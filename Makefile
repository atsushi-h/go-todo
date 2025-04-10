run:
	find . -name '*.go' | entr -r go run cmd/main.go