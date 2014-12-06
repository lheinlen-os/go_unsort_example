package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"os"
	"time"
)

// Global variable to hold the path to the inputFile from flags
var inputFile string

// Global variable to hold the path to the outputFile from flags
// Honestly, I probably wouldn't use an output file if I was writing this myself.  I prefer the unix-style methodology
// where output goes to STDOUT so you can pipe it into another command.  Then, if I want a file, I just redirect it
// like unsort file > somefile.txt.  Using flags is unnecessary in that case as there is just a single argument.  So,
// I left it like this to give an example of using the flag library and to make it more like the original.
var outputFile string

// Init functions from each source file in a go program run before the main.  It is important to note that their order
// cannot be guaranteed.  We're using this to configure the flags.
func init() {
	flag.StringVar(&inputFile, "i", "", "The input file to unsort")
	flag.StringVar(&outputFile, "o", "", "The output file into which the unsorted content of inputfile will be placed.  This file *will* be overwritten.")
}

// The main entrypoint for a program is the main() function in the main package.
func main() {
	flag.Parse()

	// If these were not passed in they will be the "zero" value for their type (string).  That is "".
	if inputFile == "" || outputFile == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Open the input file for reading
	// Go has the ability to return multiple values from function calls.  It is a very common idiom to return the expected
	// return and an error struct.  The error will be its zero value (a nil) if there was no error and an error if not.
	input, err := os.Open(inputFile)
	if err != nil {
		// Panic immediately ends the program after printing an error
		panic(err)
	}
	// The defer keyword is totally awesome.  It pushes the command after it only a list of commands all of which will
	// be executed after the surrounding function exits.  This lets you do important cleanup for a function (like closing
	// files) whether it ended in error, or quit early or quit normally.  Pretty great.
	// http://blog.golang.org/defer-panic-and-recover
	defer input.Close()

	// Open the output file for writing (this will create or truncate the file as well)
	output, err := os.Create(outputFile)
	if err != nil {
		panic(err)
	}
	defer output.Close()

	// Actually run the unsort algorithm seeding it with the current time in nanoseconds
	err = unsort(input, output, time.Now().UnixNano())
	if err != nil {
		panic(err)
	}
}

// The unsort algorithm.  I'm actually just going to load the lines into memory and randomize them.  If this is a problem
// because the files are extremely large there are other methods to achieve this but I'm assuming smallish files.
// I'm passing in a reader and writer and the main function is actually opening the files.  That allows me to much
// more easily unit test as I can pass in string (or []byte) based readers/writers and I don't need actual files for my
// unit tests.  Unit testing is also why I'm passing in the seed.  This allows me to pass in a known seed so I can check known results.
func unsort(input io.Reader, output io.Writer, seed int64) error {
	// Seed the random number library
	rand.Seed(seed)

	// Create an empty slice to hold the lines of the file.  A slice is Go's kind of special array implementation.
	// http://blog.golang.org/go-slices-usage-and-internals
	lines := []string{}

	// Create a scanner object which makes it easy to read line by line
	scanner := bufio.NewScanner(input)
	// scanner.Scan() returns true if it gets a line and false if it is at the end of the file or it gets an error.
	for scanner.Scan() {
		// The mechanics of append are interesting http://blog.golang.org/slices but the short version is that it is returning
		// a slice with the text of the current line from the scanner appended to the end.
		lines = append(lines, scanner.Text())
	}

	// Check to see if we quit looping because of an error.
	if scanner.Err() != nil {
		return scanner.Err()
	}

	// The rand.Perm function actually takes a slice and returns a randomized version of it.  However, that means that
	// we have the file content in memory twice (original slice & randomized slice).  That is wholly unnecessary.
	// So instead, we randomize ourselves using the http://en.wikipedia.org/wiki/Fisher%E2%80%93Yates_shuffle.

	// range actually returns two values the index and the value at that index.  In this case, we don't need the value
	// so we can just leave it since it's second.  If we had only wanted the value (as we will later) we would have to
	// discard the index by using the special variable _.
	for i := range lines {
		// Intn returns a random number between 0 and n-1
		// http://golang.org/pkg/math/rand/#Intn
		j := rand.Intn(i + 1)
		// parallel assignment allows us to swap without having to use a temporary variable
		lines[i], lines[j] = lines[j], lines[i]
	}

	// Loop through the now random lines and print them to the output writer.
	// Here is an example of using the _ to ignore the first return value (the index in this case).
	for _, line := range lines {
		_, err := fmt.Fprintln(output, line)
		if err != nil {
			return err
		}
	}

	// Return nil to indicate there was no error
	return nil
}
