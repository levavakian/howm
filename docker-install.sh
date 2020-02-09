set -e

STARTEDDOWN="false"
if [ ! "$(docker ps -q -f name=rowmc)" ]; then
    echo "No container detected, starting own"
    STARTEDDOWN="true"
    if [ "$(docker ps -aq -f status=exited -f name=rowmc)" ]; then
        docker rm -f rowmc
    fi
    # run your container
    docker run -itd --network=host --name rowmc -v $(pwd)/../rowm:/go/src/rowm/ golang:latest bash
fi

echo "Compiling..."
$(pwd)/compile.sh
echo "Installing to global directories"
rm -rf /usr/bin/rowm
rm -rf /usr/bin/rowmbright.sh
rm -rf /usr/share/xsessions/rowm.desktop
mkdir -p  /usr/local/share/wingo/
cp $(pwd)/rowm /usr/bin/rowm
cp $(pwd)/cmd/rowmbright/rowmbright /usr/bin/rowmbright
cp $(pwd)/resources/dejavu/DejaVuSans.ttf /usr/local/share/wingo/DejaVuSans.ttf
cp $(pwd)/resources/nofont/write-your-password-with-this-font.ttf  /usr/local/share/wingo/write-your-password-with-this-font.ttf
cp $(pwd)/resources/rowm.desktop /usr/share/xsessions/rowm.desktop
chmod u+s /usr/bin/rowmbright

if [ "$STARTEDDOWN" = "true" ]; then
  echo "Removing container rowmc"
  docker rm -f rowmc
fi