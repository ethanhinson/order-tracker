FROM golang
WORKDIR /app
COPY . /app
COPY go.mod .
COPY go.sum .
RUN go mod vendor
RUN go build -o dist/tracker main.go
RUN chmod u+x /app/dist/tracker
ENV PATH="/app/dist:${PATH}"