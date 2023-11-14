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
	buf, err := getFileReader(file)
	if err != nil {
		return "", err
	}

	buff := bytes.NewBuffer([]byte{})
	_, err = io.Copy(buff, buf)
	if err != nil {
		return "", err
	}

	rdr := bytes.NewReader(buff.Bytes())

	reader, err := pdf.NewReader(rdr, rdr.Size())
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
		for _, genders := range stg.GenderMapping {
			if strings.Contains(distanceSplit[1], genders[0]) {
				event.Gender = genders[1]
				style = trim(substr(distanceSplit[1], genders[0]))
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
				yearString := trim(substrr(substr(laneString, teamNameString), athleteNameString))
				swimTimeRestString := trim(substrr(laneString, teamNameString))

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

				if runImport() && event.RelayDistance == "" {
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

// ImportPdfResultList takes the path to a pdf file that contains a result
// list. All teams, athletes and results will be imported.
// If exclude is set, given event numbers will be excluded from import.
// If include is set, only given event numbers will be imported.
//
// For import process details see documentation on GitHub.
func ImportPdfResultList(file string, meeting string, exclude []int, include []int, stg importModel.ImportPdfResultListSettings) (*importModel.ImportFileStats, error) {
	return nil, nil
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
