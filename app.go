package main

import (
	"context"
	"fmt"
	"log"

	firebase "firebase.google.com/go"
)

func main() {
	fmt.Println("hello world")
	_, err := firebase.NewApp(context.Background(), nil)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}
}
