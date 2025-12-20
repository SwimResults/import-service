package importer

import (
	"encoding/json"
	"fmt"
	"github.com/swimresults/import-service/model"
	"github.com/swimresults/service-core/misc"
	startModel "github.com/swimresults/start-service/model"
	"os"
	"time"
)

var CurrentMeeting model.EasyWkMeeting

func SetEasyWkMeeting() {
	dat, err1 := os.ReadFile("config/live_meeting.json")
	if err1 != nil {
		println(err1.Error())
		return
	}
	err := json.Unmarshal(dat, &CurrentMeeting)
	if err != nil {
		println(err.Error())
		return
	}
	fmt.Printf("set meeting for live services to: '%s'; password is: '%s'\n", CurrentMeeting.Meeting, CurrentMeeting.Password)
}

func GetCurrentEasyWkMeetingId() string {
	return CurrentMeeting.Meeting
}

func SetHeatStartTime(meeting string, event int, heat int) error {
	_, err := hc.SetHeatStart(meeting, event, heat)
	return err
	//return SetHeatTime(event, heat, misc.TimeNow(), time.Time{})
}

func SetHeatFinishTime(meeting string, event int, heat int) error {
	return SetHeatTime(meeting, event, heat, time.Time{}, misc.TimeNow())
}

func SetHeatTime(meeting string, event int, heatNumber int, startAt time.Time, finishedAt time.Time) error {
	heat := startModel.Heat{
		Meeting:    meeting,
		Event:      event,
		Number:     heatNumber,
		StartAt:    startAt,
		FinishedAt: finishedAt,
	}

	_, _, err := hc.ImportHeat(heat)
	return err
}

func ImportResult(meeting string, event int, heat int, lane int, time time.Duration, meter int, reaction bool, finished bool) error {
	start := startModel.Start{
		Meeting:    meeting,
		Event:      event,
		HeatNumber: heat,
		Lane:       lane,
	}

	var rt string

	if reaction {
		rt = "reaction"
	} else if finished {
		rt = "livetiming_result"
	} else {
		rt = "lap"
	}

	result := startModel.Result{
		Time:       time,
		ResultType: rt,
	}

	if !reaction {
		result.LapMeters = meter
	}

	fmt.Printf("import result for E: %d H: %d L: %d\n", event, heat, lane)
	_, _, err := sc.ImportResult(start, result)
	return err
}
