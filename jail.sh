#!/usr/bin/env bash

mkdir -p /tmpserver/app/bin /tmpserver/app/etc /tmpserver/app/lib
ldd /bin/bash /bin/ls
cp /bin/bash /bin/ls /tmpserver/app/bin
cp -v /lib/aarch64-linux-gnu/libc.so.6 /lib/aarch64-linux-gnu/libdl.so.2 /lib/aarch64-linux-gnu/libtinfo.so.6 /lib/ld-linux-aarch64.so.1 /lib/aarch64-linux-gnu/libselinux.so.1 /lib/aarch64-linux-gnu/libpthread.so.0 /usr/lib/aarch64-linux-gnu/libpcre2-8.so.0 /tmpserver/app/lib
cp ./port /tmpserver/app
echo "PS1='JAIL $ ' " | tee /tmpserver/app/bash.bashrc
chroot /tmpserver/app /bin/bash