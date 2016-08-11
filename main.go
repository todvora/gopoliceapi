package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"
	"net/url"
	"os"
	"path"
)

func toSearchQuery(url url.Values) (vin string, regno string) {
	vin = url.Get("vin")
	regno = url.Get("regno")
	// fallback, no vin or regno, fill everywhere value from 'q' paramether. It should cover cases, where
	// the caller doesn't know, which type of the identificator is it.
	if len(vin) == 0 && len(regno) == 0 {
		query := url.Get("q")
		vin = query
		regno = query
	}
	return vin, regno
}

func searchHandler(searchCallback func(vin string, regno string) Results) http.Handler {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
		results := searchCallback(toSearchQuery(r.URL.Query()))
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(results)
	})
}

func getContentType(url string) string {
	extension := path.Ext(url)
	mime := mime.TypeByExtension(extension)
	return mime
}

// openshift health monitoring
func healthHandler(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprint(w, "OK")
}

//go:generate go-bindata -prefix "assets/" assets/...
func staticHandler(w http.ResponseWriter, req *http.Request) {
	var path string = req.URL.Path[1:]
	if path == "" {
		path = "index.html"
	}

	if bs, err := Asset(path); err != nil {
		http.Error(w, "Page "+path+" not found!", http.StatusNotFound)
	} else {
		w.Header().Set("Content-Type", getContentType(path))
		var reader = bytes.NewBuffer(bs)
		io.Copy(w, reader)
	}
}


func main() {
	client := defaultClient()
	http.Handle("/search", searchHandler(client.Search))
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/", staticHandler)

	ip := os.Getenv("OPENSHIFT_GO_IP")
	port := os.Getenv("OPENSHIFT_GO_PORT")

	if len(port) == 0 {
		port = "8080"
	}

	bind := fmt.Sprintf("%s:%s", ip, port)
	fmt.Printf("listening on %s...\n", bind)

	log.Fatal(http.ListenAndServe(bind, nil))
}
