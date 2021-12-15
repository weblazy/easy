package csvx

import (
	"bufio"
	"encoding/csv"
	"os"
	"strings"
)

type CSV struct {
	path         string
	wfile        *os.File
	rfile        *os.File
	w            *csv.Writer
	r            *bufio.Reader
	rowSeparator string
}

// NewCSV return a CSV
func NewCSV(path string, rowSeparator rune, lineSeparator string) (*CSV, error) {
	wfile, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	rfile, err := os.OpenFile(path, os.O_RDONLY, 0666)
	if err != nil {
		return nil, err
	}
	w := csv.NewWriter(wfile)
	w.Comma = rowSeparator
	if lineSeparator == `\r\n` {
		w.UseCRLF = true
	}
	r := bufio.NewReader(rfile)
	return &CSV{
		path:         path,
		wfile:        wfile,
		rfile:        rfile,
		w:            w,
		r:            r,
		rowSeparator: string(rowSeparator),
	}, nil
}

// Write truncate and write one line
func (this *CSV) Write(str []string) error {
	err := this.wfile.Truncate(0)
	if err != nil {
		return err
	}
	err = this.w.Write(str)
	if err != nil {
		return err
	}
	this.w.Flush()
	return nil
}

// Truncate
func (this *CSV) Truncate() error {
	return this.wfile.Truncate(0)
}

// Append append one line
func (this *CSV) Append(str []string) error {
	err := this.w.Write(str)
	if err != nil {
		return err
	}
	this.w.Flush()
	return nil
}

// Reset
func (this *CSV) Reset() (int64, error) {
	return this.rfile.Seek(0, 0)
}

// ReadLine read one line
func (this *CSV) ReadLine() ([]string, error) {
	line, _, err := this.r.ReadLine() //以'\n'为结束符读入一行
	return strings.Split(string(line), this.rowSeparator), err
}

// Close close file
func (this *CSV) Close() error {
	err := this.wfile.Close()
	if err != nil {
		return err
	}
	return this.rfile.Close()
}
