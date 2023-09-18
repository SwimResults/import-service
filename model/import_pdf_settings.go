package model

type ImportPdfStartListSettings struct {
	EventSeparator     string   `json:"event_separator,omitempty" bson:"event_separator,omitempty"`
	EventIgnoreStrings []string `json:"event_ignore_strings,omitempty" bson:"event_ignore_strings,omitempty"`
	HeatSeparator      string   `json:"heat_separator,omitempty" bson:"heat_separator,omitempty"`
	HeatIgnoreStrings  []string `json:"heat_ignore_strings,omitempty" bson:"heat_ignore_strings,omitempty"`
}
