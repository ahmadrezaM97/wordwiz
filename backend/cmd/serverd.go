package main

import "log"

func main() {
	if err := run(); err != nil {
		log.Println(err)
	}
}

func run() error {
	return nil
}
