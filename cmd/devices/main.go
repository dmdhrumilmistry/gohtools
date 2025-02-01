package main

import (
	"flag"

	_ "github.com/dmdhrumilmistry/gohtools/logging"
	"github.com/dmdhrumilmistry/gohtools/pcap"
	gpacp "github.com/google/gopacket/pcap"
	"github.com/rs/zerolog/log"
)

func main() {
	device := flag.String("i", "", "device interface name")
	snapLen := flag.Int("b", 1600, "data amount to be captured for each frame")
	promisc := flag.Bool("p", false, "enable promisc mode")
	filter := flag.String("f", "tcp and port 80", "bpf filter")
	flag.Parse()

	if *device == "" {
		log.Fatal().Msg("device interface not found. Use -h for more details")
	}

	if err := pcap.ListenAndApplyFilter(*device, *filter, int32(*snapLen), *promisc, gpacp.BlockForever, pcap.AnalyzeAppLayer); err != nil {
		log.Error().Err(err).Msg("failed to get device interfaces")
	}

}
