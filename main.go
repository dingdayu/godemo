package main

import (
	"demo/cmd"

	_ "github.com/joho/godotenv/autoload"
	_ "go.uber.org/automaxprocs"
)

func main() {
	cmd.Execute()
}
