#!/bin/bash
if [ -d "build" ]
then 
    go build -o build/goping cmd/*
else
    echo "creating a build direcotry..."
    mkdir build
    go build -o build/goping cmd/*
fi
