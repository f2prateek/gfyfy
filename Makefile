default: server

server:
	@godep go run cmd/gfyfy/main.go

.PHONY: server
