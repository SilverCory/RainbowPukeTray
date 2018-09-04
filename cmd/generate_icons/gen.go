package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

func main() {
	iconData, err := ioutil.ReadFile("assets/icon.ico")
	if err != nil {
		panic(err)
	}

	greyIconData, err := ioutil.ReadFile("assets/iconGreyscale.ico")
	if err != nil {
		panic(err)
	}

	file, err := os.Create("out.txt")
	if err != nil {
		panic(err)
	}

	defer file.Close()

	fmt.Fprintf(file, "%#v\n%#v\n", iconData, greyIconData)
	time.Sleep(15 * time.Second)
}
