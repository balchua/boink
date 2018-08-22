FROM golang:1.10

# Download and install the latest release of dep
ADD https://github.com/golang/dep/releases/download/v0.5.0/dep-linux-amd64 /usr/bin/dep
RUN chmod +x /usr/bin/dep

# Copy the code from the host and compile it
WORKDIR $GOPATH/src/boink
COPY Gopkg.toml Gopkg.lock ./
RUN dep ensure --vendor-only

COPY main.go ./

COPY handler/* handler/

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix nocgo -o /boink .
