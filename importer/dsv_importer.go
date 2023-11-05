package importer

import (
	"bytes"
	"fmt"
	"github.com/konrad2002/dsvparser/model"
	"github.com/konrad2002/dsvparser/parser"
	athleteModel "github.com/swimresults/athlete-service/model"
	importModel "github.com/swimresults/import-service/model"
	eventModel "github.com/swimresults/meeting-service/model"
	startModel "github.com/swimresults/start-service/model"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func ImportDsvDefinitionFile(file string, meeting string, exclude []int, include []int) (*importModel.ImportFileStats, error) {

	var stats importModel.ImportFileStats
	var buf io.Reader

	if strings.Contains(file, "http") {
		// get file content from url
		resp, err := http.Get(file)
		if err != nil {
			return &stats, err
		}

		buf = resp.Body
	} else {
		// get file content from local file
		dat, err := os.ReadFile(file)
		if err != nil {
			return &stats, err
		}

		buf = bytes.NewBuffer(dat)
	}

	r := parser.NewReader(buf)
	res, err := r.Read()
	if err != nil {
		return nil, err
	}

	def := res.(*model.Wettkampfdefinitionsliste)

	for _, dsvEvent := range def.Wettkaempfe {
		if !IsEventImportable(dsvEvent.Wettkampfnummer, exclude, include) {
			continue
		}

		event := eventModel.Event{
			Number:   dsvEvent.Wettkampfnummer,
			Distance: dsvEvent.Einzelstrecke,
			Meeting:  meeting,
		}

		if dsvEvent.Geschlecht == 'W' {
			event.Gender = "FEMALE"
		}
		if dsvEvent.Geschlecht == 'M' {
			event.Gender = "MALE"
		}
		if dsvEvent.Geschlecht == 'D' || dsvEvent.Geschlecht == 'X' {
			event.Gender = "MIXED"
		}

		if dsvEvent.AnzahlStarter > 1 {
			event.RelayDistance = strconv.Itoa(dsvEvent.AnzahlStarter) + "x" + strconv.Itoa(dsvEvent.Einzelstrecke)
			event.Distance = dsvEvent.AnzahlStarter * dsvEvent.Einzelstrecke
		}

		newEvent, created, err3 := ec.ImportEvent(event, string(dsvEvent.Technik), dsvEvent.Abschnittsnummer)
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

	return &stats, nil
}

func ImportDsvResultFile(file string, meeting string, exclude []int, include []int) (*importModel.ImportFileStats, error) {

	var stats importModel.ImportFileStats

	var buf io.Reader
	if strings.Contains(file, "http") {
		// get file content from url
		resp, err := http.Get(file)
		if err != nil {
			return &stats, err
		}

		buf = resp.Body
	} else {
		// get file content from local file
		dat, err := os.ReadFile(file)
		if err != nil {
			return &stats, err
		}

		buf = bytes.NewBuffer(dat)
	}

	r := parser.NewReader(buf)
	res, err := r.Read()
	if err != nil {
		return &stats, err
	}

	erg := res.(*model.Wettkampfergebnisliste)

	var unusedEvents = make(map[int]bool)

	for _, dsvEvent := range erg.Wettkaempfe {
		if dsvEvent.Ausuebung != "GL" {
			unusedEvents[dsvEvent.Wettkampfnummer] = true
		} else {
			unusedEvents[dsvEvent.Wettkampfnummer] = false
		}
	}

	// TEAMS
	for _, team := range erg.Vereine {
		team := athleteModel.Team{
			Name:    team.Vereinsbezeichnung,
			Country: team.FinaNationenkuerzel,
			DsvId:   team.Vereinskennzahl,
			StateId: team.Landesschwimmverband,
		}
		newTeam, created, err := tc.ImportTeam(team, meeting)
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
	}

	fmt.Printf(" +==============================+ \n")

	var starts []startModel.Start

	for _, dsvResult := range erg.PNErgebnisse {

		if dsvResult.GrundDerNichtwertung == "AB" {
			continue
		}

		if unusedEvents[dsvResult.Wettkampfnummer] {
			continue
		}

		if !IsEventImportable(dsvResult.Wettkampfnummer, exclude, include) {
			continue
		}

		// ATHLETE
		athlete := athleteModel.Athlete{
			Name:   dsvResult.Name,
			Year:   dsvResult.Jahrgang,
			Gender: string(dsvResult.Geschlecht),
			DsvId:  dsvResult.DsvId,
			Team: athleteModel.Team{
				DsvId: dsvResult.Vereinskennzahl,
				Name:  dsvResult.Verein,
			},
		}

		if dsvResult.Geschlecht == 'W' {
			athlete.Gender = "FEMALE"
		}
		if dsvResult.Geschlecht == 'M' {
			athlete.Gender = "MALE"
		}
		if dsvResult.Geschlecht == 'D' || dsvResult.Geschlecht == 'X' {
			athlete.Gender = "MIXED"
		}

		newAthlete, created, err := ac.ImportAthlete(athlete, meeting)
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

		// RESULT
		result := startModel.Result{
			Time:       dsvResult.Endzeit.Duration(),
			ResultType: "result_list",
		}

		start := startModel.Start{
			Meeting:          meeting,
			Event:            dsvResult.Wettkampfnummer,
			Athlete:          newAthlete.Identifier,
			AthleteMeetingId: dsvResult.VeranstaltungsIdSchwimmer,
			AthleteName:      dsvResult.Name,
			AthleteYear:      dsvResult.Jahrgang,
			AthleteTeam:      newAthlete.Team.Identifier,
			AthleteTeamName:  dsvResult.Verein,
			Rank:             dsvResult.Platz,
			Certified:        true,
		}
		newStart, c, err2 := sc.ImportStart(start)
		if err2 != nil {
			return &stats, err2
		}
		if c {
			stats.Created.Starts++
			fmt.Printf("[ ! ] start has been created: id: '%s'; event: '%d', athlete: '%s'", newStart.Identifier, newStart.Event, newStart.AthleteName)
		}
		stats.Imported.Starts++

		starts = append(starts, *newStart)

		if dsvResult.GrundDerNichtwertung != "" {
			disqType := "disqualified"
			switch dsvResult.GrundDerNichtwertung {
			case "NA":
				disqType = "dns"
				break
			case "AU":
				disqType = "dnf"
				break
			case "ZU":
				disqType = "time"
				break
			}
			disqualification, created, err4 := dq.ImportDisqualification(start, dsvResult.Disqualifikationsbemerkung, disqType, time.UnixMicro(0))
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
		} else {
			_, _, err3 := sc.ImportResult(start, result)
			if err3 != nil {
				return &stats, err3
			}
			stats.Created.Results++
			stats.Imported.Results++
		}

	}

	for _, dsvLap := range erg.PNZwischenzeiten {
		if unusedEvents[dsvLap.Wettkampfnummer] {
			continue
		}

		if !IsEventImportable(dsvLap.Wettkampfnummer, exclude, include) {
			continue
		}

		// LAP Result
		lapResult := startModel.Result{
			Time:       dsvLap.Zwischenzeit.Duration(),
			ResultType: "lap",
			LapMeters:  dsvLap.Distanz,
		}

		var lapStart startModel.Start

		found := false
		for _, start := range starts {
			if start.AthleteMeetingId == dsvLap.VeranstaltungsIdSchwimmer && start.Event == dsvLap.Wettkampfnummer {
				lapStart = start
				found = true
				break
			}
		}

		if found {
			_, _, err3 := sc.ImportResult(lapStart, lapResult)
			if err3 != nil {
				return &stats, err3
			}
			stats.Created.Results++
			stats.Imported.Results++
		}
	}

	return &stats, nil

}
