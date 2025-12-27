package hooks

import (
	"slices"
	"strings"
	"net/mail"
	"transcribee-voctoweb/utils"

	"github.com/pocketbase/pocketbase/tools/mailer"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/routine"
)

// Todo: Make the text variable
const mailTemplate = `<div style="font-family: Helvetica, Arial, sans-serif; font-size:12px;">
				<p>Your account has been created.</p>
				<p>Your username is:</p>
				<p>{{username}}</p>
				<p>This is your password:</p>
				<p>{{password}}</p>
				<p>Please login here: <a href="{{appUrl}}">{{appUrl}}</a>
				<p>Your 39c3 subtitles team</p>
			</div>`


func CreateUsers(e *core.RecordEvent) error {
	app := e.App

	record := e.Record
	filename := record.GetString("filename")
	logger := app.Logger().With("process", "userCreationFromFile", "filename", filename)

	logger.Info("Processing user file")

	recordTargetCollection, err := app.FindCollectionByNameOrId("users")
	if err != nil || recordTargetCollection == nil {
		logger.Error("Failed to find users collection", "error", err)
		return err
	}


	// Read from file
	fileKey := record.BaseFilesPath() + "/" + record.GetString("file")
	records, header, err := utils.ReadCsv(app, fileKey, ',')
	if err != nil {
		logger.Error("Failed to read CSV file", "error", err)
		return err
	}

	// Check the file for the correct columns
	expectedColumns := []string{"nick", " mail"}
	if len(header) < len(expectedColumns) || !slices.Equal(expectedColumns, header) {
		logger.Error("CSV file does not have the expected columns", "expected", expectedColumns, "got", header)
		return err
	}


	routine.FireAndForget(func() {

		success, failed := 0, 0
		app.RunInTransaction(func(txApp core.App) error {
			for _, csvRow := range records {
				username := strings.TrimSpace(csvRow[0])
				usermail := strings.TrimSpace(csvRow[1])

				userLogger := logger.With("username", username, "mail", usermail)

				if usermail == "" || username == "" {
					userLogger.Error("Username or mail not valid")
					continue
				}

				password := utils.GeneratePassword(10, true, true, true)

				userRecord, err := txApp.FindFirstRecordByData(recordTargetCollection.Id, "username", username)
				if err != nil || userRecord == nil {
					// Record not found, create a new one
					userRecord = core.NewRecord(recordTargetCollection)
				}
				userRecord.Set("username", username)
				userRecord.SetEmail(usermail)
				userRecord.SetPassword(password)

				userLogger.Info("Sending mail for user", "address", usermail)

				mailText := mailTemplate
				mailText = strings.ReplaceAll(mailText, "{{username}}", username)
				mailText = strings.ReplaceAll(mailText, "{{password}}", password)
				mailText = strings.ReplaceAll(mailText, "{{appUrl}}", txApp.Settings().Meta.AppURL)

				// Send mail first, we just want to save the user if we can send mail
				message := &mailer.Message{
					From: mail.Address{
						Address: txApp.Settings().Meta.SenderAddress,
						Name:    txApp.Settings().Meta.SenderName,
					},
					To:      []mail.Address{{Address: usermail}},
					Subject: "Subtitles account created",
					HTML:    mailText,
				}

				err = txApp.NewMailClient().Send(message)
				if err != nil {
					userLogger.Error("Sending mail to user failed", "mailaddress", usermail, "error", err)
					continue
				}


				// Mail success, save
				err = txApp.Save(userRecord)
				if err != nil {
					userLogger.Error("Saving user failed", "error", err)
					failed++
					continue
				}
				success++

			}

			logger.Info("Success! Added users", "success", success, "failed", failed)
			return nil
		})

	})

	return e.Next()
}
