package jira

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

/*

   Version struct {
           ID          string `json:"id"`
           Name        string `json:"name"`
           Description string `json:"description"`
           Project     string `json:"project"`
           ProjectID   int `json:"projectId"`
           Archived    bool   `json:"archived"`
           Released    bool   `json:"released"`
           ReleaseDate string `json:"releaseDate"`
   }
*/
func TestGetVersions(t *testing.T) {
	response := `
	[ 
	{ 
		"archived" : false,
    		"id" : "12230",
    		"name" : "1.0",
    		"projectId" : 11300,
    		"released" : false,
    		"self" : "https://example.com/rest/api/2/version/12230"
  	},
  	{ 
		"archived" : false,
    		"id" : "12227",
    		"name" : "1.1",
    		"projectId" : 11300,
    		"released" : false,
    		"self" : "https://example.com/rest/api/2/version/12227"
  	}
	]	
	`

	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Fatalf("wanted GET but found %s\n", r.Method)
		}
		url := *r.URL
		if url.Path != "/rest/api/2/project/1/versions" {
			t.Fatalf("Want /rest/api/2/project/PRJ/versions but got %s\n", url.Path)
		}
		if r.Header.Get("Accept") != "application/json" {
			t.Fatalf("Want application/json but got %s\n", r.Header.Get("Accept"))
		}
		if r.Header.Get("Authorization") != "Basic dTpw" {
			t.Fatalf("Want Basic dTpw but got %s\n", r.Header.Get("Authorization"))
		}
		fmt.Fprintln(w, response)
	}))
	defer testServer.Close()

	url, _ := url.Parse(testServer.URL)
	client := NewClient("u", "p", url)
	r, err := client.GetVersions("1")
	if err != nil {
		t.Fatalf("Not expecting an error %v\n", err)
	}
	if len(r) != 2 {
		t.Fatalf("Want 2 but got %d\n", len(r))
	}
	if r["1.0"].Name != "1.0" {
		t.Fatalf("Want lolcats but got %s\n", r["1.0"].Name)
	}
	if r["1.0"].ID != "12230" {
		t.Fatalf("Want 12230 but got %s\n", r["1.0"].Name)
	}
	if r["1.1"].Name != "1.1" {
		t.Fatalf("Want no-nukes but got %s\n", r["1.1"].Name)
	}
	if r["1.1"].ID != "12227" {
		t.Fatalf("Want 12227 but got %s\n", r["1.1"].ID)
	}
}

func TestGetVersions404(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
	}))
	defer testServer.Close()
	url, _ := url.Parse(testServer.URL)
	client := NewClient("u", "p", url)
	_, err := client.GetVersions("1")
	if err == nil {
		t.Fatalf("Expecting an error\n")
	}
}
