run:
  go run ./cmd/hisame

test:
  go test ./...

race:
  go test -race ./...