package model

type EasyWkActionRequest struct {
	Password    string `form:"pwd"`
	Action      string `form:"action"`
	KeepSum     string `form:"keepsum,omitempty"`
	FirstLane   int    `form:"firstlane,omitempty"`
	LaneCount   int    `form:"lanecount,omitempty"`
	MeetingName string `form:"vername,omitempty"`
	Event       int    `form:"event,omitempty"`
	Heat        int    `form:"heat,omitempty"`
	MaxHeat     int    `form:"maxheat,omitempty"`
	EventName   string `form:"name,omitempty"`
	Lane        int    `form:"lane,omitempty"`
	Meter       string `form:"meter,omitempty"`
	Time        int    `form:"time,omitempty"`
	Finished    string `form:"finished,omitempty"`
	Content     string `form:"content,omitempty"`

	Athlete0 string `form:"srw0,omitempty"`
	Year0    string `form:"yob0,omitempty"`
	Team0    string `form:"club0,omitempty"`

	Athlete1 string `form:"srw1,omitempty"`
	Year1    string `form:"yob1,omitempty"`
	Team1    string `form:"club1,omitempty"`

	Athlete2 string `form:"srw2,omitempty"`
	Year2    string `form:"yob2,omitempty"`
	Team2    string `form:"club2,omitempty"`

	Athlete3 string `form:"srw3,omitempty"`
	Year3    string `form:"yob3,omitempty"`
	Team3    string `form:"club3,omitempty"`

	Athlete4 string `form:"srw4,omitempty"`
	Year4    string `form:"yob4,omitempty"`
	Team4    string `form:"club4,omitempty"`

	Athlete5 string `form:"srw5,omitempty"`
	Year5    string `form:"yob5,omitempty"`
	Team5    string `form:"club5,omitempty"`

	Athlete6 string `form:"srw6,omitempty"`
	Year6    string `form:"yob6,omitempty"`
	Team6    string `form:"club6,omitempty"`

	Athlete7 string `form:"srw7,omitempty"`
	Year7    string `form:"yob7,omitempty"`
	Team7    string `form:"club7,omitempty"`

	Athlete8 string `form:"srw8,omitempty"`
	Year8    string `form:"yob8,omitempty"`
	Team8    string `form:"club8,omitempty"`

	Athlete9 string `form:"srw9,omitempty"`
	Year9    string `form:"yob9,omitempty"`
	Team9    string `form:"club9,omitempty"`

	Athlete10 string `form:"srw10,omitempty"`
	Year10    string `form:"yob10,omitempty"`
	Team10    string `form:"club10,omitempty"`
}
