package jira

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestGetComponents(t *testing.T) {
	response := `
[
   {
      "isAssigneeTypeValid" : false,
      "lead" : {
         "avatarUrls" : {
            "48x48" : "https://example.com/secure/useravatar?ownerId=bsimpson&avatarId=11900",
            "32x32" : "https://example.com/secure/useravatar?size=medium&ownerId=bsimpson&avatarId=11900",
            "24x24" : "https://example.com/secure/useravatar?size=small&ownerId=bsimpson&avatarId=11900",
            "16x16" : "https://example.com/secure/useravatar?size=xsmall&ownerId=bsimpson&avatarId=11900"
         },     
         "active" : true,
         "name" : "bsimpson",
         "self" : "https://example.com/rest/api/2/user?username=bsimpson",
         "displayName" : "Bart Simpson",
         "key" : "bsimpson"
      },        
      "realAssigneeType" : "PROJECT_DEFAULT",
      "name" : "lolcats",
      "self" : "https://example.com/rest/api/2/component/12105",
      "description" : "lolcats rule",
      "assigneeType" : "PROJECT_DEFAULT",
      "id" : "12105"
   },
   {    
      "isAssigneeTypeValid" : false,
      "lead" : {
         "avatarUrls" : {
            "48x48" : "https://example.com/secure/useravatar?ownerId=hsimpson&avatarId=11516",
            "32x32" : "https://example.com/secure/useravatar?size=medium&ownerId=hsimpson&avatarId=11516",
            "24x24" : "https://example.com/secure/useravatar?size=small&ownerId=hsimpson&avatarId=11516",
            "16x16" : "https://example.com/secure/useravatar?size=xsmall&ownerId=hsimpson&avatarId=11516"
         },
         "active" : true,
         "name" : "hsimpson",
         "self" : "https://example.com/rest/api/2/user?username=hsimpson",
         "displayName" : "Homer Simpson",
         "key" : "hsimpson"
      },        
      "realAssigneeType" : "PROJECT_DEFAULT",
      "name" : "no-nukes",
      "self" : "https://example.com/rest/api/2/component/13400",
      "description" : "Demonstration of an nuclear power plant server",
      "assigneeType" : "PROJECT_DEFAULT",
      "id" : "13400"
   }    
]
`
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Fatalf("wanted GET but found %s\n", r.Method)
		}
		url := *r.URL
		if url.Path != "/rest/api/2/project/1/components" {
			t.Fatalf("Want /rest/api/2/project/PRJ/components but got %s\n", url.Path)
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
	r, err := client.GetComponents(1)
	if err != nil {
		t.Fatalf("Not expecting an error %v\n", err)
	}
	if len(r) != 2 {
		t.Fatalf("Want 2 but got %d\n", len(r))
	}
	if r["lolcats"].Name != "lolcats" {
		t.Fatalf("Want lolcats but got %s\n", r["12105"].Name)
	}
	if r["lolcats"].ID != "12105" {
		t.Fatalf("Want 12105 but got %s\n", r["12105"].Name)
	}
	if r["no-nukes"].Name != "no-nukes" {
		t.Fatalf("Want no-nukes but got %s\n", r["13400"].Name)
	}
	if r["no-nukes"].ID != "13400" {
		t.Fatalf("Want 13400 but got %s\n", r["13400"].ID)
	}
}

func TestGetComponents404(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
	}))
	defer testServer.Close()

	url, _ := url.Parse(testServer.URL)
	client := NewClient("u", "p", url)
	_, err := client.GetComponents(1)
	if err == nil {
		t.Fatalf("Expecting an error\n")
	}
}
