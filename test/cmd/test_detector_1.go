package cmd

import (
	"context"
	"fmt"
	"math/big"

	log "github.com/sirupsen/logrus"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/spf13/cobra"

	"github.com/rocky2015aaa/ethdefender/internal/config"
	"github.com/rocky2015aaa/ethdefender/internal/services"
)

var detectorTest1Cmd = &cobra.Command{
	Use:   "dt1",
	Short: "Detector Test 1",
	Long:  "Test Command if the detector catches a big amount ETH sending tx and send the notification",
	Run: func(cmd *cobra.Command, args []string) {
		auth, err := services.GetTransactor(config.Default.Ethereum.TestAccountPrivateKey)
		if err != nil {
			log.Fatalf("failed to get auth: %v", err)
		}
		gasPrice, err := ethClient.SuggestGasPrice(context.Background())
		if err != nil {
			log.Fatalf("failed to get gas price: %v", err)
		}
		err = sendETH(auth, big.NewInt(110000000000000), gasPrice) // Send 0.00011 ETH
		if err != nil {
			log.Fatalf("failed to send eth: %v", err)
		}
		fmt.Println("Sending 0.00011 ETH is successful")
	},
}

func sendETH(auth *bind.TransactOpts, amount, gasPrice *big.Int) error {
	log.Printf("ContractA ddress: %s\n", config.Default.Ethereum.ContractAddress)
	tx, err := sendTx(ethClient, common.HexToAddress(config.Default.Ethereum.ContractAddress),
		auth, amount, gasPrice, nil)
	if err != nil {
		return err
	}
	log.Printf("Transaction hash: %s\n", tx.Hash().Hex())
	receipt, err := services.WaitForReceipt(ethClient, tx.Hash())
	if err != nil {
		return err
	}
	if receipt.Status == types.ReceiptStatusFailed {
		return fmt.Errorf("Transaction failed")
	}
	return nil
}
