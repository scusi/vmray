// vmrDownloadAnalysis.go - can be used to download a complete analysis from vmray
//
// Example Usage:
//  go run vmrDownloadAnalysis.go -analysis_id=12345
//
package main

import (
	"flag"
	"fmt"
	"github.com/scusi/vmray"
	"io/ioutil"
	"log"
	"os"
	//"net/http"
	//"crypto/tls"
)

var analysis_id string // analysis_id

func init() {
	flag.StringVar(&analysis_id, "analysis_id", "", "ID of Sample to retrieve Infos for")
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	flag.Parse()
	/*
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		httpclient := &http.Client{Transport: tr}
	*/
	c, err := vmray.New(
		vmray.SetErrorLog(log.New(os.Stderr, "vmray error: ", log.Lshortfile)),
		//vmray.SetTraceLog(log.New(os.Stderr, "vmray trace: ", log.Lshortfile)),
		vmray.SetBasicAuth(os.Getenv("VMRAY_EMAIL"), os.Getenv("VMRAY_PASSWD")),
		//vmray.SetHttpClient(httpclient),
	)
	check(err)
	//fmt.Printf("vmray client: %#v\n", c)

	data, err := c.DownloadAnalysis(analysis_id)
	check(err)
	err = ioutil.WriteFile(analysis_id+".zip", data, 0750)
	check(err)
	fmt.Printf("downloaded analysis to: %s\n", analysis_id+".zip")
}
