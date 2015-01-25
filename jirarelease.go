package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/xoom/jira"
)

var (
	projectID     = flag.Int("project-id", -1, "JIRA project ID.")
	baseURL       = flag.String("jira-base-url", "http://localhost:8080", "JIRA base REST URL.")
	componentName = flag.String("component-name", "", "JIRA project component name.")
	username      = flag.String("jira-username", "", "JIRA admin user.")
	password      = flag.String("jira-password", "", "JIRA admin password.")
	versionName   = flag.String("version-name", "", "JIRA component version name.")
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
	if err != nil {
		log.Fatalf("Cannot parse JIRA base URL %s: %v\n", *baseURL, err)
	}

	jiraClient := jira.NewClient(*username, *password, url)

	versions, err := jiraClient.GetVersions(*projectID)
	if err != nil {
		log.Fatalf("Cannot get versions: %+v\n", err)
	}

	components, err := jiraClient.GetComponents(*projectID)
	if err != nil {
		log.Fatalf("Cannot get versions%v\n", err)
	}

	if _, present := versions[*versionName]; present {
		log.Fatalf("Version %s already exists.\n", *versionName)
	}

	if _, present := components[*componentName]; !present {
		log.Fatalf("Component %s does not exist.\n", *componentName)
	}

	err = jiraClient.CreateVersion(*projectID, *versionName)
	fmt.Printf("%+v\n", err)
}

func validate() []error {
	errors := make([]error, 0)
	if *projectID == -1 {
		errors = append(errors, fmt.Errorf("project-id must be provided"))
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
