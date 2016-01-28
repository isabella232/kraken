package main

import "testing"

func TestValidate(t *testing.T) {
	errors := validate()
	if len(errors) != 4 {
		t.Fatalf("Want 4 but got %d\n", len(errors))
	}

	u := "foo"
	v := "foo"
	nextVersionName = &u
	releaseVersionName = &v
	errors = validate()
	if len(errors) != 4 {
		t.Fatalf("Want 4 but got %d\n", len(errors))
	}
}
