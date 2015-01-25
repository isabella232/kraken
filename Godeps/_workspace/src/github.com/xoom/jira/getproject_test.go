package jira

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestGetProject(t *testing.T) {
	response := `
{
   "roles" : {
      "Developers" : "https://example.com/rest/api/2/project/11300/role/10001",
      "Users" : "https://example.com/rest/api/2/project/11300/role/10000",
      "Watcher" : "https://example.com/rest/api/2/project/11300/role/10502",
      "Service Desk Team" : "https://example.com/rest/api/2/project/11300/role/10501",
      "Service Desk Customers" : "https://example.com/rest/api/2/project/11300/role/10500",
      "Administrators" : "https://example.com/rest/api/2/project/11300/role/10002",
      "Tempo Project Managers" : "https://example.com/rest/api/2/project/11300/role/10600"
   },
   "lead" : {
      "avatarUrls" : {
         "48x48" : "https://example.com/secure/useravatar?ownerId=tbombadil&avatarId=11900",
         "16x16" : "https://example.com/secure/useravatar?size=xsmall&ownerId=tbombadil&avatarId=11900"
      },
      "active" : true,
      "name" : "tbombadil",
      "self" : "https://example.com/rest/api/2/user?username=tbombadil",
      "displayName" : "Tom Bombadil",
      "key" : "tbombadil"
   },
   "avatarUrls" : {
      "24x24" : "https://example.com/secure/projectavatar?size=small&pid=11300&avatarId=10011",
      "16x16" : "https://example.com/secure/projectavatar?size=xsmall&pid=11300&avatarId=10011"
   },
   "name" : "CoolProj",
   "components" : [
      {
         "isAssigneeTypeValid" : false,
         "name" : "amqp",
         "self" : "https://example.com/rest/api/2/component/12105",
         "id" : "12105",
         "description" : "AMQP utilities"
      },
      {
         "isAssigneeTypeValid" : false,
         "name" : "redis",
         "self" : "https://example.com/rest/api/2/component/13400",
         "id" : "13400",
         "description" : "Demonstration of Redis"
      }
   ],
   "self" : "https://example.com/rest/api/2/project/11300",
   "description" : "",
   "key" : "COOL",
   "versions" : [
      {
         "projectId" : 11300,
         "name" : "1.0",
         "self" : "https://example.com/rest/api/2/version/12230",
         "released" : false,
         "id" : "12230",
         "archived" : false
      },
      {
         "projectId" : 11300,
         "name" : "1.1",
         "self" : "https://example.com/rest/api/2/version/12227",
         "released" : false,
         "id" : "12227",
         "archived" : false
      }
   ],
   "assigneeType" : "UNASSIGNED",
   "expand" : "projectKeys",
   "issueTypes" : [
      {
         "name" : "Epic",
         "iconUrl" : "https://example.com/images/icons/issuetypes/epic.png",
         "self" : "https://example.com/rest/api/2/issuetype/6",
         "id" : "6",
         "subtask" : false,
         "description" : "Created by GreenHopper - do not edit or delete. Issue type for a big user story that needs to be broken down."
      },
      {
         "name" : "Story",
         "iconUrl" : "https://example.com/images/icons/issuetypes/story.png",
         "self" : "https://example.com/rest/api/2/issuetype/7",
         "id" : "7",
         "subtask" : false,
         "description" : ""
      }
   ],
   "id" : "11300"
}
`
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Fatalf("wanted GET but found %s\n", r.Method)
		}
		url := *r.URL
		if url.Path != "/rest/api/2/project/COOL" {
			t.Fatalf("Want /rest/api/2/project/COOL but got %s\n", url.Path)
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
	project, err := client.GetProject("COOL")
	if err != nil {
		t.Fatalf("Unexpected error: %s\n", err)
	}
	if project.ID != 11300 {
		t.Fatalf("Want 11300 but got %d\n", project.ID)
	}
}

func TestGetProjectNot200(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
	}))
	defer testServer.Close()

	url, _ := url.Parse(testServer.URL)
	client := NewClient("u", "p", url)
	_, err := client.GetProject("COOL")
	if err == nil {
		t.Fatalf("Expected an error\n")
	}
}
