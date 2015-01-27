package jira

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

/*
  req, err := http.NewRequest("PUT", fmt.Sprintf("%s/rest/com.deniz.jira.mapping/latest/releaseDate/%d?releaseDate=%s", client.baseURL, mappingID, releaseDate), nil)
        if err != nil {
                return err
        }
        if debug {
                log.Printf("jira.GetVersionsForComponent URL %s\n", req.URL)
        }
        req.Header.Set("Accept", "application/json")
        req.SetBasicAuth(client.username, client.password)
        responseCode, _, err := client.consumeResponse(req)
        if err != nil {
                return err
        }
        if responseCode != http.StatusOK {
                return fmt.Errorf("error updating mapping release date.  Status code: %d.\n", responseCode)
        }

*/

func TestUpdateReleaseDate(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Fatalf("wanted PUT but found %s\n", r.Method)
		}
		url := *r.URL
		params := url.Query()
		if params.Get("releaseDate") != "2006-12-06" {
			t.Fatalf("Want 2006-12-06 but got %s\n", params.Get("releaseDate"))
		}
		if url.Path != "/rest/com.deniz.jira.mapping/latest/releaseDate/4" {
			t.Fatalf("Want /rest/com.deniz.jira.mapping/latest/releaseDate/4 but got %s\n", url.Path)
		}
		if r.Header.Get("Accept") != "application/json" {
			t.Fatalf("Want application/json but got %s\n", r.Header.Get("Accept"))
		}
		if r.Header.Get("Authorization") != "Basic dTpw" {
			t.Fatalf("Want Basic dTpw but got %s\n", r.Header.Get("Authorization"))
		}
		w.WriteHeader(200)
	}))
	defer testServer.Close()

	url, _ := url.Parse(testServer.URL)
	client := NewClient("u", "p", url)
	err := client.UpdateReleaseDate(4, "2006-12-06")
	if err != nil {
		t.Fatalf("Not expecting an error %v\n", err)
	}
}

func TestUpdateReleaseDateNon200(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
	}))
	defer testServer.Close()

	url, _ := url.Parse(testServer.URL)
	client := NewClient("u", "p", url)
	err := client.UpdateReleaseDate(4, "2006-12-06")
	if err == nil {
		t.Fatalf("Expecting an error\n")
	}
}
