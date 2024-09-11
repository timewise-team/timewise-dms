FROM golang:1.22.5

WORKDIR /app

ARG GITHUB_TOKEN

RUN git config --global url."https://${GITHUB_TOKEN}@github.com/".insteadOf "https://github.com/"

ENV GOPRIVATE=github.com/timewise-team/timewise-models/models
ENV GONOSUMDB=github.com/timewise-team/timewise-models/models

COPY go.mod go.sum /app/
RUN go mod download

COPY . /app/

RUN CGO_ENABLED=0 GOOS=linux go build -o timewise-dms .

EXPOSE 8080

CMD ["./timewise-dms"]