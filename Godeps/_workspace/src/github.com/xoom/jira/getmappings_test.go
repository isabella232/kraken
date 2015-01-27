package jira

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestGetMappings(t *testing.T) {
	response := `
[
   {
      "projectName" : "The Project",
      "versionName" : "2.4",
      "versionId" : 13008,
      "componentId" : 12326,
      "released" : false,
      "releaseDateStr" : "",
      "projectId" : 11300,
      "id" : 137,
      "componentName" : "service-one",
      "projectKey" : "PRJ"
   },
   {
      "projectName" : "The Project",
      "versionName" : "2.3",
      "versionId" : 13007,
      "componentId" : 12326,
      "released" : false,
      "releaseDateStr" : "",
      "projectId" : 11300,
      "id" : 136,
      "componentName" : "service-one",
      "projectKey" : "PRJ"
   }
   ]
`

	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Fatalf("wanted GET but found %s\n", r.Method)
		}
		url := *r.URL
		if url.Path != "/rest/com.deniz.jira.mapping/latest/mappings" {
			t.Fatalf("Want /rest/com.deniz.jira.mapping/latest/mappings but got %s\n", url.Path)
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
	r, err := client.GetMappings()
	if err != nil {
		t.Fatalf("Not expecting an error %v\n", err)
	}
	if len(r) != 2 {
		t.Fatalf("Want 2 but got %d\n", len(r))
	}

	mapping, present := r[137]
	if !present {
		t.Fatalf("Expecting an entry at key 137\n")
	}
	if mapping.ID != 137 {
		t.Fatalf("Want 137 but got %d\n", mapping.ID)
	}
	if mapping.ProjectName != "The Project" {
		t.Fatalf("Want The Project but got %s\n", mapping.ProjectName)
	}
	if mapping.VersionName != "2.4" {
		t.Fatalf("Want 2.4 but got %s\n", mapping.VersionName)
	}
	if mapping.VersionID != 13008 {
		t.Fatalf("Want 13008 but got %d\n", mapping.VersionID)
	}
	if mapping.ComponentID != 12326 {
		t.Fatalf("Want 12326 but got %d\n", mapping.ComponentID)
	}
	if mapping.Released {
		t.Fatalf("Want false\n")
	}
	if mapping.ReleaseDateStr != "" {
		t.Fatalf("Want empty string but got %s\n", mapping.ReleaseDateStr)
	}
	if mapping.ProjectID != 11300 {
		t.Fatalf("Want 11300 but got %d\n", mapping.ProjectID)
	}
	if mapping.ComponentName != "service-one" {
		t.Fatalf("Want service-one but got %s\n", mapping.ComponentName)
	}
	if mapping.ProjectKey != "PRJ" {
		t.Fatalf("Want PRJ", mapping.ProjectKey)
	}
}

func TestGetMappingNon200(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(400)
	}))
	defer testServer.Close()

	url, _ := url.Parse(testServer.URL)
	client := NewClient("u", "p", url)
	_, err := client.GetMappings()
	if err == nil {
		t.Fatalf("Expecting an error\n")
	}
}
