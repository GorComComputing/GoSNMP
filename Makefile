.PHONY: main
main: *.go deps
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o gosnmp .


.PHONY:deps
deps:
	go get github.com/gosnmp/gosnmp
	

