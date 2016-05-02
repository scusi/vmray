// vmray api module for go
//
// vmray allows to communicate with the API of VmRay.
// VmRay is a 3rd generation malware execution and analysis environment.
// For more Information see: http://www.vmray.com/
//
// This module has been written by Florian 'scusi' Walther.
//
// For examples how to use this module see Examples directory.
//
package vmray

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

const (
	// DefaultURL specifies the default URL for VmRay API
	DefaultURL = "https://cloud.vmray.com/api/"
	// root CA cert for the DefaultURL
	GlobalSignRootCA = ` 
-----BEGIN CERTIFICATE-----
MIIDdTCCAl2gAwIBAgILBAAAAAABFUtaw5QwDQYJKoZIhvcNAQEFBQAwVzELMAkG
A1UEBhMCQkUxGTAXBgNVBAoTEEdsb2JhbFNpZ24gbnYtc2ExEDAOBgNVBAsTB1Jv
b3QgQ0ExGzAZBgNVBAMTEkdsb2JhbFNpZ24gUm9vdCBDQTAeFw05ODA5MDExMjAw
MDBaFw0yODAxMjgxMjAwMDBaMFcxCzAJBgNVBAYTAkJFMRkwFwYDVQQKExBHbG9i
YWxTaWduIG52LXNhMRAwDgYDVQQLEwdSb290IENBMRswGQYDVQQDExJHbG9iYWxT
aWduIFJvb3QgQ0EwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQDaDuaZ
jc6j40+Kfvvxi4Mla+pIH/EqsLmVEQS98GPR4mdmzxzdzxtIK+6NiY6arymAZavp
xy0Sy6scTHAHoT0KMM0VjU/43dSMUBUc71DuxC73/OlS8pF94G3VNTCOXkNz8kHp
1Wrjsok6Vjk4bwY8iGlbKk3Fp1S4bInMm/k8yuX9ifUSPJJ4ltbcdG6TRGHRjcdG
snUOhugZitVtbNV4FpWi6cgKOOvyJBNPc1STE4U6G7weNLWLBYy5d4ux2x8gkasJ
U26Qzns3dLlwR5EiUWMWea6xrkEmCMgZK9FGqkjWZCrXgzT/LCrBbBlDSgeF59N8
9iFo7+ryUp9/k5DPAgMBAAGjQjBAMA4GA1UdDwEB/wQEAwIBBjAPBgNVHRMBAf8E
BTADAQH/MB0GA1UdDgQWBBRge2YaRQ2XyolQL30EzTSo//z9SzANBgkqhkiG9w0B
AQUFAAOCAQEA1nPnfE920I2/7LqivjTFKDK1fPxsnCwrvQmeU79rXqoRSLblCKOz
yj1hTdNGCbM+w6DjY1Ub8rrvrTnhQ7k4o+YviiY776BQVvnGCv04zcQLcFGUl5gE
38NflNUVyRRBnMRddWQVDf9VMOyGj/8N7yy5Y0b2qvzfvGn9LhJIZJrglfCm7ymP
AbEVtQwdpf5pLGkkeB6zpxxxYu7KyJesF12KwvhHhm4qxFYxldBniYUr+WymXUad
DKqC5JlR3XC321Y9YeRq4VzW9v493kHMB65jUr9TU/Qr6cf9tveCX4XSQRjbgbME
HMUfpIBvFSDJ3gyICh3WZlXi/EjJKSZp4A==
-----END CERTIFICATE-----`
)

// Client type, holds all data we need for a api client
type Client struct {
	url               string
	basicAuthUsername string
	basicAuthPassword string
	errorlog          *log.Logger
	tracelog          *log.Logger
	c                 *http.Client
}

// generic error specific to vmray
type ClientError struct {
	msg string
}

// Error returns a string representation of the error condition
func (self ClientError) Error() string {
	return self.msg
}

// OptionFunc configures a client, used by New
type OptionFunc func(*Client) error

// errorf log to the error log
func (self *Client) errorf(format string, args ...interface{}) {
	if self.errorlog != nil {
		self.errorlog.Printf(format, args...)
	}
}

// tracef logs to the trace log
func (self *Client) tracef(format string, args ...interface{}) {
	if self.tracelog != nil {
		self.tracelog.Printf(format, args...)
	}
}

// New configures a new vmray client.
//
// Example on how to use vmray.New:
//
//  c, err := vmray.New(
//      vmray.SetBasicAuth(os.Getenv("VMRAY_EMAIL"), os.Getenv("VMRAY_PASSWD")),
//  )
//
// Example with custom http client and URL, error logging and request tracing
//
//  c, err := vmray.New(
//      vmray.SetHttpClient(httpclient),
//      vmray.SetUrl("https://vmray.mydomain.com/api/"),
//      vmray.SetErrorLog(log.New(os.Stderr, "vmray error: ", log.Lshortfile)),
//      vmray.SetTraceLog(log.New(os.Stderr, "vmray trace: ", log.Lshortfile)),
//      vmray.SetBasicAuth(os.Getenv("VMRAY_EMAIL"), os.Getenv("VMRAY_PASSWD")),
//  )
//
func New(options ...OptionFunc) (*Client, error) {
	// setup and configure a tls cert pool
	tlsConf := new(tls.Config)
	certPool := x509.NewCertPool()
	ok := certPool.AppendCertsFromPEM([]byte(GlobalSignRootCA))
	if !ok {
		err := fmt.Errorf("failed to parse root CA cert")
		return nil, err
	}
	tlsConf.RootCAs = certPool
	// set up transport for http
	tr := &http.Transport{
		TLSClientConfig: tlsConf,
	}
	// setup http client with transport defined above
	httpclient := &http.Client{Transport: tr}
	// Set up the client
	c := &Client{
		url: "",
		c:   httpclient,
	}
	// run options on it
	for _, option := range options {
		if err := option(c); err != nil {
			return nil, err
		}
	}
	// set default url if no url was set
	if c.url == "" {
		c.url = DefaultURL
	}
	// make sure our url is correct and ends with a slash
	if !strings.HasSuffix(c.url, "/") {
		c.url += "/"
	}
	// print to the trace log what url we are actually using
	c.tracef("Using URL [%s]\n", c.url)

	return c, nil
}

// SetHttpClient can be used to specify the http.Client to use when making
// HTTP requests to vmray.
func SetHttpClient(httpClient *http.Client) OptionFunc {
	return func(self *Client) error {
		if httpClient != nil {
			self.c = httpClient
			self.tracef("HttpClient: %+v\n", httpClient)

		} else {
			// setup and configure a tls cert pool
			tlsConf := new(tls.Config)
			certPool := x509.NewCertPool()
			ok := certPool.AppendCertsFromPEM([]byte(GlobalSignRootCA))
			if !ok {
				err = fmt.Errorf("failed to parse root CA cert")
				return err
			}
			tlsConf.RootCAs = certPool
			// set up transport for http
			tr := &http.Transport{
				TLSClientConfig: tlsConf,
			}
			// setup http client with transport defined above
			httpclient := &http.Client{Transport: tr}
			// Set up the client
			self.c = httpclient
		}
		return nil
	}
}

// SetUrl defines the URL endpoint of vmray
func SetUrl(rawurl string) OptionFunc {
	return func(self *Client) error {
		if rawurl == "" {
			rawurl = DefaultURL
		}
		u, err := url.Parse(rawurl)
		if err != nil {
			self.errorf("Invalid URL [%s] - %v\n", rawurl, err)
			return err
		}
		if u.Scheme != "http" && u.Scheme != "https" {
			msg := fmt.Sprintf("Invalid schema specified [%s]", rawurl)
			self.errorf(msg)
			return ClientError{msg: msg}
		}
		self.url = rawurl
		return nil
	}
}

// Set basic auth
func SetBasicAuth(username, password string) OptionFunc {
	return func(self *Client) error {
		self.basicAuthUsername = username
		self.basicAuthPassword = password
		return nil
	}
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

// SetErrorLog sets the logger for critical messages. It is nil by default.
func SetErrorLog(logger *log.Logger) func(*Client) error {
	return func(c *Client) error {
		c.errorlog = logger
		return nil
	}
}

// SetTraceLog specifies the logger to use for output of trace messages like
// HTTP requests and responses. It is nil by default.
func SetTraceLog(logger *log.Logger) func(*Client) error {
	return func(c *Client) error {
		c.tracelog = logger
		return nil
	}
}

// dumpRequest dumps a request to the debug logger if it was defined
func (self *Client) dumpRequest(req *http.Request) {
	if self.tracelog != nil {
		out, err := httputil.DumpRequestOut(req, true)
		if err == nil {
			self.tracef("%s\n", string(out))
		}
	}
}

// dumpResponse dumps a response to the debug logger if it was defined
func (self *Client) dumpResponse(resp *http.Response) {
	if self.tracelog != nil {
		out, err := httputil.DumpResponse(resp, true)
		if err == nil {
			self.tracef("%s\n", string(out))
		}
	}
}

// Request handling functions

// handleError will handle responses with status code different from 200
func (self *Client) handleError(resp *http.Response) error {
	if resp.StatusCode != http.StatusOK {
		if self.errorlog != nil {
			out, err := httputil.DumpResponse(resp, true)
			if err == nil {
				self.errorf("%s\n", string(out))
			}
		}
		if resp.Body != nil {
			resp.Body.Close()
		}
		msg := fmt.Sprintf("Unexpected status code: %d (%s)", resp.StatusCode, http.StatusText(resp.StatusCode))
		self.errorf(msg)
		return ClientError{msg: msg}
	}
	return nil
}

// makeApiGetRequest fetches a URL with querystring via HTTP GET and
//  returns the response if the status code is HTTP 200
// `parameters` should not include the apikey.
// The caller must call `resp.Body.Close()`.
func (self *Client) makeApiGetRequest(fullurl string, parameters Parameters) (resp *http.Response, err error) {
	values := url.Values{}
	//values.Set("apikey", self.apikey)
	for k, v := range parameters {
		values.Add(k, v)
	}
	// check if fullurl already ends in '?', append '?' if not.
	if !strings.HasSuffix(fullurl, "?") {
		fullurl += "?"
	}
	req, err := http.NewRequest("GET", fullurl+values.Encode(), nil)
	if err != nil {
		return resp, err
	}
	if self.basicAuthUsername != "" {
		req.SetBasicAuth(self.basicAuthUsername, self.basicAuthPassword)
	}
	self.dumpRequest(req)
	resp, err = self.c.Do(req)
	if err != nil {
		return resp, err
	}

	self.dumpResponse(resp)

	if err = self.handleError(resp); err != nil {
		return resp, err
	}

	return resp, nil
}

// makeApiPostRequest fetches a URL with querystring via HTTP POST and
//  returns the response if the status code is HTTP 200
// `parameters` should not include the apikey.
// The caller must call `resp.Body.Close()`.
func (self *Client) makeApiPostRequest(fullurl string, parameters Parameters) (resp *http.Response, err error) {
	bodyReader, bodyWriter := io.Pipe()
	// create a multipat/mime writer
	writer := multipart.NewWriter(bodyWriter)
	// get the Content-Type of our form data
	fdct := writer.FormDataContentType()
	errChan := make(chan error, 1)
	go func() {
		defer bodyWriter.Close()
		for k, v := range parameters {
			if err := writer.WriteField(k, v); err != nil {
				errChan <- err
				return
			}
		}
		errChan <- writer.Close()
	}()

	req, err := http.NewRequest("POST", fullurl, bodyReader)
	if err != nil {
		return resp, err
	}
	req.Header.Add("Content-Type", fdct)
	if self.basicAuthUsername != "" {
		req.SetBasicAuth(self.basicAuthUsername, self.basicAuthPassword)
	}

	self.dumpRequest(req)
	resp, err = self.c.Do(req)
	if err != nil {
		return resp, err
	}

	self.dumpResponse(resp)

	if err = self.handleError(resp); err != nil {
		return resp, err
	}

	return resp, nil
}

// makeApiUploadRequest uploads a file via multipart/mime POST and
//  returns the response if the status code is HTTP 200
// `parameters` should not include the apikey.
// The caller must call `resp.Body.Close()`.
func (self *Client) makeApiUploadRequest(fullurl string, parameters Parameters, paramName, path string) (resp *http.Response, err error) {
	// open the file
	self.tracef("makeApiUploadRequest: opening file '%s'\n", path)
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	// Pipe the file so as not to read it into memory
	bodyReader, bodyWriter := io.Pipe()
	// create a multipat/mime writer
	writer := multipart.NewWriter(bodyWriter)
	// get the Content-Type of our form data
	fdct := writer.FormDataContentType()
	// Read file errors from the channel
	errChan := make(chan error, 1)
	go func() {
		defer bodyWriter.Close()
		defer file.Close()
		part, err := writer.CreateFormFile(paramName, filepath.Base(path))
		if err != nil {
			errChan <- err
			return
		}
		if _, err := io.Copy(part, file); err != nil {
			errChan <- err
			return
		}
		for k, v := range parameters {
			if err := writer.WriteField(k, v); err != nil {
				errChan <- err
				return
			}
		}
		errChan <- writer.Close()
	}()

	// create a HTTP request with our body, that contains our file
	postReq, err := http.NewRequest("POST", fullurl, bodyReader)
	if err != nil {
		return resp, err
	}
	// add the Content-Type we got earlier to the request header.
	postReq.Header.Add("Content-Type", fdct)
	if self.basicAuthUsername != "" {
		postReq.SetBasicAuth(self.basicAuthUsername, self.basicAuthPassword)
	}
	self.dumpRequest(postReq)

	// send our request off, get response and/or error
	resp, err = self.c.Do(postReq)
	if err = <-errChan; err != nil {
		return resp, err
	}
	if err != nil {
		return resp, err
	}

	self.dumpResponse(resp)

	if err = self.handleError(resp); err != nil {
		return resp, err
	}
	// we made it, let's return
	return resp, nil
}

type Parameters map[string]string

// fetchApiJson makes a request to the API and decodes the response.
// `method` is one of "GET", "POST", or "FILE"
// `actionurl` is the final path component that specifies the API call
// `parameters` for request
// `result` is modified as an output parameter.
// 'result' must be a pointer to a vmray JSON structure.
func (self *Client) fetchApiJson(method string, actionurl string, parameters Parameters, result interface{}) (err error) {
	theurl := self.url + actionurl
	var resp *http.Response
	switch method {
	case "GET":
		resp, err = self.makeApiGetRequest(theurl, parameters)
	case "POST":
		resp, err = self.makeApiPostRequest(theurl, parameters)
	case "FILE":
		// get the path to our file from parameters["filename"]
		path := parameters["filename"]
		self.tracef("fetchApiJson FILE '%s'\n", path)
		// call makeApiUploadRequest with fresh/empty Parameters
		newparameters := Parameters{
			"comment":          "",
			"email":            self.basicAuthUsername,
			"password":         self.basicAuthPassword,
			"cmdline":          "",
			"archive_password": "",
			"type":             "api",
			"name":             "sample_file",
		}
		resp, err = self.makeApiUploadRequest(theurl, newparameters, "sample_file", path)
	}
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)
	if err = dec.Decode(result); err != nil {
		return err
	}

	return nil
}

// fetchApiFile makes a get request to the API and returns the file content
func (self *Client) fetchApiFile(actionurl string, parameters Parameters) (data []byte, err error) {
	theurl := self.url + actionurl
	var resp *http.Response
	resp, err = self.makeApiPostRequest(theurl, parameters)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// JobsInfoResult a datastructure to hold the results from GetJobsInfo
type JobInfoResult struct {
	Jobs map[string]JobInfoDetail `json:"jobs"`
}

type JobInfoDetail struct {
	Status string  `json:"status"`
	Slot   float64 `json:"slot"`
}

// GetJobsInfo queries pending and in progress jobs from vmray
func (self *Client) GetJobsInfo() (r *JobInfoResult, err error) {
	r = &JobInfoResult{}
	parameters := Parameters{"email": self.basicAuthUsername, "password": self.basicAuthPassword, "type": "api"}
	err = self.fetchApiJson("POST", "get_jobs_info", parameters, r)
	return r, err
}

// FindSampleResult a datastructure to hold the results from FindSample
type FindSampleResult struct {
	SampleId int `json:"sample_id"`
}

// FindSample finds a Sample in VmRay by its sha1, sha2 or md5 hash
func (self *Client) FindSample(hash string) (r *FindSampleResult, err error) {
	r = &FindSampleResult{}
	parameters := Parameters{"email": self.basicAuthUsername, "password": self.basicAuthPassword, "type": "api", "hash": hash}
	err = self.fetchApiJson("POST", "find_sample", parameters, r)
	return r, err
}

// SampleInfoResult a datastructure to hold the result of GetSampleInfo
type SampleInfoResult struct {
	Filesize  int    `json:"sample_filesize"`
	Priority  int    `json:"sample_priority"`
	Sha1      string `json:"sample_sha1hash"`
	Type      string `json:"sample_type"`
	Filename  string `json:"sample_filename"`
	Md5       string `json:"sample_md5hash"`
	Password  string `json:"sample_password"`
	Shareable bool   `json:"sample_shareable"`
	SampleId  int    `json:"sample_id"`
	Sha2      string `json:"sample_sha256hash"`
	Url       string `json:"sample_url"`
	Created   string `json:"sample_created"`
}

// GetSampleInfo queries Information about a given Sample from VmRay
func (self *Client) GetSampleInfo(id string) (r *SampleInfoResult, err error) {
	r = &SampleInfoResult{}
	parameters := Parameters{"email": self.basicAuthUsername, "password": self.basicAuthPassword, "type": "api", "id": id}
	err = self.fetchApiJson("POST", "get_sample_info", parameters, r)
	return r, err
}

// AnalysisInfoResult is a datastructure to hold results from GetAnalysisInfo
type AnalysisInfoResults struct {
	Analyses map[string]AnalysisInfoDetails `json:"Analyses"`
	Jobs     JobInfoResult                  `json:"jobs"`
}

type AnalysisInfoDetails struct {
	AnalyzerType              string  `json:"analyzer_type"`
	AnalysisSnapshotId        float64 `json:"analysis_snapshot_id"`
	VmhostName                string  `json:"vmhost_name"`
	AnalysisCreated           string  `json:"analysis_created"`
	AnalysisSize              float64 `json:"analysis_size"`
	AnalysisJobStarted        string  `json:"analysis_job_started"`
	SnapshotName              string  `json:"snapshot_name"`
	AnalysisResult            string  `json:"analysis_result"`
	AnalysisJobId             float64 `json:"analysis_job_id"`
	AnalysisCmdlineId         float64 `json:"analysis_cmdline_id"`
	AnalysisConfigurationID   float64 `json:"analysis_configuration_id"`
	AnalysisUserConfigID      float64 `json:"analysis_user_config_id"`
	AnalyzerName              string  `json:"analyzer_name"`
	AnalysisJobruleId         float64 `json:"analysis_jobrule_id"`
	AnalysisPriority          float64 `json:"analysis_priority"`
	Target                    string  `json:"target"`
	AnalysisHint              float64 `json:"analysis_hint"`
	AnalysisAnalyzerVersion   string  `json:"analysis_analyzer_version"`
	ConfigurationName         string  `json:"configuration_name"`
	AnalysisUserId            float64 `json:"analysis_user_id"`
	AnalysisId                float64 `json:"analysis_id"`
	AnalysisExternalReference string  `json:"analysis_external_reference"`
	AnalysisVmhostId          float64 `json:"analysis_vmhost_id"`
	VmName                    string  `json:"vm_name"`
	AnalysisAnalyzerId        float64 `json:"analysis_analyzer_id"`
	AnalysisPrescriptId       float64 `json:"analysis_prescript_id"`
	AnalysisSampleId          float64 `json:"analysis_sample_id"`
	AnalysisServerity         float64 `json:"analysis_severity"`
	AnalysisVmId              float64 `json:"analysis_vm_id"`
}

// GetAnalysisInfo queries Information about an analysis performed by VmRay
func (self *Client) GetAnalysisInfo(id string) (r *AnalysisInfoResults, err error) {
	r = &AnalysisInfoResults{}
	parameters := Parameters{"email": self.basicAuthUsername, "password": self.basicAuthPassword, "type": "api", "sample_id": id}
	err = self.fetchApiJson("POST", "get_analysis_info", parameters, r)
	return r, err
}

// DownloadAnalysis downloads results of an VmRay analysis as zip file.
func (self *Client) DownloadAnalysis(id string) (data []byte, err error) {
	parameters := Parameters{"email": self.basicAuthUsername, "password": self.basicAuthPassword, "type": "api", "analysis_id": id}
	data, err = self.fetchApiFile("download_analysis", parameters)
	return data, err
}

// UploadResult is a datastructure to hold results from UploadSample API call
type UploadResultDetails struct {
	Submission_id   int    `json:"submission_id"`
	Sample_id       int    `json:"sample_id"`
	Webif_url       string `json:"webif_url"`
	Sample_filename string `json:"sample_filename"`
	Sample_url      string `json:"sample_url"`
	Job_ids         []int  `json:"job_ids"`
}

// UploadSample uploads a given file to VmRay and returns the UploadResultDetails and error
func (self *Client) UploadSample(file string) (r *map[string]UploadResultDetails, err error) {
	self.tracef("UploadSample file: '%s'\n", file)
	var parsed map[string]UploadResultDetails
	//r = &UploadResult{}
	parameters := Parameters{
		"email":    self.basicAuthUsername,
		"password": self.basicAuthPassword,
		"type":     "api",
		"filename": file,
	}
	self.tracef("UploadSample parameters: %v", parameters)
	err = self.fetchApiJson("FILE", "upload_sample", parameters, &parsed)
	return &parsed, err
}
