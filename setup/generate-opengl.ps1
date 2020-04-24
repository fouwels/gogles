$PROJDIR="$PWD/.."
$GLOWDIR="$PROJDIR/glow"
$glversion="2.1"

if ((Get-ComputerInfo -Property "os*").OsType -eq "WINNT") {
    $glapi="gl"
    $glfwapi="gl"
} else {
    $glapi="gles"
    $glfwapi="gles2"
}

Write-Output "Getting and building GLFW for $glfwapi" 
& "go" "get -u -tags=$glfwapi github.com/go-gl/glfw/v3.3/glfw".Split(" ")

Write-Output "Installing glow to generate OpenGL Bindings"
& "go" "get github.com/go-gl/glow".Split(" ")
& "go" "install github.com/go-gl/glow".Split(" ")

Write-Output "Removing old glow folder"
Remove-Item -Recurse $GLOWDIR
Write-Output "Generating OpenGL Bindings for $glapi version $glversion"
& "glow" "generate -api=$glapi -xml=$PWD/opengl/xml -version=$glversion -out $GLOWDIR/gl".Split(" ")

Write-Output "Fix for MACOS"
New-Item -Path "$GLOWDIR/gl/" -Name "KHR" -ItemType "directory" | Out-Null
New-Item -Path "$GLOWDIR/gles/" -Name "KHR" -ItemType "directory" | Out-Null
Copy-Item "$PWD/opengl/lib/khrplatform.h" -Destination "$GLOWDIR/gl/KHR/khrplatform.h"
Copy-Item "$PWD/opengl/lib/khrplatform.h" -Destination "$GLOWDIR/gles/KHR/khrplatform.h"

# This is horrible, never do this. Required as gl and gles are generated in different packages names.
Write-Output "Rewriting go imports for $glapi"
$configFiles = Get-ChildItem $PROJDIR *.go -rec
foreach ($file in $configFiles)
{
    if ($glapi -eq "gl"){
        (Get-Content $file.PSPath) | Foreach-Object { $_ -replace 'gl "github.com/kaelanfouwels/gogles/glow/gles"', '"github.com/kaelanfouwels/gogles/glow/gl"' } |  Set-Content $file.PSPath
    } else {
        (Get-Content $file.PSPath) | Foreach-Object { $_ -replace '"github.com/kaelanfouwels/gogles/glow/gl"', 'gl "github.com/kaelanfouwels/gogles/glow/gles"' } |  Set-Content $file.PSPath
    }
}