package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/xoom/jira"
)

var (
	baseURL            = flag.String("jira-base-url", "http://localhost:8080", "JIRA base REST URL.  Required.")
	username           = flag.String("jira-username", "", "JIRA admin user.  Required.")
	password           = flag.String("jira-password", "", "JIRA admin password.  Required.")
	projectKey         = flag.String("project-key", "", "JIRA project key.  For example, PLAT.  Required.")
	releaseVersionName = flag.String("release-version-name", "", "JIRA release version name. For example, 1.1.  Required.")
	componentName      = flag.String("component-name", "", "JIRA project component name.  For example, rest-server.  Required if stashkins-job-name is not provided.")
	jobName            = flag.String("stashkins-job-name", "", "Stashkins job name.  For example, eng-abcd-release, which extracts abcd as a component name  Required if component-name is not provided.")

	nextVersionName = flag.String("next-version-name", "", "JIRA next version name. For example, 1.2.  Optional.")
	versionFlag     = flag.Bool("version", false, "Print version and exit.")

	Log = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)

	buildInfo string
)

func init() {
	flag.Parse()
	Log.Printf("%s\n", buildInfo)
	if *versionFlag {
		os.Exit(0)
	}
}

func main() {
	if err := validate(); len(err) != 0 {
		for _, i := range err {
			Log.Printf("Error: %+v\nExiting.", i)
		}
		os.Exit(-1)
	}

	if *componentName == "" {
		*componentName = componentNameFromJobname(*jobName)
	}

	*nextVersionName = nextVersion(*nextVersionName)

	Log.Printf("Specified component: <%s>\n", *componentName)
	Log.Printf("Specified release version: <%s>\n", *releaseVersionName)
	if *nextVersionName != "" {
		Log.Printf("Specified next version name: <%s>\n", *nextVersionName)
	}

	url, err := url.Parse(*baseURL)
	if err != nil {
		Log.Printf("Error parsing Jira base URL: %v\n", err)
		os.Exit(-1)
	}

	jiraClient := jira.NewClient(*username, *password, url)

	project, err := jiraClient.GetProject(*projectKey)
	if err != nil {
		Log.Printf("Error getting projects: %v\n", err)
		os.Exit(-1)
	}
	Log.Printf("Found project: %s\n", *projectKey)

	versions, err := jiraClient.GetVersions(project.ID)
	if err != nil {
		Log.Printf("Error getting project versions: %v\n", err)
		os.Exit(-1)
	}
	Log.Printf("Found %d project versions\n", len(versions))

	components, err := jiraClient.GetComponents(project.ID)
	if err != nil {
		Log.Printf("Error getting project components: %v\n", err)
		os.Exit(-1)
	}
	Log.Printf("Found %d project components\n", len(components))

	component, present := components[*componentName]
	if !present {
		Log.Printf("Component %s does not exist.\nExiting.\n", *componentName)
		os.Exit(0)
	}

	// get mappings for all projects
	mappings, err := jiraClient.GetMappings()
	if err != nil {
		Log.Printf("Error getting mappings: %v\n", err)
		os.Exit(-1)
	}

	// fetch or create release-version
	releaseVersion, err := getOrCreateVersion(project.ID, *releaseVersionName, versions, jiraClient)
	if err != nil {
		Log.Printf("Error getOrCreateVersion(): %v\n", err)
		os.Exit(-1)
	}

	// Create the release-version mapping if it does not exist
	releaseMapping, err := getOrCreateMapping(project.ID, component.ID, releaseVersion.ID, mappings, jiraClient)
	if err != nil {
		Log.Printf("Error getOrCreateMapping(): %v\n", err)
		os.Exit(-1)
	}

	// Do not update a mapping that is already released.
	if !releaseMapping.Released {
		err = jiraClient.UpdateReleasedFlag(releaseMapping.ID, true)
		if err != nil {
			Log.Printf("Error updating release flag for release-version: %v\n", err)
		}

		if err != nil {
			Log.Printf("Error updating release flag for release-version: %v\n", err)
			return
		}

		err = jiraClient.UpdateReleaseDate(releaseMapping.ID, today())
		if err != nil {
			Log.Printf("Error updating release data for release-version: %v\n", err)
			return
		}
		Log.Printf("Updated release date for release mapping %+v\n", releaseMapping)
	} else {
		Log.Printf("Skipping already released release mapping: %+v\n", releaseMapping)
	}

	// next-version
	if *nextVersionName != "" {
		nextVersion, err := getOrCreateVersion(project.ID, *nextVersionName, versions, jiraClient)
		if err != nil {
			Log.Printf("Error creating next version: %v\n", err)
			return
		}

		// Create the next-version mapping if it does not exist.
		_, err = getOrCreateMapping(project.ID, component.ID, nextVersion.ID, mappings, jiraClient)
		if err != nil {
			Log.Printf("Error creating next-version mapping: %v\n", err)
			return
		}
	}
}

func getOrCreateVersion(projectID, versionName string, versions map[string]jira.Version, client jira.Core) (jira.Version, error) {
	var err error
	var present bool
	var version jira.Version
	version, present = versions[versionName]
	if !present {
		Log.Printf("Creating project version %s ...\n", versionName)
		version, err = client.CreateVersion(projectID, versionName)
		if err != nil {
			return jira.Version{}, err
		}
		Log.Printf("Created project version %s\n", version.Name)
	} else {
		Log.Printf("Retrieved existing version %s, no need to create it.\n", version.Name)
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
		Log.Printf("Created version mapping ID: %d\n", mapping.ID)
	} else {
		Log.Printf("Retrieved existing version mapping ID %d, no need to create it.\n", mapping.ID)
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

// For inputs not ending in -SNAPSHOT, return the input.  For inputs ending -SNAPSHOT, remove that suffix and return the result.
func nextVersion(version string) string {
	if strings.HasSuffix(version, "-SNAPSHOT") {
		return version[:strings.Index(version, "-SNAPSHOT")]
	}
	return version
}

// Jobs are assumed to be of the form proj-component-release, where proj- and -release are discarded and component is returned.
func componentNameFromJobname(jobName string) string {
	firstDash := strings.Index(jobName, "-")
	lastDash := strings.LastIndex(jobName, "-")
	return jobName[firstDash+1 : lastDash]
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
	if *releaseVersionName == "" {
		errors = append(errors, fmt.Errorf("release-version-name must be provided"))
	}
	if *nextVersionName != "" && *releaseVersionName == *nextVersionName {
		errors = append(errors, fmt.Errorf("release-version-name and next-version-name must be different"))
	}
	if *jobName != "" && *componentName != "" {
		errors = append(errors, fmt.Errorf("only one of component-naem or stashkins-job-name may be provided"))
	}
	if *jobName == "" && *componentName == "" {
		errors = append(errors, fmt.Errorf("one of component-naem or stashkins-job-name must be provided"))
	}

	return errors
}
