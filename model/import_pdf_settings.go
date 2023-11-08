package model

type ImportPdfStartListSettings struct {
	EventSeparator         string       `json:"event_separator,omitempty" bson:"event_separator,omitempty"`                     // Wettkampf
	EventSkipStrings       []string     `json:"event_skip_strings,omitempty" bson:"event_skip_strings,omitempty"`               //
	EventRequiredStrings   []string     `json:"event_required_strings,omitempty" bson:"event_required_strings,omitempty"`       // Meldezeit
	EventNumberSeparator   string       `json:"event_number_separator,omitempty" bson:"event_number_separator,omitempty"`       // -
	DistanceSeparator      string       `json:"distance_separator,omitempty" bson:"distance_separator,omitempty"`               // m
	GenderMapping          [3][2]string `json:"gender_mapping,omitempty" bson:"gender_mapping,omitempty"`                       // "männlich" -> "MALE"
	StyleNameSkipString    []string     `json:"style_name_skip_string,omitempty" bson:"style_name_skip_string,omitempty"`       // staffel, beine
	HeatSeparator          string       `json:"heat_separator,omitempty" bson:"heat_separator,omitempty"`                       // Lauf
	HeatSkipStrings        []string     `json:"heat_skip_strings,omitempty" bson:"heat_skip_strings,omitempty"`                 // Finalabschnittes
	HeatRequiredStrings    []string     `json:"heat_required_strings,omitempty" bson:"heat_required_strings,omitempty"`         // Uhr, Meldezeit
	HeatNumberSeparator    string       `json:"heat_number_separator,omitempty" bson:"heat_number_separator,omitempty"`         // /
	HeatTimeLeftSeparator  string       `json:"heat_time_left_separator,omitempty" bson:"heat_time_left_separator,omitempty"`   // (ca.
	HeatTimeRightSeparator string       `json:"heat_time_right_separator,omitempty" bson:"heat_time_right_separator,omitempty"` // Uhr
	HeatTimeLayout         string       `json:"heat_time_layout,omitempty" bson:"heat_time_layout,omitempty"`                   // 15:04
	LaneSeparator          string       `json:"lane_separator,omitempty" bson:"lane_separator,omitempty"`                       // Bahn
	LaneSkipStrings        []string     `json:"lane_skip_strings,omitempty" bson:"lane_skip_strings,omitempty"`                 // Meldezeit, Uhr
	LaneNumberPattern      string       `json:"lane_number_pattern,omitempty" bson:"lane_number_pattern,omitempty"`             // [0-9]+
	YearPattern            string       `json:"year_pattern,omitempty" bson:"year_pattern,omitempty"`                           // [0-9]{4}
	SwimTimePattern        string       `json:"swim_time_pattern,omitempty" bson:"swim_time_pattern,omitempty"`                 // [0-9]{2}:[0-9]{2},[0-9]{2}
}

type ImportPdfResultListSettings struct {
	EventSeparator     string   `json:"event_separator,omitempty" bson:"event_separator,omitempty"`
	EventSkipStrings   []string `json:"event_skip_strings,omitempty" bson:"event_skip_strings,omitempty"`
	EventNoSkipStrings []string `json:"event_no_skip_strings,omitempty" bson:"event_no_skip_strings,omitempty"`
}
