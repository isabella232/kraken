package jira

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestDeleteMapping(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Fatalf("wanted DELETE but found %s\n", r.Method)
		}
		url := *r.URL
		if url.Path != "/rest/com.deniz.jira.mapping/latest/4" {
			t.Fatalf("Want /rest/com.deniz.jira.mapping/latest/4 but got %s\n", url.Path)
		}
		if r.Header.Get("Accept") != "application/json" {
			t.Fatalf("Want application/json but got %s\n", r.Header.Get("Accept"))
		}
		if r.Header.Get("Authorization") != "Basic dTpw" {
			t.Fatalf("Want Basic dTpw but got %s\n", r.Header.Get("Authorization"))
		}
		w.WriteHeader(204)
	}))
	defer testServer.Close()
	url, _ := url.Parse(testServer.URL)
	client := NewClient("u", "p", url)
	err := client.DeleteMapping(4)
	if err != nil {
		t.Fatalf("Not expecting an error %v\n", err)
	}
}

func TestDeleteMappingNon204(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
	}))
	defer testServer.Close()

	url, _ := url.Parse(testServer.URL)
	client := NewClient("u", "p", url)
	err := client.DeleteMapping(4)
	if err == nil {
		t.Fatalf("Expecting an error\n")
	}
}
