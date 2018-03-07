package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// ParseFunc handles the trajectories created from files.
type ParseFunc func(filepath string, trajectory Trajectory)

// Parser describes a type that parses data into trajectories.
type Parser interface {
	Parse(ParseFunc)
}

// PLTParser parses .plt files into trajectories.
type PLTParser struct {
	rootDir      string
	fn           ParseFunc
}

// Parse parses all .plt files in all subdirectories.
func (p *PLTParser) Parse(fn ParseFunc) {
	p.fn = fn
	filepath.Walk(p.rootDir, func(path string, info os.FileInfo, err error) error {
		split := strings.Split(info.Name(), ".")
		if len(split) == 2 && split[1] == "plt" {
			p.process(path, split[0])
		}
		return nil
	})
}

func (p *PLTParser) process(path, simpleFileName string) {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("could not open file", path)
		return
	}
	defer file.Close()
	p.parseFile(file, simpleFileName)
}

func (p *PLTParser) parseFile(file *os.File, simpleFileName string) {
	reader := newSimplePltReader(file)
	trajectory := p.buildTrajectory(reader)
	trajectory.ID = simpleFileName
	p.fn(file.Name(), trajectory)
}

func (*PLTParser) buildTrajectory(r *simplePltReader) (trajectory Trajectory) {
	for r.HasNext() {
		lat, lng, alt := r.ReadLine()
		trajectory.Path = append(trajectory.Path, Point{
			Lat: lat,
			Lng: lng,
			Alt: alt,
		})
	}
	return
}

// simplePltReader reads a subset of data from each plt file.
type simplePltReader struct {
	reader      *csv.Reader
	currentLine []string
}

// HasNext returns true if the reader can read another line.
func (spr *simplePltReader) HasNext() bool {
	var err error
	spr.currentLine, err = spr.reader.Read()
	return err == nil
}

// ReadLine returns the (latitude, longitude, altitude) of a line.
func (spr *simplePltReader) ReadLine() (float64, float64, float64) {
	var cols [4]float64 // cols[2] is always 0 and not needed
	for i := range cols {
		cols[i], _ = strconv.ParseFloat(spr.currentLine[i], 64)
	}
	return cols[0], cols[1], cols[3]
}

func newSimplePltReader(handle *os.File) *simplePltReader {
	br := bufio.NewReader(handle)
	// The first six lines contain generic information.
	for i := 0; i < 6; i++ {
		br.ReadBytes('\n')
	}
	return &simplePltReader{reader: csv.NewReader(br)}
}

// NewParser creates a Parser that operates on a given directory.
func NewParser(dir string) Parser {
	return &PLTParser{rootDir: dir}
}
