package main

func main() {
	server := NewServer()
	panic(server.ListenAndServe(":6969"))
}
