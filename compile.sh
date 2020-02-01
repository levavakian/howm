set -e
docker exec -w /go/src/howm howmc go get
docker exec -w /go/src/howm howmc go build 
docker exec -w /go/src/howm/cmd/howmbright howmc go build 