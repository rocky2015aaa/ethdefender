package reports

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/rocky2015aaa/ethdefender/internal/config"
	"github.com/rocky2015aaa/ethdefender/internal/repository"
	"github.com/rocky2015aaa/ethdefender/internal/repository/models"
	"github.com/rocky2015aaa/ethdefender/internal/services"
	httplib "github.com/rocky2015aaa/ethdefender/pkg/http"
)

const (
	pauseEventName  = "Paused"
	pauseStatusName = "paused"
)

type Controller struct {
	DB        repository.Storage
	EthClient *ethclient.Client
	TxLogs    chan types.Log
	EthSub    ethereum.Subscription
}

func NewController(db repository.Storage, ethClient *ethclient.Client, txLogs chan types.Log, ethSub ethereum.Subscription) *Controller {
	return &Controller{
		DB:        db,
		EthClient: ethClient,
		TxLogs:    txLogs,
		EthSub:    ethSub,
	}
}

func (i *Controller) CreateReports(ctx context.Context, vLog types.Log) *httplib.Error {
	// Get the transaction receipt
	receipt, err := i.EthClient.TransactionReceipt(context.Background(), vLog.TxHash)
	if err != nil {
		return httplib.NewError(http.StatusInternalServerError, err.Error())
	}

	// Query the block details for the timestamp
	block, err := i.EthClient.BlockByNumber(context.Background(), receipt.BlockNumber)
	if err != nil {
		return httplib.NewError(http.StatusInternalServerError, err.Error())
	}

	executionTime := time.Unix(int64(block.Time()), 0).UTC()
	// Format time.Time object to string
	executionTimeStr := executionTime.Format(time.RFC3339)

	gasUsed := receipt.GasUsed // Gas used in the transaction

	// Convert data to JSON
	transactionReport := models.TransactionReport{
		TransactionID: vLog.TxHash.Hex(),
		ExecutionTime: executionTimeStr,
		GasUsed:       gasUsed,
	}

	err = i.DB.CreateTransactionReport(ctx, &transactionReport)
	if err != nil {
		return httplib.NewError(http.StatusInternalServerError, err.Error())
	}
	if vLog.Topics[0] == services.ParsedABI.Events[pauseEventName].ID {
		log.WithField("event", pauseEventName).Info("Paused event detected")

		isPaused, err := getPauseStatus(services.ParsedABI, i.EthClient, common.HexToAddress(config.Default.Ethereum.ContractAddress))
		if err != nil {
			return httplib.NewError(http.StatusInternalServerError, err.Error())
		}

		pauseReport := models.PauseReport{
			EventType:   pauseEventName,
			PauseStatus: isPaused,
		}

		err = i.DB.CreatePauseReport(ctx, &pauseReport)
		if err != nil {
			return httplib.NewError(http.StatusInternalServerError, err.Error())
		}
	}

	return nil
}

func getPauseStatus(parsedABI abi.ABI, ethClient *ethclient.Client, contractAddress common.Address) (bool, error) {
	var isPaused bool
	callData, err := parsedABI.Pack(pauseStatusName)
	if err != nil {
		return isPaused, err
	}
	msg := ethereum.CallMsg{
		To:   &contractAddress,
		Data: callData,
	}
	result, err := ethClient.CallContract(context.Background(), msg, nil)
	if err != nil {
		return isPaused, err
	}
	err = parsedABI.UnpackIntoInterface(&isPaused, pauseStatusName, result)
	if err != nil {
		return isPaused, err
	}
	return isPaused, nil
}

func (i *Controller) GetTransactionReport(ctx context.Context) (*Resp, *httplib.Error) {
	reportData, err := i.DB.GetTransactionReport(ctx)
	if err != nil {
		log.WithError(err).WithField("report", reportData).Error("Get slither report error")
		return nil, httplib.ErrInternalServer
	}
	return &Resp{Data: reportData, StatusCode: http.StatusOK, Message: "OK"}, nil
}

func (i *Controller) GetPauseReport(ctx context.Context) (*Resp, *httplib.Error) {
	reportData, err := i.DB.GetPauseReport(ctx)
	if err != nil {
		log.WithError(err).WithField("report", reportData).Error("Get slither report error")
		return nil, httplib.ErrInternalServer
	}
	return &Resp{Data: reportData, StatusCode: http.StatusOK, Message: "OK"}, nil
}

func (i *Controller) GetSlitherReport(ctx context.Context) (*Resp, *httplib.Error) {
	contractReportData, err := i.DB.GetSlitherReport(ctx)
	if err != nil {
		log.WithError(err).WithField("report", contractReportData).Error("Get slither report error")
		return nil, httplib.ErrInternalServer
	}
	contractReportResp := []*ContractReportResp{}
	for _, data := range contractReportData {
		var report interface{}
		err := json.Unmarshal(data.Report, &report)
		if err != nil {
			fmt.Println("Error unmarshalling JSON:", err)
			return nil, httplib.ErrInternalServer
		}
		reportResp := ContractReportResp{
			ContractName: data.ContractName,
			Report:       report,
			CreatedAt:    data.CreatedAt,
			UpdatedAt:    data.UpdatedAt,
		}
		contractReportResp = append(contractReportResp, &reportResp)
	}
	return &Resp{Data: contractReportResp, StatusCode: http.StatusOK, Message: "OK"}, nil
}

func (i *Controller) PostSlitherReport(ctx context.Context, contractName string, reportData []byte) (*Resp, *httplib.Error) {
	contractReport := models.ContractReport{
		ContractName: contractName,
		Report:       reportData,
	}
	err := i.DB.UpsertSlitherReport(ctx, &contractReport)
	if err != nil {
		log.WithError(err).WithField("report", contractReport).Error("Upsert slither report error")
		return nil, httplib.ErrInternalServer
	}
	return &Resp{Data: nil, StatusCode: http.StatusOK, Message: "OK"}, nil
}
