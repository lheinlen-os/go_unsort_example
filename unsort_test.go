package main

import (
	"bytes"
	"fmt"
	"strconv"
	"testing"
)

func Test_unsort_1(t *testing.T) {
	// Create two buffers (which are in memory ReadWriters backed by slices of bytes)
	var b, output bytes.Buffer

	// Populate input with 10 lines each containg a single number (as a string) 0-9
	// The & means to pass a pointer to the b instead of b to the Println method.  This is because the "Writer" interface
	// is implemented on the pointer to a Buffer and not on Buffer itself.  This is one of the few things that I feel
	// gets confusing in go and probably doesn't make much sense.  The links below should be helpful.
	// https://golang.org/doc/effective_go.html#interfaces
	// https://golang.org/doc/effective_go.html#pointers_vs_values
	for i := 0; i < 10; i++ {
		fmt.Fprintln(&b, strconv.Itoa(i))
	}

	inputString := b.String()
	t.Logf("Input string is %q", inputString)
	// Create a reader from b's contents
	input := bytes.NewBufferString(inputString)

	// unsort input with a fixed seed.
	unsort(input, &output, 0)

	// I ran the index choosing part of the shuffling algorithm with this seed to divine which indexes it would choose
	// https://play.golang.org/p/WCcH3WhzBu
	// https://docs.google.com/spreadsheets/d/1dedMtWwanUycEQBG3NaEU-GwNHqvqAy7hK20xJgJ4bc
	// This is not the greatest test since the test is based on part of the actual algorithm (so we're testing what we coded instead of what it should do)
	// But, the most important thing is not the order but the fact it shuffled them and didn't lose anything, so I feel ok about it.
	expected := "8\n2\n3\n0\n5\n7\n1\n6\n9\n4\n"
	if output.String() != expected {
		t.Errorf("Expected randomized writer to contain...\n%q\n...but was...\n%q", expected, output.String())
	}
}
