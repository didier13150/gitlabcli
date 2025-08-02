FROM alpine AS build

RUN apk add --no-cache git make musl-dev go

ENV GOROOT=/usr/lib/go
ENV GOPATH=/go
ENV PATH=/go/bin:$PATH

RUN mkdir build

COPY main.go build/main.go
COPY env.go build/env.go
COPY var.go build/var.go
COPY project.go build/project.go
COPY go.mod build/go.mod
COPY go.sum build/go.sum
RUN cd build && go install && go build -o glvars

#-----------------------------------------------------------------------------

FROM alpine

LABEL org.opencontainers.image.authors="Didier FABERT <didier.fabert@gmail.com>"
LABEL eu.tartarefr.glvars.version=1.0.3

RUN mkdir -p /usr/local/share/glvars
COPY --from=build build/glvars /usr/local/bin/glvars
COPY LICENSE /usr/local/share/glvars/LICENSE
COPY README.fr.md /usr/local/share/glvars/README.fr.md
COPY README.md /usr/local/share/glvars/README.md
RUN chmod 0755 \
  /usr/local/bin/glvars

ENTRYPOINT [ "/usr/local/bin/glvars" ]
