package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/boltdb/bolt"
)

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, world!")
}

func main() {
	mux := defaultMux()

	// Connect to the BoltDB
	db, err := bolt.Open("my.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	fmt.Println(db)

	// Create a bucket in the database
	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket([]byte("MyBucket"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})

	// Insert the key value pair into the database
	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("MyBucket"))
		err := b.Put([]byte("/gd"), []byte("https://www.google.com/drive"))
		return err
	})

	dbHandler := DBHandler(db, mux)
	fmt.Println("Starting the server on :8080")
	fmt.Println(http.ListenAndServe(":8080", dbHandler))

	// variable to store the yaml file.
	yamlFile := flag.String("yamlFile", "links.yaml", "a yaml file containing short urls")

	// variable to store the json file
	jsonFile := flag.String("jsonFile", "links1.json", "a json file which contains short urls")

	//check whether the yaml file is present or not
	_, err = os.Open(*yamlFile)
	if err != nil {
		log.Fatal(err)
	}

	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := MapHandler(pathsToUrls, mux)

	Yaml := `
  - path: /urlshort
    url: https://github.com/gophercises/urlshort
  - path: /urlshort-final
    url: https://github.com/gophercises/urlshort/tree/final
  `

	yamlHandler, err := YAMLHandler([]byte(Yaml), mapHandler)
	if err != nil {
		fmt.Printf("Error parsing YAML file: %s\n", err)
	}

	// reads  the file and convert the yaml file into slice of bytes format.
	YamlFile, err := ioutil.ReadFile(*yamlFile)

	//fmt.Println(YamlFile)
	YamlHandler, err := YAMLHandler([]byte(YamlFile), mapHandler)
	if err != nil {
		fmt.Printf("Error parsing YAML file: %s\n", err)
	}

	// check whether the json file is present or not.
	_, err = os.Open(*jsonFile)
	if err != nil {
		fmt.Println(err)
	}

	// reads the file and converts the json file into slice of bytes format
	byteValue, _ := ioutil.ReadFile(*jsonFile)
	jsonHandler, err := JSONHandler(byteValue, mapHandler)
	if err != nil {
		fmt.Printf("Error in parsing the json file: %s \n", err)
	}

	fmt.Println("Starting the server on :8080")
	fmt.Println(http.ListenAndServe(":8080", jsonHandler))

	fmt.Println("Starting the server on :8080")
	fmt.Println(http.ListenAndServe(":8080", YamlHandler))

	fmt.Println("Starting the server on :8080")
	fmt.Println(http.ListenAndServe(":8080", yamlHandler))

	fmt.Println("Starting the server on :8080")
	fmt.Println(http.ListenAndServe(":8080", mapHandler))

	fmt.Println(http.ListenAndServe(":8080", nil))
}
