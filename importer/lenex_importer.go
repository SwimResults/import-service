package importer

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/konrad2002/lenexparser/model/elements"
	"github.com/konrad2002/lenexparser/model/enums"
	athleteModel "github.com/swimresults/athlete-service/model"
	importModel "github.com/swimresults/import-service/model"
	meetingModel "github.com/swimresults/meeting-service/model"
	startModel "github.com/swimresults/start-service/model"
	"io"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"time"
)

type ProgressCallback func(progress float64, message string)

func ImportLenexFile(file string, meeting string, exclude []int, include []int, features []string, stg importModel.ImportSetting, progressCallback ProgressCallback) (*importModel.ImportFileStats, error) {
	var stats importModel.ImportFileStats

	// Helper function to handle nil callback
	progress := func(pct float64, msg string) {
		if progressCallback != nil {
			progressCallback(pct, msg)
		}
	}

	// Read content via getFileReader for both local and remote sources
	buf, err1 := getFileReader(file)
	if err1 != nil {
		return nil, err1
	}

	data, err := io.ReadAll(buf)
	if err != nil {
		return nil, err
	}

	// Read LEF content directly or extract from LXF (zip) first
	var xmlString []byte
	ext := strings.ToLower(filepath.Ext(file))
	if ext == ".lxf" {
		fmt.Printf("[ unzip ] detected .lxf archive: %s\n", file)
		zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
		if err != nil {
			return nil, err
		}

		found := false
		for _, zf := range zr.File {
			if strings.HasSuffix(strings.ToLower(zf.Name), ".lef") {
				rc, err := zf.Open()
				if err != nil {
					return nil, err
				}
				b, err := io.ReadAll(rc)
				rc.Close()
				if err != nil {
					return nil, err
				}
				xmlString = b
				fmt.Printf("[ unzip ] using LEF from archive entry: %s\n", zf.Name)
				found = true
				break
			}
		}
		if !found {
			fmt.Printf("[ unzip ] no .lef file found inside archive: %s\n", file)
			return nil, fmt.Errorf("no .lef file found inside archive: %s", file)
		}
	} else {
		fmt.Printf("[ direct ] reading LEF file directly: %s\n", file)
		xmlString = data
	}

	var lenex elements.Lenex
	err = xml.Unmarshal(xmlString, &lenex)
	if err != nil {
		return nil, err
	}

	meet := lenex.Meets[0]

	eventOrdering := 1
	meetingYear := meet.AgeDate.Value.Year()

	heats := map[int]startModel.Heat{}  // map heat id to heat
	ranks := map[int]elements.Ranking{} // map result id to rank (first occurrence)

	loc, err := time.LoadLocation(stg.TimeZone)
	if err != nil {
		return nil, errors.New("timezone " + stg.TimeZone + " is not a valid timezone")
	}

	// helper to process one athlete including starts, results, splits and disqualifications
	processAthleteForTeam := func(teamName string, teamImported athleteModel.Team, athlete elements.Athlete) error {
		dsvAthlete, err := strconv.Atoi(athlete.License)
		if err != nil {
			dsvAthlete = 0
		}

		importAthlete := athleteModel.Athlete{
			Name:      athlete.Firstname + " " + athlete.Lastname,
			Firstname: athlete.Firstname,
			Lastname:  athlete.Lastname,
			Year:      athlete.Birthdate.Year(),
			DsvId:     dsvAthlete,
			Team: athleteModel.Team{
				Identifier: teamImported.Identifier,
			},
		}

		switch athlete.Gender {
		case enums.GenderFemale:
			importAthlete.Gender = "FEMALE"
		case enums.GenderMale:
			importAthlete.Gender = "MALE"
		default:
			importAthlete.Gender = "MIXED"
		}

		newAthlete, created, err := ac.ImportAthlete(importAthlete, meeting)
		if err != nil {
			return err
		}
		cs := 'o'
		if created {
			cs = '+'
			stats.Created.Athletes++
		}
		stats.Imported.Athletes++
		fmt.Printf("[ %c ] > id: %s, name: %s, part: %s\n", cs, newAthlete.Identifier.String(), newAthlete.Name, newAthlete.Participation)

		for _, entry := range athlete.Entries {
			heat := heats[entry.HeatId]

			if heat.Number == 0 {
				continue
			}

			if !IsEventImportable(heat.Event, exclude, include) {
				fmt.Printf("entry of '%s' for event: '%d' => no import\n", newAthlete.Name, heat.Event)
				continue
			}

			start := startModel.Start{
				Meeting:         meeting,
				Event:           heat.Event,
				HeatNumber:      heat.Number,
				Lane:            entry.Lane,
				Athlete:         newAthlete.Identifier,
				AthleteName:     athlete.Firstname + " " + athlete.Lastname,
				AthleteYear:     athlete.Birthdate.Year(),
				AthleteTeam:     newAthlete.Team.Identifier,
				AthleteTeamName: teamName,
			}
			newStart, c, err2 := sc.ImportStart(start)
			if err2 != nil {
				return err2
			}
			if c {
				stats.Created.Starts++
				fmt.Printf("[ ! ] start has been created from entry: id: '%s'; event: '%d', athlete: '%s'\n", newStart.Identifier, newStart.Event, newStart.AthleteName)
			}
			stats.Imported.Starts++

			if entry.EntryTime.Milliseconds() > 0 {
				resultModel := startModel.Result{
					Time:       entry.EntryTime.Duration,
					ResultType: "registration",
				}

				_, _, err3 := sc.ImportResult(*newStart, resultModel)
				if err3 != nil {
					return err3
				}
				stats.Created.Results++
				stats.Imported.Results++
			}
		}

		for _, result := range athlete.Results {
			heat := heats[result.HeatId]
			rank := ranks[result.ResultId]

			if heat.Number == 0 {
				continue
			}

			if !IsEventImportable(heat.Event, exclude, include) {
				fmt.Printf("result of '%s' for event: '%d' => no import\n", newAthlete.Name, heat.Event)
				continue
			}

			rankValue := 0
			if rank.ResultId == result.ResultId {
				rankValue = rank.Place
			}

			start := startModel.Start{
				Meeting:         meeting,
				Event:           heat.Event,
				HeatNumber:      heat.Number,
				Lane:            result.Lane,
				Athlete:         newAthlete.Identifier,
				AthleteName:     athlete.Firstname + " " + athlete.Lastname,
				AthleteYear:     athlete.Birthdate.Year(),
				AthleteTeam:     newAthlete.Team.Identifier,
				AthleteTeamName: teamName,
				Rank:            rankValue,
				Certified:       true,
			}
			fmt.Printf("[   ] import start from result: event: '%d', athlete: '%s', rank: %d\n", start.Event, start.AthleteName, start.Rank)

			newStart, c, err2 := sc.ImportStart(start)
			if err2 != nil {
				return err2
			}
			if c {
				stats.Created.Starts++
				fmt.Printf("[ ! ] start has been created from result: id: '%s'; event: '%d', athlete: '%s'\n", newStart.Identifier, newStart.Event, newStart.AthleteName)
			}
			stats.Imported.Starts++

			if result.EntryTime.Milliseconds() > 0 {
				resultModel := startModel.Result{
					Time:       result.EntryTime.Duration,
					ResultType: "registration",
				}

				if slices.Contains(features, "result") {
					_, _, err3 := sc.ImportResult(*newStart, resultModel)
					if err3 != nil {
						return err3
					}
					stats.Created.Results++
					stats.Imported.Results++
				}
			}

			for _, split := range result.Splits {
				lapResult := startModel.Result{
					Time:       split.SwimTime.Duration,
					ResultType: "lap",
					LapMeters:  split.Distance,
				}

				if slices.Contains(features, "result") {
					_, _, err3 := sc.ImportResult(*newStart, lapResult)
					if err3 != nil {
						return err3
					}
					stats.Created.Results++
					stats.Imported.Results++
				}
			}

			if result.SwimTime.Milliseconds() > 0 {
				resultModel := startModel.Result{
					Time:       result.SwimTime.Duration,
					ResultType: "result_list",
				}

				if slices.Contains(features, "result") {
					_, _, err3 := sc.ImportResult(*newStart, resultModel)
					if err3 != nil {
						return err3
					}
					stats.Created.Results++
					stats.Imported.Results++
				}
			}

			disqType := ""
			switch result.Status {
			case enums.ResultStatusDSQ:
				disqType = "disqualified"
			case enums.ResultStatusDNS:
				disqType = "dns"
			case enums.ResultStatusDNF:
				disqType = "dnf"
			case enums.ResultStatusSICK:
				disqType = "sick"
			case enums.ResultStatusWDR:
				disqType = "withdrawn"
			}

			if disqType != "" {
				if slices.Contains(features, "disqualification") {
					disqualification, created, err4 := dc.ImportDisqualification(*newStart, result.Comment, disqType, time.UnixMicro(0))
					if err4 != nil {
						return err4
					}
					cs := 'o'
					if created {
						cs = '+'
						stats.Created.Disqualifications++
					}
					stats.Imported.Disqualifications++
					fmt.Printf("[ %c ] > id: %s, type: %s, reason: %s\n", cs, disqualification.Identifier, disqualification.Type, disqualification.Reason)
				}
			}
		}

		return nil
	}

	// CALCULATE TOTAL ITEMS FOR PROGRESS TRACKING
	totalTeams := len(meet.Clubs)
	totalAthletes := 0
	totalEntries := 0
	totalResults := 0
	totalHeats := 0
	totalEvents := 0
	totalAgeGroups := 0

	for _, session := range meet.Sessions {
		for _, event := range session.Events {
			totalEvents++
			totalAgeGroups += len(event.AgeGroups)
			totalHeats += len(event.Heats)
		}
	}

	for _, team := range meet.Clubs {
		totalAthletes += len(team.Athletes)
		for _, athlete := range team.Athletes {
			totalEntries += len(athlete.Entries)
			totalResults += len(athlete.Results)
		}
	}

	totalItems := totalEvents + totalAgeGroups + totalHeats + totalTeams + totalAthletes + totalEntries + totalResults
	processedItems := 0

	progress(20, fmt.Sprintf("Starting import with %d total items to process", totalItems))

	for _, session := range meet.Sessions {
		for _, event := range session.Events {
			processedItems++

			// EVENT IMPORT
			fmt.Printf("%d", event.Number)

			println("event import")

			importEvent := meetingModel.Event{
				Number:   event.Number,
				Distance: event.SwimStyle.Distance,
				Meeting:  meeting,
				Ordering: eventOrdering,
			}

			eventOrdering++

			switch event.Gender {
			case enums.EventGenderFemale:
				importEvent.Gender = "FEMALE"
				break
			case enums.EventGenderMale:
				importEvent.Gender = "MALE"
				break
			default:
				importEvent.Gender = "MIXED"
			}

			if event.SwimStyle.RelayCount > 1 {
				importEvent.RelayDistance = fmt.Sprintf("%dx%d", event.SwimStyle.RelayCount, event.SwimStyle.Distance)
				importEvent.Distance = event.SwimStyle.RelayCount * event.SwimStyle.Distance
			}

			if IsEventImportable(event.Number, exclude, include) { // only import if in import list, but do not skip, heats need to be set
				if slices.Contains(features, "event") {
					newEvent, created, err3 := ec.ImportEvent(importEvent, string(event.SwimStyle.Stroke), session.Number)
					if err3 != nil {
						return nil, err3
					}

					if created {
						stats.Created.Events++
						print("(+) ")
					} else {
						print("( ) ")
					}
					println(newEvent.Number)

					stats.Imported.Events++
				}

				// AGE GROUP IMPORT
				for _, ageGroup := range event.AgeGroups {
					processedItems++

					minAge := meetingYear - ageGroup.AgeMin
					maxAge := meetingYear - ageGroup.AgeMax

					if ageGroup.AgeMax <= 0 {
						maxAge = 1900
					}

					importAgeGroup := meetingModel.AgeGroup{
						Meeting: meeting,
						Event:   event.Number,
						Default: false,
						MinAge:  strconv.Itoa(minAge),
						MaxAge:  strconv.Itoa(maxAge),
						IsYear:  true,
						Name:    ageGroup.Name,
					}

					switch ageGroup.Gender {
					case enums.AgeGroupGenderFemale:
						importAgeGroup.Gender = "FEMALE"
						break
					case enums.AgeGroupGenderMale:
						importAgeGroup.Gender = "MALE"
						break
					case enums.AgeGroupGenderMixed:
						importAgeGroup.Gender = "MIXED"
					default:
						importAgeGroup.Gender = "UNSET"
					}

					if slices.Contains(features, "age_group") {

						newAgeGroup, created, err5 := gc.ImportAgeGroup(importAgeGroup)
						if err5 != nil {
							return nil, err5
						}

						if created {
							stats.Created.AgeGroups++
							print("(+) ")
						} else {
							print("( ) ")
						}
						println(newAgeGroup.Name)

						stats.Imported.AgeGroups++
					}

					// COLLECT RANKS
					for _, ranking := range ageGroup.Rankings {
						if ranks[ranking.ResultId].Place < ranking.Place {
							ranks[ranking.ResultId] = ranking
						}
					}

				}

			} else {
				print(" => no import")
			}

			// HEATS
			for _, heat := range event.Heats {
				processedItems++

				startTime := heat.Daytime.Time

				// if Daytime does not include but date only hours, add date of the session
				if startTime.Year() < 1980 {
					startTime = time.Date(session.Date.Year(), session.Date.Month(), session.Date.Day(), startTime.Hour(), startTime.Minute(), startTime.Second(), startTime.Nanosecond(), loc)
				}

				heatModel := startModel.Heat{
					Meeting:         meeting,
					Event:           event.Number,
					Number:          heat.Number,
					StartEstimation: startTime,
				}

				stats.Found.Heats++

				// TODO: for EasyWk heat 0 contains withdrawn starts
				if heat.Number == 0 {
					heats[heat.HeatId] = heatModel // set to heat model so later can check if number == 0
					continue
				}

				if IsEventImportable(event.Number, exclude, include) { // only import if in import list, but do not skip, heats need to be set
					// IMPORT HEAT

					if slices.Contains(features, "heat") {
						newHeat, c, err := hc.ImportHeat(heatModel)
						if err != nil {
							importError(fmt.Sprintf("import heat request failed for heat %d/%d!", event.Number, heat.Number), err)
							continue
						}
						fmt.Printf("[ o ] import heat: event: '%d', number: '%d', start: '%s', start before: '%s'\n", newHeat.Event, newHeat.Number, newHeat.StartEstimation, startTime)

						if c {
							stats.Created.Heats++
						}
						stats.Imported.Heats++

						heatModel = *newHeat
					}
				}

				heats[heat.HeatId] = heatModel
			}
		}
	}

	// TEAMS
	for _, team := range meet.Clubs {
		processedItems++

		stateId, err := strconv.Atoi(team.Region)
		if err != nil {
			stateId = 0
		}

		dsv, err := strconv.Atoi(team.Code)
		if err != nil {
			dsv = 0
		}

		importTeam := athleteModel.Team{ // TODO: check if Contact, Address and Website overwrites
			Name:    team.Name,
			Country: string(team.Nation),
			DsvId:   dsv,
			StateId: stateId,
			Contact: athleteModel.Contact{
				Name:  team.Contact.Name,
				EMail: team.Contact.Email,
				Phone: team.Contact.Phone,
				Fax:   team.Contact.Fax,
			},
			Address: athleteModel.Address{
				Street:     team.Contact.Street,
				Number:     team.Contact.Street2,
				City:       team.Contact.City,
				PostalCode: team.Contact.Zip,
			},
			Website: team.Contact.Internet,
		}

		newTeam, created, err := tc.ImportTeam(importTeam, meeting)
		if err != nil {
			return &stats, err
		}
		cs := 'o'
		if created {
			cs = '+'
			stats.Created.Teams++
		}
		stats.Imported.Teams++
		fmt.Printf("[ %c ] > id: %s, name: %s, part: %s\n", cs, newTeam.Identifier.String(), newTeam.Name, newTeam.Participation)

		// ATHLETES
		for _, athlete := range team.Athletes {
			processedItems++
			if err := processAthleteForTeam(team.Name, *newTeam, athlete); err != nil {
				return &stats, err
			}
			progressPct := 20 + (float64(processedItems)/float64(totalItems))*80
			progress(progressPct, fmt.Sprintf("Processing athletes: %d / %d", processedItems, totalItems))
		}
	}

	fmt.Printf(" +==============================+ \n")

	progress(95, "Finalizing import...")

	return &stats, nil
}
