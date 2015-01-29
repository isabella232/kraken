package main

import (
	"errors"
	"testing"

	"github.com/xoom/jira"
)

type componentVersions struct {
	err      error
	versions map[int]jira.CVVersion
	mapping  jira.Mapping
	jira.ComponentVersions
}

func (r componentVersions) GetVersionsForComponent(projectID, componentID string) (map[int]jira.CVVersion, error) {
	return r.versions, r.err
}

func (r componentVersions) UpdateReleaseDate(mappingID int, releaseDate string) error {
	return r.err
}

func (r componentVersions) UpdateReleasedFlag(mappingID int, released bool) error {
	return r.err
}

func (r componentVersions) CreateMapping(projectID, componentID, versionID string) (jira.Mapping, error) {
	return r.mapping, r.err
}

func (r componentVersions) DeleteMapping(mappingID int) error {
	return r.err
}

func TestGetOrCreateMappingMapHit(t *testing.T) {
	mappings := map[int]jira.Mapping{1: jira.Mapping{ID: 1, ProjectID: 1, ComponentID: 2, VersionID: 3}}
	m, err := getOrCreateMapping("1", "2", "3", mappings, componentVersions{})
	if err != nil {
		t.Fatalf("Unexpected error: %v\n", err)
	}
	if m.ID != 1 {
		t.Fatalf("Want 1 but got: %v\n", m.ID)
	}
}

func TestGetOrCreateMappingMapMiss(t *testing.T) {
	mappings := map[int]jira.Mapping{}
	m, err := getOrCreateMapping("1", "2", "3", mappings, componentVersions{mapping: jira.Mapping{ID: 2, ProjectID: 4, ComponentID: 6, VersionID: 8}})
	if err != nil {
		t.Fatalf("Unexpected error: %v\n", err)
	}
	if m.ID != 2 {
		t.Fatalf("Want 2 but got: %v\n", m.ID)
	}
}

func TestGetOrCreateMappingMapMissAndError(t *testing.T) {
	client := componentVersions{err: errors.New("Boom")}
	_, err := getOrCreateMapping("1", "2", "3", map[int]jira.Mapping{}, client)
	if err == nil {
		t.Fatalf("Expected an error\n")
	}
}
