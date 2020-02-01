package main

import (
	"os"
	"log"
	"fmt"
	"io/ioutil"
	"strconv"
)

func main() {
	args := os.Args[1:]
	if len(args) != 2 {
		log.Println("wrong number of arguments")
		os.Exit(1)
	}

	bright, err := strconv.Atoi(args[1])
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	err = ioutil.WriteFile(fmt.Sprintf("/sys/class/backlight/%s/brightness", args[0]), []byte(strconv.Itoa(bright)), 0444)
	if err != nil {
		log.Println(err)
	}
}
