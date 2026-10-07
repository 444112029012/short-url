package main

import (
	"fmt"
	"os"

	"github.com/444112029012/short-url/internal/health"
)

func main() {
	fmt.Println(health.Status())
	os.Exit(0)
}
