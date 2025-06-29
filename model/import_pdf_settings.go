package model

type ImportPdfStartListSettings struct {
	OmitFirst              []string     `json:"omit_first,omitempty" bson:"omit_first,omitempty"`                               // strings, that are remove before starting the process
	EventSeparator         string       `json:"event_separator,omitempty" bson:"event_separator,omitempty"`                     // Wettkampf
	EventSkipStrings       []string     `json:"event_skip_strings,omitempty" bson:"event_skip_strings,omitempty"`               //
	EventRequiredStrings   []string     `json:"event_required_strings,omitempty" bson:"event_required_strings,omitempty"`       // Meldezeit
	EventNumberSeparator   string       `json:"event_number_separator,omitempty" bson:"event_number_separator,omitempty"`       // -
	DistanceSeparator      string       `json:"distance_separator,omitempty" bson:"distance_separator,omitempty"`               // m
	GenderMapping          [3][2]string `json:"gender_mapping,omitempty" bson:"gender_mapping,omitempty"`                       // "männlich" -> "MALE"
	StyleNameSkipStrings   []string     `json:"style_name_skip_strings,omitempty" bson:"style_name_skip_strings,omitempty"`     // staffel, beine
	HeatSeparator          string       `json:"heat_separator,omitempty" bson:"heat_separator,omitempty"`                       // Lauf
	HeatSkipStrings        []string     `json:"heat_skip_strings,omitempty" bson:"heat_skip_strings,omitempty"`                 // Finalabschnittes
	HeatRequiredStrings    []string     `json:"heat_required_strings,omitempty" bson:"heat_required_strings,omitempty"`         // Uhr, Meldezeit
	HeatNumberSeparator    string       `json:"heat_number_separator,omitempty" bson:"heat_number_separator,omitempty"`         // /
	HeatHasNoTime          bool         `json:"heat_has_time,omitempty" bson:"heat_has_time,omitempty"`                         // false
	HeatTimeLeftSeparator  string       `json:"heat_time_left_separator,omitempty" bson:"heat_time_left_separator,omitempty"`   // (ca.
	HeatTimeRightSeparator string       `json:"heat_time_right_separator,omitempty" bson:"heat_time_right_separator,omitempty"` // Uhr
	HeatTimeLayout         string       `json:"heat_time_layout,omitempty" bson:"heat_time_layout,omitempty"`                   // 15:04
	LaneSeparator          string       `json:"lane_separator,omitempty" bson:"lane_separator,omitempty"`                       // Bahn
	LaneSkipStrings        []string     `json:"lane_skip_strings,omitempty" bson:"lane_skip_strings,omitempty"`                 // Meldezeit, Uhr
	LaneNumberPattern      string       `json:"lane_number_pattern,omitempty" bson:"lane_number_pattern,omitempty"`             // [0-9]+
	YearPattern            string       `json:"year_pattern,omitempty" bson:"year_pattern,omitempty"`                           // [0-9]{4}
	YearOpenString         string       `json:"year_open_string,omitempty" bson:"year_open_string,omitempty"`                   // Offen
	SwimTimePattern        string       `json:"swim_time_pattern,omitempty" bson:"swim_time_pattern,omitempty"`                 // [0-9]{2}:[0-9]{2},[0-9]{2}
}

type ImportPdfResultListSettings struct {
	OmitFirst                     []string     `json:"omit_first,omitempty" bson:"omit_first,omitempty"`                                             // strings, that are remove before starting the process
	OncePerEvent                  bool         `json:"once_per_event,omitempty" bson:"once_per_event,omitempty"`                                     // true
	EventSeparator                string       `json:"event_separator,omitempty" bson:"event_separator,omitempty"`                                   // Wettkampf
	EventSkipStrings              []string     `json:"event_skip_strings,omitempty" bson:"event_skip_strings,omitempty"`                             //
	EventRequiredStrings          []string     `json:"event_required_strings,omitempty" bson:"event_required_strings,omitempty"`                     // Meldezeit
	EventNumberSeparator          string       `json:"event_number_separator,omitempty" bson:"event_number_separator,omitempty"`                     // -
	EventMultipleNumbersSeparator string       `json:"event_multiple_numbers_separator,omitempty" bson:"event_multiple_numbers_separator,omitempty"` // / for "05/205" TODO: not implemented
	DistanceSeparator             string       `json:"distance_separator,omitempty" bson:"distance_separator,omitempty"`                             // m
	GenderMapping                 [3][2]string `json:"gender_mapping,omitempty" bson:"gender_mapping,omitempty"`                                     // "männlich" -> "MALE"
	StyleNameSkipStrings          []string     `json:"style_name_skip_strings,omitempty" bson:"style_name_skip_strings,omitempty"`                   // staffel, beine
	RatingSeparators              []string     `json:"rating_separators,omitempty" bson:"rating_separators,omitempty"`                               // Jahrgang, Jahrgänge, Offene Wertung
	RatingRightSeparators         []string     `json:"rating_right_separators,omitempty" bson:"rating_right_separators,omitempty"`                   // Jahrgang, Jahrgänge, Offene Wertung
	ResultSeparator               string       `json:"result_separator,omitempty" bson:"result_separator,omitempty"`                                 // Endzeit
	DisqualificationSeparator     string       `json:"disqualification_separator,omitempty" bson:"disqualification_separator,omitempty"`             // disqualifiziert
	DisqualificationTimePattern   string       `json:"disqualification_time_pattern,omitempty" bson:"disqualification_time_pattern,omitempty"`       // [0-9]{2}:[0-9]{2}
	DnsSeparator                  string       `json:"dns_separator,omitempty" bson:"dns_separator,omitempty"`                                       // nicht am Start
	LoggedOutSeparator            string       `json:"logged_out_separator,omitempty" bson:"logged_out_separator,omitempty"`                         // abgemeldet
	ResultEndCutStrings           []string     `json:"result_end_cut_strings,omitempty" bson:"result_end_cut_strings,omitempty"`                     // disqualifiziert, abgemeldet, nicht am Start, erzeugt mit EasyWK
	ResultPattern                 string       `json:"result_pattern,omitempty" bson:"result_pattern,omitempty"`                                     // [0-9]*\..*[0-9]{4}.*[0-9]{2}:[0-9]{2},[0-9]{2}
	YearPattern                   string       `json:"year_pattern,omitempty" bson:"year_pattern,omitempty"`                                         // [0-9]{4}
	YearOpenString                string       `json:"year_open_string,omitempty" bson:"year_open_string,omitempty"`                                 // Offen
	AthleteReplaceStrings         []string     `json:"athlete_replace_strings,omitempty" bson:"athlete_replace_strings,omitempty"`                   // Bezirksmeister
	SwimTimePattern               string       `json:"swim_time_pattern,omitempty" bson:"swim_time_pattern,omitempty"`                               // [0-9]{2}:[0-9]{2},[0-9]{2}
	ReasonRightSeparator          string       `json:"reason_right_separator,omitempty" bson:"reason_right_separator,omitempty"`                     // Uhrzeit der Bekanntgabe
}
