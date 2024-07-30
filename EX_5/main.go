package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

type GrepOptions struct {
	after      int
	before     int
	context    int
	count      bool
	ignoreCase bool
	invert     bool
	fixed      bool
	lineNum    bool
}

func main() {
	after := flag.Int("A", 0, "Print +N lines after match")
	before := flag.Int("B", 0, "Print +N lines before match")
	context := flag.Int("C", 0, "Print ±N lines around match")
	count := flag.Bool("c", false, "Count lines")
	ignoreCase := flag.Bool("i", false, "Ignore case")
	invert := flag.Bool("v", false, "Invert match")
	fixed := flag.Bool("F", false, "Fixed string match")
	lineNum := flag.Bool("n", false, "Print line number")
	flag.Parse()

	if flag.NArg() < 2 {
		fmt.Println("Usage: grep [OPTIONS] PATTERN FILE")
		flag.Usage()
		return
	}

	pattern := flag.Arg(0)
	filePath := flag.Arg(1)

	options := GrepOptions{
		after:      *after,
		before:     *before,
		context:    *context,
		count:      *count,
		ignoreCase: *ignoreCase,
		invert:     *invert,
		fixed:      *fixed,
		lineNum:    *lineNum,
	}

	err := grep(pattern, filePath, options)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	}
}

func grep(pattern, filePath string, options GrepOptions) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	if options.context > 0 {
		options.after = options.context
		options.before = options.context
	}

	var matches []int
	for i, line := range lines {
		if match(line, pattern, options) {
			matches = append(matches, i)
		}
	}

	if options.count {
		fmt.Println(len(matches))
		return nil
	}

	printMatches(lines, matches, options)
	return nil
}

func match(line, pattern string, options GrepOptions) bool {
	if options.ignoreCase {
		line = strings.ToLower(line)
		pattern = strings.ToLower(pattern)
	}
	if options.fixed {
		return options.invert != strings.Contains(line, pattern)
	}
	return options.invert != strings.Contains(line, pattern)
}

func printMatches(lines []string, matches []int, options GrepOptions) {
	printed := make(map[int]bool)
	for _, match := range matches {
		start := max(0, match-options.before)
		end := min(len(lines), match+options.after+1)

		for i := start; i < end; i++ {
			if printed[i] {
				continue
			}
			printed[i] = true
			if options.lineNum {
				fmt.Printf("%d:", i+1)
			}
			fmt.Println(lines[i])
		}
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
