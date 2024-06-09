package model

type ImportFileRequest struct {
	Url           string `json:"url"`            // url where file is located
	Text          string `json:"text"`           // text to import
	FileExtension string `json:"file_extension"` // file extension, either PDF or DSV or PDF_TXT
	FileType      string `json:"file_type"`      // file type: DEFINITION; START_LIST; RESULT_LIST;
	ExcludeEvents []int  `json:"exclude_events"` // events to exclude from import as array
	IncludeEvents []int  `json:"include_events"` // events to include in import process as array
	Meeting       string `json:"meeting"`        // meeting in which to import the data
}
