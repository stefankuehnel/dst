FROM golang:1 AS build-stage

ARG VERSION

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY ./ ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /dst -ldflags="-X 'main.Version=$VERSION'" -trimpath .

FROM gcr.io/distroless/base-debian11 AS release-stage

WORKDIR /mnt/local
VOLUME /mnt/local

COPY --from=build-stage /dst /dst

ENTRYPOINT ["/dst"]