package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

// Version (semantic)
const Version = "0.0.3"

// environmental variable
const nextwordDataPath = "NEXTWORD_DATA_PATH"

// flags
var versionFlag = flag.Bool("v", false, "show version")
var dataPath = flag.String("d", os.Getenv(nextwordDataPath), "path to the data directory")
var candidateNum = flag.Int("c", 10, "max candidates number")
var helpFlag = flag.Bool("h", false, "show this message")
var greedyFlag = flag.Bool("g", false, "show as many result as possible")

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	flag.Parse()

	// version
	if *versionFlag {
		showVersion()
		return nil
	}

	// help
	if *helpFlag {
		showHelpMessage()
		return nil
	}

	// new nextword
	params := &NextwordParams{
		DataPath:     *dataPath,
		CandidateNum: *candidateNum,
		Greedy:       *greedyFlag,
	}
	nw, err := NewNextword(params)
	if err != nil {
		return err
	}

	// loop
	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		candidates, err := nw.Suggest(sc.Text())
		if err != nil {
			return err
		}
		fmt.Println(strings.Join(candidates, " "))
	}
	if sc.Err() != nil {
		return sc.Err()
	}

	return nil
}

func showVersion() {
	fmt.Fprintln(os.Stderr, fmt.Sprintf("nextword version %s", Version))
}

func showHelpMessage() {
	fmt.Fprintln(os.Stderr, "Nextword prints the most likely English words that follow the stdin sentence.")
	fmt.Fprintln(os.Stderr)
	flag.Usage()
	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, `You need to install nextword-data and set "NEXTWORD_DATA_PATH" environment variable.`)
	fmt.Fprintln(os.Stderr, `It available at https://github.com/high-moctane/nextword-data`)
	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, `The result depends on whether the input string ends with a space character.`)
	fmt.Fprintln(os.Stderr, `If the string does not end with a space, nextword will print candidate words which`)
	fmt.Fprintln(os.Stderr, `begin the last word in the sentence.`)
	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, `Example:`)
	fmt.Fprintln(os.Stderr, `	input:  "Alice was "`)
	fmt.Fprintln(os.Stderr, `	output: "not a the in still born very so to beginning too at sitting ..."`)
	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, `	input:  "Alice w"`)
	fmt.Fprintln(os.Stderr, `	output: "was would were went with will who wrote when wants ..."`)
}
