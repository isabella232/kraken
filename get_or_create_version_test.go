package main

import (
	"errors"
	"testing"

	"github.com/xoom/jira"
)

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

func TestGetOrCreateVersionMapHit(t *testing.T) {
	versions := map[string]jira.Version{"v1": jira.Version{Name: "v1"}, "v2": jira.Version{Name: "v2"}}
	v, err := getOrCreateVersion("1", "v1", versions, core{})
	if err != nil {
		t.Fatalf("Unexpected error: %v\n", err)
	}
	if v.Name != "v1" {
		t.Fatalf("Want v1 but got: %v\n", v.Name)
	}
}

func TestGetOrCreateVersionMapMiss(t *testing.T) {
	client := core{version: jira.Version{Name: "v1"}}
	v, err := getOrCreateVersion("1", "v1", map[string]jira.Version{}, client)
	if err != nil {
		t.Fatalf("Unexpected error: %v\n", err)
	}
	if v.Name != "v1" {
		t.Fatalf("Want v1 but got: %v\n", v.Name)
	}
}

func TestGetOrCreateVersionMapMissAndError(t *testing.T) {
	client := core{err: errors.New("Boom")}
	_, err := getOrCreateVersion("1", "v1", map[string]jira.Version{}, client)
	if err == nil {
		t.Fatalf("Expected an error\n")
	}
}
