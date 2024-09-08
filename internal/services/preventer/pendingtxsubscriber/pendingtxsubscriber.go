package pendingtxsubscriber

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	log "github.com/sirupsen/logrus"

	"github.com/rocky2015aaa/ethdefender/internal/config"
	"github.com/rocky2015aaa/ethdefender/internal/services"
)

const pauseContractFunction = "pause"

type PendingTxSubscriber struct {
	rpcClient           *rpc.Client
	ethClient           *ethclient.Client
	pendingTxHashCh     chan common.Hash
	pendingRpcClientSub *rpc.ClientSubscription
}

func New(rpcClient *rpc.Client, ethClient *ethclient.Client, pendingTxHashCh chan common.Hash, pendingTxSub *rpc.ClientSubscription) *PendingTxSubscriber {
	return &PendingTxSubscriber{
		rpcClient:           rpcClient,
		ethClient:           ethClient,
		pendingTxHashCh:     pendingTxHashCh,
		pendingRpcClientSub: pendingTxSub,
	}
}

func (pendingTxSub *PendingTxSubscriber) Run(ctx context.Context) {
	for {
		if err := pendingTxSub.runWithReconnect(ctx); err != nil {
			log.WithError(err).Error("Failed to handle subscription")
			pendingTxSub.pendingRpcClientSub, err = pendingTxSub.rpcClient.Subscribe(context.Background(), "eth", pendingTxSub.pendingTxHashCh, "newPendingTransactions")
			if err != nil {
				log.WithError(err).Error("Failed to subscribe to pending transactions")
			}
		}
		log.Info("Reconnecting...")
		time.Sleep(5 * time.Second) // Delay before attempting to reconnect
	}
}

func (pendingTxSub *PendingTxSubscriber) runWithReconnect(ctx context.Context) error {
	// Monitor the pending transactions
	for {
		select {
		case <-ctx.Done():
			log.Info("Preventer has stopped")
			return nil
		case err := <-pendingTxSub.pendingRpcClientSub.Err():
			log.WithError(err).Error("Pending transaction subscription error")
			return err
		case txHash := <-pendingTxSub.pendingTxHashCh:
			// Fetch the transaction details using its hash
			tx, isPending, err := pendingTxSub.ethClient.TransactionByHash(context.Background(), txHash)
			if err != nil && err.Error() != "not found" {
				log.WithError(err).Error("Failed to fetch transaction")
				continue
			}
			// Only process the transaction if it's pending
			if isPending {
				// Example: Check if the transaction is sent to a specific contract address
				if tx.To() != nil {
					contractAddress := common.HexToAddress(config.Default.Ethereum.ContractAddress)
					if *tx.To() == contractAddress {
						if isSuspicious(tx) {
							log.WithFields(
								log.Fields{
									"tx_value": tx.Hash(),
								}).Info("Suspicious pending transaction has detected.")
							err := pauseContract(pendingTxSub.ethClient, tx.GasPrice())
							if err != nil {
								log.WithError(err).Error("Pausing contract error")
							} else {
								log.WithFields(
									log.Fields{
										"contract_address": config.Default.Ethereum.ContractAddress,
									}).Info("The contract has paused.")
							}
						}
					}
				}
			}
		}
	}
}

func isSuspicious(tx *types.Transaction) bool {
	// Logic to detect suspicious transactions
	// For example, high gas price or unusual transaction value
	return tx.GasPrice().Cmp(big.NewInt(1000000000)) > 0 // Example: gas price greater than 1 Gwei
}

func pauseContract(ethClient *ethclient.Client, competitorGasPrice *big.Int) error {
	auth, err := services.GetTransactor(config.Default.Ethereum.ContractOwnerPrivateKey)
	if err != nil {
		return err
	}

	data, err := services.ParsedABI.Pack(pauseContractFunction)
	if err != nil {
		return err
	}

	nonce, err := ethClient.PendingNonceAt(context.Background(), auth.From)
	if err != nil {
		return err
	}

	adjustedGasPrice := new(big.Int).Set(competitorGasPrice)
	adjustedGasPrice.Mul(adjustedGasPrice, big.NewInt(300)) // 3 times the gas fee
	adjustedGasPrice.Div(adjustedGasPrice, big.NewInt(100))

	// Define the gas limit (can be made dynamic)
	gasLimit := uint64(90000)

	tx := types.NewTransaction(nonce, common.HexToAddress(config.Default.Ethereum.ContractAddress),
		big.NewInt(0), gasLimit, adjustedGasPrice, data)

	signedTx, err := auth.Signer(auth.From, tx)
	if err != nil {
		return err
	}

	// Set up a context with timeout for sending the transaction
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Send the transaction asynchronously
	go func() {
		err := ethClient.SendTransaction(ctx, signedTx)
		if err != nil {
			log.Printf("failed to send transaction: %v", err)
			return
		}
		log.Printf("Transaction sent successfully: %s", signedTx.Hash().Hex())
	}()

	receipt, err := services.WaitForReceipt(ethClient, signedTx.Hash())
	if err != nil {
		return err
	}

	if receipt.Status == types.ReceiptStatusFailed {
		return fmt.Errorf("transaction failed with out of gas or other errors")
	}

	return nil
}
