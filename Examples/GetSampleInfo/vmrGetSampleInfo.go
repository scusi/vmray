// vmrGetSampleInfo.go - takes a vmray sampleId, provides Information about the sample
//
package main

import (
	"encoding/json"
	"fmt"
	"github.com/scusi/vmray"
	"log"
	"os"
	"flag"
)

// sample_id variable
var sample_id string 

func init() {
	flag.StringVar(&sample_id, "sample_id", "", "ID of Sample to retrieve Infos for")
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	flag.Parse()
	c, err := vmray.New(
		vmray.SetErrorLog(log.New(os.Stderr, "vmray error: ", log.Lshortfile)),
		//vmray.SetTraceLog(log.New(os.Stderr, "vmray trace: ", log.Lshortfile)),
		vmray.SetBasicAuth(os.Getenv("VMRAY_EMAIL"), os.Getenv("VMRAY_PASSWD")),
	)
	check(err)

	f, err := c.GetSampleInfo(sample_id)
	check(err)
	j, err := json.MarshalIndent(f, "", "    ")
	fmt.Printf("Sample Information for %s\n", sample_id)
	os.Stdout.Write(j)
}
