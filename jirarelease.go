package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"runtime"

	"github.com/xoom/jira"
)

var (
	projectKey    = flag.String("project-key", "", "JIRA project key.  For example, PLAT.")
	baseURL       = flag.String("jira-base-url", "http://localhost:8080", "JIRA base REST URL.")
	componentName = flag.String("component-name", "", "JIRA project component name.  For example, rest-server.")
	username      = flag.String("jira-username", "", "JIRA admin user.")
	password      = flag.String("jira-password", "", "JIRA admin password.")
	versionName   = flag.String("version-name", "", "JIRA component version name. For example, some-version.")
	versionFlag   = flag.Bool("version", false, "Print version and exit.")

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
	check(err)

	jiraClient := jira.NewClient(*username, *password, url)

	project, err := jiraClient.GetProject(*projectKey)
	check(err)

	versions, err := jiraClient.GetVersions(project.ID)
	check(err)

	components, err := jiraClient.GetComponents(project.ID)
	check(err)

	component, present := components[*componentName]
	if !present {
		log.Fatalf("Component %s does not exist.\n", *componentName)
	}

	version, present := versions[*versionName]
	if !present {
		version, err = jiraClient.CreateVersion(project.ID, *versionName)
		check(err)
	}

	_, err = jiraClient.CreateMapping(project.ID, component.ID, version.ID)
	check(err)
}

func check(err error) {
	trace := make([]byte, 10*1024)
	_ = runtime.Stack(trace, false)
	if err != nil {
		log.Fatalf("Error: %+v\n", err)
		log.Printf("%s", trace)
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
	if *versionName == "" {
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
