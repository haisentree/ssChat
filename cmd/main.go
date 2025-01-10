package main

import "ssChat/internal"

func main() {
	s := internal.NewSSChatServer()
	s.Run()
}
