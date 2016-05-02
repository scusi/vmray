// vmrDoanloadAnalysesByHash.go - downloads all vmray analyses for a given sample, identified by its hash.
// file hash can be md5, sha1 or sha2 checksum of sample
//
//  go run vmrDownloadAnaysesByHash.go -hash d41d8cd98f00b204e9800998ecf8427e
//
package main

import (
	"flag"
	"github.com/scusi/vmray"
	"io/ioutil"
	"log"
	"os"
	"strconv"
)

var hash string
var email string
var passwd string

func init() {
	flag.StringVar(&hash, "hash", "", "hash (md5, sha1, sha2) of sample file to load analysis for from vmray")
	if os.Getenv("VMRAY_EMAIL") == "" {
		log.Fatal("Environment variable 'VMRAY_EMAIL' is not set or empty")
	}
	if os.Getenv("VMRAY_PASSWD") == "" {
		log.Fatal("Environment variable 'VMRAY_PASSWD' is not set or empty")
	}
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	flag.Parse()
	client, err := vmray.New(
		vmray.SetBasicAuth(os.Getenv("VMRAY_EMAIL"), os.Getenv("VMRAY_PASSWD")),
		vmray.SetErrorLog(log.New(os.Stderr, "vmray error: ", log.Lshortfile)),
		//vmray.SetTraceLog(log.New(os.Stderr, "vmray trace: ", log.Lshortfile)),
	)
	check(err)
	r, err := client.FindSample(hash)
	check(err)
	log.Printf("Found Sample '%d' for hash '%s'\n", r.SampleId, hash)
	sample_id := strconv.Itoa(r.SampleId)
	// GetAnalysisInfo to get analyses ids
	a, err := client.GetAnalysisInfo(sample_id)
	check(err)
	for k, _ := range a.Analyses {
		filename, err := downloadAndSaveAnalysis(client, k)
		check(err)
		//filename := k + ".zip"
		log.Printf("downloaded analysis to: %s\n", filename)
	}
}

func downloadAndSaveAnalysis(client *vmray.Client, analysis_id string) (filename string, err error) {
	// Download analysis file as zip
	data, err := client.DownloadAnalysis(analysis_id)
	if err != nil {
		return filename, err
	}
	filename = analysis_id + ".zip"
	err = ioutil.WriteFile(filename, data, 0750)
	if err != nil {
		return filename, err
	}
	return filename, err
}
