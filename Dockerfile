FROM golang:1.18-alpine AS build
WORKDIR /root/project
COPY libp2p-private-relay ./
RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=cache,target=/root/.cache/go-build \
    go build -v -o /bin/app

FROM alpine
WORKDIR /bin
COPY --from=build /bin/app ./
COPY config.json ./
CMD [ "./app" ]
