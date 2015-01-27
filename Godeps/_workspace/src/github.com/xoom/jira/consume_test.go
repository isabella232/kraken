package jira

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"testing"
)

type rt struct {
	err      error
	response *http.Response
}

func (r rt) RoundTrip(req *http.Request) (*http.Response, error) {
	if r.err != nil {
		return nil, r.err
	}
	return r.response, nil
}

func TestConsumeResponseWithError(t *testing.T) {
	r := rt{err: errors.New("boom!")}
	client := DefaultClient{httpClient: &http.Client{Transport: r}}
	req, err := http.NewRequest("GET", "http://localhost:8080", nil)
	if err != nil {
		t.Fatalf("Unexpected error: %v\n", err)
	}
	_, _, err = client.consumeResponse(req)
	if err == nil {
		t.Fatalf("Expected error\n")
	}
}

func TestConsumeResponse(t *testing.T) {
	buffer := bytes.NewBufferString("hello")
	body := ioutil.NopCloser(buffer)
	response := &http.Response{Status: "200 OK", StatusCode: 200, Proto: "HTTP/1.0", ProtoMajor: 1, ProtoMinor: 0, Body: body}
	r := rt{response: response}

	client := DefaultClient{httpClient: &http.Client{Transport: r}}

	req, err := http.NewRequest("GET", "http://localhost:8080", nil)
	if err != nil {
		t.Fatalf("Unexpected error: %v\n", err)
	}
	rc, data, err := client.consumeResponse(req)
	if err != nil {
		t.Fatalf("Unexpected error: %v\n", err)
	}
	if rc != 200 {
		t.Fatalf("Want 200 but got: %v\n", rc)
	}
	if string(data) != "hello" {
		t.Fatalf("Want hello but got: %v\n", string(data))
	}
}
