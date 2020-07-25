#!/bin/sh
set -e

SCRIPT=$(readlink -f "$0")
DIR=$(dirname "$SCRIPT")

STARTEDDOWN="false"
if [ ! "$(docker ps -q -f name=rowmc)" ]; then
    echo "No container detected, starting own"
    STARTEDDOWN="true"
    if [ "$(docker ps -aq -f status=exited -f name=rowmc)" ]; then
        docker rm -f rowmc
    fi
    # run your container
    docker run -itd --network=host --name rowmc -v $DIR/../rowm:/go/src/github.com/levavakian/rowm/ golang:latest bash
fi

echo "Compiling..."
$DIR/compile.sh
echo "Installing to global directories"
rm /usr/bin/rowmbinaryfinder | true
rm /usr/bin/rowmbright | true
rm /usr/bin/rowm | true
rm /usr/share/xsessions/rowm.desktop | true
mkdir -p /usr/local/share/wingo/
mkdir -p /usr/share/xsessions
cp $DIR/rowm /usr/bin/rowm
cp $DIR/cmd/rowmbright/rowmbright /usr/bin/rowmbright
cp $DIR/cmd/rowmbinaryfinder /usr/bin/rowmbinaryfinder
cp $DIR/resources/dejavu/DejaVuSans.ttf /usr/local/share/wingo/DejaVuSans.ttf
cp $DIR/resources/nofont/write-your-password-with-this-font.ttf  /usr/local/share/wingo/write-your-password-with-this-font.ttf
cp $DIR/resources/rowm.desktop /usr/share/xsessions/rowm.desktop
chmod u+s /usr/bin/rowmbright

if [ "$STARTEDDOWN" = "true" ]; then
  echo "Removing container rowmc"
  docker rm -f rowmc
fi
