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
	bw := bufio.NewWriter(w)
	defer bw.Flush()
	nothing := []byte("")

	paths := [][]string{}
	for _, f := range fields {
		paths = append(paths, []string{f})
	}

	out := make([][]byte, numFields)
	for {
		line, err := br.ReadSlice('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		jsonparser.EachKey(line, func(idx int, value []byte, vt jsonparser.ValueType, err error) {
			out[idx] = value
		}, paths...)
		for i = 0; i < numFields-1; i++ {
			bw.Write(out[i])
			out[i] = nothing
			bw.Write([]byte("\t"))
		}
		bw.Write(out[numFields-1])
		out[numFields-1] = nothing
		bw.Write([]byte("\n"))
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
