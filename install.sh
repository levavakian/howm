set -e

STARTEDDOWN="false"
if [ ! "$(docker ps -q -f name=howmc)" ]; then
    echo "No container detected, starting own"
    STARTEDDOWN="true"
    if [ "$(docker ps -aq -f status=exited -f name=howmc)" ]; then
        docker rm -f howmc
    fi
    # run your container
    docker run -itd --network=host --name howmc -v $(pwd)/../howm:/go/src/howm/ golang:latest bash
fi

echo "Compiling..."
$(pwd)/compile.sh
echo "Installing to global directories"
rm -rf /usr/bin/howm
rm -rf /usr/bin/howmbright.sh
rm -rf /usr/share/xsessions/howm.desktop
mkdir -p  /usr/local/share/wingo/
cp $(pwd)/howm /usr/bin/howm
cp $(pwd)/cmd/howmbright/howmbright /usr/bin/howmbright
cp $(pwd)/resources/dejavu/DejaVuSans.ttf  /usr/local/share/wingo/DejaVuSans.ttf
cp $(pwd)/resources/nofont/write-your-password-with-this-font.ttf  /usr/local/share/wingo/write-your-password-with-this-font.ttf
cp $(pwd)/resources/howm.desktop /usr/share/xsessions/howm.desktop
chmod u+s /usr/bin/howmbright

if [ "$STARTEDDOWN" = "true" ]; then
  echo "Removing container howmc"
  docker rm -f howmc
fi