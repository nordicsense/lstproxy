package main

import (
	"log"
	"os"

	"github.com/nordicsense/modis"
	"github.com/nordicsense/modis/dataset"
)

func main() {
	if err := convert(os.Args[1], os.Args[2]); err != nil {
		log.Fatal(err)
	}
}

func convert(src string, tgt string) error {
	r, err := dataset.Open(src)
	if err != nil {
		return err
	}
	defer r.Close()
	w, err := dataset.New(tgt, dataset.GTiff, r.ImageParams())
	if err != nil {
		return err
	}

	defer w.Close()
	box := modis.Box{0, 0, r.ImageParams().XSize(), r.ImageParams().YSize()}
	data, err := r.ReadBlock(0, 0, box)
	if err != nil {
		return err
	}
	return w.WriteBlock(0, 0, box, data)
}
