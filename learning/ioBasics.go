package learning

import (
	"fmt"
	"io"
	"log"
	"os"
)

type Reader struct{}

func (reader *Reader) Read(buffer []byte) (int, error) {
	fmt.Print("input > ")
	return os.Stdin.Read(buffer)
}

type Writer struct{}

func (writer *Writer) Write(buffer []byte) (int, error) {
	fmt.Print("output >")
	return os.Stdout.Write(buffer)
}

func ManualCopyMethod() {
	var (
		reader Reader
		writer Writer
	)

	buffer := make([]byte, 4096) // create buffer to hold input and output data

	size, err := reader.Read(buffer)
	if err != nil {
		fmt.Println("Unable to read data")
	}

	fmt.Printf("Read %d bytes from Stdin\n", size)

	_, err = writer.Write(buffer)
	if err != nil {
		fmt.Println("Unable to write data")
	}
	fmt.Printf("Wrote %d bytes from Stdout\n", size)

}

func IoCopyMethod() {
	var (
		reader Reader
		writer Writer
	)

	if _, err := io.Copy(&writer, &reader); err != nil {
		log.Fatalln("Unable to read/write data!")
	}
}
