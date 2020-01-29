set -e
$(pwd)/compile.sh
rm -rf /usr/bin/howm
cp howm /usr/bin/howm