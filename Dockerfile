# Building the binary of the App
FROM fedora:35 AS build

#RUN dnf install golang -y

RUN curl -LO https://dl.google.com/go/go1.17.5.linux-amd64.tar.gz

RUN tar -C /usr/local -xzf go1.17.2.linux-amd64.tar.gz

RUN export GOROOT=/usr/local/go

RUN mkdir -p /go

RUN export GOPATH=/go

RUN export PATH=$GOPATH/bin:$GOROOT/bin:$PATH

# install task file
RUN 'sh -c "$(curl --location https://taskfile.dev/install.sh)" -- -d'

# install libinjection

WORKDIR /src/go/ackycdn-node

COPY . .

RUN task setup_fedora

RUN go mod download

RUN go mod tidy

RUN rm -rf ./dist/

ENV CGO_ENABLED=0

ENV GOOS=linux

ENV GOARCH=amd64

RUN go build -a -installsuffix cgo -o ./dist/ackycdn-node_linux_arm64.bin ackycdn.go

FROM fedora:35

WORKDIR /ackycdn

RUN mkdir ./conf

COPY ./conf ./conf

COPY --from=build /src/go/ackycdn-node/dist/ackycdn-node_linux_arm64.bin .

EXPOSE [80,443]

CMD ["./ackycdn-node_linux_arm64.bin"]




