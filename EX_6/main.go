package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	// Определение флагов
	var fields string
	var delimiter string
	var separated bool

	flag.StringVar(&fields, "f", "", "Выберите поля (колонки)")
	flag.StringVar(&delimiter, "d", "\t", "Использовать другой разделитель")
	flag.BoolVar(&separated, "s", false, "Только строки с разделителем")
	flag.Parse()

	// Обработка ввода
	if len(flag.Args()) > 0 {
		fmt.Println("Usage: cut [-f fields] [-d delimiter] [-s] < input_file")
		return
	}

	// Создание мапы для выбранных полей
	selectedFields := parseFields(fields)

	// Чтение строк из стандартного ввода
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()

		// Проверка на наличие разделителя
		if separated && !strings.Contains(line, delimiter) {
			continue
		}

		// Разделение строки на колонки
		columns := strings.Split(line, delimiter)
		output := []string{}

		for _, index := range selectedFields {
			if index < len(columns) {
				output = append(output, columns[index])
			}
		}

		fmt.Println(strings.Join(output, delimiter))
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Error reading input:", err)
	}
}

// Функция для парсинга полей
func parseFields(fields string) []int {
	var selectedFields []int
	for _, field := range strings.Split(fields, ",") {
		field = strings.TrimSpace(field)
		if field != "" {
			if index, err := strconv.Atoi(field); err == nil {
				selectedFields = append(selectedFields, index-1) // Поля считаем с 1
			}
		}
	}
	return selectedFields
}
