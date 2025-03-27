FROM golang:1.24-alpine AS build

WORKDIR /app

ARG bin_to_build

COPY go.mod .

RUN go mod download

COPY . .

RUN go build -o server cmd/${bin_to_build}/main.go

FROM scratch
COPY --from=build /app/server /server
EXPOSE 8080
CMD [ "./server" ]
