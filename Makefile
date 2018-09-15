all: prepare 
	go build -o ${GOPATH}/bin/gomeboycolor

prepare:
	go mod tidy
