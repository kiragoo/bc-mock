package main

import (
	"bufio"
	"fmt"
	"os"
)

func Prompt(p string) {
	fmt.Printf("[> %s\n", p)
	fmt.Printf("-> Press Return key to continue.")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		break
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	fmt.Println()
}
