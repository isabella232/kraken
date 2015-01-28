package main

import "github.com/xoom/jira"

type core struct {
	err        error
	project    jira.Project
	version    jira.Version
	components map[string]jira.Component
	versions   map[string]jira.Version
	jira.Core
}

func (r core) GetProject(projectKey string) (jira.Project, error) {
	return r.project, r.err
}

func (r core) GetComponents(projectKey string) (map[string]jira.Component, error) {
	return r.components, r.err
}

func (r core) GetVersions(projectKey string) (map[string]jira.Version, error) {
	return r.versions, r.err
}

func (r core) CreateVersion(projectID, versionName string) (jira.Version, error) {
	return r.version, r.err
}
