COMPOSE=docker-compose
BUILDFILE=build.yml
DOCKER=docker

#Docker
build: Dockerfile install
	$(COMPOSE) -f $(BUILDFILE) build
push:
	$(COMPOSE) -f $(BUILDFILE) push

run-local:
	go run .

build-local:
	go build .

debug:
	go build -gcflags="-N -l" .
	gdb ./gogles