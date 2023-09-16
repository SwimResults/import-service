package importer

import (
	athleteClient "github.com/swimresults/athlete-service/client"
	"github.com/swimresults/meeting-service/client"
	startClient "github.com/swimresults/start-service/client"
)

//var startServiceUrl = "https://api.swimresults.de/start/v1/"
//var athleteServiceUrl = "https://api.swimresults.de/athlete/v1/"
//var meetingServiceUrl = "https://api.swimresults.de/meeting/v1/"

var startServiceUrl = "http://localhost:8087/"
var athleteServiceUrl = "http://localhost:8086/"
var meetingServiceUrl = "http://localhost:8089/"

var ec = client.NewEventClient(meetingServiceUrl)
var hc = startClient.NewHeatClient(startServiceUrl)
var sc = startClient.NewStartClient(startServiceUrl)
var dq = startClient.NewDisqualificationClient(startServiceUrl)
var ac = athleteClient.NewAthleteClient(athleteServiceUrl)
var tc = athleteClient.NewTeamClient(athleteServiceUrl)

func IsEventImportable(ev int, ex []int, in []int) bool {
	if ex != nil {
		for _, e := range ex {
			if ev == e { // in exclude list -> next
				return false
			}
		}
	}

	if in != nil {
		for _, e := range in {
			if ev == e {
				return true
			}
		}
		return false
	}

	return true

}
