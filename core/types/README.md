# Types

Types and interfaces are defined with [protobuf](https://developers.google.com/protocol-buffers/) (version 3).

## Prerequisites

* Make sure [Protobuf](http://brewformulas.org/Protobuf) is installed (version 3.5 or later).

  > Install with `brew install protobuf`
  
  > Verify with `protoc --version`

* Make sure [Go workspace bin](https://stackoverflow.com/questions/42965673/cant-run-go-bin-in-terminal) is in your path.
  
  > Install with ``export PATH=$PATH:`go env GOPATH`/bin``
  
  > Verify with `echo $PATH`

* Make sure [Go protobuf plugin](https://developers.google.com/protocol-buffers/docs/gotutorial) is installed.
  
  > Install with `go get -u github.com/golang/protobuf/protoc-gen-go`
  
  > Verify with `protoc-gen-go` (it will wait for stdin, just close it)

* Make sure [Square goprotowrap](https://github.com/square/goprotowrap) is installed.
  
  > Install with `go get -u github.com/square/goprotowrap/cmd/protowrap`
  
  > Verify with `protowrap --version`

## Build

* Compile all `.proto` files to `.go` files from this directory:
```
./generate_go_types.sh
```

* **NOTE:** Commit generated `.go` files to git.

