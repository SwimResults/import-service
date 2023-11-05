package importer

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/swimresults/import-service/model"
	startModel "github.com/swimresults/start-service/model"
	"os"
	"time"
)

var CurrentMeeting model.EasyWkMeeting

func SetEasyWkMeeting() {
	dat, _ := os.ReadFile("./config/live_meeting.json")
	err := json.Unmarshal(dat, &CurrentMeeting)
	if err != nil {
		println(err.Error())
		return
	}
	fmt.Printf("set meeting for live services to: '%s'; password is: '%s'\n", CurrentMeeting.Meeting, CurrentMeeting.Password)
}

func SetHeatStartTime(event int, heat int) error {
	if CurrentMeeting.Meeting == "" {
		return errors.New("no meeting for live services declared")
	}

	return SetHeatTime(event, heat, time.Now(), time.Time{})
}

func SetHeatFinishTime(event int, heat int) error {
	if CurrentMeeting.Meeting == "" {
		return errors.New("no meeting for live services declared")
	}

	return SetHeatTime(event, heat, time.Time{}, time.Now())
}

func SetHeatTime(event int, heatNumber int, startAt time.Time, finishedAt time.Time) error {
	heat := startModel.Heat{
		Meeting:    CurrentMeeting.Meeting,
		Event:      event,
		Number:     heatNumber,
		StartAt:    startAt,
		FinishedAt: finishedAt,
	}

	_, _, err := hc.ImportHeat(heat)
	return err
}

func ImportResult(event int, heat int, lane int, time time.Duration, meter int, reaction bool, finished bool) error {
	if CurrentMeeting.Meeting == "" {
		return errors.New("no meeting for live services declared")
	}

	start := startModel.Start{
		Meeting:    CurrentMeeting.Meeting,
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
