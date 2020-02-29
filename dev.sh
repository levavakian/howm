#!/bin/bash
docker run -itd --network=host --name rowmc -v $(pwd)/../rowm:/go/src/github.com/levavakian/rowm/ golang:latest bash
