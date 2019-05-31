#!/bin/bash

package=beacon
platforms=("windows/amd64" "linux/amd64" "darwin/amd64")
ldflags="-s -w"
version="v1.0.0"

for platform in "${platforms[@]}"
do
    platform_split=(${platform//\// })
    GOOS=${platform_split[0]}
    GOARCH=${platform_split[1]}
    folder_name=$package'-'$GOOS'-'$GOARCH-$version
    output_name=$package
    if [ $GOOS = "windows" ]; then
        output_name+='.exe'
    fi  

    env GOOS=$GOOS GOARCH=$GOARCH go build -ldflags "$ldflags" -o bin/$folder_name/$output_name $package
    if [ $? -ne 0 ]; then
        echo 'An error has occurred! Aborting the script execution...'
        exit 1
    fi

    cp -f {LICENSE,README.md} bin/$folder_name/
    tar -czf bin/$folder_name.tgz -C bin $folder_name
done
