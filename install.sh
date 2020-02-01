set -e
$(pwd)/compile.sh
rm -rf /usr/bin/howm
rm -rf /usr/bin/howmbright.sh
mkdir -p /usr/share/fonts/truetype/dejavu
cp $(pwd)/howm /usr/bin/howm
cp $(pwd)/cmd/howmbright/howmbright /usr/bin/howmbright
cp $(pwd)/resources/DejaVuSans.ttf /usr/share/fonts/truetype/dejavu/DejaVuSans.ttf
chmod u+s /usr/bin/howmbright