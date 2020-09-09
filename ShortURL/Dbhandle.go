package main

import (
	"fmt"
	"net/http"

	"github.com/boltdb/bolt"
)

// This function will return the http.HandlerFunc which will map the path
// to urls from BoltDB
func DBHandler(db *bolt.DB, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//db, _ := bolt.Open("my.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
		path := r.URL.Path
		var path1 string
		fmt.Println(path)
		db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("MyBucket"))
			v := b.Get([]byte(path))
			path1 = string(v)
			fmt.Printf("The url is: %s\n", path1)
			//tx.DeleteBucket([]byte("MyBucket"))
			return nil
		})
		fmt.Println(path1)
		desturl := path1
		http.Redirect(w, r, desturl, http.StatusFound)
	}
}
