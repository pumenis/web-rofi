package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	defaultCGI := filepath.Join(os.Getenv("HOME"), ".cgi-bin")
	cgiPath := flag.String("cgi", defaultCGI, "Path to CGI scripts directory")
	port := flag.String("port", "1121", "Port to run the HTTP server on")
	flag.Parse()

	http.HandleFunc("/rofi", func(wr http.ResponseWriter, r *http.Request) {
		script := r.URL.Query().Get("script")
		if script == "" {
			http.Error(wr, "Missing 'script' query parameter", http.StatusBadRequest)
			return
		}

		scriptPath := filepath.Join(*cgiPath, script)
		if _, err := os.Stat(scriptPath); err != nil {
			http.Error(wr, "Script not found: "+err.Error(), http.StatusNotFound)
			return
		}

		args := r.URL.Query()["arg"]
		cmd := exec.Command(scriptPath, args...)
		output, err := cmd.CombinedOutput()
		if err != nil {
			http.Error(wr, "Script execution error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		lines := strings.Split(string(output), "\n")
		if len(lines) == 0 {
			http.Error(wr, "No output from script", http.StatusInternalServerError)
			return
		}

		rawurl := lines[0]
		options := lines[1:]

		parsed, err := url.Parse(rawurl)
		if err != nil {
			http.Error(wr, "Invalid URL: "+err.Error(), http.StatusInternalServerError)
			return
		}

		queryValues, err := url.ParseQuery(parsed.RawQuery)
		if err != nil {
			http.Error(wr, "Invalid query string: "+err.Error(), http.StatusInternalServerError)
			return
		}

		parsed.RawQuery = queryValues.Encode()
		fmt.Fprintln(wr, GetHTML(options, parsed.String()))
	})

	http.HandleFunc("/view", func(wr http.ResponseWriter, r *http.Request) {
		script := r.URL.Query().Get("script")
		if script == "" {
			http.Error(wr, "Missing 'script' query parameter", http.StatusBadRequest)
			return
		}

		scriptPath := filepath.Join(*cgiPath, script)
		if _, err := os.Stat(scriptPath); err != nil {
			http.Error(wr, "Script not found: "+err.Error(), http.StatusNotFound)
			return
		}

		args := r.URL.Query()["arg"]
		cmd := exec.Command(scriptPath, args...)
		output, err := cmd.CombinedOutput()
		if err != nil {
			http.Error(wr, "Script execution error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Fprintln(wr, string(output))
	})

	log.Printf("Starting HTTP server on :%s\n", *port)
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}
