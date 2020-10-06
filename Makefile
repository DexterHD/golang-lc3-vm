
build:
	go build .

lint:
	golangci-lint run

test:
	go test -v -race -cover ./...

ci-coverage-dependencies:
	go get github.com/axw/gocov/...
	go get github.com/AlekSi/gocov-xml
	go mod tidy

ci-coverage-report: ci-coverage-dependencies
	go test -race -covermode=atomic -coverprofile=coverage.txt ./...
	gocov convert coverage.txt | gocov-xml > coverage.xml

clean:
	rm -f ./coverage.txt
	rm -f ./coverage.xml
	rm -f ./golang-lc3-vm
