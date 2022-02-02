package main

import (
	"fmt"
	"log"

	"github.com/coffemanfp/chat/server"
)

const PORT = 8080
const STATIC_FILES_PATH = "./public"

func main() {
	fmt.Println("Starting...")
	server := server.NewServer(PORT, STATIC_FILES_PATH)

	fmt.Printf("Listening on port: %d\n", PORT)
	log.Fatal(server.Run())
}
