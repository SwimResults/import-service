package importer

import (
	"fmt"
	athleteClient "github.com/swimresults/athlete-service/client"
	athleteModel "github.com/swimresults/athlete-service/model"
	"github.com/swimresults/meeting-service/client"
	"github.com/swimresults/meeting-service/model"
	startClient "github.com/swimresults/start-service/client"
	startModel "github.com/swimresults/start-service/model"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var doTheImport = true
var meeting = "ESS23F"

func DoTheMagic() error {
	str, err1 := ReadPdf("../assets/ME_KKJS.pdf")
	if err1 != nil {
		return err1
	}

	c := client.NewEventClient("https://api.swimresults.de/meeting/v1/")
	hc := startClient.NewHeatClient("https://api.swimresults.de/start/v1/")
	sc := startClient.NewStartClient("https://api.swimresults.de/start/v1/")
	ac := athleteClient.NewAthleteClient("https://api.swimresults.de/athlete/v1/")
	tc := athleteClient.NewTeamClient("https://api.swimresults.de/athlete/v1/")

	wks := strings.Split(str, "Wettkampf ")
	for _, s1 := range wks {
		if !strings.Contains(s1, "Meldezeit") {
			continue
		}

		event := model.Event{
			Meeting: meeting,
		}

		var err error

		sa1 := strings.SplitN(s1, "-", 2)
		event.Number, err = strconv.Atoi(strings.Trim(sa1[0], " "))
		if err != nil {
			panic(err)
		}

		sr1 := sa1[1]

		sa2 := strings.SplitN(sr1, "m", 2)
		event.Distance, err = strconv.Atoi(strings.Trim(sa2[0], " "))

		sr2 := sa2[0]
		sr3 := sa2[1]

		if err != nil {
			if strings.Contains(sr2, "x") {
				event.RelayDistance = strings.Trim(sr2, " ")
			} else {
				panic(err)
			}
		}

		if strings.Contains(sr3, "weiblich") {
			event.Gender = "FEMALE"
			sr3 = strings.Replace(sr3, "weiblich", "", 1)
		}

		if strings.Contains(sr3, "männlich") {
			event.Gender = "MALE"
			sr3 = strings.Replace(sr3, "männlich", "", 1)
		}

		if strings.Contains(sr3, "mixed") {
			event.Gender = "MIXED"
			sr3 = strings.Replace(sr3, "mixed", "", 1)
		}

		if strings.Contains(sr3, "Bein") {
			continue
		}

		if strings.Contains(sr3, "staffel") {
			continue
		}

		sr3 = substr(sr3, "(")
		sr3 = substr(sr3, "Final")
		sr3 = substr(sr3, "Lauf")

		style := strings.Trim(sr3, " ")

		if doTheImport {
			newEvent, created, err3 := c.ImportEvent(event, style, 1)
			if err3 != nil {
				println("----------")
				println(err3.Error())
				continue
			}

			if created {
				print("(+) ")
			} else {
				print("( ) ")
			}

			println(newEvent.Number)
		}

		// +----=====[ HEAT ]=====----+

		heats := strings.Split(s1, "Lauf")

		for _, s2 := range heats {
			if !strings.Contains(s2, "Uhr") {
				continue
			}
			heat := startModel.Heat{
				Meeting: meeting,
				Event:   event.Number,
			}

			s3 := substrr(s2, "ca.")
			s3 = substr(s3, "Uhr")

			heat.StartEstimation, _ = time.Parse("15:04", strings.Trim(s3, " "))

			heat.Number, _ = strconv.Atoi(substr(s2, "/"))

			if doTheImport {
				newHeat, created, err5 := hc.ImportHeat(heat)
				if err5 != nil {
					panic(err5)
				}

				if created {
					print("(+) ")
				} else {
					print("( ) ")
				}

				println(newHeat.Number)
				heat = *newHeat
			}

			lanes := strings.Split(s2, "Bahn")
			for _, s4 := range lanes {
				if strings.Contains(s4, "Meldezeit") || strings.Contains(s4, "Uhr") {
					continue
				}

				athlete := athleteModel.Athlete{
					Gender: event.Gender,
				}

				re := regexp.MustCompile("[0-9]+")

				sa3 := re.Split(s4, 3)
				if sa3[1] == "" {
					continue
				}
				athlete.Name = sa3[1]

				y1 := substrr(s4, athlete.Name)
				y1 = strings.Trim(y1, " ")
				y1 = y1[:4]

				athlete.Year, err = strconv.Atoi(y1)
				if err != nil {
					panic(err)
				}

				s5 := substrr(s4, y1)
				s5 = substr(s5, "Erzgebirgsspiele")
				s5 = substr(s5, "---")
				s5 = s5[:len(s5)-8]

				s5 = strings.Trim(s5, " ")

				team := athleteModel.Team{
					Name:    s5,
					Country: "GER",
				}

				s6 := substr(s4, sa3[1])
				s6 = strings.Trim(s6, " ")
				lane, err := strconv.Atoi(s6)
				if err != nil {
					panic(err)
				}

				println(lane)

				if doTheImport {
					newTeam, created, err3 := tc.ImportTeam(team, meeting)
					if err3 != nil {
						panic(err3)
					}

					if created {
						print("(+) ")
						fmt.Println(newTeam)
					}
					athlete.Team.Name = newTeam.Name

					newAthlete, created, err4 := ac.ImportAthlete(athlete, meeting)
					if err4 != nil {
						panic(err4)
					}

					if created {
						print("(+) ")
						fmt.Println(newAthlete)
					}

					start := startModel.Start{
						Meeting:         meeting,
						Event:           event.Number,
						HeatNumber:      heat.Number,
						Lane:            lane,
						Athlete:         newAthlete.Identifier,
						AthleteName:     newAthlete.Name,
						AthleteYear:     newAthlete.Year,
						AthleteTeam:     newTeam.Identifier,
						AthleteTeamName: newTeam.Name,
					}

					newStart, created, err5 := sc.ImportStart(start)
					if err5 != nil {
						panic(err5)
					}

					if created {
						print("(+) ")
						fmt.Println(newStart)
					}

					s7 := substrr(s4, s5)
					s7 = strings.Trim(s7, " ")
					s7 = s7[:8]

					s7 = strings.Replace(s7, ":", "m", 1)
					s7 = strings.Replace(s7, ",", "s", 1)
					s7 = s7 + "0ms"
					dur, err7 := time.ParseDuration(s7)
					if err7 != nil {
						panic(err7)
					}

					result := startModel.Result{
						Time:       dur,
						ResultType: "registration",
					}

					_, created, err8 := sc.ImportResult(*newStart, result)
					if err8 != nil {
						panic(err8)
					}

					if created {
						fmt.Printf("(+) result: %s\n", s7)
					}
				}

			}
		}
	}
	return nil
}

func substr(s string, substr string) string {
	return strings.Trim(strings.SplitN(s, substr, 2)[0], " ")
}

func substrr(s string, substr string) string {
	return strings.Trim(strings.SplitN(s, substr, 2)[1], " ")
}
