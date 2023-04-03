package main

import "os"

func main() {
	notMain()
}
func notMain() {
	os.Exit(1)
}
