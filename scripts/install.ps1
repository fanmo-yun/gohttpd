$buildDir = "./build/gohttpd-windows"
if (-Not (Test-Path -Path $buildDir)) {
    New-Item -ItemType Directory -Force -Path $buildDir
}

Write-Host "Running 'go mod tidy'..."
go mod tidy
if ($?) {
    Write-Host "'go mod tidy' completed successfully."
} else {
    Write-Host "Error: 'go mod tidy' failed."
    exit 1
}

Write-Host "Building the project..."
go build -o "$buildDir/gohttpd.exe" ./cmd/main.go
if ($?) {
    Write-Host "Build completed successfully."
} else {
    Write-Host "Error: Build failed."
    exit 1
}

$resources = @(
    @{Path = "./conf"; Dest = "$buildDir\conf"},
    @{Path = "./html"; Dest = "$buildDir\html"},
    @{Path = "./banner.txt"; Dest = "$buildDir\banner.txt"}
)

foreach ($resource in $resources) {
    if (Test-Path -Path $resource.Path) {
        if (Test-Path -Path $resource.Path -PathType Container) {
            Write-Host "Copying directory $($resource.Path)..."
            Copy-Item -Recurse -Force -Path $resource.Path -Destination $resource.Dest
        } else {
            Write-Host "Copying file $($resource.Path)..."
            Copy-Item -Force -Path $resource.Path -Destination $resource.Dest
        }
    } else {
        Write-Host "Warning: $($resource.Path) not found. Skipping..."
    }
}

Write-Host "Build and installation completed successfully!"