package main

func main() {
	// StartEchoServer(4444, "ioCopyEcho")
	StartTcpProxy(4444, "localhost:8000")
}
