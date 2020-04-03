COMPOSE=docker-compose
BUILDFILE=build.yml
DOCKER=docker

init:
	cd ./setup
	./generate-opengl.sh
	#./install-glfw.sh
	
#Docker
build: Dockerfile install
	$(COMPOSE) -f $(BUILDFILE) build
push:
	$(COMPOSE) -f $(BUILDFILE) push

run-local:
	go run .

build-local:
	go build .