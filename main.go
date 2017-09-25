package main

import (
	"bufio"
	"flag"
	"io"
	"os"

	"github.com/buger/jsonparser"
)

func cut(r io.Reader, w io.Writer, fields []string) error {
	var i int
	numFields := len(fields)
	br := bufio.NewReader(r)
	nothing := []byte("")

	paths := [][]string{}
	for _, f := range fields {
		paths = append(paths, []string{f})
	}

	out := make([][]byte, numFields)
	for {
		line, err := br.ReadBytes('\n')
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		jsonparser.EachKey(line, func(idx int, value []byte, vt jsonparser.ValueType, err error) {
			out[idx] = value
		}, paths...)
		for i = 0; i < numFields-1; i++ {
			w.Write(out[i])
			out[i] = nothing
			w.Write([]byte("\t"))
		}
		w.Write(out[numFields-1])
		out[numFields-1] = nothing
		w.Write([]byte("\n"))
	}
	return nil
}

func main() {
	flag.Parse()
	fields := flag.Args()
	err := cut(os.Stdin, os.Stdout, fields)
	if err != nil {
		panic(err)
	}
}
