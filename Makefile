.PHONY: ejecutar probar cobertura formato revisar limpiar

ejecutar:
	go run ./cmd/butaca

probar:
	go test ./...

cobertura:
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out | tail -1

formato:
	gofmt -l -w .

revisar:
	go vet ./...

limpiar:
	rm -f coverage.out coverage.html
