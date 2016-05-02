VmRay API module for go
=======================

_vmray.go_ allows to communicate with the API of VmRay.

VmRay is a 3rd generation malware execution and analysis environment.
For more Information see [VmRay Website](http://www.vmray.com/)

## Disclaimers

This code is based on the old API of VmRay. 
Since beginning of 2016 or version 1.9 of VmRay there is a new API which is not yet coverd in this code. See [Issue #1](https://github.com/scusi/vmray/issues/1)

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

Please see the [Examples directory](https://github.com/scusi/vmray/tree/master/Examples)

## Commits

If you want to commit to this code feel free to send me pull requests.
I prefer lots of small commits that do change one thing rather than 
one huge commit with a dozen of changes hard to follow.

## Author

This module has been written by _Florian 'scusi' Walther_.

