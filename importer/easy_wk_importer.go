package importer

import (
	"errors"
	"fmt"
	startModel "github.com/swimresults/start-service/model"
	"os"
	"time"
)

var currentMeeting string

func SetEasyWkMeeting() {
	dat, _ := os.ReadFile("config/live_meeting.txt")
	currentMeeting = string(dat)
	fmt.Printf("set meeting for live services to: %s\n", currentMeeting)
}

func SetHeatStartTime(event int, heat int) error {
	if currentMeeting == "" {
		return errors.New("no meeting for live services declared")
	}
	// TODO set heat started at time

	return nil
}

func SetHeatFinishTime(event int, heat int) error {
	if currentMeeting == "" {
		return errors.New("no meeting for live services declared")
	}
	// TODO set heat finished at time

	return nil
}

func ImportResult(event int, heat int, lane int, time time.Duration, meter int, finalResult bool) error {
	if currentMeeting == "" {
		return errors.New("no meeting for live services declared")
	}

	start := startModel.Start{
		Meeting:    currentMeeting,
		Event:      event,
		HeatNumber: heat,
		Lane:       lane,
	}

	var rt string

	if finalResult {
		rt = "result_list"
	} else {
		rt = "lap"
	}

	result := startModel.Result{
		Time:       time,
		ResultType: rt,
	}

	if !finalResult {
		result.LapMeters = meter
	}

	fmt.Printf("import result for E: %d H: %d L: %d\n", event, heat, lane)
	_, _, err := sc.ImportResult(start, result)
	return err
}
