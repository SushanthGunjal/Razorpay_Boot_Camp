package main

import (
	"encoding/json"
	"net/http"
)

// user defined type to store path and url of yaml file
type pathurl struct {
	Path string `json:"path"`
	URL  string `json:"url"`
}

// This funciton will return http.HandlerFunc which will map path to
// their urls
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if dest, ok := pathsToUrls[path]; ok {
			http.Redirect(w, r, dest, http.StatusFound)
			return
		}
		fallback.ServeHTTP(w, r)

	}
}

// This funciton will return http.HandlerFunc which will map path to
// their urls from the json file
func JSONHandler(jsonBytes []byte, fallback http.Handler) (http.HandlerFunc, error) {
	pathUrls, err := parseJson(jsonBytes)
	if err != nil {
		return nil, err
	}
	pathsToUrls := buildMap1(pathUrls)

	return MapHandler(pathsToUrls, fallback), nil
}

// Parse the json file and convert it into path and url format
func parseJson(data []byte) ([]pathurl, error) {
	var pathUrls []pathurl
	err := json.Unmarshal(data, &pathUrls)
	if err != nil {
		return nil, err
	}
	return pathUrls, nil
}

// build a map to store path and url
func buildMap1(pathUrls []pathurl) map[string]string {
	pathsToUrls := make(map[string]string)
	for _, pu := range pathUrls {
		pathsToUrls[pu.Path] = pu.URL
	}
	return pathsToUrls
}
