package main

import (
	"fmt"
	"os"
)

func main() {
	envs := []string{"GOFILE", "GOLINE", "GOPACKAGE", "DOLLAR"}
	for _, e := range envs {
		fmt.Printf("%s: %s\n", e, os.Getenv(e))
	}
}
