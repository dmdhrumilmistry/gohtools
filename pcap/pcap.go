package pcap

import (
	"fmt"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"github.com/rs/zerolog/log"
)

type PacketProcessor func(packet gopacket.Packet)

func AnalyzeAppLayer(packet gopacket.Packet) {
	appLayer := packet.ApplicationLayer()
	if appLayer != nil {
		payload := appLayer.Payload()
		log.Info().Msg(string(payload))
	}
}

func ListenAndApplyFilter(iface, filter string, snaplen int32, promisc bool, timeout time.Duration, processorFunc PacketProcessor) error {
	devices, err := pcap.FindAllDevs()
	if err != nil {
		log.Error().Err(err).Msg("failed to get devices")
		return err
	}

	deviceFound := false
	for _, device := range devices {
		if device.Name == iface {
			deviceFound = true
		}
	}

	if !deviceFound {
		return fmt.Errorf("%s device not found", iface)
	}

	log.Info().Msgf("Starting to listen on interface %s with filter %s", iface, filter)

	handler, err := pcap.OpenLive(iface, snaplen, promisc, timeout)
	if err != nil {
		log.Error().Err(err).Msg("failed to open device handler")
		return err
	}
	defer handler.Close()

	if err := handler.SetBPFFilter(filter); err != nil {
		log.Error().Err(err).Msg("failed to set bpf filter")
		return err
	}

	source := gopacket.NewPacketSource(handler, handler.LinkType())
	for packet := range source.Packets() {
		// log.Info().Msgf("%v", packet)
		processorFunc(packet)
	}

	return nil
}
