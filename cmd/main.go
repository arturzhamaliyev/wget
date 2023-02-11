package main

import (
	"fmt"

	"wget/internal/transporter"
)

func main() {
	// file, _ := os.Open("hmm-dot-dot-dot-stick-figure-intresting-gif-23376515")
	// scanner := bufio.NewScanner(file)
	// for scanner.Scan() {
	// 	fmt.Println(scanner.Text())
	// 	fmt.Println("BRUH")
	// }

	if err := transporter.Switcher(); err != nil {
		fmt.Println(err)
		return
	}
}
