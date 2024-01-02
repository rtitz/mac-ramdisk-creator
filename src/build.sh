#!/bin/bash

platforms=( "darwin/arm64" "darwin/amd64" )

cd $( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
package_name=$(cd .. && basename $(pwd) && cd - >/dev/null 2>&1)
#version=$(git tag | tail -n1)
output_directory="../bin/"

mkdir -p $output_directory >/dev/null 2>&1

echo "Downloading required modules..."
go get -u && go mod tidy

for platform in "${platforms[@]}"
do
    platform_split=(${platform//\// })
    GOOS=${platform_split[0]}
    GOARCH=${platform_split[1]}
    
    if [ $GOOS = "darwin" ]; then
        #output_name=$package_name'-'$version'_'$GOOS'-'$GOARCH
        output_name=$package_name'_macos-'$GOARCH
    else
        #output_name=$package_name'-'$version'_'$GOOS'-'$GOARCH
        output_name=$package_name'_'$GOOS'-'$GOARCH
    fi

    if [ $GOOS = "windows" ]; then
        output_name+='.exe'
    fi

    echo "Building $GOOS/$GOARCH output: $output_name"

    env GOOS=$GOOS GOARCH=$GOARCH go build -ldflags "-s -w" -o $output_name $package
    if [ $? -ne 0 ]; then
           echo 'An error has occurred! Aborting the script execution...'
        exit 1
    fi
    mv $output_name $output_directory
done
