#! /bin/bash


name=txmMma.exe
echo "###########" compile [ $name ] try "###########"

go build -o ./$name main.go

echo "compile end..."