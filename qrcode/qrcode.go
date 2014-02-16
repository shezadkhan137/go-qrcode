package qrcode

// #cgo LDFLAGS: -lzbar -lpng -ljpeg -lz -lrt -lm -pthread
// #include <stdio.h>
// #include <stdlib.h>
// #include <png.h>
// #include <zbar.h>
// #include "get_data.h"
// typedef void (*zbar_image_set_data_callback)(zbar_image_t *  image);
import "C"
import (
	"errors"
	"fmt"
	"unsafe"
)

type Result struct {
	symbolType string
	data       string
}

func GetDataFromPng(pngPath string) (results []Result, err error) {

	pth := C.CString(pngPath)
	scanner := C.zbar_image_scanner_create()
	C.zbar_image_scanner_set_config(scanner, 0, C.ZBAR_CFG_ENABLE, 1)

	defer C.zbar_image_scanner_destroy(scanner)

	var width, height C.int = 0, 0
	var raw unsafe.Pointer = nil
	errorCode := C.get_data(pth, &width, &height, &raw)
	if int(errorCode) != 0 {
		err = errors.New(fmt.Sprintf("Error reading from png file. Error code %d", errorCode))
		return
	}

	//defer C.free(raw)

	image := C.zbar_image_create()

	defer C.zbar_image_destroy(image)

	C.zbar_image_set_format(image, C.ulong(808466521))
	C.zbar_image_set_size(image, C.uint(width), C.uint(height))

	f := C.zbar_image_set_data_callback(C.zbar_image_free_data)
	C.zbar_image_set_data(image, raw, C.ulong(width*height), f)

	C.zbar_scan_image(scanner, image)

	symbol := C.zbar_image_first_symbol(image)

	for ; symbol != nil; symbol = C.zbar_symbol_next(symbol) {
		typ := C.zbar_symbol_get_type(symbol)
		data := C.zbar_symbol_get_data(symbol)
		symbolType := C.GoString(C.zbar_get_symbol_name(typ))
		dataString := C.GoString(data)
		results = append(results, Result{symbolType, dataString})
	}

	return
}
