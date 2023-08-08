package importer

import (
	athleteClient "github.com/swimresults/athlete-service/client"
	"github.com/swimresults/meeting-service/client"
	startClient "github.com/swimresults/start-service/client"
)

var ec = client.NewEventClient("https://api.swimresults.de/meeting/v1/")
var hc = startClient.NewHeatClient("https://api.swimresults.de/start/v1/")
var sc = startClient.NewStartClient("https://api.swimresults.de/start/v1/")
var dq = startClient.NewDisqualificationClient("https://api.swimresults.de/start/v1/")
var ac = athleteClient.NewAthleteClient("https://api.swimresults.de/athlete/v1/")
var tc = athleteClient.NewTeamClient("https://api.swimresults.de/athlete/v1/")
