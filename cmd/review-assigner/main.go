package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("server port:", os.Getenv("SERVER_PORT"))
}
