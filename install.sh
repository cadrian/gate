#!/bin/sh

export ROOTDIR=$(dirname $(readlink -f $0))

mkdir -p $DESTDIR/usr/bin
mkdir -p $DESTDIR/usr/lib/gate
mkdir -p $DESTDIR/etc/gate

for bin in console menu server; do
    cp $ROOTDIR/bin/$bin $DESTDIR/usr/lib/gate/
    chmod +x $DESTDIR/usr/lib/gate/$bin
done

for script in gate_console gate_menu; do
    sed 's|^dist=.*$|exe=$(dirname $(readlink -f $0))/../lib/gate|;s| \$prop$| \$1|g' < $ROOTDIR/script/$script > $DESTDIR/usr/bin/$script
    chmod +x $DESTDIR/usr/bin/$script
done

for rc in $ROOTDIR/conf/*.rc; do
    cp $rc $DESTDIR/etc/gate/
done
