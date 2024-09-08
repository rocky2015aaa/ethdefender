package txsubscriber

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rocky2015aaa/ethdefender/internal/config"
	log "github.com/sirupsen/logrus"
	"gopkg.in/gomail.v2"
)

// Global map to track tx and a mutex for thread safety
var (
	txHistory = make(map[common.Address][]time.Time)
	mu        sync.Mutex
)

// Maximum allowed txs in a given time period (e.g., 3 in 5 minute)
const (
	maxTxs     = 3
	timeWindow = 1 * time.Minute // 1 minute window

	emailFromHeader     = "From"
	emailToHeader       = "To"
	emailSubjectHeader  = "Subject"
	emailSubjectContent = "Suspicious Transaction Detected"
	emailBodyType       = "text/plain"

	notificationType1 = "much ETH"
	notificationType2 = "many TX"
)

type TxSubscriber struct {
	ethClient *ethclient.Client
	txLogs    chan types.Log
	txSub     ethereum.Subscription
}

func New(ethClient *ethclient.Client, txLogs chan types.Log, txSub ethereum.Subscription) *TxSubscriber {
	return &TxSubscriber{
		ethClient: ethClient,
		txLogs:    txLogs,
		txSub:     txSub,
	}
}

// Run starts the subscription and handles reconnections
func (txSub *TxSubscriber) Run(ctx context.Context) {
	for {
		if err := txSub.runWithReconnect(ctx); err != nil {
			log.WithError(err).Error("Failed to handle subscription")
			contractAddress := common.HexToAddress(config.Default.Ethereum.ContractAddress)
			query := ethereum.FilterQuery{
				Addresses: []common.Address{contractAddress},
			}
			txSub.txSub, err = txSub.ethClient.SubscribeFilterLogs(context.Background(), query, txSub.txLogs)
			if err != nil {
				log.WithError(err).Error("Failed to subscribe to contract logs")
			}
		}
		log.Info("Reconnecting...")
		time.Sleep(5 * time.Second) // Delay before attempting to reconnect
	}
}

// runWithReconnect manages subscription and reconnection logic
func (txSub *TxSubscriber) runWithReconnect(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			log.Info("Detector has stopped")
			return nil
		case err := <-txSub.txSub.Err():
			log.WithError(err).Error("Transaction subscription error")
			return err // Exit to trigger reconnection
		case txLog := <-txSub.txLogs:
			notificationType, isSuspiciousTx := isSuspiciousTransaction(txSub.ethClient, txLog)
			if isSuspiciousTx {
				sendNotification(notificationType, txLog)
			}
		}
	}
}

func isSuspiciousTransaction(ethClient *ethclient.Client, txLog types.Log) (string, bool) {
	// Analyze the log for suspicious activity (large withdrawals or multiple withdrawals)
	// Get the transaction details
	tx, _, err := ethClient.TransactionByHash(context.Background(), txLog.TxHash)
	if err != nil {
		log.WithError(err).Error("Failed to get transaction")
		return "", false
	}
	// Get the transaction value in Wei
	txValueWei := tx.Value()
	thresholdWei := new(big.Int).SetUint64(100000000000000) // 0.0001 Ether

	// Check if transaction value exceeds threshold
	if txValueWei.Cmp(thresholdWei) > 0 {
		log.WithFields(log.Fields{
			"tx_value": txValueWei.String(),
		}).Info("Suspicious transaction: Value exceeds 0.0001 Ether")
		return notificationType1, true
	}

	// Get the address from which the transaction originated
	fromAddress := common.HexToAddress(txLog.Topics[0].Hex()) // Extract the sender's address from the log

	// Check for multiple withdrawals within a short period
	if checkMultipleTxs(fromAddress) {
		log.WithFields(
			log.Fields{
				"from_address": fromAddress,
			}).Info("Suspicious transaction: Multiple withdrawals detected in a short time")
		return notificationType2, true
	}

	return "", false
}

// Function to check if an address has made multiple txs within a short time frame
func checkMultipleTxs(fromAddress common.Address) bool {
	mu.Lock()
	defer mu.Unlock()

	now := time.Now()

	// Get the tx history for the address
	history := txHistory[fromAddress]

	// Filter the history to remove timestamps older than the time window
	recentTx := []time.Time{}
	for _, timestamp := range history {
		if now.Sub(timestamp) <= timeWindow {
			recentTx = append(recentTx, timestamp)
		}
	}

	// Update the history with only recent tx
	txHistory[fromAddress] = recentTx
	// Check if the number of recent tx exceeds the allowed limit
	if len(recentTx) >= maxTxs {
		return true
	}

	// Add the current tx timestamp to the history
	txHistory[fromAddress] = append(txHistory[fromAddress], now)

	return false
}

func sendNotification(notificationType string, ethLog types.Log) {
	m := gomail.NewMessage()
	m.SetHeader(emailFromHeader, config.Default.Notification.EmailFrom)
	m.SetHeader(emailToHeader, config.Default.Notification.EmailTo)
	m.SetHeader(emailSubjectHeader, emailSubjectContent)
	m.SetBody(emailBodyType,
		fmt.Sprintf("Suspicious transaction detected: \nType: %s\nAddress: %s\nBlockNumber: %d\nTx Hash: %s",
			notificationType, ethLog.Address, ethLog.BlockNumber, ethLog.TxHash))

	d := gomail.NewDialer(config.Default.Notification.SMTPDomain,
		config.Default.Notification.SMTPPort,
		config.Default.Notification.SMTPUser,
		config.Default.Notification.SMTPKey)
	if err := d.DialAndSend(m); err != nil {
		log.WithError(err).Error("Failed to send email")
	}
	log.WithFields(
		log.Fields{
			"mail_to": config.Default.Notification.EmailTo,
			"subject": emailSubjectContent,
		}).Info("The notification has sent")
}
