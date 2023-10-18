package service

import (
	"fmt"
	"github.com/swimresults/import-service/importer"
	"github.com/swimresults/import-service/model"
	"strconv"
	"strings"
	"time"
)

var currentEvent int
var currentHeat int

func EasyWkLivetimingRequest(request model.EasyWkActionRequest) (string, error) {
	if request.Password != importer.CurrentMeeting.Password || request.Action == "" {
		fmt.Printf("password or no action error with password '%s' and action '%s'\n", request.Password, request.Action)
		return "ERROR: Passwort nicht korrekt oder keine Aktion definiert", nil
	}

	var err error

	switch request.Action {
	case "ping", "clearsum", "init", "disq", "text":
		return "OK", nil
	case "newrace", "ready":
		// store current event and heat for later timings import
		currentEvent = request.Event
		currentHeat = request.Heat

		// start heat by setting start time
		err = importer.SetHeatStartTime(currentEvent, currentHeat)

		return "OK", err
	case "time":
		// import result time
		if currentEvent == 0 || currentHeat == 0 {
			return "OK", fmt.Errorf("[EasyWk Time Import] no event or heat set, skipping import")
		}

		if request.Meter == "" {
			return "OK", fmt.Errorf("[EasyWk Time Import] meter not set for E: %d, H: %d, L: %d", currentEvent, currentHeat, request.Lane)
		}

		reaction := false
		m := 0

		if request.Meter == "RT" {
			reaction = true
		} else {
			mStr := strings.ReplaceAll(request.Meter, "m", "")
			mStr = strings.Trim(mStr, " ")
			m, err = strconv.Atoi(mStr)
			if err != nil {
				return "OK", fmt.Errorf("[EasyWk Time Import] meter to int conversion failed for: %s", mStr)
			}
		}

		var t time.Duration
		var err2 error

		if request.Meter == "RT" {
			t, err2 = EasyWkReactionToDuration(request.Time)
		} else {
			t, err2 = EasyWkTimeToDuration(request.Time)
		}

		if err2 != nil {
			return "OK", fmt.Errorf("[EasyWk Time Import] time to duration conversion failed for: %d", request.Time)
		}

		err = importer.ImportResult(currentEvent, currentHeat, request.Lane, t, m, reaction, request.Finished == "yes")

		return "OK", err
	case "raceresult":
		// set heat to finished
		err = importer.SetHeatFinishTime(currentEvent, currentHeat)
		return "OK", err
	default:
		return "ERROR: Unbekannte Aktion", nil
	}
}

func EasyWkTimeToDuration(t int) (time.Duration, error) {
	tStr := fmt.Sprintf("%08d", t)
	h, _ := strconv.Atoi(tStr[6:8])
	fmt.Println(
		tStr[0:2] + "h" +
			tStr[2:4] + "m" +
			tStr[4:6] + "s" +
			fmt.Sprintf("%03d", h*10) + "ms")
	d, err := time.ParseDuration(
		tStr[0:2] + "h" +
			tStr[2:4] + "m" +
			tStr[4:6] + "s" +
			fmt.Sprintf("%03d", h*10) + "ms")
	return d, err
}

func EasyWkReactionToDuration(t int) (time.Duration, error) {
	tStr := fmt.Sprintf("%03d", t)
	h, _ := strconv.Atoi(tStr[1:3])
	fmt.Println(
		tStr[0:1] + "s" +
			fmt.Sprintf("%03d", h*10) + "ms")
	d, err := time.ParseDuration(
		tStr[0:1] + "s" +
			fmt.Sprintf("%03d", h*10) + "ms")
	return d, err
}
