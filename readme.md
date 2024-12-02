


install first:
```
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

run below cmd

```
make clean
make gen
make grpc-server
```

open another terminal and run below cmd
```
make grpc-client
```


