package main

import (
	"encoding/csv"
	"flag"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// Usage:
//   csv [options] cut <field>...
//   csv [options] (prefix|match|regex) (<field> <filter>)...
//
// Options:
//   -d ','  Delimiter character.
//   -q '"'  Quote character.
//   -c      Comment character.
//   -a      Match all, instead of any (AND instead of OR) doesn't apply to cut
//   -s      Return string, instead of csv. (only applies to cut with 1 field)
//   -f      Input file. (default: stdin)
//   -o      Output file. (default: stdout)

var (
	// TODO: use flag.Value interface to populate Var of type rune
	delim = flag.String("d", ",", "Delimiter character")
	// q = flag.String("q", "\"", "Quote character")
	// c = flag.String("c", "", "Comment character")

	single = flag.Bool("s", false, "Return single field as string instead of csv")
	all    = flag.Bool("a", false, "Apply all filters (AND, instead of OR)")

	// f = flag.String("f", "", "Input file")
	// o = flag.String("o", "", "Output file")
)

func cut(input []string, fields []int) []string {
	var output []string

	for _, i := range fields {
		if i >= len(input) {
			continue // TODO: test index out of range
		}
		output = append(output, input[i])
	}

	return output
}

func and(input []string, filters []filter) []string {
	for _, filter := range filters {
		if filter.field >= len(input) {
			return nil
		}
		if !filter.match(input[filter.field]) {
			return nil
		}
	}
	return input
}

func or(input []string, filters []filter) []string {
	for _, filter := range filters {
		if filter.field >= len(input) {
			continue
		}
		if filter.match(input[filter.field]) {
			return input
		}
	}
	return nil
}

type filter struct {
	field int
	match func(string) bool
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
	case "c", "cut":
		if len(args) < 1 {
			log.Fatal("TODO: usage string for subcommand cut here")
		}

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
	case "p", "prefix":
		if len(args) < 2 || len(args)%2 > 0 {
			log.Fatal("TODO: usage string for subcommand prefix here")
		}

		filters := make([]filter, len(args)/2)
		for n := range filters {
			u, err := strconv.ParseUint(args[n*2], 10, 0)
			if err != nil {
				log.Fatal("TODO: usage string for subcommand prefix here")
			}

			field, prefix := int(u), args[n*2+1]

			filters[n] = filter{
				field: field,
				match: func(input string) bool {
					return strings.HasPrefix(input, prefix)
				},
			}
		}
		if *all {
			command = func(input []string) []string {
				return and(input, filters)
			}
		} else {
			command = func(input []string) []string {
				return or(input, filters)
			}
		}
	case "m", "match":
		if len(args) < 2 || len(args)%2 > 0 {
			log.Fatal("TODO: usage string for subcommand match here")
		}

		filters := make([]filter, len(args)/2)
		for n := range filters {
			u, err := strconv.ParseUint(args[n*2], 10, 0)
			if err != nil {
				log.Fatal("TODO: usage string for subcommand match here")
			}

			field, match := int(u), args[n*2+1]

			filters[n] = filter{
				field: field,
				match: func(input string) bool {
					return input == match
				},
			}
		}
		if *all {
			command = func(input []string) []string {
				return and(input, filters)
			}
		} else {
			command = func(input []string) []string {
				return or(input, filters)
			}
		}
	case "r", "re", "regex", "regexp":
		if len(args) < 2 || len(args)%2 > 0 {
			log.Fatal("TODO: usage string for subcommand regex here")
		}

		filters := make([]filter, len(args)/2)
		for n := range filters {
			u, err := strconv.ParseUint(args[n*2], 10, 0)
			if err != nil {
				log.Fatal("TODO: usage string for subcommand regexp here")
			}

			field, expr := int(u), args[n*2+1]

			re, err := regexp.Compile(expr)
			if err != nil {
				log.Fatal("TODO: usage string for subcommand regex here (bad regex)")
			}

			filters[n] = filter{
				field: field,
				match: func(input string) bool {
					return re.Match([]byte(input))
				},
			}
		}
		if *all {
			command = func(input []string) []string {
				return and(input, filters)
			}
		} else {
			command = func(input []string) []string {
				return or(input, filters)
			}
		}
	default:
		log.Fatal("TODO: usage string here")
	}

	r := csv.NewReader(os.Stdin)
	r.Comma = rune((*delim)[0]) // TODO: should be rune
	var line int
	w := csv.NewWriter(os.Stdout)
	w.Comma = rune((*delim)[0]) // TODO: should be rune

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
