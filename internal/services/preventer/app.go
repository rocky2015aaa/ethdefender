package preventer

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	log "github.com/sirupsen/logrus"

	"github.com/rocky2015aaa/ethdefender/internal/config"
	"github.com/rocky2015aaa/ethdefender/internal/services"
	"github.com/rocky2015aaa/ethdefender/internal/services/preventer/pendingtxsubscriber"
	"github.com/rocky2015aaa/ethdefender/pkg/service"
)

type App struct {
	rpcClient           *rpc.Client
	ethClient           *ethclient.Client
	pendingTxHashCh     chan common.Hash
	pendingRpcClientSub *rpc.ClientSubscription
}

func NewApp() *App {
	services.Setup()

	rpcClient, err := rpc.Dial(config.Default.Ethereum.SubscriptionUrl)
	if err != nil {
		log.WithError(err).Fatal("Failed to connect to the Ethereum client")
	}

	// Create a new ethclient instance
	ethClient, err := ethclient.Dial(config.Default.Ethereum.SubscriptionUrl)
	if err != nil {
		log.WithError(err).Fatal("Failed to connect to the Ethereum client")
	}

	// Channel to receive pending transaction hashes
	pendingTxHashCh := make(chan common.Hash)

	// Subscribe to pending transactions
	pendingTxSub, err := rpcClient.Subscribe(context.Background(), "eth", pendingTxHashCh, "newPendingTransactions")
	if err != nil {
		log.WithError(err).Fatal("Failed to subscribe to pending transactions")
	}

	return &App{
		rpcClient:           rpcClient,
		ethClient:           ethClient,
		pendingTxHashCh:     pendingTxHashCh,
		pendingRpcClientSub: pendingTxSub,
	}
}

func (a *App) Run(ctx context.Context) {
	service.RunWithGracefulShutdown(ctx, func(ctx context.Context) {
		pendingtxsubscriber.New(a.rpcClient, a.ethClient, a.pendingTxHashCh, a.pendingRpcClientSub).Run(ctx)
		defer a.pendingRpcClientSub.Unsubscribe()
	})
}
