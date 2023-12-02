package importer

import (
	"fmt"
	athleteClient "github.com/swimresults/athlete-service/client"
	"github.com/swimresults/meeting-service/client"
	startClient "github.com/swimresults/start-service/client"
	"os"
	"strings"
)

//var startServiceUrl = "https://api.swimresults.de/start/v1/"
//var athleteServiceUrl = "https://api.swimresults.de/athlete/v1/"
//var meetingServiceUrl = "https://api.swimresults.de/meeting/v1/"

//var startServiceUrl = "http://localhost:8087/"
//var athleteServiceUrl = "http://localhost:8086/"
//var meetingServiceUrl = "http://localhost:8089/"

var startServiceUrl = os.Getenv("SR_IMPORT_START_URL")
var athleteServiceUrl = os.Getenv("SR_IMPORT_ATHLETE_URL")
var meetingServiceUrl = os.Getenv("SR_IMPORT_MEETING_URL")

var ec = client.NewEventClient(meetingServiceUrl)
var gc = client.NewAgeGroupClient(meetingServiceUrl)
var hc = startClient.NewHeatClient(startServiceUrl)
var sc = startClient.NewStartClient(startServiceUrl)
var dc = startClient.NewDisqualificationClient(startServiceUrl)
var ac = athleteClient.NewAthleteClient(athleteServiceUrl)
var tc = athleteClient.NewTeamClient(athleteServiceUrl)

func IsEventImportable(ev int, ex []int, in []int) bool {
	if ex != nil && len(ex) > 0 {
		for _, e := range ex {
			if ev == e { // in exclude list -> next
				return false
			}
		}
	}

	if in != nil && len(in) > 0 {
		for _, e := range in {
			if ev == e {
				return true
			}
		}
		return false
	}

	return true

}

// substr returns everything from s until sep appears
func substr(s string, sep string) string {
	return strings.Trim(strings.SplitN(s, sep, 2)[0], " ")
}

// substrr returns everything from s that comes after sep; if sep does not occur, s is returned
func substrr(s string, substr string) string {
	s1 := strings.SplitN(s, substr, 2)
	s2 := ""
	if len(s1) > 1 {
		s2 = s1[1]
	}
	return strings.Trim(s2, " ")
}

func trim(s string) string {
	return strings.Trim(s, " ")
}

func runImport() bool {
	return os.Getenv("SR_NO_IMPORT") == ""
}

func importError(text string, err error) {
	fmt.Println("\033[31m" + text + "\033[33m")
	fmt.Println(err)
	fmt.Println("\033[0m")
}
