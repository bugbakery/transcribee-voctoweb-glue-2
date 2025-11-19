package cron

import (
	"log"
	"transcribee-voctoweb/transcribee_api"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
)

func monitorTranscriptionProgressCron(app *pocketbase.PocketBase, transcribeeApiBaseUrl string) error {
	conferenceRecords, err := app.FindAllRecords("conferences",
		dbx.HashExp{"autocreate_active": true},
	)
	if err != nil {
		return err
	}

	for _, conference := range conferenceRecords {
		transcribeeToken := conference.GetString("transcribee_user_token")

		if transcribeeToken == "" {
			continue
		}

		transcribeeApi := transcribee_api.New(transcribeeApiBaseUrl, transcribeeToken)

		talkRecords, err := app.FindRecordsByFilter(
			"talks",
			"transcribee_id != '' && transcribee_state != 'done' && transcribee_state != 'failed' && conference = {:conference}",
			"date",
			0, 0,
			dbx.Params{"conference": conference.Id},
		)
		if err != nil {
			log.Printf("Error finding talks to monitor for conference %s: %v", conference.GetString("name"), err)
			continue
		}
		log.Printf("Found %d talks for conference '%s' in flight", len(talkRecords), conference.GetString("name"))

		for _, talkRecord := range talkRecords {
			transcribeeId := talkRecord.GetString("transcribee_id")
			currentState := talkRecord.GetString("transcribee_state")

			tasks, err := transcribeeApi.GetTasksForDocument(transcribeeId)
			if err != nil {
				log.Printf("Error getting tasks for talk '%s' (transcribee ID: %s): %v",
					talkRecord.GetString("title"), transcribeeId, err)
				continue
			}

			newState := determineTranscribeeState(tasks)

			if newState != currentState {
				log.Printf("Updating transcribee_state for talk '%s' from '%s' to '%s'",
					talkRecord.GetString("title"), currentState, newState)

				talkRecord.Set("transcribee_state", newState)
				err = app.Save(talkRecord)
				if err != nil {
					log.Printf("Error updating transcribee_state for talk '%s': %v",
						talkRecord.GetString("title"), err)
				}
			}
		}
	}

	return nil
}

func determineTranscribeeState(tasks []transcribee_api.TaskResponse) string {
	// return "done" if the document is transcribed and speaker diarization is complete
	// return "in_progress" if the document is being transcribed or speaker diarization is in progress
	// return "creating" if the transcription task is not yet created

	if len(tasks) == 0 {
		return "created"
	}

	var transcribeTask *transcribee_api.TaskResponse
	var speakerIdentificationTask *transcribee_api.TaskResponse

	// Find the transcribe and speaker identification tasks
	for i := range tasks {
		task := &tasks[i]
		if task.TaskType == transcribee_api.TaskTypeModelTRANSCRIBE {
			transcribeTask = task
		} else if task.TaskType == transcribee_api.TaskTypeModelIDENTIFY_SPEAKERS {
			speakerIdentificationTask = task
		}
	}

	// If there's no transcribe task yet, we're still creating
	if transcribeTask == nil {
		return "created"
	}

	// Check if any task has failed
	for _, task := range tasks {
		if task.State == transcribee_api.TaskStateFAILED {
			return "failed"
		}
	}

	// If transcription is not completed yet, we're in progress
	if transcribeTask.State != transcribee_api.TaskStateCOMPLETED {
		return "in_progress"
	}

	// If there's no speaker identification task, but transcription is done, we're done
	if speakerIdentificationTask == nil {
		return "done"
	}

	// If speaker identification is not completed yet, we're still in progress
	if speakerIdentificationTask.State != transcribee_api.TaskStateCOMPLETED {
		return "in_progress"
	}

	// Both transcription and speaker identification are complete
	return "done"
}

func RegisterMonitorTranscriptionProgressCron(app *pocketbase.PocketBase, transcribeeApiBaseUrl string) error {
	job_finished := true
	app.Cron().MustAdd("monitor_transcription_progress", "* * * * *", func() {
		if !job_finished {
			log.Println("monitor_transcription_progress Job already running, skipping")
			return
		}
		job_finished = false
		err := monitorTranscriptionProgressCron(app, transcribeeApiBaseUrl)
		if err != nil {
			log.Println(err)
		}
		job_finished = true
	})

	return nil
}
