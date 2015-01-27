package jira

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestGetVersionsForComponent(t *testing.T) {
	response := `
[
   {
      "name" : "Unknown",
      "isReleased" : false,
      "id" : -1,
      "description" : "Unknown"
   },
   {
      "name" : "1.1",
      "isReleased" : true,
      "id" : 12227
   }
]
`

	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Fatalf("wanted GET but found %s\n", r.Method)
		}
		url := *r.URL
		if url.Path != "/rest/com.deniz.jira.mapping/latest/applicable_versions" {
			t.Fatalf("Want /rest/com.deniz.jira.mapping/latest/applicable_versions but got %s\n", url.Path)
		}
		queryParams := url.Query()
		if queryParams.Get("projectId") != "11300" {
			t.Fatalf("Want 11300 but got %s\n", queryParams.Get("projectId"))
		}
		if queryParams.Get("selectedComponentIds") != "200" {
			t.Fatalf("Want 200 but got %s\n", queryParams.Get("selectedComponentIds"))
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
	r, err := client.GetVersionsForComponent("11300", "200")
	if err != nil {
		t.Fatalf("Not expecting an error %v\n", err)
	}
	if len(r) != 2 {
		t.Fatalf("Want 2 but got %d\n", len(r))
	}

	version, present := r[12227]
	if !present {
		t.Fatalf("Want true\n")
	}
	if version.ID != 12227 {
		t.Fatalf("Want 12227 but got %d\n", version.ID)
	}
	if version.Name != "1.1" {
		t.Fatalf("Want 1.1 but got %s\n", version.Name)
	}
	if !version.Released {
		t.Fatalf("Want true\n")
	}
}

func TestGetVersionsForComponentNon200(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
	}))
	defer testServer.Close()
	url, _ := url.Parse(testServer.URL)
	client := NewClient("u", "p", url)
	_, err := client.GetVersionsForComponent("11300", "200")
	if err == nil {
		t.Fatalf("Expecting an error")
	}
}
