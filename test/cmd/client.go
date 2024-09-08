package cmd

import (
	"log"

	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/rocky2015aaa/ethdefender/internal/config"
	"github.com/rocky2015aaa/ethdefender/internal/services"
)

var (
	ethClient *ethclient.Client
)

func init() {
	services.InitConfig()
	services.InitContractABI()

	var err error
	ethClient, err = ethclient.Dial(config.Default.Ethereum.SubscriptionUrl)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}
}
