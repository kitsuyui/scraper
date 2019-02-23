FROM gliderlabs/alpine:3.9@sha256:23b993692b943f0799b3f36042d8a1331557103eb4ac2c0b8ab36cab9f399f8b
ENV \
SCRAPER_VERSION='v0.1.0' \
SCRAPER_HASH='b46b2ff9cbfbff8c8f6772cc115362085d8ca454e34683da9971b88f71672fda'
RUN \
apk --update add --no-cache --virtual dep-download \
ca-certificates curl && \
update-ca-certificates && \
curl -fsSL "https://github.com/kitsuyui/scraper/releases/download/$SCRAPER_VERSION/scraper_linux_amd64" \
> /usr/bin/scraper && \
echo "${SCRAPER_HASH}  /usr/bin/scraper" | sha256sum -s -c - && \
chmod +x /usr/bin/scraper && \
apk del dep-download && \
mkdir /lib64 && \
ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2 && \
mkdir -p /root/server-config
ENTRYPOINT ["scraper", "server"]
CMD ["-d", "/root/server-config", "-p", "8080", "-H", "0.0.0.0"]
EXPOSE 8080
