FROM golang:1.22-bookworm AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/tracker .

FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /

COPY --from=build /out/tracker /tracker

USER nonroot:nonroot

ENTRYPOINT ["/tracker"]
