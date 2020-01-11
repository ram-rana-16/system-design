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
	bURL := os.Getenv("B_CONDITION_RL")
	defultURL := os.Getenv("DEFAULT_CONDITION_URL")
	fmt.Printf("redirect to URL a : %s", aURL)
	fmt.Printf("redirect tor URL b: %s", bURL)
	fmt.Printf("redirect to default URL: %s", defultURL)
	fmt.Printf("server is running on %s", getPort())

}
func main() {
	logSetup()

}