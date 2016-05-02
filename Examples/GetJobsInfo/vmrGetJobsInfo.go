// vmrGetJobsInfo.go - provides information about running jobs on vmray
//
//  go run vmrGetJobsInfo.go
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

func init() {
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	flag.Parse()
	c, err := vmray.New(
		//vmray.SetErrorLog(log.New(os.Stderr, "vmray error: ", log.Lshortfile)),
		//vmray.SetTraceLog(log.New(os.Stderr, "vmray trace: ", log.Lshortfile)),
		vmray.SetBasicAuth(os.Getenv("VMRAY_EMAIL"), os.Getenv("VMRAY_PASSWD")),
	)
	check(err)

	r, err := c.GetJobsInfo()
	check(err)
	j, err := json.MarshalIndent(r, "", "    ")
	fmt.Printf("GetJobsInfo:\n")
	os.Stdout.Write(j)
	fmt.Println()
	// Do something with the result, parse out interesting info
}
