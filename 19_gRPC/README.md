## Install
`sudo apt install protobuf-compiler`

`go install google.golang.org/protobuf/cmd/protoc-gen-go@latest`

`go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest`

## Generate Go Code
`protoc --go_out=. --go-grpc_out=. {name_file}.proto`
