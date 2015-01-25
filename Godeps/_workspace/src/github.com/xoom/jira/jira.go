package jira

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"strings"
	"time"
)

var debug bool

type (
	Jira interface {
		GetComponents(projectID int) (map[string]Component, error)
		GetVersions(projectID int) (map[string]Version, error)
		CreateVersion(projectID int, versionName string) error
		MapVersionToComponent(componentID, versionName string) error
	}

	DefaultClient struct {
		username   string
		password   string
		baseURL    *url.URL
		httpClient *http.Client
		Jira
	}

	Component struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	Version struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
		Project     string `json:"project"`
		ProjectID   int    `json:"projectId"`
		Archived    bool   `json:"archived"`
		Released    bool   `json:"released"`
		ReleaseDate string `json:"releaseDate"`
	}
)

func init() {
	debug = strings.ToLower(os.Getenv("JIRA_DEBUG")) == "true"
}

// NewClient returns a new default Jira client.
func NewClient(username, password string, baseURL *url.URL) Jira {
	return DefaultClient{username: username, password: password, baseURL: baseURL, httpClient: &http.Client{Timeout: 10 * time.Second}}
}

// GetComponents returns a map of Component indexed by component name.
func (client DefaultClient) GetComponents(projectID int) (map[string]Component, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/rest/api/2/project/%d/components", client.baseURL, projectID), nil)
	if err != nil {
		return nil, err
	}
	if debug {
		log.Printf("jira.GetComponents URL %s\n", req.URL)
	}
	req.Header.Set("Accept", "application/json")
	req.SetBasicAuth(client.username, client.password)

	responseCode, data, err := client.consumeResponse(req)
	if err != nil {
		return nil, err
	}
	if responseCode != http.StatusOK {
		return nil, fmt.Errorf("Error getting project versions.  Status code: %d.\n", responseCode)
	}

	var r []Component
	if err := json.Unmarshal(data, &r); err != nil {
		return nil, err
	}

	m := make(map[string]Component)
	for _, c := range r {
		m[c.Name] = c
	}

	return m, nil
}

// GetVersions returns a map of Version indexed by version name.
func (client DefaultClient) GetVersions(projectID int) (map[string]Version, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/rest/api/2/project/%d/versions", client.baseURL, projectID), nil)
	if err != nil {
		return nil, err
	}
	if debug {
		log.Printf("jira.GetVersions URL %s\n", req.URL)
	}
	req.Header.Set("Accept", "application/json")
	req.SetBasicAuth(client.username, client.password)

	responseCode, data, err := client.consumeResponse(req)
	if err != nil {
		return nil, err
	}
	if responseCode != http.StatusOK {
		return nil, fmt.Errorf("Error getting project versions.  Status code: %d.\n", responseCode)
	}

	var r []Version
	if err := json.Unmarshal(data, &r); err != nil {
		return nil, err
	}

	m := make(map[string]Version)
	for _, c := range r {
		m[c.Name] = c
	}

	return m, nil
}

// CreateVersion creates a new version in Jira for the given project ID and version name.
func (client DefaultClient) CreateVersion(projectID int, versionName string) error {
	version := Version{Name: versionName, Description: "Version " + versionName, ProjectID: projectID, Archived: false, Released: true, ReleaseDate: time.Now().Format("2006-01-02")}
	data, err := json.Marshal(&version)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/rest/api/2/version", client.baseURL), bytes.NewBuffer(data))
	if debug {
		log.Printf("jira.CreateVersion URL %s\n", req.URL)
	}
	if err != nil {
		return err
	}
	req.Header.Set("Content-type", "application/json")
	req.SetBasicAuth(client.username, client.password)

	responseCode, data, err := client.consumeResponse(req)
	if err != nil {
		return nil
	}
	if responseCode != http.StatusCreated {
		return fmt.Errorf("Error getting project versions.  Status code: %d.\n", responseCode)
	}
	return nil
}

func (client DefaultClient) MapVersionToComponent(componentID, versionName string) error {
	return nil
}

func (client DefaultClient) consumeResponse(req *http.Request) (rc int, buffer []byte, err error) {
	response, err := client.httpClient.Do(req)

	defer func() {
		if response != nil && response.Body != nil {
			response.Body.Close()
		}
		if e := recover(); e != nil {
			trace := make([]byte, 10*1024)
			_ = runtime.Stack(trace, false)
			log.Printf("%s", trace)
			err = fmt.Errorf("%v", e)
		}
	}()

	if err != nil {
		panic(err)
	}

	if data, err := ioutil.ReadAll(response.Body); err != nil {
		panic(err)
	} else {
		return response.StatusCode, data, nil
	}
}
