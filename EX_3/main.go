package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode"
)

func main() {
	inputFile := flag.String("i", "", "Input file path")
	outputFile := flag.String("o", "", "Output file path")
	column := flag.Int("k", 0, "Column number for sorting (1-based index, 0 for entire line)")
	numSort := flag.Bool("n", false, "Sort numerically")
	reverse := flag.Bool("r", false, "Sort in reverse order")
	unique := flag.Bool("u", false, "Unique lines only")
	monthSort := flag.Bool("M", false, "Sort by month name")
	ignoreTail := flag.Bool("b", false, "Ignore trailing spaces")
	checkSort := flag.Bool("c", false, "Check if sorted")
	humanSort := flag.Bool("h", false, "Sort by human readable sizes")
	flag.Parse()

	if *inputFile == "" || *outputFile == "" {
		fmt.Println("Input and output file paths must be specified.")
		flag.Usage()
		return
	}

	content, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		fmt.Printf("Error reading input file: %v\n", err)
		return
	}

	lines := strings.Split(string(content), "\n")
	if len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1] // Remove trailing empty line if any
	}

	if *checkSort {
		if isSorted(lines, *column, *numSort, *reverse, *monthSort, *ignoreTail, *humanSort) {
			fmt.Println("The file is already sorted.")
		} else {
			fmt.Println("The file is not sorted.")
		}
		return
	}

	sortedLines := sortLines(lines, *column, *numSort, *reverse, *unique, *monthSort, *ignoreTail, *humanSort)
	err = ioutil.WriteFile(*outputFile, []byte(strings.Join(sortedLines, "\n")), 0644)
	if err != nil {
		fmt.Printf("Error writing output file: %v\n", err)
	}
}

func sortLines(lines []string, column int, numSort, reverse, unique, monthSort, ignoreTail, humanSort bool) []string {
	if monthSort {
		monthMap := map[string]time.Month{
			"jan": time.January, "feb": time.February, "mar": time.March,
			"apr": time.April, "may": time.May, "jun": time.June,
			"jul": time.July, "aug": time.August, "sep": time.September,
			"oct": time.October, "nov": time.November, "dec": time.December,
		}
		sort.SliceStable(lines, func(i, j int) bool {
			return monthMap[getField(lines[i], column)] < monthMap[getField(lines[j], column)]
		})
	} else if numSort || humanSort {
		sort.SliceStable(lines, func(i, j int) bool {
			numI := extractNumber(getField(lines[i], column))
			numJ := extractNumber(getField(lines[j], column))
			if reverse {
				return numI > numJ
			}
			return numI < numJ
		})
	} else {
		sort.SliceStable(lines, func(i, j int) bool {
			fieldI := getField(lines[i], column)
			fieldJ := getField(lines[j], column)
			if ignoreTail {
				fieldI = strings.TrimRightFunc(fieldI, unicode.IsSpace)
				fieldJ = strings.TrimRightFunc(fieldJ, unicode.IsSpace)
			}
			if reverse {
				return fieldI > fieldJ
			}
			return fieldI < fieldJ
		})
	}

	if unique {
		lines = uniqueLines(lines)
	}

	return lines
}

func getField(line string, column int) string {
	fields := strings.Fields(line)
	if column > 0 && column <= len(fields) {
		return fields[column-1]
	}
	return line
}

func extractNumber(s string) int {
	num, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return num
}

func uniqueLines(lines []string) []string {
	seen := make(map[string]bool)
	result := []string{}
	for _, line := range lines {
		if !seen[line] {
			seen[line] = true
			result = append(result, line)
		}
	}
	return result
}

func isSorted(lines []string, column int, numSort, reverse, monthSort, ignoreTail, humanSort bool) bool {
	for i := 1; i < len(lines); i++ {
		if !isOrdered(lines[i-1], lines[i], column, numSort, reverse, monthSort, ignoreTail, humanSort) {
			return false
		}
	}
	return true
}

func isOrdered(a, b string, column int, numSort, reverse, monthSort, ignoreTail, humanSort bool) bool {
	fieldA := getField(a, column)
	fieldB := getField(b, column)

	if ignoreTail {
		fieldA = strings.TrimRightFunc(fieldA, unicode.IsSpace)
		fieldB = strings.TrimRightFunc(fieldB, unicode.IsSpace)
	}

	if monthSort {
		monthMap := map[string]int{
			"jan": 1, "feb": 2, "mar": 3, "apr": 4,
			"may": 5, "jun": 6, "jul": 7, "aug": 8,
			"sep": 9, "oct": 10, "nov": 11, "dec": 12,
		}
		monthA := monthMap[strings.ToLower(fieldA)]
		monthB := monthMap[strings.ToLower(fieldB)]
		if reverse {
			return monthA > monthB
		}
		return monthA < monthB
	}

	if numSort || humanSort {
		numA := extractNumber(fieldA)
		numB := extractNumber(fieldB)
		if reverse {
			return numA > numB
		}
		return numA < numB
	}

	if reverse {
		return fieldA > fieldB
	}
	return fieldA < fieldB
}
