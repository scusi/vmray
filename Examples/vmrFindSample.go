// vmrFindSample.go - can be used to find a sample by its hash on vmray
//
// EXAMPLE USAGE:
//  go run vmrFindSample.go -rsrc="07bd860cf26e56a02bbf1b0ab6874578"

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/scusi/vmray"
	"log"
	"os"
)

var rsrc string // resource / hash

func init() {
	flag.StringVar(&rsrc, "rsrc", "", "resource to search for; Hash like MD5, SHA1,... for sample")
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

	f, err := c.FindSample(rsrc)
	check(err)
	j, err := json.MarshalIndent(f, "", "    ")
	fmt.Printf("FindSample:\n")
	os.Stdout.Write(j)
	fmt.Println()
}
