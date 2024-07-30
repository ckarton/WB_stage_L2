package main

import (
	"fmt"
	"os"
	"time"

	"github.com/beevik/ntp"
)

func main() {
	// Попытка получить точное время с NTP сервера
	currentTime, err := ntp.Time("pool.ntp.org")
	if err != nil {
		// Вывод ошибки в STDERR и завершение программы с кодом 1
		fmt.Fprintln(os.Stderr, "Error fetching time:", err)
		os.Exit(1)
	}

	// Вывод текущего времени в стандартный вывод
	fmt.Println("Current time:", currentTime.Format(time.RFC3339))
}
