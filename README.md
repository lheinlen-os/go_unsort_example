# Go Unsort Example

This is meant to be a simple example to help illustrate some core Go concepts.  It is not meant to be a robust tool.

## Installation

To install this example, please use `go get`.

```
$ go get github.com/lheinlen/go_unsort_example
```

## Usage
To run the example utility itself...

```
$ go_unsort_example
  -i="": The input file to unsort
  -o="": The output file into which the unsorted content of inputfile will be placed.  This file *will* be overwritten.
```

To run the tests...

```
$ go test github.com/lheinlen/go_unsort_example
ok      github.com/lheinlen/go_unsort_example   0.293s
```
