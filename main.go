package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

type fieldMapping map[string]int
type extractor struct {
	fields    []string
	fm        fieldMapping
	numFields int
	outBuffer []string
}

func extractKey(line string) (string, string) {
	for idx, ch := range line {
		if ch == '"' {
			return line[:idx], line[idx+1:]
		}
	}
	panic("wtf")
}
func skipToValue(line string) string {
	previous := '_'
	inString := false
	for idx, ch := range line {
		if ch == '"' && previous != '\\' {
			inString = !inString
		}
		if ch == ':' && !inString {
			line = line[idx+1:]
			break
		}
		previous = ch
	}
	for idx, ch := range line {
		if ch == ' ' {
			line = line[idx+1:]
		} else {
			return line
		}
	}
	panic("wtf")
}

//Capture until the next comma or } that is not in a string
func captureValue(line string) (string, string) {
	inString := false
	sawString := false
	previous := '_'
	for idx, ch := range line {
		if ch == '"' && previous != '\\' {
			inString = !inString
			sawString = true
		}
		if (ch == ',' || ch == '}') && !inString {
			if !sawString {
				return line[:idx], line[idx+1:]
			} else {
				return line[1 : idx-1], line[idx+1:]
			}
		}
		previous = ch
	}
	panic("wtf")
}

func (e *extractor) Extract(line string, w io.Writer) error {
	if line[0] != '{' {
		return fmt.Errorf("Invalid line: %s", line)
	}
	for i := 0; i < e.numFields; i++ {
		e.outBuffer[i] = ""
	}

	var key string
	var value string
	depth := 0
	wantkey := false
	nextOutindex := 0

	capture := false
	remaining := e.numFields

	for len(line) > 0 && remaining > 0 {
		ch := line[0]
		if ch == '{' {
			depth++
			line = line[1:]
		}
		if ch == '}' {
			line = line[1:]
		}
		if ch == ' ' || ch == '\n' || ch == ',' || ch == '[' || ch == ']' {
			line = line[1:]
		}
		if ch == '"' {
			key, line = extractKey(line[1:])
			log.Printf("Key was: %q", key)
			line = skipToValue(line)
			if nextOutindex, wantkey = e.fm[key]; wantkey {
				log.Printf("Want this key for %d", nextOutindex)
				value, line = captureValue(line)
				log.Printf("Got value %q", value)
				e.outBuffer[nextOutindex] = value
				remaining--
			} else {
				log.Printf("Don't want this key")
				_, line = captureValue(line)
			}
		}
		log.Printf("line is %s", line)
		_ = capture
	}
	for i := 0; i < e.numFields; i++ {
		if i > 0 {
			io.WriteString(w, "\t")
		}
		io.WriteString(w, e.outBuffer[i])
	}
	io.WriteString(w, "\n")
	return nil

}

func cut(r io.Reader, w io.Writer, fields []string) error {

	fm := make(fieldMapping)
	for idx, f := range fields {
		fm[f] = idx
	}
	e := &extractor{
		fm:        fm,
		numFields: len(fields),
		outBuffer: make([]string, len(fields)),
	}

	br := bufio.NewReader(r)
	for {
		line, err := br.ReadString('\n')
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		err = e.Extract(line, w)
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {
	flag.Parse()
	fields := flag.Args()
	err := cut(os.Stdin, os.Stdout, fields)
	if err != nil {
		log.Fatal(err)
	}
}
