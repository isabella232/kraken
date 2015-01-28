package jira

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestUpdateReleaseDate(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Fatalf("wanted PUT but found %s\n", r.Method)
		}
		url := *r.URL
		params := url.Query()
		if params.Get("releaseDate") != "16/Feb/14" {
			t.Fatalf("Want 16/Feb/14 but got %s\n", params.Get("releaseDate"))
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
	err := client.UpdateReleaseDate(4, "16/Feb/14")
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
	err := client.UpdateReleaseDate(4, "16/Feb/14")
	if err == nil {
		t.Fatalf("Expecting an error\n")
	}
}
