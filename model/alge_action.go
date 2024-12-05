package model

type AlgeActionRequest struct {
	Password string `json:"password"`
	Action   string `json:"action"` // START, LANE, STOP
	Event    int    `json:"event,omitempty"`
	Heat     int    `json:"heat,omitempty"`
	Lane     int    `json:"lane,omitempty"`
	Meter    int    `json:"meter,omitempty"`
	Time     int    `json:"time,omitempty"`
	Finished string `json:"finished,omitempty"`
}
