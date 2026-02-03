.PHONY: all
all: realtime

.PHONY: review-ui
review-ui:
	@cd packages/review-ui && pnpm build:dev

.PHONY: realtime
realtime: review-ui
	@go run ./cmd/saferplace
