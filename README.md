VmRay API module for go
=======================

This module provides a client to the API of vmray.com.
_vmray.go_ allows to communicate with the API of VmRay.

VmRay is a 3rd generation malware execution and analysis environment.
For more Information see: http://www.vmray.com/

## Disclaimers

This code is based on the old API of VmRay. 
Since beginning of 2016 or Version 1.9 of VmRay there is a new API which is not yet coverd in this code. See [Issue #1](https://github.com/scusi/vmray/issues/1)

*This code is not final and may be subject to changes.*

## Usage

Go and get the code

```shell
go get github.com/scusi/vmray
```

Within your program just import the library

```go
import("github.com/scusi/vmray")
```

## Documentation

[![GoDoc](https://godoc.org/github.com/scusi/vmray?status.svg)](https://godoc.org/github.com/scusi/vmray)

Documentation is available on [GoDoc](https://godoc.org/github.com/scusi/vmray)

For TLS certificate issues please see [TlsCertReadme.md](https://github.com/scusi/vmray/blob/master/TlsCertReadme.md)

## Examples

See the [Examples directory](https://github.com/scusi/vmray/tree/master/Examples)

## Author

This module has been written by _Florian 'scusi' Walther_.

