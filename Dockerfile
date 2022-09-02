#syntax
FROM golang:alpine AS builder

WORKDIR /src
ENV CGO_ENABLED=0

RUN apk add -U --no-cache ca-certificates

COPY . .

RUN apk --update add ca-certificates

RUN go mod download

ARG TARGETOS
ARG TARGETARCH

RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o /bin/saferplace ./cmd/saferplace

######################
FROM scratch AS target

ENV PORT=8080
ENV GIN_MODE=release

EXPOSE ${PORT}

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /bin/saferplace /bin/saferplace
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENTRYPOINT [ "/bin/saferplace" ]
