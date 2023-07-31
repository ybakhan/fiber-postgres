build:
	@go build -o bin/fiber-postgres  ./main 

test:
	@go test ./... -v -count=1 -tags=unit