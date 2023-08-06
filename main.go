package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/michaelcosj/stms/cmd"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("error loading env file")
	}

	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}
