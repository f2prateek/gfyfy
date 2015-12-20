default: server

server:
	@godep go run cmd/gfycat/main.go

.PHONY: server
