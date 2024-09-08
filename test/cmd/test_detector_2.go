package cmd

import (
	"context"
	"math/big"

	log "github.com/sirupsen/logrus"

	"github.com/rocky2015aaa/ethdefender/internal/config"
	"github.com/rocky2015aaa/ethdefender/internal/services"
	"github.com/spf13/cobra"
)

var detectorTest2Cmd = &cobra.Command{
	Use:   "dt2",
	Short: "Detector Test 2",
	Long:  "Test Command if the detector catches multi txs(more than 3 txs) in a minute and send the notification",
	Run: func(cmd *cobra.Command, args []string) {
		auth, err := services.GetTransactor(config.Default.Ethereum.TestAccountPrivateKey)
		if err != nil {
			log.Fatalf("Failed to create transactor: %v", err)
		}
		// Define the number of transactions to send
		numTxs := 4

		for i := 0; i < numTxs; i++ {
			gasPrice, err := ethClient.SuggestGasPrice(context.Background())
			if err != nil {
				log.Fatalf("failed to get gas price: %v", err)
			}
			err = sendETH(auth, big.NewInt(90000000000000), gasPrice) // Send 0.00009 ETH
			if err != nil {
				log.Fatalf("failed to send eth: %v", err)
			}
			log.Printf("Sending %d time(s) tx is successful\n", i+1)
		}
		log.Println("Sending mutiple tx is successful")
	},
}
