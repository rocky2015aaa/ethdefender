package cmd

import (
	"context"
	"fmt"
	"math/big"

	log "github.com/sirupsen/logrus"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rocky2015aaa/ethdefender/internal/config"
	"github.com/rocky2015aaa/ethdefender/internal/services"
	"github.com/spf13/cobra"
)

const (
	isContractPaused = "paused"

	pauseContract   = "pause"
	unpauseContract = "unpause"
)

var offFlag bool

func init() {
	// Add the flags to the command
	unpauseTestCmd.Flags().BoolVarP(&offFlag, "off", "f", false, "Deactivate (unpause) the smart contract")
}

var unpauseTestCmd = &cobra.Command{
	Use:   "pause",
	Short: "pause the smart contract with a flag",
	Long:  "pause the smart contract with a flag",
	Run: func(cmd *cobra.Command, args []string) {
		auth, err := services.GetTransactor(config.Default.Ethereum.ContractOwnerPrivateKey)
		if err != nil {
			log.Fatalf("Failed to create transactor: %v", err)
		}
		abiName := pauseContract
		if offFlag {
			abiName = unpauseContract
		}
		data, err := services.ParsedABI.Pack(abiName)
		if err != nil {
			log.Fatalf("Failed to get function data: %v", err)
		}
		gasPrice, err := ethClient.SuggestGasPrice(context.Background())
		if err != nil {
			log.Fatalf("failed to get gas price: %v", err)
		}
		log.Printf("ContractA ddress: %s\n", config.Default.Ethereum.ContractAddress)
		tx, err := sendTx(ethClient, common.HexToAddress(config.Default.Ethereum.ContractAddress),
			auth, big.NewInt(0), gasPrice, data)
		if err != nil {
			log.Fatalf("Failed to send ETH: %v", err)
		}
		log.Printf("Transaction hash: %s\n", tx.Hash().Hex())
		receipt, err := services.WaitForReceipt(ethClient, tx.Hash())
		if err != nil {
			log.Fatalf("Failed to get receipt: %v", err)
		}
		if receipt.Status == types.ReceiptStatusFailed {
			log.Fatalf("Transaction failed")
		}
		if offFlag {
			fmt.Println("Unpausing the contract is successful")
		} else {
			fmt.Println("Pausing the contract is successful")
		}
	},
}

var checkPaused = &cobra.Command{
	Use:   "ispaused",
	Short: "Test paused function",
	Long:  "Test Command if contract is paused or not",
	Run: func(cmd *cobra.Command, args []string) {
		contractAddress := common.HexToAddress(config.Default.Ethereum.ContractAddress)
		// Call the getOwner function
		callData, err := services.ParsedABI.Pack(isContractPaused)
		if err != nil {
			log.Fatalf("Failed to get abi data: %v", err)
		}
		// Prepare the call message
		msg := ethereum.CallMsg{
			To:   &contractAddress,
			Data: callData,
		}
		// Call the contract function
		result, err := ethClient.CallContract(context.Background(), msg, nil)
		if err != nil {
			log.Fatalf("Failed to get abi function call: %v", err)
		}
		var isPaused bool
		err = services.ParsedABI.UnpackIntoInterface(&isPaused, isContractPaused, result)
		if err != nil {
			log.Fatalf("Failed to unpack abi call: %v", err)
		}
		fmt.Println("The contract is paused?", isPaused)
	},
}
