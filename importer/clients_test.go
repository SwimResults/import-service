package importer

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/swimresults/swimresults-api/go/athleteapi"
	"github.com/swimresults/swimresults-api/go/meetingapi"
	"github.com/swimresults/swimresults-api/go/startapi"
)

type recorded struct {
	method, path, serviceKey string
	body                     map[string]any
}

// stub answers every request with status/body and records the last request.
func stub(t *testing.T, status int, body string) (*httptest.Server, *recorded) {
	t.Helper()
	rec := &recorded{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rec.method, rec.path, rec.serviceKey = r.Method, r.URL.Path, r.Header.Get("X-Swimresults-Service")
		b, _ := io.ReadAll(r.Body)
		rec.body = nil
		_ = json.Unmarshal(b, &rec.body)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		_, _ = io.WriteString(w, body)
	}))
	t.Cleanup(srv.Close)
	return srv, rec
}

func TestImportStartCreatedUpdatedAndError(t *testing.T) {
	t.Setenv("SR_SERVICE_KEY", "test-key")
	for _, tc := range []struct {
		status      int
		wantCreated bool
		wantErr     bool
	}{{201, true, false}, {200, false, false}, {500, false, true}} {
		srv, rec := stub(t, tc.status, `{"_id":"650c2e653df25edacb51bd3b","meeting":"IESC19","lane":4}`)
		c := &startClient{newStartAPI(srv.URL + "/")} // env URLs end with "/"
		got, created, err := c.ImportStart(startapi.Start{Meeting: "IESC19", Lane: 4})

		if rec.method != http.MethodPost || rec.path != "/start/import" {
			t.Errorf("request = %s %s, want POST /start/import", rec.method, rec.path)
		}
		if rec.serviceKey != "test-key" {
			t.Errorf("X-Swimresults-Service = %q", rec.serviceKey)
		}
		if start, ok := rec.body["start"].(map[string]any); !ok || start["meeting"] != "IESC19" || start["lane"] != float64(4) {
			t.Errorf("body = %v, want {start: {meeting: IESC19, lane: 4, ...}}", rec.body)
		}
		if (err != nil) != tc.wantErr || created != tc.wantCreated {
			t.Errorf("status %d: created=%v err=%v", tc.status, created, err)
		}
		if !tc.wantErr && (got == nil || got.Identifier != "650c2e653df25edacb51bd3b") {
			t.Errorf("status %d: got %+v", tc.status, got)
		}
	}
}

func TestImportResultSendsStartAndResult(t *testing.T) {
	srv, rec := stub(t, 201, `{"time":61230000000,"result_type":"result_list"}`)
	c := &startClient{newStartAPI(srv.URL)}
	got, created, err := c.ImportResult(startapi.Start{Meeting: "M"}, startapi.Result{Time: 61230 * time.Millisecond, ResultType: "result_list"})
	if err != nil || !created || got.Time != 61230*time.Millisecond {
		t.Fatalf("got=%+v created=%v err=%v", got, created, err)
	}
	if rec.path != "/result/import" || rec.body["result"].(map[string]any)["time"] != float64(61230000000) {
		t.Errorf("path=%s body=%v", rec.path, rec.body)
	}
}

func TestSetHeatStartPath(t *testing.T) {
	srv, rec := stub(t, 200, `{"meeting":"M","event":3,"number":2}`)
	c := &heatClient{newStartAPI(srv.URL)}
	h, err := c.SetHeatStart("M", 3, 2)
	if err != nil || h.Number != 2 {
		t.Fatalf("h=%+v err=%v", h, err)
	}
	if rec.method != http.MethodPost || rec.path != "/heat/meet/M/event/3/heat/2/start" {
		t.Errorf("request = %s %s", rec.method, rec.path)
	}
}

func TestImportDisqualificationReportsCreatedOn200(t *testing.T) {
	srv, rec := stub(t, 200, `{"reason":"false start","type":"DSQ"}`)
	c := &disqualificationClient{newStartAPI(srv.URL)}
	d, created, err := c.ImportDisqualification(startapi.Start{Meeting: "M"}, "false start", "DSQ", time.Time{})
	if err != nil || !created || d.Reason != "false start" {
		t.Fatalf("d=%+v created=%v err=%v", d, created, err)
	}
	if rec.path != "/disqualification/import" || rec.body["disqualification"].(map[string]any)["type"] != "DSQ" {
		t.Errorf("path=%s body=%v", rec.path, rec.body)
	}
}

func TestGetAthletesByMeeting404IsEmpty(t *testing.T) {
	srv, _ := stub(t, 404, `{"message":"not found"}`)
	c := &athleteClient{newAthleteAPI(srv.URL)}
	athletes, err := c.GetAthletesByMeeting("M")
	if err != nil || athletes != nil {
		t.Fatalf("athletes=%v err=%v, want nil,nil", athletes, err)
	}
}

func TestGetAthletesByMeetingNoPaging(t *testing.T) {
	rec := &recorded{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rec.path = r.URL.RequestURI()
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `[{"_id":"650c2e653df25edacb51bd3b","name":"A B","year":2010}]`)
	}))
	defer srv.Close()
	c := &athleteClient{newAthleteAPI(srv.URL)}
	athletes, err := c.GetAthletesByMeeting("IESC19")
	if err != nil || len(athletes) != 1 || athletes[0].Year != 2010 {
		t.Fatalf("athletes=%v err=%v", athletes, err)
	}
	if rec.path != "/athlete/meet/IESC19" { // no limit/offset: the service returns all athletes
		t.Errorf("request URI = %s", rec.path)
	}
}

func TestGetTeamsByMeeting(t *testing.T) {
	srv, _ := stub(t, 404, `{"message":"not found"}`)
	teams, found, err := (&teamClient{newAthleteAPI(srv.URL)}).GetTeamsByMeeting("M")
	if err != nil || found || teams != nil {
		t.Fatalf("404: teams=%v found=%v err=%v", teams, found, err)
	}
	srv2, _ := stub(t, 200, `[{"name":"SC Test"}]`)
	teams, found, err = (&teamClient{newAthleteAPI(srv2.URL)}).GetTeamsByMeeting("M")
	if err != nil || !found || len(*teams) != 1 || (*teams)[0].Name != "SC Test" {
		t.Fatalf("200: teams=%v found=%v err=%v", teams, found, err)
	}
}

func TestImportCertificateOnly201(t *testing.T) {
	srv, rec := stub(t, 201, `{"name":"Urkunde","path":"/files/x.pdf"}`)
	c := &athleteClient{newAthleteAPI(srv.URL)}
	cert, created, err := c.ImportCertificate("Urkunde", "650c2e653df25edacb51bd3b", "M", "/files/x.pdf")
	if err != nil || !created || cert.Name != "Urkunde" {
		t.Fatalf("cert=%+v created=%v err=%v", cert, created, err)
	}
	if rec.path != "/certificate/import" || rec.body["athlete_id"] != "650c2e653df25edacb51bd3b" {
		t.Errorf("path=%s body=%v", rec.path, rec.body)
	}
}

func TestImportTeamAndAthlete(t *testing.T) {
	srv, rec := stub(t, 200, `{"name":"SC Test"}`)
	team, created, err := (&teamClient{newAthleteAPI(srv.URL)}).ImportTeam(athleteapi.Team{Name: "SC Test"}, "M")
	if err != nil || created || team.Name != "SC Test" || rec.path != "/team/import" || rec.body["meeting"] != "M" {
		t.Fatalf("team=%+v created=%v err=%v path=%s body=%v", team, created, err, rec.path, rec.body)
	}
	srv2, rec2 := stub(t, 201, `{"name":"A B","year":2010}`)
	ath, created, err := (&athleteClient{newAthleteAPI(srv2.URL)}).ImportAthlete(athleteapi.Athlete{Name: "A B", Year: 2010}, "M")
	if err != nil || !created || ath.Year != 2010 || rec2.path != "/athlete/import" {
		t.Fatalf("ath=%+v created=%v err=%v path=%s", ath, created, err, rec2.path)
	}
}

func TestEventClient(t *testing.T) {
	srv, rec := stub(t, 201, `{"number":7,"distance":100}`)
	ev, created, err := newEventClient(srv.URL).ImportEvent(meetingapi.Event{Number: 7, Distance: 100}, "FREESTYLE", 2)
	if err != nil || !created || ev.Number != 7 {
		t.Fatalf("ev=%+v created=%v err=%v", ev, created, err)
	}
	if rec.path != "/event/import" || rec.body["style_name"] != "FREESTYLE" || rec.body["part_number"] != float64(2) {
		t.Errorf("path=%s body=%v", rec.path, rec.body)
	}
	srv2, rec2 := stub(t, 200, `{"number":7}`)
	ev, err = newEventClient(srv2.URL).GetEventByMeetingAndNumber("M", 7)
	if err != nil || ev.Number != 7 || rec2.method != http.MethodGet || rec2.path != "/event/meet/M/event/7" {
		t.Fatalf("ev=%+v err=%v req=%s %s", ev, err, rec2.method, rec2.path)
	}
}

func TestIsZeroID(t *testing.T) {
	for id, want := range map[string]bool{"": true, "000000000000000000000000": true, "650c2e653df25edacb51bd3b": false} {
		if isZeroID(id) != want {
			t.Errorf("isZeroID(%q) != %v", id, want)
		}
	}
}
