package main

import (
	"fmt"
	"sort"
	"strings"
)

// findAnagramSets находит все множества анаграмм в словаре.
func findAnagramSets(words []string) map[string][]string {
	anagrams := make(map[string][]string)
	seen := make(map[string]bool)

	for _, word := range words {
		word = strings.ToLower(word)
		sortedWord := sortString(word)
		anagrams[sortedWord] = append(anagrams[sortedWord], word)
	}

	result := make(map[string][]string)
	for _, group := range anagrams {
		if len(group) > 1 {
			sort.Strings(group)
			key := group[0]
			if !seen[key] {
				result[key] = group
				for _, w := range group {
					seen[w] = true
				}
			}
		}
	}

	return result
}

// sortString сортирует символы в строке в алфавитном порядке.
func sortString(s string) string {
	runes := []rune(s)
	sort.Slice(runes, func(i, j int) bool { return runes[i] < runes[j] })
	return string(runes)
}

func main() {
	words := []string{"пятак", "пятка", "тяпка", "листок", "слиток", "столик", "кипяток", "питокя", "слово"}
	anagramSets := findAnagramSets(words)
	for key, group := range anagramSets {
		fmt.Printf("Key: %s, Group: %v\n", key, group)
	}
}
