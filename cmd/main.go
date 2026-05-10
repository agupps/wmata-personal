package main

import (
	"os"
	"wmata/internal/app"
)


func main() {

	os.Exit(app.New().Run())
}
