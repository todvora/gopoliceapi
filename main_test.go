package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"net/url"
	"github.com/stretchr/testify/assert"
	"encoding/json"
	"errors"
)

func TestToSearchQueryOnlyQ(t *testing.T) {
	values := make(url.Values)
	values["q"] =  []string{"FOOBAR"}
	vin, regno := toSearchQuery(values)
	assert.Equal(t, "FOOBAR", vin);
	assert.Equal(t, "FOOBAR", regno);
}

func TestToSearchQueryOnlyVin(t *testing.T) {
	values := make(url.Values)
	values["vin"] =  []string{"FOOBAR"}
	vin, regno := toSearchQuery(values)
	assert.Equal(t, "FOOBAR", vin);
	assert.Equal(t, "", regno);
}

func TestToSearchQueryOnlyRegno(t *testing.T) {
	values := make(url.Values)
	values["regno"] =  []string{"FOOBAR"}
	vin, regno := toSearchQuery(values)
	assert.Equal(t, "", vin);
	assert.Equal(t, "FOOBAR", regno);
}

func TestHealthHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(healthHandler)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := "OK"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestStaticHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(staticHandler)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expectedContentType := "text/html; charset=utf-8"
	actualContentType := rr.Header().Get("content-type")
	if actualContentType != expectedContentType {
		t.Errorf("handler returned unexpected content-type: got %v want %v", actualContentType, expectedContentType)
	}

	if !strings.Contains(rr.Body.String(), "Unofficial API to the official Czech police database") {
		t.Errorf("handler returned unexpected body: got %v", rr.Body.String())
	}
}

func TestStaticHandlerNonexistent(t *testing.T) {
	req, err := http.NewRequest("GET", "/foobar.html", nil)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(staticHandler)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.

	assert.Equal(t, http.StatusNotFound, rr.Code, "handler returned wrong status code")
}


func TestSearchHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/search?q=ABCD", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	searchCallback := func(vin string, regno string) Results {
		return toErrorResults(errors.New("Failed to load PCR page"))
	}

	handler := searchHandler(searchCallback)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	jsonResponse := &Results{}
	err = json.Unmarshal(rr.Body.Bytes(), jsonResponse)

	assert.Equal(t, "Failed to load PCR page", jsonResponse.Error)
	assert.Equal(t, 200, rr.Code)
}