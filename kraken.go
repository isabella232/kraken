package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/xoom/jira"
)

var (
	// required
	baseURL            = flag.String("jira-base-url", "http://localhost:8080", "JIRA base REST URL.  Required.")
	username           = flag.String("jira-username", "", "JIRA admin user.  Required.")
	password           = flag.String("jira-password", "", "JIRA admin password.  Required.")
	projectKey         = flag.String("project-key", "", "JIRA project key.  For example, PLAT.  Required.")
	releaseVersionName = flag.String("release-version-name", "", "JIRA release version name. For example, 1.1.  Required.")
	componentName      = flag.String("component-name", "", "JIRA project component name.  For example, rest-server.  Required.")

	// optional
	nextVersionName = flag.String("next-version-name", "", "JIRA next version name. For example, 1.2.  Optional.")
	versionFlag     = flag.Bool("version", false, "Print version and exit.")

	buildInfo string
)

func init() {
	flag.Parse()
	log.Printf("%s\n", buildInfo)
	if *versionFlag {
		os.Exit(0)
	}
}

func main() {
	if err := validate(); len(err) != 0 {
		for _, i := range err {
			log.Printf("Error: %+v\nExiting.", i)
		}
		os.Exit(0)
	}

	url, err := url.Parse(*baseURL)
	check("Error parsing Jira base URL", err)

	jiraClient := jira.NewClient(*username, *password, url)

	project, err := jiraClient.GetProject(*projectKey)
	check("Error getting projects", err)
	log.Printf("Retrieved project: %s\n", *projectKey)

	versions, err := jiraClient.GetVersions(project.ID)
	check("Error getting project versions", err)
	log.Printf("Retrieved %d project versions\n", len(versions))

	components, err := jiraClient.GetComponents(project.ID)
	check("Error getting project components", err)
	log.Printf("Retrieved %d project components\n", len(components))

	component, present := components[*componentName]
	if !present {
		log.Printf("Component %s does not exist.\nExiting.\n", *componentName)
		os.Exit(0)
	}

	// get mappings for all projects
	mappings, err := jiraClient.GetMappings()
	check("Error getting mappings", err)

	// fetch or create release-version
	releaseVersion, err := getOrCreateVersion(project.ID, *releaseVersionName, versions, jiraClient)
	check("Error getOrCreateVersion()", err)

	// Create the release-version mapping if it does not exist
	releaseMapping, err := getOrCreateMapping(project.ID, component.ID, releaseVersion.ID, mappings, jiraClient)
	check("Error getOrCreateMapping()", err)

	// Do not update a mapping that is already released.
	if !releaseMapping.Released {
		err = jiraClient.UpdateReleasedFlag(releaseMapping.ID, true)
		check("Error updating release flag for release-version", err)
		log.Printf("Updated released flag for release mapping %+v\n", releaseMapping)

		err = jiraClient.UpdateReleaseDate(releaseMapping.ID, today())
		check("Error updating release date for release-version", err)
		log.Printf("Updated release date for release mapping %+v\n", releaseMapping)
	} else {
		log.Printf("Skipping already released release mapping: %+v\n", releaseMapping)
	}

	// next-version
	if *nextVersionName != "" {
		nextVersion, err := getOrCreateVersion(project.ID, *nextVersionName, versions, jiraClient)
		check("Error creating next-version", err)

		// Create the next-version mapping if it does not exist.
		_, err = getOrCreateMapping(project.ID, component.ID, nextVersion.ID, mappings, jiraClient)
		check("Error creating next-version mapping", err)
	}
}

func getOrCreateVersion(projectID, versionName string, versions map[string]jira.Version, client jira.Core) (jira.Version, error) {
	var err error
	var present bool
	var version jira.Version
	version, present = versions[versionName]
	if !present {
		log.Printf("Creating project version %s ...\n", versionName)
		version, err = client.CreateVersion(projectID, versionName)
		if err != nil {
			return jira.Version{}, err
		}
		log.Printf("Created project version %s\n", version.Name)
	} else {
		log.Printf("Retrieved existing version %s\n", version.Name)
	}
	return version, nil
}

func getOrCreateMapping(projectID, componentID, releaseVersionID string, mappings map[int]jira.Mapping, client jira.ComponentVersions) (jira.Mapping, error) {
	var mapping jira.Mapping
	var present bool
	var err error

	mapping, present = findMapping(mappings, projectID, componentID, releaseVersionID)
	if !present {
		mapping, err = client.CreateMapping(projectID, componentID, releaseVersionID)
		if err != nil {
			return jira.Mapping{}, err
		}
		log.Printf("Created version mapping: %d\n", mapping.ID)
	} else {
		log.Printf("Retrieved existing version mapping: %d\n", mapping.ID)
	}
	return mapping, nil
}

func findMapping(mappings map[int]jira.Mapping, projectID, componentID, versionID string) (jira.Mapping, bool) {
	for _, mapping := range mappings {
		pID := fmt.Sprintf("%d", mapping.ProjectID)
		cID := fmt.Sprintf("%d", mapping.ComponentID)
		vID := fmt.Sprintf("%d", mapping.VersionID)
		if pID == projectID && cID == componentID && vID == versionID {
			return mapping, true
		}
	}
	return jira.Mapping{}, false
}

func today() string {
	t := time.Now()
	return fmt.Sprintf("%d/%s/%d", t.Day(), t.Month().String()[:3], t.Year()%100)
}

func check(message string, err error) {
	if err != nil {
		log.Printf("%s: %+v\nExiting.", message, err)
		os.Exit(0)
	}
}

func validate() []error {
	errors := make([]error, 0)
	if *baseURL == "" {
		errors = append(errors, fmt.Errorf("jira-base-url must be provided"))
	}
	if *username == "" {
		errors = append(errors, fmt.Errorf("jira-username must be provided"))
	}
	if *password == "" {
		errors = append(errors, fmt.Errorf("jira-password must be provided"))
	}
	if *projectKey == "" {
		errors = append(errors, fmt.Errorf("project-key must be provided"))
	}
	if *componentName == "" {
		errors = append(errors, fmt.Errorf("component-name must be provided"))
	}
	if *releaseVersionName == "" {
		errors = append(errors, fmt.Errorf("release-version-name must be provided"))
	}
	if *nextVersionName != "" && *releaseVersionName == *nextVersionName {
		errors = append(errors, fmt.Errorf("release-version-name and next-version-name must be different"))
	}
	return errors
}
