// vmrGetAnalysisInfo.go - take a sample_id and returns information about analysis available for that sample
//
//  vmrGetAnalysisInfo -sample_id 12345
//
package main

import (
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/scusi/vmray"
	"log"
	"net/http"
	"os"
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
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	httpclient := &http.Client{Transport: tr}
	c, err := vmray.New(
		//vmray.SetUrl("https://live.vmray.com/api/"),
		//vmray.SetErrorLog(log.New(os.Stderr, "vmray error: ", log.Lshortfile)),
		//vmray.SetTraceLog(log.New(os.Stderr, "vmray trace: ", log.Lshortfile)),
		vmray.SetBasicAuth(os.Getenv("VMRAY_EMAIL"), os.Getenv("VMRAY_PASSWD")),
		vmray.SetHttpClient(httpclient),
	)
	check(err)
	//fmt.Printf("vmray client: %#v\n", c)

	r, err := c.GetAnalysisInfo(sample_id)
	check(err)
	//fmt.Printf("GetAnalysis Result (raw): %#v\n", r)

	// Print out JSON indeted
	j, err := json.MarshalIndent(r, "", "    ")
	fmt.Printf("GetAnalysisInfo for SampleId: %s\n", sample_id)
	os.Stdout.Write(j)

	// Do something with the result, parse out interesting info
	/*
		for name, itfc := range r.Analyses {
			log.Printf("name: %v, interface: %v\n", name, itfc)
			for key, value := range itfc.(map[string]interface{}) {
				//log.Printf("key: %s value: %v, type of value: %T\n", key, value, value)
				fmt.Printf("%T `json:\"%s\"`\n", value, key)
			}

		}
	*/

	// der folgende Teil funktioniert.
	/*
			var ana MyAnalyses
			var analysisList []MyAnalyses

			for _, i := range r.Analyses {
				byt, err := json.Marshal(i.(map[string]interface{}))
				if err != nil {
					log.Println(err)
				}
				//fmt.Printf("%s\n", byt)
				err = json.Unmarshal(byt, &ana)
				if err != nil {
					log.Println(err)
				}
				//fmt.Printf("%#v\n", analysis)
				analysisList = append(analysisList, ana)
			}
			//fmt.Printf("%#v\n", analysisList)
			j, err := json.MarshalIndent(analysisList, "", "    ")
		    fmt.Printf("GetAnalysisInfo for SampleId: %s\n", sample_id)
			os.Stdout.Write(j)
	*/
	/*
	   	var ana Ana
	   	j, err := json.MarshalIndent(r, "", "    ")
	   	err = json.Unmarshal(j, &ana)
	   	if err != nil {
	   		log.Fatal(err)
	   	}
	   	//fmt.Printf("%#v\n", ana)
	   	for k, v := range ana.Analyses {
	   		//log.Printf("%s %T %.0f\n", k, v, v.AnalysisServerity)
	   		log.Printf("AnalysisId: %s Result: %s, Serverity: %.0f Target: %s\n", k, v.AnalysisResult, v.AnalysisServerity, v.Target)
	   	}
	   }

	   type Ana struct {
	   	Analyses map[string]MyAnalyses	`json:"Analyses"`
	   }
	   type MyAnalyses struct {
	   	AnalyzerType		string `json:"analyzer_type"`
	   	AnalysisSnapshotId	float64 `json:"analysis_snapshot_id"`
	   	VmhostName			string `json:"vmhost_name"`
	   	AnalysisCreated		string `json:"analysis_created"`
	   	AnalysisSize		float64 `json:"analysis_size"`
	   	AnalysisJobStarted	string `json:"analysis_job_started"`
	   	SnapshotName		string `json:"snapshot_name"`
	   	AnalysisResult		string `json:"analysis_result"`
	   	AnalysisJobId		float64 `json:"analysis_job_id"`
	   	AnalysisCmdlineId	float64 `json:"analysis_cmdline_id"`
	   	AnalysisConfigurationID	float64 `json:"analysis_configuration_id"`
	   	AnalysisUserConfigID	float64 `json:"analysis_user_config_id"`
	   	AnalyzerName		string `json:"analyzer_name"`
	   	AnalysisJobruleId	float64 `json:"analysis_jobrule_id"`
	   	AnalysisPriority	float64 `json:"analysis_priority"`
	   	Target				string `json:"target"`
	   	AnalysisHint		float64 `json:"analysis_hint"`
	   	AnalysisAnalyzerVersion	string `json:"analysis_analyzer_version"`
	   	ConfigurationName	string `json:"configuration_name"`
	   	AnalysisUserId		float64 `json:"analysis_user_id"`
	   	AnalysisId			float64 `json:"analysis_id"`
	   	AnalysisExternalReference	string `json:"analysis_external_reference"`
	   	AnalysisVmhostId	float64 `json:"analysis_vmhost_id"`
	   	VmName				string `json:"vm_name"`
	   	AnalysisAnalyzerId	float64 `json:"analysis_analyzer_id"`
	   	AnalysisPrescriptId	float64 `json:"analysis_prescript_id"`
	   	AnalysisSampleId	float64 `json:"analysis_sample_id"`
	   	AnalysisServerity	float64 `json:"analysis_severity"`
	   	AnalysisVmId		float64 `json:"analysis_vm_id"`
	   }
	*/
}
