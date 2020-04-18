FROM golang:1.13.5-alpine AS build

# In order to debug remotely with GoLand, the paths of the source code relative to GOPATH must match.
WORKDIR ${GOPATH}/src/github.com/bellwood4486/sample-go-containerized-debug/1_dlv
COPY . .
# `CGO_ENABLED=0` is required to run on scratch image.
# see: https://qiita.com/katoken-0215/items/f3a502fe0c2044709012
RUN CGO_ENABLED=0 go build -o /bin/sample-go-server ./app
# To debug with Delve, disable optimization.
# see: https://github.com/derekparker/delve/blob/master/Documentation/usage/dlv_exec.md
RUN CGO_ENABLED=0 go build -gcflags "all=-N -l" -o /bin/sample-go-server-debug ./app

FROM build AS debug
RUN set -ex && \
  apk update && \
  apk add --no-cache git && \
  go get -v -u github.com/go-delve/delve/cmd/dlv
# Port 8080 belongs to our application, 2345 belongs to Delve
EXPOSE 8080 2345
# Run Delve with the uninterrupted setting in the main function.
# see: https://github.com/derekparker/delve/blob/master/Documentation/faq.md#how-do-i-use-delve-with-docker
CMD /go/bin/dlv exec /bin/sample-go-server-debug \
  --listen=:2345 \
  --headless \
  --api-version=2 \
  --continue \
  --accept-multiclient

FROM scratch
COPY --from=build /bin/sample-go-server /bin/sample-go-server
EXPOSE 8080
ENTRYPOINT ["/bin/sample-go-server"]
