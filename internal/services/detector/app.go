package detector

import (
	"context"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	log "github.com/sirupsen/logrus"

	"github.com/rocky2015aaa/ethdefender/internal/config"
	"github.com/rocky2015aaa/ethdefender/internal/services"
	"github.com/rocky2015aaa/ethdefender/internal/services/detector/txsubscriber"
	"github.com/rocky2015aaa/ethdefender/pkg/service"
)

type App struct {
	ethClient *ethclient.Client
	txLogs    chan types.Log
	txSub     ethereum.Subscription
}

func NewApp() *App {
	services.Setup()

	ethClient, err := ethclient.Dial(config.Default.Ethereum.SubscriptionUrl)
	if err != nil {
		log.WithError(err).Fatal("Ethereum client init error")
	}

	txLogs := make(chan types.Log)
	contractAddress := common.HexToAddress(config.Default.Ethereum.ContractAddress)
	query := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddress},
	}
	txSub, err := ethClient.SubscribeFilterLogs(context.Background(), query, txLogs)
	if err != nil {
		log.Fatalf("Failed to subscribe to contract logs: %v", err)
	}

	return &App{
		ethClient: ethClient,
		txLogs:    txLogs,
		txSub:     txSub,
	}
}

func (a *App) Run(ctx context.Context) {
	service.RunWithGracefulShutdown(ctx, func(ctx context.Context) {
		txsubscriber.New(a.ethClient, a.txLogs, a.txSub).Run(ctx)
	})
}
