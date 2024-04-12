package model

type EasyWkActionRequest struct {
	Password    string `json:"pwd"`
	Action      string `json:"action"`
	KeepSum     string `json:"keepsum,omitempty"`
	FirstLane   int    `json:"firstlane,omitempty"`
	LaneCount   int    `json:"lanecount,omitempty"`
	MeetingName string `json:"vername,omitempty"`
	Event       int    `json:"event,omitempty"`
	Heat        int    `json:"heat,omitempty"`
	MaxHeat     int    `json:"maxheat,omitempty"`
	EventName   string `json:"name,omitempty"`
	Lane        int    `json:"lane,omitempty"`
	Meter       string `json:"meter,omitempty"`
	Time        int    `json:"time,omitempty"`
	Finished    string `json:"finished,omitempty"`
	Content     string `json:"content,omitempty"`

	Athlete0 string `json:"srw0,omitempty"`
	Year0    string `json:"yob0,omitempty"`
	Team0    string `json:"club0,omitempty"`

	Athlete1 string `json:"srw1,omitempty"`
	Year1    string `json:"yob1,omitempty"`
	Team1    string `json:"club1,omitempty"`

	Athlete2 string `json:"srw2,omitempty"`
	Year2    string `json:"yob2,omitempty"`
	Team2    string `json:"club2,omitempty"`

	Athlete3 string `json:"srw3,omitempty"`
	Year3    string `json:"yob3,omitempty"`
	Team3    string `json:"club3,omitempty"`

	Athlete4 string `json:"srw4,omitempty"`
	Year4    string `json:"yob4,omitempty"`
	Team4    string `json:"club4,omitempty"`

	Athlete5 string `json:"srw5,omitempty"`
	Year5    string `json:"yob5,omitempty"`
	Team5    string `json:"club5,omitempty"`

	Athlete6 string `json:"srw6,omitempty"`
	Year6    string `json:"yob6,omitempty"`
	Team6    string `json:"club6,omitempty"`

	Athlete7 string `json:"srw7,omitempty"`
	Year7    string `json:"yob7,omitempty"`
	Team7    string `json:"club7,omitempty"`

	Athlete8 string `json:"srw8,omitempty"`
	Year8    string `json:"yob8,omitempty"`
	Team8    string `json:"club8,omitempty"`

	Athlete9 string `json:"srw9,omitempty"`
	Year9    string `json:"yob9,omitempty"`
	Team9    string `json:"club9,omitempty"`

	Athlete10 string `json:"srw10,omitempty"`
	Year10    string `json:"yob10,omitempty"`
	Team10    string `json:"club10,omitempty"`
}
