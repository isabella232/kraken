package jira

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestCreateMapping(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Fatalf("wanted POST but found %s\n", r.Method)
		}
		url := *r.URL
		if url.Path != "/jira/rest/com.deniz.jira.mapping/latest/" {
			t.Fatalf("Want /jira/rest/com.deniz.jira.mapping/latest/ but got %s\n", url.Path)
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

		/*
		   Mapping struct {
		                ProjectID   int  `json:"projectId"`
		                ComponentID int  `json:"componentId"`
		                VersionID   int  `json:"versionId"`
		                Released    bool `json:"released"`
		        }
		*/
		var v Mapping
		if err := json.Unmarshal(data, &v); err != nil {
			t.Fatalf("Unexpected error: %v\n", err)
		}
		if v.ProjectID != 1 {
			t.Fatalf("Want 1 but got %d\n", v.ProjectID)
		}
		if v.ComponentID != 2 {
			t.Fatalf("Want 2 but got %d\n", v.ComponentID)
		}
		if v.VersionID != 3 {
			t.Fatalf("Want 3 but got %d\n", v.VersionID)
		}
		if v.Released {
			t.Fatalf("Want false\n")
		}

		data, err = json.Marshal(&v)
		w.WriteHeader(201)
		fmt.Fprintf(w, "%s", string(data))
	}))
	defer testServer.Close()

	url, _ := url.Parse(testServer.URL)
	client := NewClient("u", "p", url)
	_, err := client.CreateMapping("1", "2", "3")
	if err != nil {
		t.Fatalf("Unexpected error:  %v\n", err)
	}
}

func TestCreateMappingNon201(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(401)
	}))
	defer testServer.Close()

	url, _ := url.Parse(testServer.URL)
	client := NewClient("u", "p", url)
	_, err := client.CreateMapping("1", "2", "3")
	if err == nil {
		t.Fatalf("Expected error\n")
	}
}
