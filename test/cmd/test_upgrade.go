package cmd

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	log "github.com/sirupsen/logrus"

	"github.com/rocky2015aaa/ethdefender/internal/config"
	"github.com/rocky2015aaa/ethdefender/internal/services"
	"github.com/spf13/cobra"
)

var upgradeTestCmd = &cobra.Command{
	Use:   "getowner",
	Short: "Test getOwner function",
	Long:  "Test Command if contract upgrade is successful or not",
	Run: func(cmd *cobra.Command, args []string) {
		contractAddress := common.HexToAddress(config.Default.Ethereum.ContractAddress)
		// Call the getOwner function
		callData, err := services.ParsedABI.Pack(services.AddedContractFunction)
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
		var owner common.Address
		err = services.ParsedABI.UnpackIntoInterface(&owner, services.AddedContractFunction, result)
		if err != nil {
			log.Fatalf("Failed to unpack abi call: %v", err)
		}
		fmt.Println("The owner is", owner.Hex())
	},
}
