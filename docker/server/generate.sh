#!/usr/bin/env bash
version="$(curl https://api.github.com/repos/kitsuyui/scraper/releases/latest | jq -r .tag_name)"
homepage='https://github.com/kitsuyui/scraper'

gethash() {
  curl -fsSL "${homepage}/releases/download/${version}/$1" 2>/dev/null \
  | shasum -a 256 \
  | awk '{print $1}'
}

sha256_amd64=$(gethash scraper_linux_amd64)

cd "${0%/*}"
cat <<EOF > Dockerfile
FROM gliderlabs/alpine:3.6
ENV \\
SCRAPER_VERSION='${version}' \\
SCRAPER_HASH='${sha256_amd64}'
RUN \\
apk --update add --no-cache --virtual dep-download \\
ca-certificates curl && \\
update-ca-certificates && \\
curl -fsSL "https://github.com/kitsuyui/scraper/releases/download/\$SCRAPER_VERSION/scraper_linux_amd64" \\
> /usr/bin/scraper && \\
echo "\${SCRAPER_HASH}  /usr/bin/scraper" | sha256sum -s -c - && \\
chmod +x /usr/bin/scraper && \\
apk del dep-download && \\
mkdir /lib64 && \\
ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2 && \\
mkdir -p /root/server-config
ENTRYPOINT ["scraper", "server"]
CMD ["-d", "/root/server-config", "-p", "8080", "-H", "0.0.0.0"]
EXPOSE 8080
EOF
