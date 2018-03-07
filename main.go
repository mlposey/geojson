package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	srcDirName := flag.String("in", ".", "the root directory of the unzipped dataset")
	dstFileName := flag.String("out", "geolife_simple.json", "the output filename")
	flag.Parse()

	parser := NewParser(*srcDirName)
	dstFile, err := os.Create(*dstFileName)
	if err != nil {
		log.Fatal(err)
	}
	defer dstFile.Close()

	w := bufio.NewWriter(dstFile)	
	parser.Parse(func (path string, trajectory Trajectory) {
		if b, err := json.Marshal(trajectory); err == nil {
			fmt.Fprintln(w, string(b))
		}
	})
	w.Flush()
}
