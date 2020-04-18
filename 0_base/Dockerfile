FROM golang:1.13.5-alpine AS build

WORKDIR /go/src/sample-go-server
COPY ./app /go/src/sample-go-server
# `CGO_ENABLED=0` is required to run on scratch image.
# see: https://qiita.com/katoken-0215/items/f3a502fe0c2044709012
RUN CGO_ENABLED=0 go build -o /bin/sample-go-server

FROM scratch
COPY --from=build /bin/sample-go-server /bin/sample-go-server
EXPOSE 8080
ENTRYPOINT ["/bin/sample-go-server"]
