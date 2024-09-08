package services

import (
	"context"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rocky2015aaa/ethdefender/internal/config"
)

func WaitForReceipt(ethClient *ethclient.Client, txHash common.Hash) (*types.Receipt, error) {
	for {
		receipt, err := ethClient.TransactionReceipt(context.Background(), txHash)
		if err == nil {
			return receipt, nil
		}
		if err != ethereum.NotFound {
			return nil, err
		}
		// Wait for a while before retrying
		time.Sleep(5 * time.Second)
	}
}

func GetTransactor(privateKey string) (*bind.TransactOpts, error) {
	key, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		return nil, err
	}

	auth, err := bind.NewKeyedTransactorWithChainID(key, big.NewInt(int64(config.Default.Ethereum.ChainID)))
	if err != nil {
		return nil, err
	}
	return auth, nil
}
