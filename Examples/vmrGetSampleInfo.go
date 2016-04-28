package main

import (
	"encoding/json"
	"fmt"
	"github.com/scusi/vmray"
	"log"
	"os"
	//"net/http"
	//"crypto/tls"
	"flag"
)

var sample_id string // sample_id

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
	//tr := &http.Transport{
	//TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	//}
	//httpclient := &http.Client{Transport: tr}
	c, err := vmray.New(
		//vmray.SetUrl("https://live.vmray.com/api/"),
		vmray.SetErrorLog(log.New(os.Stderr, "vmray error: ", log.Lshortfile)),
		//vmray.SetTraceLog(log.New(os.Stderr, "vmray trace: ", log.Lshortfile)),
		vmray.SetBasicAuth(os.Getenv("VMRAY_EMAIL"), os.Getenv("VMRAY_PASSWD")),
		//vmray.SetHttpClient(httpclient),
		//vmray.SetHttpClient(&http.Client{}),
	)
	check(err)
	//fmt.Printf("vmray client: %#v\n", c)

	f, err := c.GetSampleInfo(sample_id)
	check(err)
	j, err := json.MarshalIndent(f, "", "    ")
	fmt.Printf("Sample Information for %s\n", sample_id)
	os.Stdout.Write(j)
}
