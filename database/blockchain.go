package database

import (
	"fmt"
	"github.com/kaleido-io/kaleido-fabric-go/fabric"
	"github.com/kaleido-io/kaleido-fabric-go/kaleido"
	"os"
)

var (
	conf    map[string]interface{}
	network *kaleido.KaleidoNetwork
)

func BlockchainConnect() {
	networknew := kaleido.NewNetwork()
	networknew.Initialize()
	config, err := fabric.BuildConfig(networknew)
	if err != nil {
		fmt.Printf("Failed to generate network configuration for the SDK: %s\n", err)
		os.Exit(1)
	}

	SetBlockchainConfig(config)
	SetBlockchainNetwork(networknew)
}

func SetBlockchainConfig(config map[string]interface{}) {
	conf = config
}

func GetBlockchainConfig() map[string]interface{} {
	return conf
}

func SetBlockchainNetwork(net *kaleido.KaleidoNetwork) {
	network = net
}

func GetBlockchainNetwork() *kaleido.KaleidoNetwork {
	return network
}
