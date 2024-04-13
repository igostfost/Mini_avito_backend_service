FROM golang:latest

RUN go version
ENV GOPATH=/

COPY ./ ./

##RUN apt-get update
#RUN apt-get -y install postgresql-client
#
#RUN chmod +x postgres-whait.sh


RUN go mod download
RUN go build -o avito_backend_trainee ./cmd/main.go

CMD ["./avito_backend_trainee"]