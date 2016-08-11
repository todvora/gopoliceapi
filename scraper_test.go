package main

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestApiHandler(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/patrani-vozidla/default.aspx") {
			http.ServeFile(w, r, "testfiles/list.html")
		} else if strings.HasPrefix(r.URL.Path, "/patrani-vozidla/Detail.aspx") {
			http.ServeFile(w, r, "testfiles/detail.html")
		} else {
			http.Error(w, "Mocked URL not found!", http.StatusNotFound)
		}
	}))
	defer server.Close()
	// Make a transport that reroutes all traffic to the example server
	transport := &http.Transport{
		Proxy: func(req *http.Request) (*url.URL, error) {
			return url.Parse(server.URL)
		},
	}

	// Make a http.Client with the transport
	httpClient := &http.Client{Transport: transport}
	apiClient := newClient(httpClient)

	results := apiClient.Search("WUAMMM4F58N901800", "")
	assert.Equal(t, 1, results.Count, "API didn't got all the results form the search page")

	item := results.Results[0]

	assert.Equal(t, item["class"], "osobní vozidlo")
	assert.Equal(t, item["manufacturer"], "AUDI")
	assert.Equal(t, item["type"], "RS 6")
	assert.Equal(t, item["color"], "červená metalíza")
	assert.Equal(t, item["regno"], "9Q91234")
	assert.Equal(t, item["vin"], "WUAMMM4F58N901800")
	assert.Equal(t, item["engine"], "ABC123DEF")
	assert.Equal(t, item["productionyear"], "2008")
	assert.Equal(t, item["stolendate"], "1.3.2012")
	assert.Equal(t, item["url"], "http://aplikace.policie.cz/patrani-vozidla/Detail.aspx?id=987654321")
}
