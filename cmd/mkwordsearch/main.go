package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/gilramir/argparse/v2"
)

type ProgramOptions struct {
	Verbose          bool
	NumCols          int
	NumRows          int
	WordListFilename string
}

type Program struct {
	ProgramOptions
}

func (s *Program) runCli() {

	opts := &s.ProgramOptions
	ap := argparse.New(&argparse.Command{
		Description: "Development testing tool golang-WoSeCon",
		Values:      opts,
	})

	ap.Add(&argparse.Argument{
		Switches: []string{"-v", "--verbose"},
		Help:     "Set verbose mode",
	})

	ap.Add(&argparse.Argument{
		Name: "num_cols",
		Help: "The number of columns in the word search puzzle",
	})

	ap.Add(&argparse.Argument{
		Name: "num_rows",
		Help: "The number of rows in the word search puzzle",
	})

	ap.Add(&argparse.Argument{
		Name: "word_list_filename",
		Help: "The file which has the list of words, one per line. Or '-' for STDIN",
	})

	ap.Parse()

	if s.Verbose {
		setLogger("-")
	} else {
		setLogger("")
	}

	err := s.run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s", err)
		os.Exit(1)
	}
}

func setLogger(logfileName string) {
	switch logfileName {
	case "":
		log.SetOutput(ioutil.Discard)
		return
	case "-":
		log.SetOutput(os.Stderr)
	default:
		fh, err := os.Create(logfileName)
		if err != nil {
			log.Fatalf("Cannot create %s for logging: %s", logfileName, err)
		}
		log.SetOutput(fh)
	}
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Lshortfile)
}

func main() {
	var p Program
	p.runCli()
}
