build:
	@go build -o ./bin/app

test:
	@go test -v -cover ./..

run: build
	@./bin/app

clean:
	@rm -rf ./bin
