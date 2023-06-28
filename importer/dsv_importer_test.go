package importer

import (
	"bytes"
	"fmt"
	"github.com/konrad2002/dsvparser/model"
	"github.com/konrad2002/dsvparser/parser"
	"github.com/swimresults/athlete-service/client"
	athleteModel "github.com/swimresults/athlete-service/model"
	"os"
	"testing"
)

func Test(t *testing.T) {
	dat, err := os.ReadFile("../assets/Ergebnisdatei.dsv6")
	if err != nil {
		panic(err)
	}
	buf := bytes.NewBuffer(dat)
	r := parser.NewReader(buf)
	res, err := r.Read()
	if err != nil {
		panic(err)
	}
	def := res.(*model.Wettkampfergebnisliste)

	c := client.NewTeamClient("https://api.swimresults.de/athlete/v1/")

	for _, verein := range def.Vereine {
		team := athleteModel.Team{
			Name:    verein.Vereinsbezeichnung,
			Country: verein.FinaNationenkuerzel,
			DsvId:   verein.Vereinskennzahl,
			StateId: verein.Landesschwimmverband,
		}
		newTeam, created, err := c.ImportTeam(team, "IESC19")
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

	c2 := client.NewAthleteClient("https://api.swimresults.de/athlete/v1/")

	for _, sportler := range def.PNErgebnisse {
		athlete := athleteModel.Athlete{
			Name:   sportler.Name,
			Year:   sportler.Jahrgang,
			Gender: string(sportler.Geschlecht),
			DsvId:  sportler.DsvId,
			Team: athleteModel.Team{
				DsvId: sportler.Vereinskennzahl,
				Name:  sportler.Verein,
			},
		}
		newAthlete, created, err := c2.ImportAthlete(athlete, "IESC19")
		if err != nil {
			fmt.Printf(err.Error())
		}
		cs := 'o'
		if created {
			cs = '+'
		}
		fmt.Printf("[ %c ] > id: %s, name: %s, part: %s\n", cs, newAthlete.Identifier.String(), newAthlete.Name, newAthlete.Participation)
	}
}
