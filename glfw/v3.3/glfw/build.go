package glfw

/*
// Windows Build Tags
// ----------------
// GLFW Options:
#cgo windows CFLAGS: -D_GLFW_WIN32 -Iglfw/deps/mingw

// Linker Options:
#cgo windows LDFLAGS: -lopengl32 -lgdi32


// Darwin Build Tags
// ----------------
// GLFW Options:
#cgo darwin CFLAGS: -D_GLFW_COCOA -Wno-deprecated-declarations

// Linker Options:
#cgo darwin LDFLAGS: -framework Cocoa -framework OpenGL -framework IOKit -framework CoreVideo


// Linux Build Tags
// ----------------
// GLFW Options:
#cgo linux CFLAGS: -D_GLFW_X11 -D_GNU_SOURCE

// Linker Options:
#cgo linux LDFLAGS: -lGL -lX11 -lXrandr -lXxf86vm -lXi -lXcursor -lm -lXinerama -ldl -lrt

// FreeBSD Build Tags
// ----------------
// GLFW Options:
#cgo freebsd pkg-config: glfw3
#cgo freebsd CFLAGS: -D_GLFW_HAS_DLOPEN
#cgo freebsd,!wayland CFLAGS: -D_GLFW_X11 -D_GLFW_HAS_GLXGETPROCADDRESSARB
#cgo freebsd,wayland CFLAGS: -D_GLFW_WAYLAND

// Linker Options:
#cgo freebsd,!wayland LDFLAGS: -lm -lGL -lX11 -lXrandr -lXxf86vm -lXi -lXcursor -lXinerama
#cgo freebsd,wayland LDFLAGS: -lm -lGL -lwayland-client -lwayland-cursor -lwayland-egl -lxkbcommon
*/
import "C"
