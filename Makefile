.PHONY: all dev frontend server

all: dev

dev:
	@echo "Starting frontend and Go server..."
	@trap 'kill 0' EXIT; \
	$(MAKE) frontend & \
	$(MAKE) server & \
	wait

ui:
	npm run dev

sv:
	cd server && go run .

build:
	npm run build
