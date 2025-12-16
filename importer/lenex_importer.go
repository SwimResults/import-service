package importer

import (
	"encoding/xml"
	"fmt"
	"github.com/konrad2002/lenexparser/model/elements"
	"github.com/konrad2002/lenexparser/model/enums"
	athleteModel "github.com/swimresults/athlete-service/model"
	importModel "github.com/swimresults/import-service/model"
	meetingModel "github.com/swimresults/meeting-service/model"
	startModel "github.com/swimresults/start-service/model"
	"io"
	"strconv"
	"time"
)

func ImportLenexFile(file string, meeting string, exclude []int, include []int, stg importModel.ImportSetting) (*importModel.ImportFileStats, error) {
	var stats importModel.ImportFileStats

	buf, err1 := getFileReader(file)
	if err1 != nil {
		return nil, err1
	}

	xmlString, err := io.ReadAll(buf)
	if err != nil {
		return nil, err
	}

	var lenex elements.Lenex
	err = xml.Unmarshal(xmlString, &lenex)
	if err != nil {
		return nil, err
	}

	meet := lenex.Meets[0]

	eventOrdering := 1
	meetingYear := meet.EntryStartDate.Year()

	heats := map[int]startModel.Heat{}  // map heat id to heat
	ranks := map[int]elements.Ranking{} // map result id to rank (first occurrence)

	for _, session := range meet.Sessions {
		for _, event := range session.Events {

			// EVENT IMPORT
			fmt.Printf("%d", event.Number)
			if !IsEventImportable(event.Number, exclude, include) {
				print(" => no import")
				continue
			}

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

			// AGE GROUP IMPORT
			for _, ageGroup := range event.AgeGroups {
				minAge := meetingYear - ageGroup.AgeMin
				maxAge := meetingYear - ageGroup.AgeMax

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
					importEvent.Gender = "FEMALE"
					break
				case enums.AgeGroupGenderMale:
					importEvent.Gender = "MALE"
					break
				default:
					importEvent.Gender = "MIXED"
				}

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

				// COLLECT RANKS
				for _, ranking := range ageGroup.Rankings {
					if ranks[ranking.ResultId].Place > ranking.Place {
						ranks[ranking.ResultId] = ranking
					}
				}
			}

			// HEATS
			for _, heat := range event.Heats {
				startTime := heat.Daytime.Time

				// if Daytime does not include but date only hours, add date of the session
				if startTime.Year() < 1980 {
					startTime = time.Date(session.Date.Year(), session.Date.Month(), session.Date.Day(), startTime.Hour(), startTime.Minute(), startTime.Second(), startTime.Nanosecond(), stg.TimeZone.Location())
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

				// IMPORT HEAT
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

				heats[heat.HeatId] = *newHeat
			}
		}
	}

	// TEAMS
	for _, team := range meet.Clubs {
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
					Identifier: newTeam.Identifier,
				},
			}

			switch athlete.Gender {
			case enums.GenderFemale:
				importAthlete.Gender = "FEMALE"
				break
			case enums.GenderMale:
				importAthlete.Gender = "MALE"
				break
			default:
				importAthlete.Gender = "MIXED"
			}

			newAthlete, created, err := ac.ImportAthlete(importAthlete, meeting)
			if err != nil {
				return &stats, err
			}
			cs := 'o'
			if created {
				cs = '+'
				stats.Created.Athletes++
			}
			stats.Imported.Athletes++
			fmt.Printf("[ %c ] > id: %s, name: %s, part: %s\n", cs, newAthlete.Identifier.String(), newAthlete.Name, newAthlete.Participation)

			// STARTS + RESULTS + DISQUALIFICATIONS
			// TODO: support entry lists (no heat)
			for _, entry := range athlete.Entries {
				heat := heats[entry.HeatId]

				if !IsEventImportable(heat.Event, exclude, include) {
					fmt.Printf("entry of '%s' for event: '%d' => no import\n", newAthlete.Name, heat.Event)
					continue
				}

				// START
				start := startModel.Start{
					Meeting:         meeting,
					Event:           heat.Event,
					HeatNumber:      heat.Number,
					Lane:            entry.Lane,
					Athlete:         newAthlete.Identifier,
					AthleteName:     athlete.Firstname + " " + athlete.Lastname,
					AthleteYear:     athlete.Birthdate.Year(),
					AthleteTeam:     newAthlete.Team.Identifier,
					AthleteTeamName: team.Name,
				}
				newStart, c, err2 := sc.ImportStart(start)
				if err2 != nil {
					return &stats, err2
				}
				if c {
					stats.Created.Starts++
					fmt.Printf("[ ! ] start has been created from entry: id: '%s'; event: '%d', athlete: '%s'\n", newStart.Identifier, newStart.Event, newStart.AthleteName)
				}
				stats.Imported.Starts++

				// IMPORT REGISTRATION TIME
				if entry.EntryTime.Milliseconds() > 0 {
					resultModel := startModel.Result{
						Time:       entry.EntryTime.Duration,
						ResultType: "registration",
					}

					_, _, err3 := sc.ImportResult(start, resultModel)
					if err3 != nil {
						return &stats, err3
					}
					stats.Created.Results++
					stats.Imported.Results++
				}
			}

			for _, result := range athlete.Results {
				heat := heats[result.HeatId]
				rank := ranks[result.ResultId]

				if !IsEventImportable(heat.Event, exclude, include) {
					fmt.Printf("result of '%s' for event: '%d' => no import\n", newAthlete.Name, heat.Event)
					continue
				}

				rankValue := 0
				if rank.ResultId == result.ResultId {
					rankValue = rank.Place
				}

				// START
				start := startModel.Start{
					Meeting:         meeting,
					Event:           heat.Event,
					HeatNumber:      heat.Number,
					Lane:            result.Lane,
					Athlete:         newAthlete.Identifier,
					AthleteName:     athlete.Firstname + " " + athlete.Lastname,
					AthleteYear:     athlete.Birthdate.Year(),
					AthleteTeam:     newAthlete.Team.Identifier,
					AthleteTeamName: team.Name,
					Rank:            rankValue,
					Certified:       true,
				}
				newStart, c, err2 := sc.ImportStart(start)
				if err2 != nil {
					return &stats, err2
				}
				if c {
					stats.Created.Starts++
					fmt.Printf("[ ! ] start has been created from result: id: '%s'; event: '%d', athlete: '%s'\n", newStart.Identifier, newStart.Event, newStart.AthleteName)
				}
				stats.Imported.Starts++

				// IMPORT TIME SPLITS
				for _, split := range result.Splits {
					// LAP Result
					lapResult := startModel.Result{
						Time:       split.SwimTime.Duration,
						ResultType: "lap",
						LapMeters:  split.Distance,
					}

					_, _, err3 := sc.ImportResult(*newStart, lapResult)
					if err3 != nil {
						return &stats, err3
					}
					stats.Created.Results++
					stats.Imported.Results++
				}

				// IMPORT RESULT
				if result.SwimTime.Milliseconds() > 0 {
					resultModel := startModel.Result{
						Time:       result.SwimTime.Duration,
						ResultType: "result_list",
					}

					_, _, err3 := sc.ImportResult(start, resultModel)
					if err3 != nil {
						return &stats, err3
					}
					stats.Created.Results++
					stats.Imported.Results++
				}

				// DISQUALIFICATION
				disqType := ""
				switch result.Status {
				case enums.ResultStatusDSQ:
					disqType = "disqualified"
					break
				case enums.ResultStatusDNS:
					disqType = "dns"
					break
				case enums.ResultStatusDNF:
					disqType = "dnf"
					break
				case enums.ResultStatusSICK:
					disqType = "sick"
					break
				case enums.ResultStatusWDR:
					disqType = "withdrawn"
					break
				}

				if disqType != "" {
					disqualification, created, err4 := dc.ImportDisqualification(*newStart, result.Comment, disqType, time.UnixMicro(0))
					if err4 != nil {
						return &stats, err4
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
	}

	fmt.Printf(" +==============================+ \n")

	//var starts []startModel.Start
	//
	//for _, dsvResult := range erg.PNErgebnisse {
	//
	//	if dsvResult.GrundDerNichtwertung == "AB" {
	//		continue
	//	}
	//
	//	if unusedEvents[dsvResult.Wettkampfnummer] {
	//		continue
	//	}
	//
	//	if !IsEventImportable(dsvResult.Wettkampfnummer, exclude, include) {
	//		continue
	//	}
	//
	//	// RESULT
	//	result := startModel.Result{
	//		Time:       dsvResult.Endzeit.Duration(),
	//		ResultType: "result_list",
	//	}
	//
	//	starts = append(starts, *newStart)
	//
	//	if dsvResult.GrundDerNichtwertung != "" {
	//		disqType := "disqualified"
	//		switch dsvResult.GrundDerNichtwertung {
	//		case "NA":
	//			disqType = "dns"
	//			break
	//		case "AU":
	//			disqType = "dnf"
	//			break
	//		case "ZU":
	//			disqType = "time"
	//			break
	//		}
	//		disqualification, created, err4 := dc.ImportDisqualification(start, dsvResult.Disqualifikationsbemerkung, disqType, time.UnixMicro(0))
	//		if err4 != nil {
	//			return &stats, err4
	//		}
	//		cs := 'o'
	//		if created {
	//			cs = '+'
	//			stats.Created.Disqualifications++
	//		}
	//		stats.Imported.Disqualifications++
	//		fmt.Printf("[ %c ] > id: %s, type: %s, reason: %s\n", cs, disqualification.Identifier, disqualification.Type, disqualification.Reason)
	//	} else {
	//		_, _, err3 := sc.ImportResult(start, result)
	//		if err3 != nil {
	//			return &stats, err3
	//		}
	//		stats.Created.Results++
	//		stats.Imported.Results++
	//	}
	//
	//}
	//
	//for _, dsvLap := range erg.PNZwischenzeiten {
	//	if unusedEvents[dsvLap.Wettkampfnummer] {
	//		continue
	//	}
	//
	//	if !IsEventImportable(dsvLap.Wettkampfnummer, exclude, include) {
	//		continue
	//	}
	//
	//	// LAP Result
	//	lapResult := startModel.Result{
	//		Time:       dsvLap.Zwischenzeit.Duration(),
	//		ResultType: "lap",
	//		LapMeters:  dsvLap.Distanz,
	//	}
	//
	//	var lapStart startModel.Start
	//
	//	found := false
	//	for _, start := range starts {
	//		if start.AthleteMeetingId == dsvLap.VeranstaltungsIdSchwimmer && start.Event == dsvLap.Wettkampfnummer {
	//			lapStart = start
	//			found = true
	//			break
	//		}
	//	}
	//
	//	if found {
	//		_, _, err3 := sc.ImportResult(lapStart, lapResult)
	//		if err3 != nil {
	//			return &stats, err3
	//		}
	//		stats.Created.Results++
	//		stats.Imported.Results++
	//	}
	//}

	return &stats, nil
}
