package importer

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/ledongthuc/pdf"
	athleteModel "github.com/swimresults/athlete-service/model"
	importModel "github.com/swimresults/import-service/model"
	"github.com/swimresults/meeting-service/model"
	startModel "github.com/swimresults/start-service/model"
	"io"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// ReadPdf opens pdf under given path and reads plain text to a buffer and returns content as string
func ReadPdf(path string) (string, error) {
	f, r, err := pdf.Open(path)
	defer f.Close()
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	b, err := r.GetPlainText()
	if err != nil {
		return "", err
	}

	buf.ReadFrom(b)
	return buf.String(), nil
}

// GetPdfFileContent fetches pdf from source (local path or online source) and returns content as text string
func GetPdfFileContent(file string) (string, error) {
	reader, err := GetPdfReader(file)
	if err != nil {
		return "", err
	}

	var bf bytes.Buffer
	b, err := reader.GetPlainText()
	if err != nil {
		return "", err
	}

	_, err = bf.ReadFrom(b)
	if err != nil {
		return "", err
	}

	return bf.String(), nil
}

// GetPdfReader fetches pdf from source (local path or online source) and returns pdf.Reader
func GetPdfReader(file string) (*pdf.Reader, error) {
	buf, err := getFileReader(file)
	if err != nil {
		return nil, err
	}

	buff := bytes.NewBuffer([]byte{})
	_, err = io.Copy(buff, buf)
	if err != nil {
		return nil, err
	}

	rdr := bytes.NewReader(buff.Bytes())

	return pdf.NewReader(rdr, rdr.Size())
}

// ImportPdfStartList takes the path to a pdf file that contains a start
// list. All events, teams, athletes, heats and starts will be imported.
// If exclude is set, given event numbers will be excluded from import.
// If include is set, only given event numbers will be imported.
//
// For import process details see documentation on GitHub.
func ImportPdfStartList(file string, meeting string, exclude []int, include []int, stg importModel.ImportPdfStartListSettings) (*importModel.ImportFileStats, error) {

	text, err := GetPdfFileContent(file)
	if err != nil {
		return nil, err
	}

	var stats importModel.ImportFileStats

	for _, omit := range stg.OmitFirst {
		text = strings.ReplaceAll(text, omit, "")
	}

	// split by event
	eventSplit := strings.Split(text, stg.EventSeparator)
	for _, eventString := range eventSplit {

		// eliminate some events
		if shouldSkip(eventString, stg.EventSkipStrings, stg.EventRequiredStrings) {
			continue
		}

		event := model.Event{
			Meeting: meeting,
		}

		eventNumberSplit := strings.SplitN(eventString, stg.EventNumberSeparator, 2)
		event.Number, err = strconv.Atoi(trim(eventNumberSplit[0]))
		if err != nil {
			return &stats, err
		}

		if event.Number <= 0 {
			continue
		}

		if !IsEventImportable(event.Number, exclude, include) {
			continue
		}

		distanceSplit := strings.SplitN(eventNumberSplit[1], stg.DistanceSeparator, 2)
		distance := trim(distanceSplit[0])
		event.Distance, err = strconv.Atoi(distance)
		if err != nil {
			if strings.Contains(distance, "x") {
				event.RelayDistance = distance
			} else {
				return &stats, err
			}
		}

		var style string
		pos := 1000
		for _, gender := range stg.GenderMapping {
			genderPos := strings.Index(distanceSplit[1], gender[0])
			if genderPos > 0 && genderPos < pos {
				pos = genderPos
				event.Gender = gender[1]
				style = trim(substr(distanceSplit[1], gender[0]))
			}
		}

		if shouldSkip(style, stg.StyleNameSkipStrings, []string{}) {
			importError(fmt.Sprintf("skipped event %d with style '%s'", event.Number, style), errors.New(""))
			continue
		}

		fmt.Printf("WK %d - %dm %s (%s)\n", event.Number, event.Distance, style, event.Gender)

		// +===========================+
		//        EVENT IMPORT
		// +===========================+

		stats.Found.Events++

		if runImport() {
			_, c, err := ec.ImportEvent(event, style, 1)
			if err != nil {
				importError(fmt.Sprintf("import event request failed for event %d!", event.Number), err)
				continue
			}
			if c {
				stats.Created.Events++
			}
			stats.Imported.Events++
		}

		// +===========================+
		//        HEAT PARSING
		// +===========================+

		heatSplit := strings.Split(eventString, stg.HeatSeparator)

		for _, heatString := range heatSplit {

			// eliminate some heats
			if shouldSkip(heatString, stg.HeatSkipStrings, stg.HeatRequiredStrings) {
				continue
			}

			heat := startModel.Heat{
				Meeting: meeting,
				Event:   event.Number,
			}

			heatNumberString := trim(substr(heatString, stg.HeatNumberSeparator))
			heat.Number, err = strconv.Atoi(heatNumberString)
			if err != nil {
				importError(fmt.Sprintf("failed to parse number for heat %d/%s!", event.Number, heatNumberString), err)
				continue
			}

			heatTimeString := trim(substrr(substr(heatString, stg.HeatTimeRightSeparator), stg.HeatTimeLeftSeparator))
			heat.StartEstimation, err = time.Parse(stg.HeatTimeLayout, heatTimeString)
			if err != nil {
				importError(fmt.Sprintf("failed to parse time for heat %d/%d!", event.Number, heat.Number), err)
				continue
			}

			fmt.Printf("\tHeat %d (%s)\n", heat.Number, heat.StartEstimation.Format("15:04"))

			// +===========================+
			//          HEAT IMPORT
			// +===========================+

			stats.Found.Heats++
			if runImport() {
				_, c, err := hc.ImportHeat(heat)
				if err != nil {
					importError(fmt.Sprintf("import heat request failed for heat %d/%d!", event.Number, heat.Number), err)
					continue
				}
				if c {
					stats.Created.Heats++
				}
				stats.Imported.Heats++
			}

			// +=================================+
			//    LANE / ATHLETE / TEAM PARSING
			// +=================================+

			laneSplit := strings.Split(heatString, stg.LaneSeparator)
			for _, laneString := range laneSplit {
				// eliminate some heats
				if shouldSkip(laneString, stg.LaneSkipStrings, []string{}) {
					continue
				}

				laneNumberRegex := regexp.MustCompile(stg.LaneNumberPattern)
				yearRegex := regexp.MustCompile(stg.YearPattern)
				swimTimeRegex := regexp.MustCompile(stg.SwimTimePattern)

				laneNumberSplit := laneNumberRegex.Split(laneString, 2)

				isOpen := false
				var yearSplit []string
				if strings.Contains(laneNumberSplit[1], stg.YearOpenString) {
					yearSplit = strings.SplitN(laneNumberSplit[1], stg.YearOpenString, 2)
					isOpen = true
				} else {
					yearSplit = yearRegex.Split(laneNumberSplit[1], 2)
				}

				swimTimeSplit := swimTimeRegex.Split(yearSplit[1], 2)

				athleteNameString := trim(yearSplit[0])
				teamNameString := trim(swimTimeSplit[0])
				laneNumberString := trim(substr(laneString, athleteNameString))
				yearString := trim(substrr(substr(laneString, yearSplit[1]), athleteNameString))
				swimTimeRestString := trim(substrr(yearSplit[1], teamNameString))

				laneNumber, err := strconv.Atoi(laneNumberString)
				if err != nil {
					importError(fmt.Sprintf("failed to parse number for lane %d/%d/%s!", event.Number, heat.Number, laneNumberString), err)
					continue
				}

				athleteYear := 0
				if !isOpen {
					athleteYear, err = strconv.Atoi(yearString)
					if err != nil {
						importError(fmt.Sprintf("failed to parse athlete year for lane %d/%d/%d!", event.Number, heat.Number, laneNumber), err)
						continue
					}
				}

				team := athleteModel.Team{
					Name: teamNameString,
				}

				// +===========================+
				//          TEAM IMPORT
				// +===========================+

				stats.Found.Teams++
				if runImport() {
					newTeam, c, err := tc.ImportTeam(team, meeting)
					if err != nil {
						importError(fmt.Sprintf("import team request failed for start %d/%d/%d and team '%s'!", event.Number, heat.Number, laneNumber, team.Name), err)
						continue
					}
					if c {
						stats.Created.Teams++
					}
					stats.Imported.Teams++

					team = *newTeam
				}

				athlete := athleteModel.Athlete{
					Name:   athleteNameString,
					Year:   athleteYear,
					Gender: event.Gender,
					Team:   team,
				}

				// +===========================+
				//        ATHLETE IMPORT
				// +===========================+
				if event.RelayDistance == "" {

					stats.Found.Athletes++
					if runImport() {
						newAthlete, c, err := ac.ImportAthlete(athlete, meeting)
						if err != nil {
							importError(fmt.Sprintf("import athlete request failed for start %d/%d/%d and athlete '%s'!", event.Number, heat.Number, laneNumber, athlete.Name), err)
							continue
						}
						if c {
							stats.Created.Athletes++
						}
						stats.Imported.Athletes++

						athlete = *newAthlete
					}
				}

				start := startModel.Start{
					Meeting:         meeting,
					Event:           event.Number,
					HeatNumber:      heat.Number,
					Lane:            laneNumber,
					Athlete:         athlete.Identifier,
					AthleteName:     athlete.Name,
					AthleteYear:     athleteYear,
					AthleteTeam:     team.Identifier,
					AthleteTeamName: team.Name,
				}

				if event.RelayDistance != "" {
					start.IsRelay = true
				}

				// +===========================+
				//         START IMPORT
				// +===========================+

				stats.Found.Starts++
				if runImport() {
					newStart, c, err := sc.ImportStart(start)
					if err != nil {
						importError(fmt.Sprintf("import start request failed for start %d/%d/%d!", event.Number, heat.Number, start.Lane), err)
						continue
					}
					if c {
						stats.Created.Starts++
					}
					stats.Imported.Starts++

					start = *newStart
				}

				swimTimeString := swimTimeRestString[:8]
				dur, err := swimTimeToDuration(swimTimeString)
				if err != nil {
					importError(fmt.Sprintf("failed to parse duration for start %d/%d/%d with content '%s'", event.Number, heat.Number, start.Lane, swimTimeString), err)
					continue
				}

				result := startModel.Result{
					Time:       dur,
					ResultType: "registration",
				}

				// +===========================+
				//        RESULT IMPORT
				// +===========================+

				stats.Found.Results++
				if runImport() {
					_, c, err := sc.ImportResult(start, result)
					if err != nil {
						importError(fmt.Sprintf("import result request failed for start %d/%d/%d!", event.Number, heat.Number, start.Lane), err)
						continue
					}
					if c {
						stats.Created.Results++
					}
					stats.Imported.Results++
				}

				fmt.Printf("\t\tLane %d - %s (%d) %s -> %s\n", start.Lane, start.AthleteName, start.AthleteYear, start.AthleteTeamName, dur.String())
			}
		}
	}
	return &stats, nil
}

type importedKey struct {
	event   int
	athlete string
}

// ImportPdfResultList takes the path to a pdf file that contains a result
// list. All teams, athletes and results will be imported.
// If exclude is set, given event numbers will be excluded from import.
// If include is set, only given event numbers will be imported.
//
// For import process details see documentation on GitHub.
func ImportPdfResultList(file string, meeting string, exclude []int, include []int, stg importModel.ImportPdfResultListSettings) (*importModel.ImportFileStats, error) {

	text, err := GetPdfFileContent(file)
	if err != nil {
		return nil, err
	}

	teams, _, err := tc.GetTeamsByMeeting(meeting)
	if err != nil {
		return nil, err
	}

	var stats importModel.ImportFileStats

	importedAthletes := make(map[importedKey]bool)

	lastEvent := 0
	results := make(map[int][]string)
	disqualifications := make(map[int][]string)
	var rankCount int
	var rankRepetition int

	// split by event
	eventSplit := strings.Split(text, stg.EventSeparator)
	for _, eventString := range eventSplit {

		// eliminate some events
		if shouldSkip(eventString, stg.EventSkipStrings, stg.EventRequiredStrings) {
			continue
		}

		event := model.Event{
			Meeting: meeting,
		}

		eventNumberSplit := strings.SplitN(eventString, stg.EventNumberSeparator, 2)

		eventNumberString := trim(eventNumberSplit[0])

		event.Number, err = strconv.Atoi(eventNumberString)
		if err != nil {
			importWarning(fmt.Sprintf("event number is not a number, try to separate multiple event numbers using '%s' at '%s'", stg.EventMultipleNumbersSeparator, eventNumberString))
			importError("ATTENTION! This feature is not implemented yet!", nil)
			continue
		}

		if event.Number <= 0 {
			continue
		}

		if !IsEventImportable(event.Number, exclude, include) {
			continue
		}

		distanceSplit := strings.SplitN(eventNumberSplit[1], stg.DistanceSeparator, 2)
		distance := trim(distanceSplit[0])
		event.Distance, err = strconv.Atoi(distance)
		if err != nil {
			if strings.Contains(distance, "x") {
				event.RelayDistance = distance
			} else {
				return &stats, err
			}
		}

		var style string
		pos := 1000
		for _, gender := range stg.GenderMapping {
			genderPos := strings.Index(distanceSplit[1], gender[0])
			if genderPos > 0 && genderPos < pos {
				pos = genderPos
				event.Gender = gender[1]
				style = trim(substr(distanceSplit[1], gender[0]))
			}
		}

		if shouldSkip(style, stg.StyleNameSkipStrings, []string{}) {
			importError(fmt.Sprintf("skipped event %d with style '%s'", event.Number, style), errors.New(""))
			continue
		}

		fmt.Printf("WK %d - %dm %s (%s)\n", event.Number, event.Distance, style, event.Gender)

		// +===========================+
		//        EVENT IMPORT
		// +===========================+

		stats.Found.Events++
		if runImport() {
			_, c, err := ec.ImportEvent(event, style, 1)
			if err != nil {
				importError(fmt.Sprintf("import event request failed for event %d!", event.Number), err)
				continue
			}
			if c {
				stats.Created.Events++
			}
			stats.Imported.Events++
		}

		// +===========================+
		//           RATING
		// +===========================+

		var ratingSplit []string
		ratingSplit = append(ratingSplit, eventString)
		for _, separator := range stg.RatingSeparators {
			var newSplits []string
			for _, split := range ratingSplit {
				newSplit := strings.Split(split, separator)
				newSplits = append(newSplits, newSplit...)
			}
			ratingSplit = newSplits
		}

		if lastEvent != event.Number {
			lastEvent = event.Number
			rankCount = 1
			rankRepetition = 0
		}

		for _, ratingString := range ratingSplit {

			for _, separator := range stg.RatingRightSeparators {
				if strings.Contains(ratingString, separator) {
					ratingString = substr(ratingString, separator)
				}
			}

			// extract result rows
			if strings.Contains(ratingString, stg.ResultSeparator) {
				resultsString := substrr(ratingString, stg.ResultSeparator)

				for _, cutString := range stg.ResultEndCutStrings {
					resultsString = substr(resultsString, cutString)
				}

				if strings.Index(resultsString, "1.") == 0 {
					rankCount = 1
					rankRepetition = 0
				}

				for resultsString != "" {
					var rs string

					// check what is coming first
					nextRank := strings.Index(resultsString, strconv.Itoa(rankCount+rankRepetition)+".")
					sameRank := strings.Index(resultsString, strconv.Itoa(rankCount)+".")

					if nextRank == -1 && sameRank == -1 { // nextRank = -1; sameRank -1
						rs = resultsString
						resultsString = ""
						// next higher rank
					} else if sameRank == -1 || (nextRank != -1 && nextRank < sameRank) {
						rs = substr(resultsString, strconv.Itoa(rankCount+rankRepetition)+".")
						resultsString = strconv.Itoa(rankCount+rankRepetition) + "# " + substrr(resultsString, strconv.Itoa(rankCount+rankRepetition)+".")
						rankCount += rankRepetition
						rankRepetition = 1
						// still the same rank
					} else {
						rs = substr(resultsString, strconv.Itoa(rankCount)+".")
						resultsString = strconv.Itoa(rankCount) + "# " + substrr(resultsString, strconv.Itoa(rankCount)+".")
						rankRepetition++
					}

					results[event.Number] = append(results[event.Number], strings.Replace(rs, "#", ".", 1))
				}

			}

			// extract disqualification
			if strings.Contains(ratingString, "disqualifiziert") {
				disqualificationString := substrr(ratingString, "disqualifiziert")

				for _, cutString := range stg.ResultEndCutStrings {
					disqualificationString = substr(disqualificationString, cutString)
				}

				timeRegex := regexp.MustCompile(stg.DisqualificationTimePattern)

				for timeRegex.Match([]byte(disqualificationString)) && disqualificationString != "" {
					disqualificationSplit := timeRegex.Split(disqualificationString, 2)

					var appendString string
					if len(disqualificationSplit) < 2 || disqualificationSplit[1] == "" {
						appendString = disqualificationString
						disqualificationString = ""
					} else {
						appendString = substr(disqualificationString, disqualificationSplit[1])
						disqualificationString = disqualificationSplit[1]
					}

					disqualifications[event.Number] = append(disqualifications[event.Number], appendString)
				}
			}

			// extract dns
			// TODO collect dns starts (not easily possible, won't be done now; use DSV7)
			//if strings.Contains(ratingString, "nicht am Start") {
			//	dnsString := substrr(ratingString, "nicht am Start")
			//
			//	for _, cutString := range stg.ResultEndCutStrings {
			//		dnsString = substr(dnsString, cutString)
			//	}
			//	println("n\t\t" + dnsString)
			//}

			// extract canceled starts
			// TODO collect logged out starts (not easily possible, won't be done now; use DSV7)
			//if strings.Contains(ratingString, "abgemeldet") {
			//	canceledString := substrr(ratingString, "abgemeldet")
			//
			//	for _, cutString := range stg.ResultEndCutStrings {
			//		canceledString = substr(canceledString, cutString)
			//	}
			//	println("a\t\t" + canceledString)
			//}
		}

	}

	for ev, eventResults := range results {
		println("WK: " + strconv.Itoa(ev))
		event, err := ec.GetEventByMeetingAndNumber(meeting, ev)
		if err != nil {
			importError(fmt.Sprintf("failed to fetch event %d for result import", ev), err)
			continue
		}
		for _, result := range eventResults {

			for _, cutString := range stg.ResultEndCutStrings {
				result = substr(result, cutString)
			}

			resultRegex := regexp.MustCompile(stg.ResultPattern)

			if !resultRegex.Match([]byte(result)) {
				continue
			}

			// result like "7. Vazanska, Aneta2008Plavecký klub Litvínov03:31,06259  50m: 00:48,10 | 100m: 01:43,30 | 150m: 02:38,35"

			rankingSplit := strings.SplitN(result, ".", 2)
			ranking, err := strconv.Atoi(trim(rankingSplit[0]))
			if err != nil {
				importError(fmt.Sprintf("failed to parse ranking for result in event %d with content '%s'", ev, result), err)
				continue
			}

			yearRegex := regexp.MustCompile(stg.YearPattern)
			yearSplit := yearRegex.Split(rankingSplit[1], 2)

			athleteName := trim(yearSplit[0])
			for _, replaceString := range stg.AthleteReplaceStrings {
				athleteName = strings.Replace(athleteName, replaceString, "", -1)
			}
			athleteYearString := trim(substrr(rankingSplit[1], athleteName))[:4]
			athleteYear, err := strconv.Atoi(athleteYearString)
			if err != nil {
				importError(fmt.Sprintf("failed to parse year for result e: %d a: %s with content '%s'", ev, athleteName, athleteYearString), err)
				continue
			}

			swimTimeRegex := regexp.MustCompile(stg.SwimTimePattern)
			swimTimeSplit := swimTimeRegex.Split(rankingSplit[1], 2)
			swimTimeString := trim(substrr(rankingSplit[1], swimTimeSplit[0]))[:8]

			swimTime, err := swimTimeToDuration(swimTimeString)
			if err != nil {
				importError(fmt.Sprintf("failed to parse swimtime for result e: %d a: %s (%d) with content '%s'", ev, athleteName, athleteYear, swimTimeString), err)
				continue
			}

			athleteTeam := trim(substrr(substr(result, swimTimeString), athleteYearString))

			start := startModel.Start{
				Meeting:         meeting,
				Event:           ev,
				AthleteName:     athleteName,
				AthleteYear:     athleteYear,
				AthleteTeamName: athleteTeam,
				Rank:            ranking,
				Certified:       true,
			}

			if event.RelayDistance != "" {
				start.IsRelay = true
			}

			athleteKey := importedKey{
				event:   start.Event,
				athlete: start.AthleteName,
			}

			if stg.OncePerEvent && importedAthletes[athleteKey] {
				importWarning(fmt.Sprintf("do not import again: e: %d - %s", athleteKey.event, athleteKey.athlete))
				continue
			}

			fmt.Printf("\t\tResult %d. - %s (%d) %s -> %s\n", start.Rank, start.AthleteName, start.AthleteYear, start.AthleteTeamName, swimTime.String())

			// +===========================+
			//         START IMPORT
			// +===========================+

			stats.Found.Starts++
			if runImport() {
				newStart, c, err := sc.ImportStart(start)
				if err != nil {
					importError(fmt.Sprintf("import start request failed for start e: %d %d. %s (%d)!", start.Event, start.Rank, start.AthleteName, start.AthleteYear), err)
					continue
				}
				if c {
					stats.Created.Starts++
				}
				stats.Imported.Starts++

				start = *newStart

				importedAthletes[athleteKey] = true
			}

			result := startModel.Result{
				Time:       swimTime,
				ResultType: "result_list",
				LapMeters:  event.Distance,
			}

			// +===========================+
			//        RESULT IMPORT
			// +===========================+

			stats.Found.Results++
			if runImport() {
				_, c, err := sc.ImportResult(start, result)
				if err != nil {
					importError(fmt.Sprintf("import result request failed for start %d/%d/%d!", start.Event, start.HeatNumber, start.Lane), err)
					continue
				}
				if c {
					stats.Created.Results++
				}
				stats.Imported.Results++
			}
		}
	}

	// +===========================+
	//       DISQUALIFICATION
	// +===========================+

	for ev, eventDisqualifications := range disqualifications {
		for _, disqualification := range eventDisqualifications {

			start := startModel.Start{
				Meeting:   meeting,
				Event:     ev,
				Certified: true,
			}

			yearRegex := regexp.MustCompile(stg.YearPattern)
			nameSplit := yearRegex.Split(disqualification, 2)
			year := substr(substrr(disqualification, nameSplit[0]), nameSplit[1])
			reasonTime := trim(nameSplit[1])

			reasonSplit := strings.SplitN(reasonTime, stg.ReasonRightSeparator, 2)
			reason := reasonSplit[0]

			// remove team name from reason since might be included
			for _, team := range *teams {
				reason = strings.ReplaceAll(reason, team.Name, "")
				for _, alias := range team.Alias {
					reason = strings.ReplaceAll(reason, alias, "")
				}
			}

			timeHour := 0
			timeMin := 0
			if len(reasonSplit) > 1 {
				timeRegex := regexp.MustCompile(stg.DisqualificationTimePattern)
				beforeTime := timeRegex.Split(reasonSplit[1], 2)
				clockTime := substrr(reasonSplit[1], beforeTime[0])
				if len(clockTime) >= 5 {
					timeHour, _ = strconv.Atoi(clockTime[:2])
					timeMin, _ = strconv.Atoi(clockTime[3:5])
				}
			}

			now := timeNow()

			clockTime := time.Date(now.Year(), now.Month(), now.Day(), timeHour, timeMin, 0, 0, time.Local)

			start.AthleteName = trim(nameSplit[0])
			start.AthleteYear, _ = strconv.Atoi(trim(year))

			// +===========================+
			//         START IMPORT
			// +===========================+

			stats.Found.Starts++
			if runImport() {
				newStart, c, err := sc.ImportStart(start)
				if err != nil {
					importError(fmt.Sprintf("import start request failed for start e: %d %s (%d)!", start.Event, start.AthleteName, start.AthleteYear), err)
					continue
				}
				if c {
					stats.Created.Starts++
				}
				stats.Imported.Starts++

				start = *newStart
			}

			// +===========================+
			//    DISQUALIFICATION IMPORT
			// +===========================+

			stats.Found.Disqualifications++
			if runImport() {
				_, c, err := dc.ImportDisqualification(start, reason, "disqualified", clockTime)
				if err != nil {
					importError(fmt.Sprintf("import disqualification request failed for %s (%d) - %s (%s)!", start.AthleteName, start.AthleteYear, reason, clockTime), err)
					continue
				}
				if c {
					stats.Created.Disqualifications++
				}
				stats.Imported.Disqualifications++
			}
		}
	}

	return &stats, nil
}

func shouldSkip(s string, skipStrings []string, requiredStrings []string) bool {
	skip := false
	for _, skipString := range skipStrings {
		if strings.Contains(s, skipString) {
			skip = true
			break
		}
	}
	for _, requiredString := range requiredStrings {
		if !strings.Contains(s, requiredString) {
			skip = true
			break
		}
	}
	return skip
}

func swimTimeToDuration(tm string) (time.Duration, error) {
	tm = strings.Replace(tm, ":", "m", 1)
	tm = strings.Replace(tm, ",", "s", 1)
	tm += "0ms"
	return time.ParseDuration(tm)
}
