package importer

// Clients for the meeting, athlete and start services, built on the generated packages from
// github.com/swimresults/swimresults-api. They keep the method set and semantics of the per-service
// client packages they replace (same endpoints, same "created" flag, same 404 handling), so the importers
// only changed their model imports.

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/swimresults/swimresults-api/go/athleteapi"
	"github.com/swimresults/swimresults-api/go/meetingapi"
	"github.com/swimresults/swimresults-api/go/startapi"
)

const requestTimeout = 60 * time.Second

var httpClient = &http.Client{Timeout: requestTimeout}

// serviceKey authenticates import-service against the other services (X-Swimresults-Service),
// like the service-core client did.
func serviceKey(_ context.Context, req *http.Request) error {
	if key := os.Getenv("SR_SERVICE_KEY"); key != "" {
		req.Header.Set("X-Swimresults-Service", key)
	}
	return nil
}

// baseURL accepts the existing env values ("http://start-service:8080/"); the generated clients
// append "/<path>" themselves.
func baseURL(u string) string { return strings.TrimRight(u, "/") }

func ctx() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), requestTimeout)
}

// importResult maps the services' import responses: 201 = created, 200 = updated.
func importResult[T any](status int, created, updated *T, what string) (*T, bool, error) {
	switch {
	case status == http.StatusCreated && created != nil:
		return created, true, nil
	case status == http.StatusOK && updated != nil:
		return updated, false, nil
	default:
		return nil, false, fmt.Errorf("import %s request returned: %d", what, status)
	}
}

// isZeroID reports an unset ObjectId: "" or the zero ObjectID, which the services both treat as nil.
func isZeroID(id string) bool { return id == "" || id == "000000000000000000000000" }

func mustClient[C any](c C, err error) C {
	if err != nil {
		panic(err) // only fails for unparsable base URLs
	}
	return c
}

// ---------------------------------------------------------------- meeting-service

type eventClient struct {
	c *meetingapi.ClientWithResponses
}

func newEventClient(url string) *eventClient {
	return &eventClient{mustClient(meetingapi.NewClientWithResponses(baseURL(url),
		meetingapi.WithHTTPClient(httpClient), meetingapi.WithRequestEditorFn(serviceKey)))}
}

func (e *eventClient) ImportEvent(event meetingapi.Event, styleName string, partNumber int) (*meetingapi.Event, bool, error) {
	c, cancel := ctx()
	defer cancel()
	res, err := e.c.ImportEventWithResponse(c, meetingapi.ImportEventRequestDto{Event: event, StyleName: styleName, MeetingPartNumber: partNumber})
	if err != nil {
		return nil, false, err
	}
	return importResult(res.StatusCode(), res.JSON201, res.JSON200, "event")
}

func (e *eventClient) GetEventByMeetingAndNumber(meeting string, number int) (*meetingapi.Event, error) {
	c, cancel := ctx()
	defer cancel()
	res, err := e.c.GetEventByMeetingAndNumberWithResponse(c, meeting, number)
	if err != nil {
		return nil, err
	}
	if res.JSON200 == nil {
		return nil, fmt.Errorf("GetEventByMeetingAndNumber received error: %d", res.StatusCode())
	}
	return res.JSON200, nil
}

// ---------------------------------------------------------------- athlete-service

type athleteClient struct {
	c *athleteapi.ClientWithResponses
}
type teamClient struct {
	c *athleteapi.ClientWithResponses
}

func newAthleteAPI(url string) *athleteapi.ClientWithResponses {
	return mustClient(athleteapi.NewClientWithResponses(baseURL(url),
		athleteapi.WithHTTPClient(httpClient), athleteapi.WithRequestEditorFn(serviceKey)))
}

func (a *athleteClient) ImportAthlete(athlete athleteapi.Athlete, meeting string) (*athleteapi.Athlete, bool, error) {
	c, cancel := ctx()
	defer cancel()
	res, err := a.c.ImportAthleteWithResponse(c, athleteapi.ImportAthleteRequestDto{Athlete: athlete, Meeting: meeting})
	if err != nil {
		return nil, false, err
	}
	return importResult(res.StatusCode(), res.JSON201, res.JSON200, "athlete")
}

func (a *athleteClient) ImportCertificate(name string, athleteId athleteapi.ObjectId, meeting string, path string) (*athleteapi.Certificate, bool, error) {
	c, cancel := ctx()
	defer cancel()
	res, err := a.c.ImportCertificateWithResponse(c, athleteapi.ImportCertificateRequestDto{Name: name, AthleteId: athleteId, Meeting: meeting, Path: path})
	if err != nil {
		return nil, false, err
	}
	return importResult(res.StatusCode(), res.JSON201, (*athleteapi.Certificate)(nil), "certificate") // only 201
}

// GetAthletesByMeeting returns all athletes of a meeting (no paging), nil if the meeting has none (404).
func (a *athleteClient) GetAthletesByMeeting(meeting string) ([]athleteapi.Athlete, error) {
	c, cancel := ctx()
	defer cancel()
	res, err := a.c.GetAthletesByMeetingWithResponse(c, meeting, &athleteapi.GetAthletesByMeetingParams{})
	if err != nil {
		return nil, err
	}
	switch {
	case res.StatusCode() == http.StatusNotFound:
		return nil, nil
	case res.JSON200 == nil:
		return nil, fmt.Errorf("GetAthletesByMeeting received error: %d", res.StatusCode())
	}
	return *res.JSON200, nil
}

func (t *teamClient) ImportTeam(team athleteapi.Team, meeting string) (*athleteapi.Team, bool, error) {
	c, cancel := ctx()
	defer cancel()
	res, err := t.c.ImportTeamWithResponse(c, athleteapi.ImportTeamRequestDto{Team: team, Meeting: meeting})
	if err != nil {
		return nil, false, err
	}
	return importResult(res.StatusCode(), res.JSON201, res.JSON200, "team")
}

// GetTeamsByMeeting returns (teams, true) or (nil, false) if the meeting has none (404).
func (t *teamClient) GetTeamsByMeeting(meeting string) (*[]athleteapi.Team, bool, error) {
	c, cancel := ctx()
	defer cancel()
	res, err := t.c.GetTeamsByMeetingWithResponse(c, meeting, &athleteapi.GetTeamsByMeetingParams{})
	if err != nil {
		return nil, false, err
	}
	switch {
	case res.StatusCode() == http.StatusNotFound:
		return nil, false, nil
	case res.JSON200 == nil:
		return nil, false, fmt.Errorf("GetTeamsByMeeting received error: %d", res.StatusCode())
	}
	return res.JSON200, true, nil
}

// ---------------------------------------------------------------- start-service

type startClient struct{ c *startapi.ClientWithResponses }
type heatClient struct{ c *startapi.ClientWithResponses }
type disqualificationClient struct{ c *startapi.ClientWithResponses }
type rankingClient struct{ c *startapi.ClientWithResponses }

func newStartAPI(url string) *startapi.ClientWithResponses {
	return mustClient(startapi.NewClientWithResponses(baseURL(url),
		startapi.WithHTTPClient(httpClient), startapi.WithRequestEditorFn(serviceKey)))
}

func (s *startClient) ImportStart(start startapi.Start) (*startapi.Start, bool, error) {
	c, cancel := ctx()
	defer cancel()
	res, err := s.c.ImportStartWithResponse(c, startapi.ImportStartRequestDto{Start: start})
	if err != nil {
		return nil, false, err
	}
	return importResult(res.StatusCode(), res.JSON201, res.JSON200, "start")
}

func (s *startClient) ImportResult(start startapi.Start, result startapi.Result) (*startapi.Result, bool, error) {
	c, cancel := ctx()
	defer cancel()
	res, err := s.c.ImportResultWithResponse(c, startapi.ImportResultRequestDto{Start: start, Result: result})
	if err != nil {
		return nil, false, err
	}
	return importResult(res.StatusCode(), res.JSON201, res.JSON200, "result")
}

func (h *heatClient) ImportHeat(heat startapi.Heat) (*startapi.Heat, bool, error) {
	c, cancel := ctx()
	defer cancel()
	res, err := h.c.ImportHeatWithResponse(c, startapi.ImportHeatRequestDto{Heat: heat})
	if err != nil {
		return nil, false, err
	}
	return importResult(res.StatusCode(), res.JSON201, res.JSON200, "heat")
}

func (h *heatClient) SetHeatStart(meeting string, event int, number int) (*startapi.Heat, error) {
	c, cancel := ctx()
	defer cancel()
	res, err := h.c.TriggerHeatStartWithResponse(c, meeting, event, number)
	if err != nil {
		return nil, err
	}
	if res.JSON200 == nil {
		return nil, fmt.Errorf("set heat start request returned: %d", res.StatusCode())
	}
	return res.JSON200, nil
}

// ImportDisqualification: the service answers 200 for both new and updated disqualifications;
// like before, success is reported as created.
func (d *disqualificationClient) ImportDisqualification(start startapi.Start, reason string, disqType string, timeOfAnnouncement time.Time) (*startapi.Disqualification, bool, error) {
	c, cancel := ctx()
	defer cancel()
	res, err := d.c.ImportDisqualificationWithResponse(c, startapi.ImportDisqualificationRequestDto{
		Start: start,
		Disqualification: startapi.Disqualification{
			Reason:           reason,
			AnnouncementTime: timeOfAnnouncement,
			Type:             disqType,
		},
	})
	if err != nil {
		return nil, false, err
	}
	if res.JSON200 == nil {
		return nil, false, fmt.Errorf("import disqualification request returned: %d", res.StatusCode())
	}
	return res.JSON200, true, nil
}

func (r *rankingClient) ImportRanking(ranking startapi.Ranking) (*startapi.Ranking, bool, error) {
	c, cancel := ctx()
	defer cancel()
	res, err := r.c.ImportRankingWithResponse(c, startapi.ImportRankingRequestDto{Ranking: ranking})
	if err != nil {
		return nil, false, err
	}
	return importResult(res.StatusCode(), res.JSON201, res.JSON200, "ranking")
}
