package model

type ImportPdfStartListSettings struct {
	EventSeparator     string   `json:"event_separator,omitempty" bson:"event_separator,omitempty"`
	EventSkipStrings   []string `json:"event_skip_strings,omitempty" bson:"event_skip_strings,omitempty"`
	EventNoSkipStrings []string `json:"event_no_skip_strings,omitempty" bson:"event_no_skip_strings,omitempty"`
	HeatSeparator      string   `json:"heat_separator,omitempty" bson:"heat_separator,omitempty"`
	HeatSkipStrings    []string `json:"heat_skip_strings,omitempty" bson:"heat_skip_strings,omitempty"`
	HeatNoSkipStrings  []string `json:"heat_no_skip_strings,omitempty" bson:"heat_no_skip_strings,omitempty"`
}

type ImportPdfResultListSettings struct {
	EventSeparator     string   `json:"event_separator,omitempty" bson:"event_separator,omitempty"`
	EventSkipStrings   []string `json:"event_skip_strings,omitempty" bson:"event_skip_strings,omitempty"`
	EventNoSkipStrings []string `json:"event_no_skip_strings,omitempty" bson:"event_no_skip_strings,omitempty"`
}
