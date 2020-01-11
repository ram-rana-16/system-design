package main

import (
	"os"
	"fmt"
)

func getPort() string {
	port := os.Getenv("PORT")
	return ":"+ port
}
func logSetup() {
	aURL := os.Getenv("A_CONDITION_URL")
	bURL := os.Getenv("B_CONDITION_URL")
	defultURL := os.Getenv("DEFAULT_CONDITION_URL")
	fmt.Printf("redirect to URL a : %s\n", aURL)
	fmt.Printf("redirect tor URL b: %s\n", bURL)
	fmt.Printf("redirect to default URL: %s\n", defultURL)
	fmt.Printf("server is running on %s\n", getPort())

}
func main() {
	logSetup()

}