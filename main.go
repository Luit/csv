package main

import (
	"encoding/csv"
	"flag"
	"io"
	"log"
	"os"
	"strconv"
)

// Usage:
//   csv [options] cut <field>...
//   csv [options] (prefix|match|regex) (<field> <filter>)...
//
// Options:
//   -d ','  Delimiter character.
//   -q '"'  Quote character.
//   -c      Comment character.
//   -f      Input file. (default: stdin)
//   -o      Output file. (default: stdout)

var (
	// TODO: use flag.Value interface to populate Var of type rune
	d = flag.String("d", ",", "Delimiter character")
	// q = flag.String("q", "\"", "Quote character")
	// c = flag.String("c", "", "Comment character")

	// f = flag.String("f", "", "Input file")
	// o = flag.String("o", "", "Output file")
)

func cut(input []string, fields []int) []string {
	var output []string

	inLen := len(input)
	for _, i := range fields {
		if i >= inLen {
			continue // TODO: test index out of range
		}
		output = append(output, input[i])
	}

	return output
}

type filter struct {
	field int
}

func main() {
	flag.Parse()

	var command func([]string) []string

	if flag.NArg() < 1 {
		log.Fatal("TODO: usage string here")
	}
	args := flag.Args()
	cmd, args := args[0], args[1:]
	switch cmd {
	case "cut":
		if len(args) < 1 {
			log.Fatal("TODO: usage string for subcommand cut here")
		}

		// Turn args into ints
		fields := make([]int, len(args))
		for n, arg := range args {
			u, err := strconv.ParseUint(arg, 10, 0)
			if err != nil {
				log.Fatal("TODO: usage string for subcommand cut here")
			}
			fields[n] = int(u)
		}

		command = func(input []string) []string {
			return cut(input, fields)
		}

	case "prefix":
		if len(args) < 2 || len(args)%2 > 0 {
			log.Fatal("TODO: usage string for subcommand prefix here")
		}
	case "match":
		if len(args) < 2 || len(args)%2 > 0 {
			log.Fatal("TODO: usage string for subcommand match here")
		}
	case "regex":
		if len(args) < 2 || len(args)%2 > 0 {
			log.Fatal("TODO: usage string for subcommand regex here")
		}

		// turn every second arg into regexp.Regexp

	default:
		log.Fatal("TODO: usage string here")
	}

	r := csv.NewReader(os.Stdin)
	var line int
	w := csv.NewWriter(os.Stdout)

	for {
		line += 1
		input, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error reading line %d: %s", line, err.Error())
		}

		output := command(input)
		if len(output) > 0 {
			err := w.Write(output)
			if err != nil {
				log.Fatal("Error writing output:", err)
			}
		}
	}
	w.Flush()
}
