package oack

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// testClient creates a test server + client. The handler receives (method, path, body).
func testClient(t *testing.T, handler func(w http.ResponseWriter, r *http.Request)) *Client {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(handler))
	t.Cleanup(srv.Close)
	return New(BearerToken("tok"), WithBaseURL(srv.URL))
}

func assertPath(t *testing.T, r *http.Request, want string) {
	t.Helper()
	if r.URL.Path != want {
		t.Errorf("path: got %q, want %q", r.URL.Path, want)
	}
}

func assertMethod(t *testing.T, r *http.Request, want string) {
	t.Helper()
	if r.Method != want {
		t.Errorf("method: got %q, want %q", r.Method, want)
	}
}

// ---------------------------------------------------------------------------
// Accounts
// ---------------------------------------------------------------------------

func TestCreateAccount(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertPath(t, r, "/api/v1/accounts")
		assertMethod(t, r, http.MethodPost)
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(Account{ID: "a1", Name: "Test"})
	})
	a, err := c.CreateAccount(context.Background(), "Test")
	if err != nil {
		t.Fatal(err)
	}
	if a.ID != "a1" {
		t.Errorf("got %q", a.ID)
	}
}

func TestListAccounts(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertPath(t, r, "/api/v1/accounts")
		_ = json.NewEncoder(w).Encode([]Account{{ID: "a1"}, {ID: "a2"}})
	})
	accounts, err := c.ListAccounts(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(accounts) != 2 {
		t.Errorf("got %d", len(accounts))
	}
}

func TestGetAccount(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertPath(t, r, "/api/v1/accounts/a1")
		_ = json.NewEncoder(w).Encode(Account{ID: "a1", Name: "My Org"})
	})
	a, err := c.GetAccount(context.Background(), "a1")
	if err != nil {
		t.Fatal(err)
	}
	if a.Name != "My Org" {
		t.Errorf("got %q", a.Name)
	}
}

func TestUpdateAccount(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertPath(t, r, "/api/v1/accounts/a1")
		assertMethod(t, r, http.MethodPut)
		_ = json.NewEncoder(w).Encode(Account{ID: "a1", Name: "New Name"})
	})
	a, err := c.UpdateAccount(context.Background(), "a1", "New Name")
	if err != nil {
		t.Fatal(err)
	}
	if a.Name != "New Name" {
		t.Errorf("got %q", a.Name)
	}
}

func TestDeleteAccount(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertPath(t, r, "/api/v1/accounts/a1")
		assertMethod(t, r, http.MethodDelete)
		w.WriteHeader(http.StatusNoContent)
	})
	if err := c.DeleteAccount(context.Background(), "a1"); err != nil {
		t.Fatal(err)
	}
}

func TestListAccountMembers(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertPath(t, r, "/api/v1/accounts/a1/members")
		_ = json.NewEncoder(w).Encode([]AccountMember{{UserID: "u1", Role: "owner"}})
	})
	members, err := c.ListAccountMembers(context.Background(), "a1")
	if err != nil {
		t.Fatal(err)
	}
	if len(members) != 1 {
		t.Errorf("got %d", len(members))
	}
}

func TestCreateAccountInvite(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertPath(t, r, "/api/v1/accounts/a1/invites")
		assertMethod(t, r, http.MethodPost)
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(AccountInvite{ID: "inv1", Email: "x@y.com"})
	})
	inv, err := c.CreateAccountInvite(context.Background(), "a1", "x@y.com", "member")
	if err != nil {
		t.Fatal(err)
	}
	if inv.ID != "inv1" {
		t.Errorf("got %q", inv.ID)
	}
}

func TestCreateAccountAPIKey(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertPath(t, r, "/api/v1/accounts/acc1/api-keys")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(CreateAccountAPIKeyResult{
			Key:    "oack_acc_abc123",
			APIKey: AccountAPIKey{ID: "k1", Name: "Terraform", KeyPrefix: "oack_acc_abc"},
		})
	})
	result, err := c.CreateAccountAPIKey(context.Background(), "acc1", CreateAccountAPIKeyParams{Name: "Terraform"})
	if err != nil {
		t.Fatal(err)
	}
	if result.Key != "oack_acc_abc123" {
		t.Errorf("got %q", result.Key)
	}
}

func TestListAccountAPIKeys(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertPath(t, r, "/api/v1/accounts/acc1/api-keys")
		_ = json.NewEncoder(w).Encode([]AccountAPIKey{{ID: "k1"}, {ID: "k2"}})
	})
	keys, err := c.ListAccountAPIKeys(context.Background(), "acc1")
	if err != nil {
		t.Fatal(err)
	}
	if len(keys) != 2 {
		t.Errorf("got %d", len(keys))
	}
}

func TestDeleteAccountAPIKey(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertPath(t, r, "/api/v1/accounts/acc1/api-keys/k1")
		assertMethod(t, r, http.MethodDelete)
		w.WriteHeader(http.StatusNoContent)
	})
	if err := c.DeleteAccountAPIKey(context.Background(), "acc1", "k1"); err != nil {
		t.Fatal(err)
	}
}

// ---------------------------------------------------------------------------
// Teams
// ---------------------------------------------------------------------------

func TestCreateTeam(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertPath(t, r, "/api/v1/accounts/a1/teams")
		assertMethod(t, r, http.MethodPost)
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(Team{ID: "t1", Name: "Dev"})
	})
	team, err := c.CreateTeam(context.Background(), "a1", "Dev")
	if err != nil {
		t.Fatal(err)
	}
	if team.Name != "Dev" {
		t.Errorf("got %q", team.Name)
	}
}

func TestListTeams(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertPath(t, r, "/api/v1/teams")
		_ = json.NewEncoder(w).Encode([]Team{{ID: "t1"}, {ID: "t2"}})
	})
	teams, err := c.ListTeams(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(teams) != 2 {
		t.Errorf("got %d", len(teams))
	}
}

func TestGetTeam(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertPath(t, r, "/api/v1/teams/t1")
		_ = json.NewEncoder(w).Encode(Team{ID: "t1", Name: "Dev"})
	})
	team, err := c.GetTeam(context.Background(), "t1")
	if err != nil {
		t.Fatal(err)
	}
	if team.ID != "t1" {
		t.Errorf("got %q", team.ID)
	}
}

func TestUpdateTeam(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertPath(t, r, "/api/v1/teams/t1")
		assertMethod(t, r, http.MethodPut)
		_ = json.NewEncoder(w).Encode(Team{ID: "t1", Name: "Ops"})
	})
	team, err := c.UpdateTeam(context.Background(), "t1", "Ops")
	if err != nil {
		t.Fatal(err)
	}
	if team.Name != "Ops" {
		t.Errorf("got %q", team.Name)
	}
}

func TestDeleteTeam(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertPath(t, r, "/api/v1/teams/t1")
		assertMethod(t, r, http.MethodDelete)
		w.WriteHeader(http.StatusNoContent)
	})
	if err := c.DeleteTeam(context.Background(), "t1"); err != nil {
		t.Fatal(err)
	}
}

func TestListMembers(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertPath(t, r, "/api/v1/teams/t1/members")
		_ = json.NewEncoder(w).Encode([]TeamMember{{UserID: "u1", Role: "admin"}})
	})
	members, err := c.ListMembers(context.Background(), "t1")
	if err != nil {
		t.Fatal(err)
	}
	if len(members) != 1 {
		t.Errorf("got %d", len(members))
	}
}

func TestAddMember(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertPath(t, r, "/api/v1/teams/t1/members")
		assertMethod(t, r, http.MethodPost)
		_ = json.NewEncoder(w).Encode(TeamMember{UserID: "u2", Role: "member"})
	})
	m, err := c.AddMember(context.Background(), "t1", "u2", "member")
	if err != nil {
		t.Fatal(err)
	}
	if m.Role != "member" {
		t.Errorf("got %q", m.Role)
	}
}

func TestRemoveMember(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertPath(t, r, "/api/v1/teams/t1/members/u1")
		assertMethod(t, r, http.MethodDelete)
		w.WriteHeader(http.StatusNoContent)
	})
	if err := c.RemoveMember(context.Background(), "t1", "u1"); err != nil {
		t.Fatal(err)
	}
}

func TestCreateTeamAPIKey(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertPath(t, r, "/api/v1/teams/t1/api-keys")
		assertMethod(t, r, http.MethodPost)
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(CreateTeamAPIKeyResult{
			Key:    "oack_tk_abc",
			APIKey: TeamAPIKey{ID: "k1", Name: "CI"},
		})
	})
	result, err := c.CreateTeamAPIKey(context.Background(), "t1", &CreateTeamAPIKeyParams{Name: "CI"})
	if err != nil {
		t.Fatal(err)
	}
	if result.Key != "oack_tk_abc" {
		t.Errorf("got %q", result.Key)
	}
}

func TestListTeamAPIKeys(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertPath(t, r, "/api/v1/teams/t1/api-keys")
		_ = json.NewEncoder(w).Encode([]TeamAPIKey{{ID: "k1"}})
	})
	keys, err := c.ListTeamAPIKeys(context.Background(), "t1")
	if err != nil {
		t.Fatal(err)
	}
	if len(keys) != 1 {
		t.Errorf("got %d", len(keys))
	}
}

func TestDeleteTeamAPIKey(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertPath(t, r, "/api/v1/teams/t1/api-keys/k1")
		assertMethod(t, r, http.MethodDelete)
		w.WriteHeader(http.StatusNoContent)
	})
	if err := c.DeleteTeamAPIKey(context.Background(), "t1", "k1"); err != nil {
		t.Fatal(err)
	}
}

// ---------------------------------------------------------------------------
// Monitors
// ---------------------------------------------------------------------------

func TestCreateMonitor(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertPath(t, r, "/api/v1/teams/t1/monitors")
		assertMethod(t, r, http.MethodPost)
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(Monitor{ID: "m1", Name: "API", Type: "http"})
	})
	m, err := c.CreateMonitor(context.Background(), "t1", &CreateMonitorParams{
		Name: "API", URL: "https://example.com",
	})
	if err != nil {
		t.Fatal(err)
	}
	if m.ID != "m1" {
		t.Errorf("got %q", m.ID)
	}
}

func TestCreateMonitor_WithBrowserConfig(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		var req CreateMonitorParams
		_ = json.NewDecoder(r.Body).Decode(&req)
		if req.Type != "browser" {
			t.Errorf("type: got %q", req.Type)
		}
		if req.BrowserConfig == nil || req.BrowserConfig.Mode != "script" {
			t.Error("browser_config missing or wrong mode")
		}
		if len(req.Locations) != 2 {
			t.Errorf("locations: got %d", len(req.Locations))
		}
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(Monitor{ID: "m1", Type: "browser", AggregateFailureMode: "at_least_n"})
	})
	m, err := c.CreateMonitor(context.Background(), "t1", &CreateMonitorParams{
		Name: "Browser", URL: "https://example.com", Type: "browser",
		BrowserConfig: &BrowserConfig{Mode: "script", Script: "page.goto('https://example.com')"},
		Locations: []LocationParams{
			{CheckerRegion: "North America", Label: "Dallas"},
			{CheckerRegion: "Europe", Label: "Amsterdam"},
		},
		AggregateFailureMode: "at_least_n", AggregateFailureCount: 2,
	})
	if err != nil {
		t.Fatal(err)
	}
	if m.Type != "browser" {
		t.Errorf("got %q", m.Type)
	}
}

func TestListMonitors(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertPath(t, r, "/api/v1/teams/t1/monitors")
		_ = json.NewEncoder(w).Encode([]Monitor{{ID: "m1"}, {ID: "m2"}})
	})
	monitors, err := c.ListMonitors(context.Background(), "t1")
	if err != nil {
		t.Fatal(err)
	}
	if len(monitors) != 2 {
		t.Errorf("got %d", len(monitors))
	}
}

func TestGetMonitor(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertPath(t, r, "/api/v1/teams/t1/monitors/m1")
		_ = json.NewEncoder(w).Encode(Monitor{ID: "m1", Name: "API"})
	})
	m, err := c.GetMonitor(context.Background(), "t1", "m1")
	if err != nil {
		t.Fatal(err)
	}
	if m.Name != "API" {
		t.Errorf("got %q", m.Name)
	}
}

func TestDeleteMonitor(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertPath(t, r, "/api/v1/teams/t1/monitors/m1")
		assertMethod(t, r, http.MethodDelete)
		w.WriteHeader(http.StatusNoContent)
	})
	if err := c.DeleteMonitor(context.Background(), "t1", "m1"); err != nil {
		t.Fatal(err)
	}
}

func TestPauseMonitor(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertPath(t, r, "/api/v1/teams/t1/monitors/m1/pause")
		assertMethod(t, r, http.MethodPost)
		_ = json.NewEncoder(w).Encode(Monitor{ID: "m1", Status: "paused"})
	})
	m, err := c.PauseMonitor(context.Background(), "t1", "m1")
	if err != nil {
		t.Fatal(err)
	}
	if m.Status != "paused" {
		t.Errorf("got %q", m.Status)
	}
}

func TestUnpauseMonitor(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertPath(t, r, "/api/v1/teams/t1/monitors/m1/unpause")
		assertMethod(t, r, http.MethodPost)
		_ = json.NewEncoder(w).Encode(Monitor{ID: "m1", Status: "active"})
	})
	m, err := c.UnpauseMonitor(context.Background(), "t1", "m1")
	if err != nil {
		t.Fatal(err)
	}
	if m.Status != "active" {
		t.Errorf("got %q", m.Status)
	}
}

func TestDuplicateMonitor(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertPath(t, r, "/api/v1/teams/t1/monitors/m1/duplicate")
		assertMethod(t, r, http.MethodPost)
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(Monitor{ID: "m2", Name: "API (copy)"})
	})
	m, err := c.DuplicateMonitor(context.Background(), "t1", "m1")
	if err != nil {
		t.Fatal(err)
	}
	if m.ID != "m2" {
		t.Errorf("got %q", m.ID)
	}
}

func TestMoveMonitor(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertPath(t, r, "/api/v1/teams/t1/monitors/m1/move")
		assertMethod(t, r, http.MethodPost)
		_ = json.NewEncoder(w).Encode(Monitor{ID: "m1", TeamID: "t2"})
	})
	m, err := c.MoveMonitor(context.Background(), "t1", "m1", "t2")
	if err != nil {
		t.Fatal(err)
	}
	if m.TeamID != "t2" {
		t.Errorf("got %q", m.TeamID)
	}
}

func TestTestMonitorAlert(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertPath(t, r, "/api/v1/teams/t1/monitors/m1/test-alert")
		assertMethod(t, r, http.MethodPost)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{}`))
	})
	if err := c.TestMonitorAlert(context.Background(), "t1", "m1"); err != nil {
		t.Fatal(err)
	}
}

// ---------------------------------------------------------------------------
// Probes
// ---------------------------------------------------------------------------

func TestListProbes(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertPath(t, r, "/api/v1/teams/t1/monitors/m1/probes")
		_ = json.NewEncoder(w).Encode(ProbeList{Probes: []Probe{{ID: "p1"}}, Total: 1})
	})
	list, err := c.ListProbes(context.Background(), "t1", "m1", ProbeListOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if list.Total != 1 {
		t.Errorf("got %d", list.Total)
	}
}

func TestGetProbe(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertPath(t, r, "/api/v1/teams/t1/monitors/m1/probes/p1")
		_ = json.NewEncoder(w).Encode(Probe{ID: "p1", StatusCode: 200})
	})
	p, err := c.GetProbe(context.Background(), "t1", "m1", "p1")
	if err != nil {
		t.Fatal(err)
	}
	if p.StatusCode != 200 {
		t.Errorf("got %d", p.StatusCode)
	}
}

func TestAggregateProbes(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertPath(t, r, "/api/v1/teams/t1/monitors/m1/probes/aggregate")
		_ = json.NewEncoder(w).Encode(ProbeAggregation{
			Buckets: []ProbeAggBucket{{Timestamp: "2026-01-01T00:00:00Z", TotalCount: 10}},
		})
	})
	agg, err := c.AggregateProbes(context.Background(), "t1", "m1", 0, 0, "1h", "avg")
	if err != nil {
		t.Fatal(err)
	}
	if len(agg.Buckets) != 1 {
		t.Errorf("got %d", len(agg.Buckets))
	}
}

// ---------------------------------------------------------------------------
// Alert Channels
// ---------------------------------------------------------------------------

func TestCreateAlertChannel(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertPath(t, r, "/api/v1/teams/t1/alert-channels")
		assertMethod(t, r, http.MethodPost)
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(AlertChannel{ID: "ch1", Type: "email"})
	})
	ch, err := c.CreateAlertChannel(context.Background(), "t1", &CreateAlertChannelParams{
		Type: "email", Name: "Ops", Config: map[string]string{"email": "ops@test.com"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if ch.Type != "email" {
		t.Errorf("got %q", ch.Type)
	}
}

func TestListAlertChannels(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertPath(t, r, "/api/v1/teams/t1/alert-channels")
		_ = json.NewEncoder(w).Encode([]AlertChannel{{ID: "ch1"}, {ID: "ch2"}})
	})
	channels, err := c.ListAlertChannels(context.Background(), "t1")
	if err != nil {
		t.Fatal(err)
	}
	if len(channels) != 2 {
		t.Errorf("got %d", len(channels))
	}
}

func TestDeleteAlertChannel(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertPath(t, r, "/api/v1/teams/t1/alert-channels/ch1")
		assertMethod(t, r, http.MethodDelete)
		w.WriteHeader(http.StatusNoContent)
	})
	if err := c.DeleteAlertChannel(context.Background(), "t1", "ch1"); err != nil {
		t.Fatal(err)
	}
}

func TestLinkMonitorChannel(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertPath(t, r, "/api/v1/teams/t1/monitors/m1/alert-channels/ch1")
		assertMethod(t, r, http.MethodPost)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{}`))
	})
	if err := c.LinkMonitorChannel(context.Background(), "t1", "m1", "ch1"); err != nil {
		t.Fatal(err)
	}
}

func TestListMonitorChannels(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertPath(t, r, "/api/v1/teams/t1/monitors/m1/alert-channels")
		_ = json.NewEncoder(w).Encode(MonitorChannelsResponse{ChannelIDs: []string{"ch1", "ch2"}})
	})
	ids, err := c.ListMonitorChannels(context.Background(), "t1", "m1")
	if err != nil {
		t.Fatal(err)
	}
	if len(ids) != 2 {
		t.Errorf("got %d", len(ids))
	}
}

func TestListAlertEvents(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertPath(t, r, "/api/v1/teams/t1/monitors/m1/alert-events")
		_ = json.NewEncoder(w).Encode([]AlertEvent{{ID: "ae1", Type: "down"}})
	})
	events, err := c.ListAlertEvents(context.Background(), "t1", "m1")
	if err != nil {
		t.Fatal(err)
	}
	if len(events) != 1 {
		t.Errorf("got %d", len(events))
	}
}

// ---------------------------------------------------------------------------
// Comments
// ---------------------------------------------------------------------------

func TestCreateComment(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertPath(t, r, "/api/v1/teams/t1/monitors/m1/comments")
		assertMethod(t, r, http.MethodPost)
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(Comment{ID: "c1", Body: "test"})
	})
	comment, err := c.CreateComment(context.Background(), "t1", "m1", CreateCommentParams{Body: "test"})
	if err != nil {
		t.Fatal(err)
	}
	if comment.ID != "c1" {
		t.Errorf("got %q", comment.ID)
	}
}

func TestListComments(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertPath(t, r, "/api/v1/teams/t1/monitors/m1/comments")
		_ = json.NewEncoder(w).Encode([]Comment{{ID: "c1"}, {ID: "c2"}})
	})
	comments, err := c.ListComments(context.Background(), "t1", "m1", CommentListOptions{
		From: "2026-01-01T00:00:00Z", To: "2026-01-02T00:00:00Z",
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(comments) != 2 {
		t.Errorf("got %d", len(comments))
	}
}

func TestEditComment(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertPath(t, r, "/api/v1/teams/t1/monitors/m1/comments/c1")
		assertMethod(t, r, http.MethodPut)
		_ = json.NewEncoder(w).Encode(Comment{ID: "c1", Body: "edited"})
	})
	comment, err := c.EditComment(context.Background(), "t1", "m1", "c1", "edited")
	if err != nil {
		t.Fatal(err)
	}
	if comment.Body != "edited" {
		t.Errorf("got %q", comment.Body)
	}
}

func TestDeleteComment(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertPath(t, r, "/api/v1/teams/t1/monitors/m1/comments/c1")
		assertMethod(t, r, http.MethodDelete)
		w.WriteHeader(http.StatusNoContent)
	})
	if err := c.DeleteComment(context.Background(), "t1", "m1", "c1"); err != nil {
		t.Fatal(err)
	}
}

func TestReplyToComment(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertPath(t, r, "/api/v1/teams/t1/monitors/m1/comments/c1/replies")
		assertMethod(t, r, http.MethodPost)
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(Comment{ID: "c2", Body: "reply"})
	})
	reply, err := c.ReplyToComment(context.Background(), "t1", "m1", "c1", "reply")
	if err != nil {
		t.Fatal(err)
	}
	if reply.ID != "c2" {
		t.Errorf("got %q", reply.ID)
	}
}

func TestResolveComment(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertPath(t, r, "/api/v1/teams/t1/monitors/m1/comments/c1/resolve")
		assertMethod(t, r, http.MethodPost)
		w.WriteHeader(http.StatusNoContent)
	})
	if err := c.ResolveComment(context.Background(), "t1", "m1", "c1"); err != nil {
		t.Fatal(err)
	}
}

func TestReopenComment(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertPath(t, r, "/api/v1/teams/t1/monitors/m1/comments/c1/reopen")
		assertMethod(t, r, http.MethodPost)
		w.WriteHeader(http.StatusNoContent)
	})
	if err := c.ReopenComment(context.Background(), "t1", "m1", "c1"); err != nil {
		t.Fatal(err)
	}
}

func TestListCommentsByTeam(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertPath(t, r, "/api/v1/teams/t1/comments")
		_ = json.NewEncoder(w).Encode([]Comment{{ID: "c1"}})
	})
	comments, err := c.ListCommentsByTeam(context.Background(), "t1")
	if err != nil {
		t.Fatal(err)
	}
	if len(comments) != 1 {
		t.Errorf("got %d", len(comments))
	}
}

// ---------------------------------------------------------------------------
// External Links
// ---------------------------------------------------------------------------

func TestCreateExternalLink(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertPath(t, r, "/api/v1/teams/t1/external-links")
		assertMethod(t, r, http.MethodPost)
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(ExternalLink{ID: "el1", Name: "Grafana"})
	})
	link, err := c.CreateExternalLink(context.Background(), "t1", &CreateExternalLinkParams{
		Name: "Grafana", URLTemplate: "https://grafana.io/d/{{.MonitorID}}", TimeWindowMinutes: 60,
	})
	if err != nil {
		t.Fatal(err)
	}
	if link.Name != "Grafana" {
		t.Errorf("got %q", link.Name)
	}
}

func TestListExternalLinks(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertPath(t, r, "/api/v1/teams/t1/external-links")
		_ = json.NewEncoder(w).Encode([]ExternalLink{{ID: "el1"}})
	})
	links, err := c.ListExternalLinks(context.Background(), "t1")
	if err != nil {
		t.Fatal(err)
	}
	if len(links) != 1 {
		t.Errorf("got %d", len(links))
	}
}

func TestDeleteExternalLink(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertPath(t, r, "/api/v1/teams/t1/external-links/el1")
		assertMethod(t, r, http.MethodDelete)
		w.WriteHeader(http.StatusNoContent)
	})
	if err := c.DeleteExternalLink(context.Background(), "t1", "el1"); err != nil {
		t.Fatal(err)
	}
}

func TestAssignExternalLink(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertPath(t, r, "/api/v1/teams/t1/monitors/m1/external-links/el1")
		assertMethod(t, r, http.MethodPost)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{}`))
	})
	if err := c.AssignExternalLink(context.Background(), "t1", "m1", "el1"); err != nil {
		t.Fatal(err)
	}
}

func TestListMonitorExternalLinks(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertPath(t, r, "/api/v1/teams/t1/monitors/m1/external-links")
		_ = json.NewEncoder(w).Encode([]ExternalLink{{ID: "el1"}})
	})
	links, err := c.ListMonitorExternalLinks(context.Background(), "t1", "m1")
	if err != nil {
		t.Fatal(err)
	}
	if len(links) != 1 {
		t.Errorf("got %d", len(links))
	}
}

// ---------------------------------------------------------------------------
// Shares
// ---------------------------------------------------------------------------

func TestCreateShare(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertPath(t, r, "/api/v1/teams/t1/monitors/m1/shares")
		assertMethod(t, r, http.MethodPost)
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(Share{ID: "s1", Token: "abc"})
	})
	s, err := c.CreateShare(context.Background(), "t1", "m1", CreateShareParams{})
	if err != nil {
		t.Fatal(err)
	}
	if s.Token != "abc" {
		t.Errorf("got %q", s.Token)
	}
}

func TestListShares(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertPath(t, r, "/api/v1/teams/t1/monitors/m1/shares")
		_ = json.NewEncoder(w).Encode([]Share{{ID: "s1"}})
	})
	shares, err := c.ListShares(context.Background(), "t1", "m1")
	if err != nil {
		t.Fatal(err)
	}
	if len(shares) != 1 {
		t.Errorf("got %d", len(shares))
	}
}

func TestRevokeShare(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertPath(t, r, "/api/v1/teams/t1/monitors/m1/shares/s1")
		assertMethod(t, r, http.MethodDelete)
		w.WriteHeader(http.StatusNoContent)
	})
	if err := c.RevokeShare(context.Background(), "t1", "m1", "s1"); err != nil {
		t.Fatal(err)
	}
}

// ---------------------------------------------------------------------------
// Traces
// ---------------------------------------------------------------------------

func TestListTraces(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertPath(t, r, "/api/v1/teams/t1/monitors/m1/traces")
		_ = json.NewEncoder(w).Encode([]Trace{{ID: "tr1", Status: "completed"}})
	})
	traces, err := c.ListTraces(context.Background(), "t1", "m1")
	if err != nil {
		t.Fatal(err)
	}
	if len(traces) != 1 {
		t.Errorf("got %d", len(traces))
	}
}

func TestRequestTrace(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertPath(t, r, "/api/v1/teams/t1/monitors/m1/traces")
		assertMethod(t, r, http.MethodPost)
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(Trace{ID: "tr1", Status: "pending"})
	})
	tr, err := c.RequestTrace(context.Background(), "t1", "m1")
	if err != nil {
		t.Fatal(err)
	}
	if tr.Status != "pending" {
		t.Errorf("got %q", tr.Status)
	}
}

// ---------------------------------------------------------------------------
// Geo
// ---------------------------------------------------------------------------

func TestListCheckers(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertPath(t, r, "/api/v1/checkers")
		_ = json.NewEncoder(w).Encode([]Checker{{ID: "ck1", Region: "us-east"}})
	})
	checkers, err := c.ListCheckers(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(checkers) != 1 {
		t.Errorf("got %d", len(checkers))
	}
}

func TestListRegions(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertPath(t, r, "/api/v1/regions")
		_ = json.NewEncoder(w).Encode(GeoRegionsResponse{
			Regions: []GeoRegion{{Code: "NA", Name: "North America"}},
		})
	})
	resp, err := c.ListRegions(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(resp.Regions) != 1 {
		t.Errorf("got %d", len(resp.Regions))
	}
}

// ---------------------------------------------------------------------------
// Triggers (formerly Watchdogs)
// ---------------------------------------------------------------------------

func TestCreateTrigger(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertPath(t, r, "/api/v1/accounts/a1/status-pages/sp1/components/comp1/triggers")
		assertMethod(t, r, http.MethodPost)
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(Trigger{ID: "t1", Severity: "major"})
	})
	tr, err := c.CreateTrigger(context.Background(), "a1", "sp1", "comp1", &CreateTriggerParams{
		MonitorID: "m1", Severity: "major",
	})
	if err != nil {
		t.Fatal(err)
	}
	if tr.Severity != "major" {
		t.Errorf("got %q", tr.Severity)
	}
}

func TestListTriggers(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertPath(t, r, "/api/v1/accounts/a1/status-pages/sp1/components/comp1/triggers")
		_ = json.NewEncoder(w).Encode([]Trigger{{ID: "t1"}})
	})
	triggers, err := c.ListTriggers(context.Background(), "a1", "sp1", "comp1")
	if err != nil {
		t.Fatal(err)
	}
	if len(triggers) != 1 {
		t.Errorf("got %d", len(triggers))
	}
}

func TestDeleteTrigger(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertPath(t, r, "/api/v1/accounts/a1/status-pages/sp1/components/comp1/triggers/t1")
		assertMethod(t, r, http.MethodDelete)
		w.WriteHeader(http.StatusNoContent)
	})
	if err := c.DeleteTrigger(context.Background(), "a1", "sp1", "comp1", "t1"); err != nil {
		t.Fatal(err)
	}
}

// ---------------------------------------------------------------------------
// Metrics / Chart Events
// ---------------------------------------------------------------------------

func TestGetMonitorMetrics(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertPath(t, r, "/api/v1/teams/t1/monitors/m1/metrics")
		_ = json.NewEncoder(w).Encode(MonitorMetrics{})
	})
	_, err := c.GetMonitorMetrics(context.Background(), "t1", "m1")
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetMonitorExpiration(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertPath(t, r, "/api/v1/teams/t1/monitors/m1/expiration")
		_ = json.NewEncoder(w).Encode(Expiration{})
	})
	_, err := c.GetMonitorExpiration(context.Background(), "t1", "m1")
	if err != nil {
		t.Fatal(err)
	}
}

func TestListTimeline(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertPath(t, r, "/api/v1/teams/t1/monitors/m1/timeline")
		_ = json.NewEncoder(w).Encode([]TimelineEvent{{ID: "tl1"}})
	})
	events, err := c.ListTimeline(context.Background(), "t1", "m1", TimelineListOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if len(events) != 1 {
		t.Errorf("got %d", len(events))
	}
}

func TestCreateChartEvent_UsesEventsPath(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertPath(t, r, "/api/v1/teams/t1/events")
		assertMethod(t, r, http.MethodPost)
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(ChartEvent{ID: "ev1", Kind: "deploy"})
	})
	ev, err := c.CreateChartEvent(context.Background(), "t1", CreateChartEventParams{
		Kind: "deploy", Title: "v1.0", StartedAt: "2026-01-01T00:00:00Z",
	})
	if err != nil {
		t.Fatal(err)
	}
	if ev.ID != "ev1" {
		t.Errorf("got %q", ev.ID)
	}
}

func TestListChartEvents_WithFilters(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertPath(t, r, "/api/v1/teams/t1/events")
		if r.URL.Query().Get("kind") != "deploy" {
			t.Errorf("kind: got %q", r.URL.Query().Get("kind"))
		}
		_ = json.NewEncoder(w).Encode([]ChartEvent{{ID: "ev1"}})
	})
	events, err := c.ListChartEvents(context.Background(), "t1", ChartEventListOptions{
		From: "2026-01-01T00:00:00Z", To: "2026-01-02T00:00:00Z", Kind: "deploy",
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(events) != 1 {
		t.Errorf("got %d", len(events))
	}
}

func TestUpdateChartEvent(t *testing.T) {
	title := "updated"
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertPath(t, r, "/api/v1/teams/t1/events/ev1")
		assertMethod(t, r, http.MethodPut)
		_ = json.NewEncoder(w).Encode(ChartEvent{ID: "ev1", Title: "updated"})
	})
	ev, err := c.UpdateChartEvent(context.Background(), "t1", "ev1", UpdateChartEventParams{Title: &title})
	if err != nil {
		t.Fatal(err)
	}
	if ev.Title != "updated" {
		t.Errorf("got %q", ev.Title)
	}
}

func TestDeleteChartEvent(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertPath(t, r, "/api/v1/teams/t1/events/ev1")
		assertMethod(t, r, http.MethodDelete)
		w.WriteHeader(http.StatusNoContent)
	})
	if err := c.DeleteChartEvent(context.Background(), "t1", "ev1"); err != nil {
		t.Fatal(err)
	}
}

func TestIngestChartEvent(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertPath(t, r, "/api/v1/teams/t1/events/ingest")
		assertMethod(t, r, http.MethodPost)
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(ChartEvent{ID: "ev2", Source: "webhook"})
	})
	ev, err := c.IngestChartEvent(context.Background(), "t1", IngestChartEventParams{Kind: "deploy", Title: "v1.0"})
	if err != nil {
		t.Fatal(err)
	}
	if ev.Source != "webhook" {
		t.Errorf("got %q", ev.Source)
	}
}

// ---------------------------------------------------------------------------
// Browser Probes
// ---------------------------------------------------------------------------

func TestListBrowserProbes(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertPath(t, r, "/api/v1/teams/t1/monitors/m1/browser-probes")
		if r.URL.Query().Get("from") != "2026-01-01T00:00:00Z" {
			t.Errorf("from: got %q", r.URL.Query().Get("from"))
		}
		_ = json.NewEncoder(w).Encode(BrowserProbeList{
			Items: []BrowserProbe{{ID: "bp1", Status: 200, TotalMs: 1500, CheckedAt: "2026-01-01T00:01:00Z"}},
		})
	})
	list, err := c.ListBrowserProbes(context.Background(), "t1", "m1", BrowserProbeListOptions{
		From: "2026-01-01T00:00:00Z", To: "2026-01-02T00:00:00Z", Limit: 10,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(list.Items) != 1 {
		t.Errorf("got %d", len(list.Items))
	}
}

func TestGetBrowserProbe(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertPath(t, r, "/api/v1/teams/t1/monitors/m1/browser-probes/bp1")
		_ = json.NewEncoder(w).Encode(BrowserProbe{
			ID: "bp1", Status: 200, LcpMs: 250.5,
			ConsoleMessages: []ConsoleMessage{{Type: "error", Text: "oops"}},
			StepResults:     []StepResult{{Action: "click", Status: "passed", DurationMs: 100}},
			CheckedAt:       "2026-01-01T00:01:00Z",
		})
	})
	bp, err := c.GetBrowserProbe(context.Background(), "t1", "m1", "bp1")
	if err != nil {
		t.Fatal(err)
	}
	if bp.LcpMs != 250.5 {
		t.Errorf("got %f", bp.LcpMs)
	}
	if len(bp.ConsoleMessages) != 1 {
		t.Errorf("got %d", len(bp.ConsoleMessages))
	}
}

func TestAggregateBrowserProbes(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertPath(t, r, "/api/v1/teams/t1/monitors/m1/browser-probes/aggregate")
		_ = json.NewEncoder(w).Encode(BrowserProbeAggregation{
			Buckets: []BrowserProbeAggBucket{{Timestamp: "2026-01-01T00:00:00Z", ProbeCount: 5}},
		})
	})
	agg, err := c.AggregateBrowserProbes(context.Background(), "t1", "m1",
		"2026-01-01T00:00:00Z", "2026-01-02T00:00:00Z", "1h")
	if err != nil {
		t.Fatal(err)
	}
	if len(agg.Buckets) != 1 {
		t.Errorf("got %d", len(agg.Buckets))
	}
}

func TestDownloadBrowserScreenshot(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertPath(t, r, "/api/v1/teams/t1/monitors/m1/browser-probes/bp1/screenshot")
		_, _ = w.Write([]byte("PNG-DATA"))
	})
	data, err := c.DownloadBrowserScreenshot(context.Background(), "t1", "m1", "bp1")
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "PNG-DATA" {
		t.Errorf("got %q", string(data))
	}
}

func TestDownloadBrowserHAR(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertPath(t, r, "/api/v1/teams/t1/monitors/m1/browser-probes/bp1/har")
		_, _ = w.Write([]byte(`{"log":{}}`))
	})
	data, err := c.DownloadBrowserHAR(context.Background(), "t1", "m1", "bp1")
	if err != nil {
		t.Fatal(err)
	}
	if len(data) == 0 {
		t.Error("empty HAR")
	}
}

// ---------------------------------------------------------------------------
// Services
// ---------------------------------------------------------------------------

func TestCreateService(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertPath(t, r, "/api/v1/accounts/a1/services")
		assertMethod(t, r, http.MethodPost)
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(Service{ID: "svc1", Name: "Payment API"})
	})
	s, err := c.CreateService(context.Background(), "a1", &CreateServiceParams{
		Name: "Payment API",
	})
	if err != nil {
		t.Fatal(err)
	}
	if s.ID != "svc1" {
		t.Errorf("got %q", s.ID)
	}
}

func TestListServices(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertPath(t, r, "/api/v1/accounts/a1/services")
		assertMethod(t, r, http.MethodGet)
		_ = json.NewEncoder(w).Encode([]Service{{ID: "svc1"}, {ID: "svc2"}})
	})
	services, err := c.ListServices(context.Background(), "a1")
	if err != nil {
		t.Fatal(err)
	}
	if len(services) != 2 {
		t.Errorf("got %d", len(services))
	}
}

func TestGetService(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertPath(t, r, "/api/v1/accounts/a1/services/svc1")
		assertMethod(t, r, http.MethodGet)
		_ = json.NewEncoder(w).Encode(Service{ID: "svc1", Name: "Payment API"})
	})
	s, err := c.GetService(context.Background(), "a1", "svc1")
	if err != nil {
		t.Fatal(err)
	}
	if s.Name != "Payment API" {
		t.Errorf("got %q", s.Name)
	}
}

func TestUpdateService(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertPath(t, r, "/api/v1/accounts/a1/services/svc1")
		assertMethod(t, r, http.MethodPut)
		_ = json.NewEncoder(w).Encode(Service{ID: "svc1", Name: "Updated"})
	})
	name := "Updated"
	s, err := c.UpdateService(context.Background(), "a1", "svc1", &UpdateServiceParams{
		Name: &name,
	})
	if err != nil {
		t.Fatal(err)
	}
	if s.Name != "Updated" {
		t.Errorf("got %q", s.Name)
	}
}

func TestDeleteService(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertPath(t, r, "/api/v1/accounts/a1/services/svc1")
		assertMethod(t, r, http.MethodDelete)
		w.WriteHeader(http.StatusNoContent)
	})
	err := c.DeleteService(context.Background(), "a1", "svc1")
	if err != nil {
		t.Fatal(err)
	}
}

func TestLinkServiceMonitors(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertPath(t, r, "/api/v1/accounts/a1/services/svc1/monitors")
		assertMethod(t, r, http.MethodPost)
		w.WriteHeader(http.StatusNoContent)
	})
	err := c.LinkServiceMonitors(context.Background(), "a1", "svc1", []string{"m1", "m2"})
	if err != nil {
		t.Fatal(err)
	}
}

func TestUnlinkServiceMonitor(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertPath(t, r, "/api/v1/accounts/a1/services/svc1/monitors/m1")
		assertMethod(t, r, http.MethodDelete)
		w.WriteHeader(http.StatusNoContent)
	})
	err := c.UnlinkServiceMonitor(context.Background(), "a1", "svc1", "m1")
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetServiceAnalytics(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertPath(t, r, "/api/v1/accounts/a1/services/svc1/analytics")
		assertMethod(t, r, http.MethodGet)
		_ = json.NewEncoder(w).Encode(ServiceAnalytics{IncidentCount: 5})
	})
	a, err := c.GetServiceAnalytics(context.Background(), "a1", "svc1")
	if err != nil {
		t.Fatal(err)
	}
	if a.IncidentCount != 5 {
		t.Errorf("got %d", a.IncidentCount)
	}
}

// ---------------------------------------------------------------------------
// Incidents
// ---------------------------------------------------------------------------

func TestCreateIncident(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertPath(t, r, "/api/v1/accounts/a1/incidents")
		assertMethod(t, r, http.MethodPost)
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(AccountIncident{ID: "inc1", Name: "API down"})
	})
	inc, err := c.CreateAccountIncident(context.Background(), "a1", &CreateAccountIncidentParams{
		Name:     "API down",
		Severity: "major",
	})
	if err != nil {
		t.Fatal(err)
	}
	if inc.ID != "inc1" {
		t.Errorf("got %q", inc.ID)
	}
}

func TestListIncidents(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertPath(t, r, "/api/v1/accounts/a1/incidents")
		assertMethod(t, r, http.MethodGet)
		_ = json.NewEncoder(w).Encode([]AccountIncident{{ID: "inc1"}, {ID: "inc2"}})
	})
	incidents, err := c.ListAccountIncidents(context.Background(), "a1", nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(incidents) != 2 {
		t.Errorf("got %d", len(incidents))
	}
}

func TestDeleteIncident(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertPath(t, r, "/api/v1/accounts/a1/incidents/inc1")
		assertMethod(t, r, http.MethodDelete)
		w.WriteHeader(http.StatusNoContent)
	})
	if err := c.DeleteAccountIncident(context.Background(), "a1", "inc1"); err != nil {
		t.Fatal(err)
	}
}

func TestGetIncident(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertPath(t, r, "/api/v1/accounts/a1/incidents/inc1")
		assertMethod(t, r, http.MethodGet)
		_ = json.NewEncoder(w).Encode(AccountIncidentWithDetails{
			AccountIncident: AccountIncident{ID: "inc1", Name: "API down"},
			Updates:         []AccountIncidentUpdate{{ID: "u1", Status: "identified"}},
		})
	})
	inc, err := c.GetAccountIncident(context.Background(), "a1", "inc1")
	if err != nil {
		t.Fatal(err)
	}
	if inc.ID != "inc1" {
		t.Errorf("got %q", inc.ID)
	}
	if len(inc.Updates) != 1 {
		t.Errorf("updates: got %d", len(inc.Updates))
	}
}

func TestAcknowledgeIncident(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertPath(t, r, "/api/v1/accounts/a1/incidents/inc1/acknowledge")
		assertMethod(t, r, http.MethodPost)
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "acknowledged"})
	})
	if err := c.AcknowledgeAccountIncident(context.Background(), "a1", "inc1"); err != nil {
		t.Fatal(err)
	}
}

// ---------------------------------------------------------------------------
// On-Call Schedules
// ---------------------------------------------------------------------------

func TestCreateSchedule(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertPath(t, r, "/api/v1/accounts/a1/oncall/schedules")
		assertMethod(t, r, http.MethodPost)
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(OnCallSchedule{ID: "sch1", Name: "Primary"})
	})
	s, err := c.CreateSchedule(context.Background(), "a1", &CreateScheduleParams{
		Name:         "Primary",
		Timezone:     "UTC",
		RotationType: "weekly",
		Participants: []string{"u1", "u2"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if s.ID != "sch1" {
		t.Errorf("got %q", s.ID)
	}
}

func TestListSchedules(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertPath(t, r, "/api/v1/accounts/a1/oncall/schedules")
		assertMethod(t, r, http.MethodGet)
		_ = json.NewEncoder(w).Encode([]OnCallSchedule{{ID: "sch1"}, {ID: "sch2"}})
	})
	schedules, err := c.ListSchedules(context.Background(), "a1")
	if err != nil {
		t.Fatal(err)
	}
	if len(schedules) != 2 {
		t.Errorf("got %d", len(schedules))
	}
}

func TestGetSchedule(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertPath(t, r, "/api/v1/accounts/a1/oncall/schedules/sch1")
		assertMethod(t, r, http.MethodGet)
		_ = json.NewEncoder(w).Encode(OnCallSchedule{ID: "sch1", Name: "Primary"})
	})
	s, err := c.GetSchedule(context.Background(), "a1", "sch1")
	if err != nil {
		t.Fatal(err)
	}
	if s.Name != "Primary" {
		t.Errorf("got %q", s.Name)
	}
}

func TestDeleteSchedule(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertPath(t, r, "/api/v1/accounts/a1/oncall/schedules/sch1")
		assertMethod(t, r, http.MethodDelete)
		w.WriteHeader(http.StatusNoContent)
	})
	err := c.DeleteSchedule(context.Background(), "a1", "sch1")
	if err != nil {
		t.Fatal(err)
	}
}

// ---------------------------------------------------------------------------
// On-Call Overrides
// ---------------------------------------------------------------------------

func TestCreateOverride(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertPath(t, r, "/api/v1/accounts/a1/oncall/schedules/sch1/overrides")
		assertMethod(t, r, http.MethodPost)
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(OnCallOverride{ID: "ovr1", ScheduleID: "sch1"})
	})
	o, err := c.CreateOverride(context.Background(), "a1", "sch1", &CreateOverrideParams{
		OriginalUserID:    "u1",
		ReplacementUserID: "u2",
		StartAt:           "2026-04-01T00:00:00Z",
		EndAt:             "2026-04-02T00:00:00Z",
	})
	if err != nil {
		t.Fatal(err)
	}
	if o.ID != "ovr1" {
		t.Errorf("got %q", o.ID)
	}
}

func TestListOverrides(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertPath(t, r, "/api/v1/accounts/a1/oncall/schedules/sch1/overrides")
		assertMethod(t, r, http.MethodGet)
		_ = json.NewEncoder(w).Encode([]OnCallOverride{{ID: "ovr1"}, {ID: "ovr2"}})
	})
	overrides, err := c.ListOverrides(context.Background(), "a1", "sch1")
	if err != nil {
		t.Fatal(err)
	}
	if len(overrides) != 2 {
		t.Errorf("got %d", len(overrides))
	}
}

func TestDeleteOverride(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertPath(t, r, "/api/v1/accounts/a1/oncall/schedules/sch1/overrides/ovr1")
		assertMethod(t, r, http.MethodDelete)
		w.WriteHeader(http.StatusNoContent)
	})
	err := c.DeleteOverride(context.Background(), "a1", "sch1", "ovr1")
	if err != nil {
		t.Fatal(err)
	}
}

// ---------------------------------------------------------------------------
// Who's On Call
// ---------------------------------------------------------------------------

func TestGetWhosOnCall(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertPath(t, r, "/api/v1/accounts/a1/oncall/now")
		assertMethod(t, r, http.MethodGet)
		_ = json.NewEncoder(w).Encode([]WhosOnCall{
			{ScheduleID: "sch1", UserID: "u1"},
		})
	})
	info, err := c.GetWhosOnCall(context.Background(), "a1")
	if err != nil {
		t.Fatal(err)
	}
	if len(info) != 1 {
		t.Errorf("got %d", len(info))
	}
	if info[0].UserID != "u1" {
		t.Errorf("got %q", info[0].UserID)
	}
}

// ---------------------------------------------------------------------------
// Escalation Policies
// ---------------------------------------------------------------------------

func TestCreateEscalationPolicy(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertPath(t, r, "/api/v1/accounts/a1/oncall/escalation-policies")
		assertMethod(t, r, http.MethodPost)
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(EscalationPolicy{
			ID: "ep1", Name: "Default",
			Levels: []EscalationLevel{{ScheduleID: "sch1", AckTimeoutMinutes: 5}},
		})
	})
	p, err := c.CreateEscalationPolicy(context.Background(), "a1", &CreateEscalationPolicyParams{
		Name:   "Default",
		Levels: []EscalationLevel{{ScheduleID: "sch1", AckTimeoutMinutes: 5}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if p.ID != "ep1" {
		t.Errorf("got %q", p.ID)
	}
}

func TestListEscalationPolicies(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertPath(t, r, "/api/v1/accounts/a1/oncall/escalation-policies")
		assertMethod(t, r, http.MethodGet)
		_ = json.NewEncoder(w).Encode([]EscalationPolicy{{ID: "ep1"}, {ID: "ep2"}})
	})
	policies, err := c.ListEscalationPolicies(context.Background(), "a1")
	if err != nil {
		t.Fatal(err)
	}
	if len(policies) != 2 {
		t.Errorf("got %d", len(policies))
	}
}

func TestDeleteEscalationPolicy(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertPath(t, r, "/api/v1/accounts/a1/oncall/escalation-policies/ep1")
		assertMethod(t, r, http.MethodDelete)
		w.WriteHeader(http.StatusNoContent)
	})
	err := c.DeleteEscalationPolicy(context.Background(), "a1", "ep1")
	if err != nil {
		t.Fatal(err)
	}
}
