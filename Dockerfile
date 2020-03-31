FROM golang:1.14

RUN apt update && apt install -y libgl1-mesa-dev build-essential

WORKDIR /source
COPY . .
RUN ./prerequisites.sh
RUN go run .