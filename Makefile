run:
	go run cmd/*.go

build:
	CGO_ENABLED=1 go build -o bin/main cmd/*.go

test:
	go test ./internal...

testCoverage:
	go test -coverprofile=coverage.out  ./internal...
	go tool cover -html=coverage.out

lint:
	golangci-lint run -v

update-go-packages:
	go get -u -t ./... && go get -u=patch all && go mod tidy
