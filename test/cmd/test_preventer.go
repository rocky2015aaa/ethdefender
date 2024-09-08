package cmd

import (
	"fmt"
	"math/big"

	log "github.com/sirupsen/logrus"

	"github.com/rocky2015aaa/ethdefender/internal/config"
	"github.com/rocky2015aaa/ethdefender/internal/services"
	"github.com/spf13/cobra"
)

var preventerTestCmd = &cobra.Command{
	Use:   "pt",
	Short: "Preventer Test",
	Long:  "Test Command if the preventer catches tx used much gas(more than 1 Gwei) and pause the contract",
	Run: func(cmd *cobra.Command, args []string) {
		auth, err := services.GetTransactor(config.Default.Ethereum.TestAccountPrivateKey)
		if err != nil {
			log.Fatalf("Failed to create transactor: %v", err)
		}
		log.Printf("ContractA ddress: %s\n", config.Default.Ethereum.ContractAddress)

		err = sendETH(auth, big.NewInt(90000000000000), big.NewInt(1000100000)) // Send 0.00009 ETH with Gwei 1.0001
		if err != nil {
			log.Fatalf("failed to send eth: %v", err)
		}
		fmt.Println("Sending 0.00009 ETH is successful")
	},
}
