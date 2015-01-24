package main

import (
	"flag"
	"log"
	"os"
)

var (
	projectID     = flag.Int("project-id", -1, "JIRA project ID.")
	baseURL       = flag.String("jira-base-url", "", "JIRA base REST URL.")
	componentName = flag.String("component-name", "", "JIRA project component name.")
	versionName   = flag.String("version-name", "", "JIRA component version name.")
	versionFlag   = flag.Bool("version", false, "Print version and exit.")

	version   string
	commit    string
	buildTime string
	sdkInfo   string
)

func main() {
	log.Printf("Version: %s, CommitID: %s, build time: %s, SDK Info: %s\n", version, commit, buildTime, sdkInfo)
	if *versionFlag {
		os.Exit(0)
	}
}
