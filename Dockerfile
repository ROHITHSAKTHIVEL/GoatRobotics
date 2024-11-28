FROM golang:alpine

WORKDIR /GoatRobotics

RUN apk update && apk add --no-cache 

COPY go.mod go.sum ./

RUN go mod tidy

COPY . .

RUN go build -o goatrobotics   .

# COPY   goatrobotics .

COPY config.json .

ENTRYPOINT [ "./goatrobotics" ]