#!/bin/sh

# Get and build GLFW
go get -u -tags=gles1 github.com/go-gl/glfw/v3.3/glfw # OpenGL Mode
#go get -u -tags=gles1 github.com/go-gl/glfw/v3.3/glfw # OpenGL ES Mode

# Generate OpenGL bindings
go get github.com/go-gl/glow
go install github.com/go-gl/glow

glow generate -api=gl -version=1.1 -out ../glow/gl # OpenGL Bindings
#glow generate -api=gles1 -version=1.1 -out ../glow/gles #  OpenGL ES Bindings

mkdir -p ../glow/gl/KHR/ && cp ./lib/khrplatform.h ../glow/gl/KHR/khrplatform.h #Fix for MacOS.
mkdir -p ../glow/gles/KHR/ && cp ./lib/khrplatform.h ../glow/gles/KHR/khrplatform.h #Fix for MacOS.
