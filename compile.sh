set -e
docker exec -w /go/src/github.com/levavakian/rowm rowmc go get
docker exec -w /go/src/github.com/levavakian/rowm rowmc go build 
docker exec -w /go/src/github.com/levavakian/rowm/cmd/rowmbright rowmc go build 