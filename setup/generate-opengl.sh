#!/bin/sh

GLOWDIR="$PWD/../glow"

# Get and build GLFW
go get -u -tags=gl github.com/go-gl/glfw/v3.3/glfw # OpenGL Mode
#go get -u -tags=gles1 github.com/go-gl/glfw/v3.3/glfw # OpenGL ES Mode

# Generate OpenGL bindings
go get github.com/go-gl/glow
go install github.com/go-gl/glow

cd ./opengl && glow generate -api=gl -version=1.1 -out $GLOWDIR/gl # OpenGL Bindings
#glow generate -api=gles1 -version=1.0 -out $GLOWDIR/gles #  OpenGL ES Bindings
cd ..

# Fix for MACOS
mkdir -p $GLOWDIR/gl/KHR/ && cp ./opengl/lib/khrplatform.h $GLOWDIR/gl/KHR/khrplatform.h #Fix for MacOS.
mkdir -p $GLOWDIR/gles/KHR/ && cp ./opengl/lib/khrplatform.h $GLOWDIR/gles/KHR/khrplatform.h #Fix for MacOS.