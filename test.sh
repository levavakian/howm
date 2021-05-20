#!/bin/sh
set -e
SCRIPTPATH="$( cd "$(dirname "$0")" >/dev/null 2>&1 ; pwd -P )"
$SCRIPTPATH/compile.sh

# 2. Run.
#
# We need to specify the full path to Xephyr, as otherwise xinit will not
# interpret it as an argument specifying the X server to launch and will launch
# the default X server instead.
XEPHYR=$(whereis -b Xephyr | cut -f2 -d' ')
cd $SCRIPTPATH
xinit $SCRIPTPATH/xinitrc -- \
    "$XEPHYR" \
        :100 \
        -ac \
        -sw-cursor \
        -screen 1500x800 \
        +xinerama \
        +extension RANDR
cd -
