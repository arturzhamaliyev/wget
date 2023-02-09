package main

import (
	"fmt"

	"wget/internal/transporter"
)

func main() {
	if err := transporter.Switcher(); err != nil {
		fmt.Println(err)
		return
	}
}
