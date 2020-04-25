COMPOSE=docker-compose
BUILDFILE=build.yml
DOCKER=docker

#Docker
build: Dockerfile
	$(COMPOSE) -f $(BUILDFILE) build
push:
	$(COMPOSE) -f $(BUILDFILE) push
up:
	$(COMPOSE) -f $(BUILDFILE) up
up-d:
	$(COMPOSE) -f $(BUILDFILE) up -d
down:
	$(COMPOSE) -f $(BUILDFILE) down

run-local:
	go run .

build-local:
	go build .

debug:
	go build -gcflags="-N -l" .
	gdb ./gogles