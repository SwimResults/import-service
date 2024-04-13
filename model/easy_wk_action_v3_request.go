package model

type EasyWkActionV3Request struct {
	Password    string `json:"pwd" form:"pwd"`
	Action      string `json:"action" form:"action"`
	KeepSum     string `json:"keepsum,omitempty" form:"keepsum,omitempty"`
	FirstLane   int    `json:"firstlane,omitempty" form:"firstlane,omitempty"`
	LaneCount   int    `json:"lanecount,omitempty" form:"lanecount,omitempty"`
	MeetingName string `json:"vername,omitempty" form:"vername,omitempty"`
	Event       int    `json:"event,omitempty" form:"event,omitempty"`
	Heat        int    `json:"heat,omitempty" form:"heat,omitempty"`
	MaxHeat     int    `json:"maxheat,omitempty" form:"maxheat,omitempty"`
	EventName   string `json:"name,omitempty" form:"name,omitempty"`
	Lane        int    `json:"lane,omitempty" form:"lane,omitempty"`
	Meter       string `json:"meter,omitempty" form:"meter,omitempty"`
	Time        int    `json:"time,omitempty" form:"time,omitempty"`
	Finished    string `json:"finished,omitempty" form:"finished,omitempty"`
	Content     string `json:"content,omitempty" form:"content,omitempty"`

	Athlete0 string `json:"srw0,omitempty" form:"srw0,omitempty"`
	Year0    string `json:"yob0,omitempty" form:"yob0,omitempty"`
	Team0    string `json:"club0,omitempty" form:"club0,omitempty"`

	Athlete1 string `json:"srw1,omitempty" form:"srw1,omitempty"`
	Year1    string `json:"yob1,omitempty" form:"yob1,omitempty"`
	Team1    string `json:"club1,omitempty" form:"club1,omitempty"`

	Athlete2 string `json:"srw2,omitempty" form:"srw2,omitempty"`
	Year2    string `json:"yob2,omitempty" form:"yob2,omitempty"`
	Team2    string `json:"club2,omitempty" form:"club2,omitempty"`

	Athlete3 string `json:"srw3,omitempty" form:"srw3,omitempty"`
	Year3    string `json:"yob3,omitempty" form:"yob3,omitempty"`
	Team3    string `json:"club3,omitempty" form:"club3,omitempty"`

	Athlete4 string `json:"srw4,omitempty" form:"srw4,omitempty"`
	Year4    string `json:"yob4,omitempty" form:"yob4,omitempty"`
	Team4    string `json:"club4,omitempty" form:"club4,omitempty"`

	Athlete5 string `json:"srw5,omitempty" form:"srw5,omitempty"`
	Year5    string `json:"yob5,omitempty" form:"yob5,omitempty"`
	Team5    string `json:"club5,omitempty" form:"club5,omitempty"`

	Athlete6 string `json:"srw6,omitempty" form:"srw6,omitempty"`
	Year6    string `json:"yob6,omitempty" form:"yob6,omitempty"`
	Team6    string `json:"club6,omitempty" form:"club6,omitempty"`

	Athlete7 string `json:"srw7,omitempty" form:"srw7,omitempty"`
	Year7    string `json:"yob7,omitempty" form:"yob7,omitempty"`
	Team7    string `json:"club7,omitempty" form:"club7,omitempty"`

	Athlete8 string `json:"srw8,omitempty" form:"srw8,omitempty"`
	Year8    string `json:"yob8,omitempty" form:"yob8,omitempty"`
	Team8    string `json:"club8,omitempty" form:"club8,omitempty"`

	Athlete9 string `json:"srw9,omitempty" form:"srw9,omitempty"`
	Year9    string `json:"yob9,omitempty" form:"yob9,omitempty"`
	Team9    string `json:"club9,omitempty" form:"club9,omitempty"`

	Athlete10 string `json:"srw10,omitempty" form:"srw10,omitempty"`
	Year10    string `json:"yob10,omitempty" form:"yob10,omitempty"`
	Team10    string `json:"club10,omitempty" form:"club10,omitempty"`
}
