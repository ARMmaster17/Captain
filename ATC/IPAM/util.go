package IPAM

import (
	"github.com/rs/zerolog/log"
	"net"
)

func parseSubnetBlocks(blocks []string) []net.IPNet {
	var outBlocks []net.IPNet
	for _, block := range blocks {
		_, net, err := net.ParseCIDR(block)
		if err != nil {
			log.Debug().Msgf("unrecognized IP block: %s", block)
		} else {
			outBlocks = append(outBlocks, *net)
		}
	}
	return outBlocks
}
