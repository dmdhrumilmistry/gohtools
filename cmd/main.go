package main

import (
	"fmt"

	"github.com/dmdhrumilmistry/gohtools/bing"
)

func main() {

	// TCP Servers
	// StartEchoServer(4444, "ioCopyEcho")
	// StartTcpProxy(4444, "localhost:8000")
	// StartListener(4444, HandleNcConn) // connect with server using net.Dial or nc IP 4444

	// HTTP request
	// ip := http.GetHostMachinePublicIP() // package to import: "github.com/dmdhrumilmistry/gohtools/http"
	// fmt.Printf("[*] Host Machine Public Ip: %s\n", ip)

	// msf login
	// msfHost := os.Getenv("MSF_HOST")
	// msfUsername := os.Getenv("MSF_USERNAME")
	// msfPassword := os.Getenv("MSF_PASSWORD")

	// if msfHost == "" || msfPassword == "" || msfUsername == "" {
	// 	log.Fatalln("[ERROR] Required environment variables (MSF_HOST, MSF_USERNAME, MSF_PASSWORD) not found!")
	// }

	// msf, err := metasploit.NewMsfClient(msfHost, msfUsername, msfPassword)
	// if err != nil {
	// 	log.Fatalf("[ERROR] Unable to create MSF client due to error: %s", err)
	// }
	// defer msf.Logout()

	// sessions, err := msf.ListSessions()
	// if err != nil {
	// 	log.Panicf("[ERROR] Unable to list MSF sessions due to error: %s", err)
	// }

	// if len(sessions) < 1 {
	// 	fmt.Println("[*] No active sessions found!")
	// } else {
	// 	fmt.Println("Sessions:")
	// 	for _, session := range sessions {
	// 		fmt.Printf("%5d\t%s\n", session.Id, session.Info)
	// 	}
	// }

	// Bing Doc Search
	var docLinks []string
	site := "example.com"
	docType := "xlsx"

	bsClient := bing.NewBingSearch()
	bsClient.SearchDocument(site, docType)

	for _, result := range bsClient.Results {
		docLink, ok := result.(string)
		if !ok {
			continue
		}
		docLinks = append(docLinks, docLink)
	}

	fmt.Printf("[*] Docs found: %d\n", len(docLinks))
	bsClient.ExtractDocumentsMetaData(docLinks)
}
