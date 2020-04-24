FROM golang:1

RUN apt update && apt install -y libgl1-mesa-dev build-essential

WORKDIR /source
COPY . .
RUN ./prerequisites.sh
RUN go build . -o out
ENTRYPOINT ["/source/out"]
