package model

type ImportFileRequest struct {
	Url           string   `json:"url" form:"url"`                       // url where file is located
	Text          string   `json:"text" form:"text"`                     // text to import
	FileExtension string   `json:"file_extension" form:"file_extension"` // file extension, either PDF or DSV or PDF_TXT
	FileType      string   `json:"file_type" form:"file_type"`           // file type: DEFINITION; START_LIST; RESULT_LIST;
	ExcludeEvents []int    `json:"exclude_events" form:"exclude_events"` // events to exclude from import as array
	IncludeEvents []int    `json:"include_events" form:"include_events"` // events to include in import process as array
	Meeting       string   `json:"meeting" form:"meeting"`               // meeting in which to import the data
	SessionID     string   `json:"session_id" form:"session_id"`         // optional session ID for progress streaming
	Features      []string `json:"features" form:"features"`             // included features are imported, missing features are ignored,
	// options: event, age_group, heat, start, result, disqualification
	// only used by lenex for now
}
