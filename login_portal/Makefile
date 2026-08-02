generate:
	@templ generate

build: generate
	@go build -o ./bin/app

build-run: build
	@./bin/app

watch:
	@templ generate --watch --proxy="http://localhost:8080" --cmd="go run ."

tidy:
	@go mod tidy

test:
	@go test ./... -v -cover

clean:
	@rm -rf ./bin
