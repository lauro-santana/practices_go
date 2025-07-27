### Setup Consul web service discovery

```cmd
docker run -d \
  -p 8500:8500 \
  -p 8600:8600/udp \
  --name=dev-consul \
  hashicorp/consul agent \
    -server \
    -ui \
    -node=server-1 \
    -bootstrap-expect=1 \
    -client=0.0.0.0
```

```cmd
docker start dev-consul
```

#### Consul web UI


http://localhost:8500/


### Setup services

```cmd
go run metadata/cmd/main.go
``` 

```cmd
go run movie/cmd/main.go
```

```cmd
go run rating/cmd/main.go
```


### Generate gRPC Code by proto files

```cmd
protoc -I=api --go_out=. --go-grpc_out=. api/movie.proto
```

### Call a grpc service

```cmd
grpcurl -plaintext -d '{"record_id":"1", "record_type":
"movie"}' localhost:8082 RatingService/GetAggregatedRating
```