FROM alpine

RUN mkdir /lib64 && \
    ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2 && \
    apk add curl

COPY bin/linux/server /server

EXPOSE 10086

CMD ["/server"]