package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"runtime"
	"time"

	"github.com/xoom/jira"
)

var (
	projectKey         = flag.String("project-key", "", "JIRA project key.  For example, PLAT.")
	baseURL            = flag.String("jira-base-url", "http://localhost:8080", "JIRA base REST URL.")
	componentName      = flag.String("component-name", "", "JIRA project component name.  For example, rest-server.")
	username           = flag.String("jira-username", "", "JIRA admin user.")
	password           = flag.String("jira-password", "", "JIRA admin password.")
	releaseVersionName = flag.String("release-version-name", "", "JIRA release version name. For example, 1.1.")
	nextVersionName    = flag.String("next-version-name", "", "JIRA next version name. For example, 1.2.")
	versionFlag        = flag.Bool("version", false, "Print version and exit.")

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
	log.Printf("Retrieved %d project versions: %s\n", len(versions))

	components, err := jiraClient.GetComponents(project.ID)
	check("Error getting project components", err)
	log.Printf("Retrieved %d project components: %s\n", len(components))

	component, present := components[*componentName]
	if !present {
		log.Fatalf("Component %s does not exist.\n", *componentName)
	}

	version, present := versions[*releaseVersionName]
	log.Printf("Retrieved %d project versions: %s\n", len(versions))
	if !present {
		version, err = jiraClient.CreateVersion(project.ID, *releaseVersionName)
		check("Error creating version", err)
		log.Printf("Created project versions: %s\n", version.Name)
	}

	mapping, err := jiraClient.CreateMapping(project.ID, component.ID, version.ID)
	check("Error creating mapping", err)
	log.Printf("Created mapping: %d\n", mapping.ID)

	err = jiraClient.UpdateReleasedFlag(mapping.ID, true)
	check("Error updating release flag", err)

	err = jiraClient.UpdateReleaseDate(mapping.ID, today())
	check("Error updating release date", err)

}

func today() string {
	t := time.Now()
	return fmt.Sprintf("%d/%s/%d", t.Day(), t.Month().String()[:3], t.Year()%100)
}

func check(message string, err error) {
	trace := make([]byte, 10*1024)
	_ = runtime.Stack(trace, false)
	if err != nil {
		log.Fatalf("%s: %+v\n", message, err)
		log.Printf("%s", string(trace))
	}
}

func validate() []error {
	errors := make([]error, 0)
	if *projectKey == "" {
		errors = append(errors, fmt.Errorf("project-key must be provided"))
	}
	if *componentName == "" {
		errors = append(errors, fmt.Errorf("component-name must be provided"))
	}
	if *releaseVersionName == "" {
		errors = append(errors, fmt.Errorf("version-name must be provided"))
	}
	if *username == "" {
		errors = append(errors, fmt.Errorf("jira-username must be provided"))
	}
	if *password == "" {
		errors = append(errors, fmt.Errorf("jira-password must be provided"))
	}
	return errors
}
