#!/bin/sh
cd glfw-3.3.2
mkdir -p build
cd build

#cmake -DBUILD_SHARED_LIBS=ON -DGLFW_USE_OSMESA=TRUE ../
cmake -DBUILD_SHARED_LIBS=ON ../
make 
make install