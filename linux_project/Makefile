# Makefile for Linux + Docker + Go + React deployment

BROKER_BINARY=brokerApp
AUTH_BINARY=authApp
LOGGER_BINARY=loggerServiceApp
MAIL_BINARY=mailerApp
LISTENER_BINARY=listenerApp

up:
	@echo "Starting all Docker containers..."
	docker-compose up -d
	@echo "Docker containers started!"

up_build: build_broker build_logger build_auth build_mail build_listener
	@echo "Stopping Docker containers (if running)..."
	docker-compose down
	@echo "Building and starting Docker containers..."
	docker-compose up --build -d
	@echo "Docker containers built and started!"

down:
	@echo "Stopping Docker containers..."
	docker-compose down
	@echo "Docker containers stopped!"

# --- Go service builds (Linux binaries) ---
build_broker:
	@echo "Building Broker service..."
	cd ../broker-service && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ${BROKER_BINARY} ./cmd/api
	@echo "Broker build done."

build_listener:
	@echo "Building Listener service..."
	cd ../listener-service && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ${LISTENER_BINARY} .
	@echo "Listener build done."

build_mail:
	@echo "Building Mail service..."
	cd ../mail-service && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ${MAIL_BINARY} ./cmd/api
	@echo "Mail build done."

build_logger:
	@echo "Building Logger service..."
	cd ../logger-service && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ${LOGGER_BINARY} ./cmd/api
	@echo "Logger build done."

build_auth:
	@echo "Building Auth service..."
	cd ../authentication-service && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ${AUTH_BINARY} ./cmd/api
	@echo "Auth build done."

# --- Clean build artifacts ---
clean:
	@echo "Cleaning Go binaries..."
	rm -f ../broker-service/${BROKER_BINARY} \
	      ../listener-service/${LISTENER_BINARY} \
	      ../mail-service/${MAIL_BINARY} \
	      ../logger-service/${LOGGER_BINARY} \
	      ../authentication-service/${AUTH_BINARY}
	@echo "Cleaning frontend build..."
	rm -rf ${FRONT_END_BUILD_DIR}
	@echo "Clean done."
