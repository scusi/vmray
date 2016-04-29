# vmray examples README

## Overview

* DownloadAnalysesByHash - takes a hash of an sample (md5, sha1, sha2) retrieves analyses for that sample
* DownloadAnalysis - takes an analysis id and downloads that analysis
* FindSampleByHash - takes a hash of an sample (md5, sha1, sha2) and provides the vmray sampleId
* GetAnalysisInfo - takes a sampleId and shows information about available analyses
* GetJobsInfo - takes no argument, shows information about open/running jobs
* GetSampleInfo - takes sampleId and provides information about that sample
* UploadSample - takes a filename as argument, uploads that file to vmray and provides some information about the sample uploaded.

## vmray credentials

All example scripts assume that your vmray credentials (email, password) are 
available as environment variables called ```VMRAY_EMAIL``` and ```VMRAY_PASSWD```

So before starting any example you shoud make sure to set these environment variables accordingly.

```shell
 export VMRAY_EMAIL="me@mydomain.tld"
 export VMRAY_PASSWD="myVerySecretPassword"
```

In order to those variables permanetly you can add the two above lines to your ```~/.profile``` file.
