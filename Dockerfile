ARG GOLANG_VERS=1.25.5

ARG ALPINE_VERS=3.23

FROM golang:${GOLANG_VERS}-alpine${ALPINE_VERS}

WORKDIR /code

COPY . .

RUN go mod download && go build -o /usr/local/bin/backup && \
    apk --no-cache add binutils && strip -vs /usr/local/bin/backup

FROM alpine:${ALPINE_VERS}

RUN apk --no-cache add wget ca-certificates libc6-compat \
    libstdc++ \
    usbutils && \
    # 下载Android Platform Tools（仅提取adb二进制）
    wget -q https://dl.google.com/android/repository/platform-tools-latest-linux.zip -O /tmp/adb.zip && \
    # 解压仅保留adb（删除其他冗余文件）
    unzip -q /tmp/adb.zip platform-tools/adb platform-tools/lib64/libc++.so -d /opt/ && \
    # 赋予执行权限
    chmod +x /opt/platform-tools/adb && \
    # 清理临时文件（极致减小镜像体积）
    rm -rf /tmp/* && \
    # 卸载不需要的依赖（wget）
    apk del wget

ENV PATH=$PATH:/opt/platform-tools

ENV GIN_MODE=release

RUN mkdir -m 0750 /root/.android

COPY files/adbkey /root/.android/adbkey
COPY files/adbkey.pub /root/.android/adbkey.pub

COPY --from=0 /usr/local/bin/backup /usr/local/bin/backup
ENTRYPOINT ["/usr/local/bin/backup"]