.PHONY: templ-generate
templ-generate:
	templ generate

.PHONY: templ-watch
templ-watch:
	templ generate --watch

.PHONY: build
build:
	@templ generate api/web
	@go build -o bin/sportspazz cmd/main.go

.PHONY: run
run: build
	@./bin/sportspazz

.PHONY: clean
clean:
	@echo "Cleaning up generated files..."
	rm -rf bin/
	rm -rf pkg/
	rm -rf vendor/
	rm -f .air .air.pid
	rm -rf tmp/
	rm -f api/web/templates/*_templ.go
	rm -f api/web/templates/*_templ.txt
	@echo "Cleanup complete"
