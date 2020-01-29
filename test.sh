set -e
$(pwd)/compile.sh

# 2. Run.
#
# We need to specify the full path to Xephyr, as otherwise xinit will not
# interpret it as an argument specifying the X server to launch and will launch
# the default X server instead.
XEPHYR=$(whereis -b Xephyr | cut -f2 -d' ')
xinit ./xinitrc -- \
    "$XEPHYR" \
        :100 \
        -ac \
        -sw-cursor \
        -screen 800x800+0+0 -screen 800x800+800+0 -screen 800x800+1600+0 \
        +xinerama \
        +extension RANDR

