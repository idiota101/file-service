


protoc --proto_path=../api/proto/v1 --go_out=plugins=grpc:../pkg/api/v1 file-service.proto
