TLS Certificate README
======================

the _vmray_ go module, by default does check the certificate of the remote 
host. Therefor it needs to have at the root CA certificate to be able to 
validate the certificate presented by the remote server.

In order to archive this, the client does load the the root CA certifacte 
from a constant called _GlobalSignRootCA_. This certificate is added to
certifcate pool used internally to validate remote hosts certificates.

How to install a different root CA
==================================

In case you want to use the _vmray_ module to connect to your own instance
useing a different certificate as cloud.vmray.com the most simple approach 
is to replace the content of _GlobalSignRootCA_ constant with the - PEM 
encoded - cert of your root CA.

If you want to add your certificates along with the default one see section 
_How to add more root CAs_.

How to skip certificate checking - INSECURE
===========================================

Another - INSECURE - approach would be to disable certificate verification 
completly. This can be archived by passing your own http client to the 
NewClient function and turn off cert checking for your client by useing the 
TLS config switch _InsecureSkipVerify_.

You can archive this by useing code like the below one in your go programm.

```go
    // define a new http transport, activate 'InsecureSkipVerify'
 	tr := &http.Transport{
 		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
 	}

    // set up a new http client that uses the afore defined transport
    httpclient := &http.Client{Transport: tr}
 	
	// initialize a new vmray client and pass it your http client
	// Now the vmray client uses your http client for connections. 
	c, err := vmray.New(
 		vmray.SetUrl("https://myvmray.mydomain.tld/api/"),
 		vmray.SetBasicAuth(os.Getenv("VMRAY_EMAIL"), os.Getenv("VMRAY_PASSWD")),
 		vmray.SetHttpClient(httpclient),
 	)
```

BTW: In a similar fashion you can for example make sure to use a certain proxy 
to connect to your vmray instance.

How to add more root CAs
=======================

Imagen you have more than one instance of vmray and they use certificates 
signed by different root CAs. Here is what you can do in your go program.

```go
	// create a new TLS config
	tlsConf := new(tls.Config)
	// create a new pool of (CA) certificates
	certPool := x509.NewCertPool()

	// have alist of your files (assumed to be in the local directory)
	myRootCerts := ["rootCert1.pem", "rootCert2.pem", "rootCert3.pem"]

	// for each of your root CA certs do
	for _, fileName := range myRootCerts {
			// setup and configure a tls cert pool
			rootPEM, err := ioutil.ReadFile(fileName)
			if err != nil {
					return nil, err
			}
			ok := certPool.AppendCertsFromPEM([]byte(rootPEM))
			if !ok {
					err = fmt.Errorf("failed to parse root CA cert")
					return nil, err
			}
	}
	// assign the pool of certs as trusted root CAs to your tls config
	tlsConf.RootCAs = certPool

	// set up transport for http, useing the tls config created
	tr := &http.Transport{
			TLSClientConfig: tlsConf,
	}

	// setup http client with transport defined above
	httpclient := &http.Client{Transport: tr}

	// initialize a new vmray client and pass it your http client
	// Now the vmray client uses your http client for connections. 
	c, err := vmray.New(
 		vmray.SetUrl("https://myvmray.mydomain.tld/api/"),
 		vmray.SetBasicAuth(os.Getenv("VMRAY_EMAIL"), os.Getenv("VMRAY_PASSWD")),
 		vmray.SetHttpClient(httpclient),
 	)
```

