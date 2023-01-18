FROM golang:1.18-alpine as dev

ENV ROOT=/go/src/app
ENV CGO_ENABLED 0
WORKDIR ${ROOT}

RUN apk update && apk add git
COPY go.mod go.sum ./
RUN go mod download
EXPOSE ${PORT}

ENTRYPOINT ["go", "run", "."]


FROM golang:1.18-alpine as builder

ENV ROOT=/go/src/app
WORKDIR ${ROOT}

RUN apk update && apk add git
COPY go.mod go.sum ./
RUN go mod download

COPY . ${ROOT}
RUN CGO_ENABLED=0 GOOS=linux go build -o $ROOT/binary


FROM alpine as prod

ENV ROOT=/go/src/app
WORKDIR ${ROOT}
COPY --from=builder ${ROOT}/binary ${ROOT}

EXPOSE ${PORT}
ENTRYPOINT ["/go/src/app/binary"]