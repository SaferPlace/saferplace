#syntax
FROM golang:1.17-alpine AS builder

WORKDIR /src
ENV CGO_ENABLED=0

COPY . .

RUN go mod download

ARG TARGETOS
ARG TARGETARCH

RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o /bin/saferplace ./cmd/saferplace

######################
FROM scratch AS target

ENV GIN_MODE=release

COPY --from=builder /bin/saferplace /bin/saferplace

ENTRYPOINT [ "/bin/saferplace" ]