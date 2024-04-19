
.PHONY: dev
dev:
	air & pnpm dev

.PHONY: fmt
fmt:
	go fmt ./...
