

go build -i -v -o ../server ../cmd/server
cd ..
./server --grpc-port 8080