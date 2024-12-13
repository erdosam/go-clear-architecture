package main

import (
	"context"
	"log"

	firebase "firebase.google.com/go"
)

func init() {

}

func main() {
	_, err := firebase.NewApp(context.Background(), nil)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}
	ListenAndServe()
}
