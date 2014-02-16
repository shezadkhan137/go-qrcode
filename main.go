package main

import (
	"flag"
	"fmt"
	"github.com/shezadkhan137/go-qrcode/qrcode"
)

//go build -o test -ldflags "-linkmode external -extldflags -static"
var image *string

func init() {
	image = flag.String("i", "", "image path")
}

func main() {

	flag.Parse()

	results, err := qrcode.GetDataFromPng(*image)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", results)
}
