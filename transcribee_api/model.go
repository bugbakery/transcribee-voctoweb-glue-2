package transcribee_api

import "time"

type TaskType string

const (
	TaskTypeALIGN TaskType = "ALIGN"
)

type AlignTask struct {
	DocumentID     string                 `json:"document_id"`
	TaskParameters map[string]interface{} `json:"task_parameters"`
	TaskType       *TaskType              `json:"task_type,omitempty"`
}

type BodyAddMediaFileApiV1DocumentsDocumentIdAddMediaFilePost struct {
	File []byte   `json:"file"`
	Tags []string `json:"tags"`
}

type BodyCreateDocumentApiV1DocumentsPost struct {
	File             []byte `json:"file"`
	Language         string `json:"language"`
	Model            string `json:"model"`
	Name             string `json:"name"`
	NumberOfSpeakers *int   `json:"number_of_speakers,omitempty"`
}

type BodyImportDocumentApiV1DocumentsImportPost struct {
	MediaFile []byte `json:"media_file"`
	Name      string `json:"name"`
}

type ChangePasswordRequest struct {
	NewPassword string `json:"new_password"`
	OldPassword string `json:"old_password"`
}

type CreateShareToken struct {
	CanWrite   bool       `json:"can_write"`
	Name       string     `json:"name"`
	ValidUntil *time.Time `json:"valid_until,omitempty"`
}

type CreateUser struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

type CreateWorker struct {
	Name string `json:"name"`
}

type DeactivateWorker struct {
	ID string `json:"id"`
}

type DocumentMedia struct {
	ContentType string   `json:"content_type"`
	Tags        []string `json:"tags"`
	Url         string   `json:"url"`
}

type DocumentShareTokenBase struct {
	CanWrite   bool       `json:"can_write"`
	DocumentID string     `json:"document_id"`
	ID         string     `json:"id"`
	Name       string     `json:"name"`
	Token      string     `json:"token"`
	ValidUntil *time.Time `json:"valid_until,omitempty"`
}

type DocumentUpdateRequest struct {
	Name *string `json:"name,omitempty"`
}

type DocumentWithAccessInfo struct {
	CanWrite      bool            `json:"can_write"`
	ChangedAt     string          `json:"changed_at"`
	CreatedAt     string          `json:"created_at"`
	HasFullAccess bool            `json:"has_full_access"`
	ID            string          `json:"id"`
	MediaFiles    []DocumentMedia `json:"media_files"`
	Name          string          `json:"name"`
}

type ExportError struct {
	Error string `json:"error"`
}

type ExportFormat string

const (
	ExportFormatVTT ExportFormat = "VTT"
	ExportFormatSRT ExportFormat = "SRT"
)

type ExportResult struct {
	Result string `json:"result"`
}

type TaskType1 string

const (
	TaskType1EXPORT TaskType1 = "EXPORT"
)

type ExportTaskParameters struct {
	Format              ExportFormat `json:"format"`
	IncludeSpeakerNames bool         `json:"include_speaker_names"`
	IncludeWordTiming   bool         `json:"include_word_timing"`
	MaxLineLength       *int         `json:"max_line_length,omitempty"`
}

type KeepaliveBody struct {
	Progress *float64 `json:"progress,omitempty"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type ModelConfig struct {
	ID        string   `json:"id"`
	Languages []string `json:"languages"`
	Name      string   `json:"name"`
}

type PageConfig struct {
	FooterPosition *int   `json:"footer_position,omitempty"`
	Name           string `json:"name"`
	Text           string `json:"text"`
}

type PublicConfig struct {
	LoggedOutRedirectUrl *string                `json:"logged_out_redirect_url,omitempty"`
	Models               map[string]ModelConfig `json:"models"`
}

type SetDurationRequest struct {
	Duration float64 `json:"duration"`
}

type ShortPageConfig struct {
	FooterPosition *int   `json:"footer_position,omitempty"`
	Name           string `json:"name"`
}

type TaskType2 string

const (
	TaskType2IDENTIFY_SPEAKERS TaskType2 = "IDENTIFY_SPEAKERS"
)

type SpeakerIdentificationTask struct {
	DocumentID     string                 `json:"document_id"`
	TaskParameters map[string]interface{} `json:"task_parameters"`
	TaskType       *TaskType2             `json:"task_type,omitempty"`
}

type TaskAttemptResponse struct {
	Progress *float64 `json:"progress,omitempty"`
}

type TaskState string

const (
	TaskStateNEW       TaskState = "NEW"
	TaskStateASSIGNED  TaskState = "ASSIGNED"
	TaskStateCOMPLETED TaskState = "COMPLETED"
	TaskStateFAILED    TaskState = "FAILED"
)

type TaskTypeModel string

const (
	TaskTypeModelIDENTIFY_SPEAKERS TaskTypeModel = "IDENTIFY_SPEAKERS"
	TaskTypeModelTRANSCRIBE        TaskTypeModel = "TRANSCRIBE"
	TaskTypeModelALIGN             TaskTypeModel = "ALIGN"
	TaskTypeModelREENCODE          TaskTypeModel = "REENCODE"
	TaskTypeModelEXPORT            TaskTypeModel = "EXPORT"
)

type TaskType3 string

const (
	TaskType3TRANSCRIBE TaskType3 = "TRANSCRIBE"
)

type TranscribeTaskParameters struct {
	Lang  string `json:"lang"`
	Model string `json:"model"`
}

type UnknownTask struct {
	DocumentID     string                 `json:"document_id"`
	TaskParameters map[string]interface{} `json:"task_parameters"`
	TaskType       string                 `json:"task_type"`
}

type UserBase struct {
	Username string `json:"username"`
}

type ValidationErrorModel struct {
	Loc  []interface{} `json:"loc"`
	Msg  string        `json:"msg"`
	Type string        `json:"type"`
}

type Worker struct {
	DeactivatedAt *time.Time `json:"deactivated_at,omitempty"`
	ID            *string    `json:"id,omitempty"`
	LastSeen      *time.Time `json:"last_seen,omitempty"`
	Name          string     `json:"name"`
	Token         string     `json:"token"`
}

type WorkerWithId struct {
	DeactivatedAt *time.Time `json:"deactivated_at,omitempty"`
	ID            *string    `json:"id,omitempty"`
	LastSeen      *time.Time `json:"last_seen,omitempty"`
	Name          string     `json:"name"`
}

type Document struct {
	ChangedAt  string          `json:"changed_at"`
	CreatedAt  string          `json:"created_at"`
	ID         string          `json:"id"`
	MediaFiles []DocumentMedia `json:"media_files"`
	Name       string          `json:"name"`
}

type ExportTask struct {
	DocumentID     string               `json:"document_id"`
	TaskParameters ExportTaskParameters `json:"task_parameters"`
	TaskType       *TaskType1           `json:"task_type,omitempty"`
}

type HTTPValidationError struct {
	Detail *[]ValidationErrorModel `json:"detail,omitempty"`
}

type TaskQueueInfoTaskEntry struct {
	ID            string        `json:"id"`
	RemainingCost float64       `json:"remaining_cost"`
	State         TaskState     `json:"state"`
	TaskType      TaskTypeModel `json:"task_type"`
}

type TaskResponse struct {
	CurrentAttempt *TaskAttemptResponse   `json:"current_attempt,omitempty"`
	Dependencies   []string               `json:"dependencies"`
	DocumentID     string                 `json:"document_id"`
	ID             string                 `json:"id"`
	State          TaskState              `json:"state"`
	TaskParameters map[string]interface{} `json:"task_parameters"`
	TaskType       TaskTypeModel          `json:"task_type"`
}

type TranscribeTask struct {
	DocumentID     string                   `json:"document_id"`
	TaskParameters TranscribeTaskParameters `json:"task_parameters"`
	TaskType       *TaskType3               `json:"task_type,omitempty"`
}

type ApiDocumentWithTasks struct {
	ChangedAt  string          `json:"changed_at"`
	CreatedAt  string          `json:"created_at"`
	ID         string          `json:"id"`
	MediaFiles []DocumentMedia `json:"media_files"`
	Name       string          `json:"name"`
	Tasks      []TaskResponse  `json:"tasks"`
}

type AssignedTaskResponse struct {
	CurrentAttempt *TaskAttemptResponse   `json:"current_attempt,omitempty"`
	Dependencies   []string               `json:"dependencies"`
	Document       Document               `json:"document"`
	DocumentID     string                 `json:"document_id"`
	ID             string                 `json:"id"`
	State          TaskState              `json:"state"`
	TaskParameters map[string]interface{} `json:"task_parameters"`
	TaskType       TaskTypeModel          `json:"task_type"`
}

type TaskQueueInfoResponse struct {
	OpenTasks []TaskQueueInfoTaskEntry `json:"open_tasks"`
}
