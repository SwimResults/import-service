package importer

import (
	"bytes"
	"fmt"
	"github.com/konrad2002/dsvparser/model"
	"github.com/konrad2002/dsvparser/parser"
	athleteModel "github.com/swimresults/athlete-service/model"
	eventModel "github.com/swimresults/meeting-service/model"
	startModel "github.com/swimresults/start-service/model"
	"os"
)

func VeranstaltungsortPlz() string {
	dat, err := os.ReadFile("assets/Ergebnisdatei.dsv6")
	if err != nil {
		panic(err)
	}
	buf := bytes.NewBuffer(dat)
	r := parser.NewReader(buf)
	res, err := r.Read()
	if err != nil {
		panic(err)
	}
	def := res.(*model.Wettkampfdefinitionsliste)
	return def.Veranstaltungsort.PLZ
}

func ImportDsvDefinitionFile(file string, meeting string) error {
	dat, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	buf := bytes.NewBuffer(dat)
	r := parser.NewReader(buf)
	res, err := r.Read()
	if err != nil {
		return err
	}

	def := res.(*model.Wettkampfdefinitionsliste)

	for _, dsvEvent := range def.Wettkaempfe {
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
		if dsvEvent.Geschlecht == 'D' {
			event.Gender = "MIXED"
		}

		newEvent, created, err3 := ec.ImportEvent(event, dsvEvent.Ausuebung, dsvEvent.Abschnittsnummer)
		if err3 != nil {
			return err3
		}

		if created {
			print("(+) ")
		} else {
			print("( ) ")
		}

		println(newEvent.Number)

	}

	return nil
}

func ImportDsvResultFile(file string, meeting string) error {
	dat, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	buf := bytes.NewBuffer(dat)
	r := parser.NewReader(buf)
	res, err := r.Read()
	if err != nil {
		return err
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
			fmt.Printf(err.Error())
		}
		cs := 'o'
		if created {
			cs = '+'
		}
		fmt.Printf("[ %c ] > id: %s, name: %s, part: %s\n", cs, newTeam.Identifier.String(), newTeam.Name, newTeam.Participation)
	}

	fmt.Printf(" +==============================+ \n")

	resultsImported := 0

	for _, dsvResult := range erg.PNErgebnisse {

		if dsvResult.GrundDerNichtwertung != "" {
			continue
		}

		if unusedEvents[dsvResult.Wettkampfnummer] {
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
		newAthlete, created, err := ac.ImportAthlete(athlete, meeting)
		if err != nil {
			fmt.Printf(err.Error())
		}
		cs := 'o'
		if created {
			cs = '+'
		}
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
			return err2
		}
		if c {
			return fmt.Errorf("start has been created: id: '%s'; event: '%d', athlete: '%s'", newStart.Identifier, newStart.Event, newStart.AthleteName)
		}

		_, _, err3 := sc.ImportResult(start, result)
		if err3 != nil {
			return err3
		}
		resultsImported++

	}
	return nil

}
