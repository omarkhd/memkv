ARG SRCDIR=/go/omarkhd/memkv

# Building stage.
FROM golang:1.16.5 AS build
ARG SRCDIR

WORKDIR ${SRCDIR}
ADD cmd cmd
ADD metrics metrics
ADD server server
ADD store store
ADD go.mod go.sum ./

RUN go build -o memkv cmd/memkv/main.go

# Runtime stage.
FROM golang:1.16.5
ARG SRCDIR

WORKDIR /opt/memkv
COPY --from=build ${SRCDIR}/memkv .
