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

	version   string
	commit    string
	buildTime string
	sdkInfo   string
)

func init() {
	flag.Parse()
	log.Printf("Version: %s, CommitID: %s, build time: %s, SDK Info: %s\n", version, commit, buildTime, sdkInfo)
	if *versionFlag {
		os.Exit(0)
	}
}

func main() {
	if err := validate(); len(err) != 0 {
		for _, i := range err {
			log.Printf("Error: %+v\n", i)
		}
		os.Exit(1)
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
		log.Fatalf("Component %s does not exist.\n", *componentName)
	}

	mappings, err := jiraClient.GetMappings()

	// fetch or create release-version
	releaseVersion, err := getOrCreateVersion(project.ID, *releaseVersionName, versions, jiraClient)
	check("Error getOrCreateVersion()", err)

	// Create the release-version mapping if it does not exist
	releaseMapping, present := findMapping(mappings, project.ID, component.ID, releaseVersion.ID)
	if !present {
		// create the release-version mapping, set released, and release with today's date
		releaseMapping, err := jiraClient.CreateMapping(project.ID, component.ID, releaseVersion.ID)
		check("Error creating mapping", err)
		log.Printf("Created release-version mapping: %d\n", releaseMapping.ID)
	}

	err = jiraClient.UpdateReleasedFlag(releaseMapping.ID, true)
	check("Error updating release flag for release-version", err)

	err = jiraClient.UpdateReleaseDate(releaseMapping.ID, today())
	check("Error updating release date for release-version", err)

	// next-version
	if *nextVersionName != "" {
		nextVersion, present := versions[*nextVersionName]
		if !present {
			log.Printf("Creating project next-version %s ...\n", nextVersion.Name)
			nextVersion, err = jiraClient.CreateVersion(project.ID, *nextVersionName)
			check("Error creating version", err)
			log.Printf("Created project next-version %s\n", nextVersion.Name)
		}
		// Create the next-version mapping if it does not exist
		if nextMapping, present := findMapping(mappings, project.ID, component.ID, nextVersion.ID); !present {
			// create the release-version mapping, set released, and release with today's date
			nextMapping, err := jiraClient.CreateMapping(project.ID, component.ID, nextVersion.ID)
			check("Error creating mapping", err)
			log.Printf("Created next-version mapping: %d\n", nextMapping.ID)
		} else {
			err = jiraClient.UpdateReleasedFlag(nextMapping.ID, true)
			check("Error updating release flag for next-version", err)
		}
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
	}
	return version, nil
}

func findMapping(mappings map[int]jira.Mapping, projectID, componentID, versionName string) (jira.Mapping, bool) {
	for _, mapping := range mappings {
		if fmt.Sprintf("%d", mapping.ProjectID) == projectID && fmt.Sprintf("%d", mapping.ComponentID) == componentID && mapping.VersionName == versionName {
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
		log.Fatalf("%s: %+v\n", message, err)
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
