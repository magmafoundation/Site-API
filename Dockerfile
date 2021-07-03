FROM golang:alpine AS build

WORKDIR /src
COPY . /src
RUN go build




FROM golang:alpine

WORKDIR /app
COPY --from=build /src/MagmaAPI /app/
EXPOSE 3000

CMD ["./MagmaAPI"]



