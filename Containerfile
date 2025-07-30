FROM alpine AS build

RUN apk add --no-cache git make musl-dev go

ENV GOROOT=/usr/lib/go
ENV GOPATH=/go
ENV PATH=/go/bin:$PATH

RUN mkdir build

COPY main.go build/main.go
COPY env.go build/env.go
COPY var.go build/var.go
COPY go.mod build/go.mod
RUN cd build && go build -o glvars

#-----------------------------------------------------------------------------

FROM alpine

LABEL org.opencontainers.image.authors="Didier FABERT <didier.fabert@gmail.com>"
LABEL eu.tartarefr.glvars.version=1.0.2

COPY --from=build build/glvars /usr/local/bin/glvars
RUN chmod 0755 \
  /usr/local/bin/glvars

ENTRYPOINT [ "/usr/local/bin/glvars" ]
