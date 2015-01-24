package jira

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

func TestCreateVersion(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Fatalf("wanted POST but found %s\n", r.Method)
		}
		url := *r.URL
		if url.Path != "/rest/api/2/version" {
			t.Fatalf("Want /rest/api/2/version but got %s\n", url.Path)
		}
		if r.Header.Get("Content-type") != "application/json" {
			t.Fatalf("Want application/json but got %s\n", r.Header.Get("Content-type"))
		}
		if r.Header.Get("Authorization") != "Basic dTpw" {
			t.Fatalf("Want Basic dTpw but got %s\n", r.Header.Get("Authorization"))
		}

		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("error reading POST body: %v\n", err)
		}

		var v Version
		if err := json.Unmarshal(data, &v); err != nil {
			t.Fatalf("Unexpected error: %v\n", err)
		}
		if v.Name != "1.0" {
			t.Fatalf("Want 1.0 but got %s\n", v.Name)
		}
		if v.Description != "Version 1.0" {
			t.Fatalf("Want Version 1.0 but got %s\n", v.Description)
		}
		if v.ProjectID != 1 {
			t.Fatalf("Want 1 but got %s\n", v.ProjectID)
		}
		if v.Archived {
			t.Fatalf("Want false\n")
		}
		if !v.Released {
			t.Fatalf("Want true\n")
		}
		if v.ReleaseDate != time.Now().Format("2006-01-02") {
			t.Fatalf("Want "+time.Now().Format("2006-01-02")+" but got %s\n", v.ReleaseDate)
		}
		w.WriteHeader(201)
	}))
	defer testServer.Close()
	url, _ := url.Parse(testServer.URL)
	client := NewClient("u", "p", url)
	err := client.CreateVersion(1, "1.0")
	if err != nil {
		t.Fatalf("Unexpected error:  %v\n", err)
	}
}

func TestCreateVersionNon201(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(403)
	}))
	defer testServer.Close()
	url, _ := url.Parse(testServer.URL)
	client := NewClient("u", "p", url)
	err := client.CreateVersion(1, "1.0")
	if err == nil {
		t.Fatalf("Expecting an error\n")
	}
}
