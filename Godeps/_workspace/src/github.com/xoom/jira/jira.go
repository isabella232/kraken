// JIRA API with Oguz Component Mappings
package jira

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

var logger = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)

type (
	Jira interface {
		Core
		ComponentVersions
	}

	// https://docs.atlassian.com/jira/REST/latest/
	Core interface {
		GetProject(projectKey string) (Project, error)
		GetComponents(projectID string) (map[string]Component, error)
		GetVersions(projectID string) (map[string]Version, error)
		CreateVersion(projectID, versionName string) (Version, error)
	}

	// http://jiraplugins.denizoguz.com/wp-content/uploads/2014/09/REST-Manual-v0.1.pdf
	ComponentVersions interface {
		GetMappings() (map[int]Mapping, error)
		GetVersionsForComponent(projectID, componentID string) (map[int]CVVersion, error)
		UpdateReleaseDate(mappingID int, releaseDate string) error
		UpdateReleasedFlag(mappingID int, released bool) error
		CreateMapping(projectID string, componentID string, versionID string) (Mapping, error)
		DeleteMapping(mappingID int) error
	}

	Project struct {
		ID string `json:"id"`
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
	}

	// Component Version add-on's notion of a version
	CVVersion struct {
		ID          int    `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
		Released    bool   `json:"isReleased"`
	}

	Mapping struct {
		ID             int    `json:"id"`
		ProjectID      int    `json:"projectId"`
		ProjectKey     string `json:"projectKey"`
		ProjectName    string `json:"projectName"`
		ComponentID    int    `json:"componentId"`
		VersionID      int    `json:"versionId"`
		VersionName    string `json:"versionName"`
		ComponentName  string `json:"componentName"`
		Released       bool   `json:"released"`
		ReleaseDateStr string `json:"releaseDateStr"`
	}

	nopLogger struct {
		io.WriteCloser
	}
)

// NewClient returns a new default Jira client for the given Jira admin username/password and base REST URL.
func NewClient(username, password string, baseURL *url.URL) Jira {
	return DefaultClient{username: username, password: password, baseURL: baseURL, httpClient: &http.Client{Timeout: 10 * time.Second}}
}

// GetProject returns a representation of a Jira project for the given project key.  An example of a key is MYPROJ.
func (client DefaultClient) GetProject(projectKey string) (Project, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/rest/api/2/project/%s", client.baseURL, projectKey), nil)
	if err != nil {
		return Project{}, err
	}
	req.Header.Set("Accept", "application/json")
	req.SetBasicAuth(client.username, client.password)

	responseCode, data, err := client.consumeResponse(req)
	if err != nil {
		return Project{}, err
	}
	if responseCode != http.StatusOK {
		logger.Printf("JIRA response: %s\n", string(data))
		return Project{}, fmt.Errorf("error getting project.  Status code: %d.\n", responseCode)
	}

	var r Project
	if err := json.Unmarshal(data, &r); err != nil {
		return Project{}, err
	}

	return r, nil
}

// GetComponents returns a map of Component indexed by component name for the given project ID.
func (client DefaultClient) GetComponents(projectID string) (map[string]Component, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/rest/api/2/project/%s/components", client.baseURL, projectID), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.SetBasicAuth(client.username, client.password)

	responseCode, data, err := client.consumeResponse(req)
	if err != nil {
		return nil, err
	}
	if responseCode != http.StatusOK {
		logger.Printf("JIRA response: %s\n", string(data))
		return nil, fmt.Errorf("error getting project components.  Status code: %d.\n", responseCode)
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

// GetVersions returns a map of Version indexed by version name for the given project ID.
func (client DefaultClient) GetVersions(projectID string) (map[string]Version, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/rest/api/2/project/%s/versions", client.baseURL, projectID), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.SetBasicAuth(client.username, client.password)

	responseCode, data, err := client.consumeResponse(req)
	if err != nil {
		return nil, err
	}
	if responseCode != http.StatusOK {
		logger.Printf("JIRA response: %s\n", string(data))
		return nil, fmt.Errorf("error getting project versions.  Status code: %d.\n", responseCode)
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
func (client DefaultClient) CreateVersion(projectID, versionName string) (Version, error) {
	i, err := strconv.Atoi(projectID)
	if err != nil {
		return Version{}, err
	}

	data, err := json.Marshal(&Version{Name: versionName, Description: "Version " + versionName, ProjectID: i, Archived: false, Released: false})
	if err != nil {
		return Version{}, err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/rest/api/2/version", client.baseURL), bytes.NewBuffer(data))
	if err != nil {
		return Version{}, err
	}
	req.Header.Set("Content-type", "application/json")
	req.SetBasicAuth(client.username, client.password)

	responseCode, data, err := client.consumeResponse(req)
	if err != nil {
		return Version{}, err
	}
	if responseCode != http.StatusCreated {
		logger.Printf("JIRA response: %s\n", string(data))
		return Version{}, fmt.Errorf("error creating project version.  Status code: %d.\n", responseCode)
	}

	var v Version
	if err := json.Unmarshal(data, &v); err != nil {
		return Version{}, err
	}
	return v, nil
}

// CreateMapping creates a mapping between the given component ID and version ID in the context of the given project ID.
func (client DefaultClient) CreateMapping(projectID, componentID, versionID string) (Mapping, error) {
	pId, err := strconv.Atoi(projectID)
	if err != nil {
		return Mapping{}, err
	}
	cId, err := strconv.Atoi(componentID)
	if err != nil {
		return Mapping{}, err
	}
	vId, err := strconv.Atoi(versionID)
	if err != nil {
		return Mapping{}, err
	}

	data, err := json.Marshal(&Mapping{ProjectID: pId, ComponentID: cId, VersionID: vId, Released: false})
	if err != nil {
		return Mapping{}, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/rest/com.deniz.jira.mapping/latest/", client.baseURL), bytes.NewBuffer(data))
	if err != nil {
		return Mapping{}, err
	}
	req.Header.Set("Content-type", "application/json")
	req.SetBasicAuth(client.username, client.password)

	response, err := client.httpClient.Do(req)
	if err != nil {
		return Mapping{}, err
	}
	defer response.Body.Close()

	data, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return Mapping{}, err
	}

	if response.StatusCode != http.StatusCreated {
		logger.Printf("JIRA response: %s\n", string(data))
		return Mapping{}, fmt.Errorf("error creating mapped version.  Status code: %d.\n", response.StatusCode)
	}

	if location, err := response.Location(); err != nil {
		log.Printf("create-mapping response has no Location header.  Unable to acquire created mapping ID\n")
		return Mapping{}, nil
	} else {
		h := location.String()
		id, err := strconv.Atoi(h[strings.LastIndex(h, "/")+1:])
		if err != nil {
			log.Printf("create-mapping response Location header %s has unparseable non-integer urlPrefix/N mapping ID\n", h)
			return Mapping{}, nil
		}
		return Mapping{ID: id}, nil
	}

}

// GetMappings returns all known mappings for all projects.
func (client DefaultClient) GetMappings() (map[int]Mapping, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/rest/com.deniz.jira.mapping/latest/mappings", client.baseURL), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.SetBasicAuth(client.username, client.password)

	responseCode, data, err := client.consumeResponse(req)
	if err != nil {
		return nil, err
	}
	if responseCode != http.StatusOK {
		logger.Printf("JIRA response: %s\n", string(data))
		return nil, fmt.Errorf("error getting mappings.  Status code: %d.\n", responseCode)
	}

	var r []Mapping
	if err := json.Unmarshal(data, &r); err != nil {
		return nil, err
	}

	m := make(map[int]Mapping)
	for _, c := range r {
		m[c.ID] = c
	}
	return m, nil
}

// GetVersionsForComponent returns the versions for the given component ID in the context of the given project ID.
func (client DefaultClient) GetVersionsForComponent(projectID, componentID string) (map[int]CVVersion, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/rest/com.deniz.jira.mapping/latest/applicable_versions?projectId=%s&selectedComponentIds=%s", client.baseURL, projectID, componentID), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.SetBasicAuth(client.username, client.password)

	responseCode, data, err := client.consumeResponse(req)
	if err != nil {
		return nil, err
	}
	if responseCode != http.StatusOK {
		logger.Printf("JIRA response: %s\n", string(data))
		return nil, fmt.Errorf("error getting mappings.  Status code: %d.\n", responseCode)
	}

	var r []CVVersion
	if err := json.Unmarshal(data, &r); err != nil {
		return nil, err
	}

	m := make(map[int]CVVersion)
	for _, c := range r {
		m[c.ID] = c
	}
	return m, nil
}

// UpdateReleaseDate updates the version release date to releaseDate for the given mapping ID.
func (client DefaultClient) UpdateReleaseDate(mappingID int, releaseDate string) error {
	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/rest/com.deniz.jira.mapping/latest/releaseDate/%d?releaseDate=%s", client.baseURL, mappingID, url.QueryEscape(releaseDate)), nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/json")
	req.SetBasicAuth(client.username, client.password)
	responseCode, data, err := client.consumeResponse(req)
	if err != nil {
		return err
	}
	if responseCode != http.StatusOK {
		logger.Printf("JIRA response: %s\n", string(data))
		return fmt.Errorf("error updating mapping release date.  Status code: %d.\n", responseCode)
	}
	return nil
}

// UpdateReleasedFlag updates the version released flag for the given mapping ID.
func (client DefaultClient) UpdateReleasedFlag(mappingID int, released bool) error {
	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/rest/com.deniz.jira.mapping/latest/releaseFlag/%d?isReleased=%v", client.baseURL, mappingID, released), nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/json")
	req.SetBasicAuth(client.username, client.password)
	responseCode, data, err := client.consumeResponse(req)
	if err != nil {
		return err
	}
	if responseCode != http.StatusOK {
		logger.Printf("JIRA response: %s\n", string(data))
		return fmt.Errorf("error updating mapping is-released flag.  Status code: %d.\n", responseCode)
	}
	return nil
}

// DeleteMapping deletes the mapping for the given mapping ID.
func (client DefaultClient) DeleteMapping(mappingID int) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/rest/com.deniz.jira.mapping/latest/%d", client.baseURL, mappingID), nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/json")
	req.SetBasicAuth(client.username, client.password)
	responseCode, data, err := client.consumeResponse(req)
	if err != nil {
		return err
	}
	if responseCode != http.StatusNoContent {
		logger.Printf("JIRA response: %s\n", string(data))
		return fmt.Errorf("error deleting mapping.  Status code: %d.\n", responseCode)
	}
	return nil
}

func (client DefaultClient) consumeResponse(req *http.Request) (rc int, buffer []byte, err error) {
	response, err := client.httpClient.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer response.Body.Close()

	if data, err := ioutil.ReadAll(response.Body); err != nil {
		return 0, nil, err
	} else {
		return response.StatusCode, data, nil
	}
}
