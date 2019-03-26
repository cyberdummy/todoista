package todoist_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/cyberdummy/todoista/todoist"
)

func TestFullSyncRead(t *testing.T) {
	token := "test_token"

	// Start a local HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Test request parameters
		if !strings.HasPrefix(req.URL.String(), "/api/v7/sync?") {
			t.Errorf("URL = %s; want = /api/v7/sync", req.URL.String())
		}
		// Send response to be tested
		rw.Write([]byte(`{
			"items": [
				{
					"id": 3,
					"content": "Content",
					"project_id": 1,
					"date_string": "Mar 23"
				},
				{
					"id": 4,
					"content": "Content",
					"project_id": 2,
					"date_string": "Mar 24"
				}
			],
			"projects": [
				{
					"id": 3,
					"name": "Project 3"
				}
			]
		}`))
	}))
	// Close the server when test finishes
	defer server.Close()

	// Use Client & URL from our local test server
	api, err := todoist.New(token)
	api.SetDomain(server.URL)
	api, err = api.ReadSync()

	if err != nil {
		t.Error("Got error " + err.Error())
	}

	if len(api.Items) != 2 {
		t.Errorf("len(Items) = %d; want = 2", len(api.Items))
	}

	if api.Items[0].ID != 3 {
		t.Errorf("Items[0].ID = %d; want = 3", api.Items[0].ID)
	}

	if api.Items[1].ID != 4 {
		t.Errorf("Items[1].ID = %d; want = 4", api.Items[1].ID)
	}

	if len(api.Projects) != 1 {
		t.Errorf("len(Projects) = %d; want = 1", len(api.Projects))
	}

	if api.Projects[0].ID != 3 {
		t.Errorf("Projects[0].ID = %d; want = 3", api.Projects[0].ID)
	}
}
