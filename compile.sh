set -e
docker exec -w /go/src/rowm rowmc go get
docker exec -w /go/src/rowm rowmc go build 
docker exec -w /go/src/rowm/cmd/rowmbright rowmc go build 