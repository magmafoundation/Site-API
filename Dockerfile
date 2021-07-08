FROM golang:alpine AS build
USER root
RUN apk add build-base
WORKDIR /src
COPY . /src
RUN go build




FROM golang:alpine

WORKDIR /app
COPY --from=build /src/MagmaAPI /app/
COPY docs /app/docs
EXPOSE 3000

CMD ["./MagmaAPI"]



