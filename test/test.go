package main

import (
	"fmt"
	"os"
)

func main() {
	dirs, err := os.ReadDir("docs")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	for _, dir := range dirs {

		if !dir.IsDir() || dir.Name() == "static" {
			continue
		}

		fmt.Println(dir.Name())
	}
}
