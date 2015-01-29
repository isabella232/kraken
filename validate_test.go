package main

import "testing"

func TestValidate(t *testing.T) {
	errors := validate()
	if len(errors) != 5 {
		t.Fatalf("Want 5 but got %d\n", len(errors))
	}

	u := "foo"
	v := "foo"
	nextVersionName = &u
	releaseVersionName = &v
	errors = validate()
	if len(errors) != 5 {
		t.Fatalf("Want 5 but got %d\n", len(errors))
	}
}
