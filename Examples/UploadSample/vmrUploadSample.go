// vmrUploadSample.go - uploads a file for analysis on vmray
//
//  go run vmrUploadSample.go -f <yourSampleFile.ext>
//
// Above example shows how to upload a given file to vmray useing vmrUploadSample.go
//
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/scusi/vmray"
	"log"
	"os"
)

var fileName string

func init() {
	flag.StringVar(&fileName, "f", "", "file to upload")
}

func main() {
	flag.Parse()
	// setup a new vmray client
	c, err := vmray.New(
		vmray.SetErrorLog(log.New(os.Stderr, "vmray error: ", log.Lshortfile)),
		//vmray.SetTraceLog(log.New(os.Stderr, "vmray trace: ", log.Lshortfile)),
		vmray.SetBasicAuth(os.Getenv("VMRAY_EMAIL"), os.Getenv("VMRAY_PASSWD")),
	)
	r, err := c.UploadSample(fileName)
	if err != nil {
		log.Printf("Error uploading file: '%s'\n", err.Error())
		os.Exit(1)
	}
	log.Printf("result raw: %#v\n", r)
	// Print out JSON indeted
	j, err := json.MarshalIndent(r, "", "    ")
	if err != nil {
		log.Printf("Error Marshal JSON: %s\n", err.Error())
	}
	fmt.Printf("Upload Sample '%s' results:\n", fileName)
	os.Stdout.Write(j)
}
