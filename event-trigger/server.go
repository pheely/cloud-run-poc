package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/", process)
	log.Printf("Listening on port 8080")
	err := http.ListenAndServe(":8080", nil)
	log.Fatal(err)
}
func process(w http.ResponseWriter, r *http.Request) {
	log.Println("Serving request")
	if r.Method == "GET" {
		fmt.Fprintln(w, "I am ready to process POST requests from bucket trigger")
		return
	}

	// retrieve the file info from request body containing Cloud Storage object metadata
	gcsInputFile, err1 := readBody(r)
	if err1 != nil {
		log.Printf("Error reading POST data: %v", err1)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Problem with POST data: %v \n", err1)
		return
	}
	// create a working directory
	localDir, err2 := ioutil.TempDir("", "")
	if err2 != nil {
		log.Printf("Error creating local temp dir: %v", err2)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Could not create a temp directory on server. \n")
		return
	}
	defer os.RemoveAll(localDir)

	// download input file from the bucket
	localInputFile, err3 := download(gcsInputFile, localDir)
	if err3 != nil {
		log.Printf("Error downloading Cloud Storage file [%s] from bucket [%s]: %v",
			gcsInputFile.Name, gcsInputFile.Bucket, err3)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error downloading Cloud Storage file [%s] from bucket [%s]",
			gcsInputFile.Name, gcsInputFile.Bucket)
		return
	}

        // Upload the file to the output bucket
	targetBucket := os.Getenv("OUT_BUCKET")
	err4 := upload(localInputFile.Name(), targetBucket)
	if err4 != nil {
		log.Printf("Error uploading file to bucket [%s]: %v", targetBucket, err4)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error downloading Cloud Storage file [%s] from bucket [%s]",
			gcsInputFile.Name, gcsInputFile.Bucket)
		return
	}

        // Delete the original input file from bucket.
	err5 := deleteGCSFile(gcsInputFile.Bucket, gcsInputFile.Name)
	if err5 != nil {
		log.Printf("Error deleting file [%s] from bucket [%s]: %v", gcsInputFile.Name,
			gcsInputFile.Bucket, err5)
	}
	log.Println("Successfully moved file")
	fmt.Fprintln(w, "Successfully moved file")
}
