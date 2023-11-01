FROM golang:1.21-bullseye as build
LABEL "responsible"="Port Inc"
RUN mkdir /app
WORKDIR /app
#ARG
#ENV
COPY . .
RUN go build .
EXPOSE 8090
ENTRYPOINT ["sh", "jail.sh"]
#CMD

# chroot jail demo
#set up root at /serverid/app
#mkdir -p /tmpserver/app/{bin,etc}
#ldd /bin/bash /bin/ls
#cp /bin/bash /bin/ls /tmpserver/app/bin
#cp -v /lib/aarch64-linux-gnu/libc.so.6 /lib/aarch64-linux-gnu/libtinfo.so.6 /lib/ld-linux-aarch64.so.1 /lib/aarch64-linux-gnu/libselinux.so.1 /lib/aarch64-linux-gnu/libpcre2-8.so.0 /tmpserver/app/lib
#cp ./port /tmpserver/app
#echo "PS1='JAIL $ ' " | sudo tee /tmpserver/app/bash.bashrc
#chroot /tmpserver/app /bin/bash
##create bin
##copy over ls and bash commands and their dependent modules
##copy over server binary and env dir/files if present
##initialize the new dir as a root
##use the new root
#
#cat <<EOT >> jail.sh
#mkdir -p /tmpserver/app/{bin,etc}
#ldd /bin/bash /bin/ls
#cp /bin/bash /bin/ls /tmpserver/app/bin
#cp -v /lib/aarch64-linux-gnu/libc.so.6 /lib/aarch64-linux-gnu/libtinfo.so.6 /lib/ld-linux-aarch64.so.1 /lib/aarch64-linux-gnu/libselinux.so.1 /lib/aarch64-linux-gnu/libpcre2-8.so.0 /tmpserver/app/lib
#cp ./port /tmpserver/app
#echo "PS1='JAIL $ ' " | tee /tmpserver/app/bash.bashrc
#chroot /tmpserver/app /bin/bash
#EOT

#FROM gcr.io/distroless/base-debian11:nonroot
#EXPOSE 8090/tcp
#COPY --chown=nonroot:nonroot --from=build /app/port /
#ENTRYPOINT ["/port"]
