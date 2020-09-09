package main

import (
	"net/http"

	"gopkg.in/yaml.v2"
)

// user defined type to store path and url of yaml file
type pathUrl struct {
	Path string `yaml:"path"`
	URL  string `yaml:"url"`
}

// This funciton will return http.HandlerFunc which will map path to
// their urls from the yaml file
func YAMLHandler(yamlBytes []byte, fallback http.Handler) (http.HandlerFunc, error) {
	pathUrls, err := parseYaml(yamlBytes)
	if err != nil {
		return nil, err
	}
	pathsToUrls := buildMap(pathUrls)

	return MapHandler(pathsToUrls, fallback), nil
}

// function to convert slice of type byte to slice of type pathUrl
func parseYaml(data []byte) ([]pathUrl, error) {
	var pathUrls []pathUrl
	err := yaml.Unmarshal(data, &pathUrls)
	if err != nil {
		return nil, err
	}
	return pathUrls, nil

}

// function to convert slice of type pathUrl to map.
func buildMap(pathUrls []pathUrl) map[string]string {
	pathsToUrls := make(map[string]string)
	for _, pu := range pathUrls {
		pathsToUrls[pu.Path] = pu.URL
	}
	return pathsToUrls
}
