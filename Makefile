SERVER_ADDR=localhost:8889

INSTALL_DEPS=go get github.com/codegangsta/negroni && \
go get github.com/gorilla/mux && \
go get gopkg.in/tylerb/graceful.v1 && \
go get github.com/Sirupsen/logrus && \
go get github.com/PuerkitoBio/goquery && \
go get github.com/mitchellh/gox && \
go get github.com/boltdb/bolt

init:
	@echo "Init development environment..." && \
	make deps

deps:
	@echo "Install dependencies..."
	@$(INSTALL_DEPS)

run:
	go run cmd/horo/main.go -addr=$(SERVER_ADDR) -debug=true

build:
	cd cmd/horo && gox -osarch=linux/amd64 -output=../../bin/horo
	cd cmd/horo_up && gox -osarch=linux/amd64 -output=../../bin/horo_up