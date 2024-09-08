package handlers

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"github.com/rocky2015aaa/ethdefender/internal/config"
	"github.com/rocky2015aaa/ethdefender/internal/repository"
	"github.com/rocky2015aaa/ethdefender/internal/services/reporter/controllers/reports"
	httplib "github.com/rocky2015aaa/ethdefender/pkg/http"
)

type ReportsAPI struct {
	controller *reports.Controller
}

func NewReportsAPI(db repository.Storage) API {
	ethClient, err := ethclient.Dial(config.Default.Ethereum.SubscriptionUrl)
	if err != nil {
		log.WithError(err).Fatal("rpc client init error")
	}

	query := ethereum.FilterQuery{
		Addresses: []common.Address{common.HexToAddress(config.Default.Ethereum.ContractAddress)},
	}

	txLogs := make(chan types.Log)
	ethSub, err := ethClient.SubscribeFilterLogs(context.Background(), query, txLogs)
	if err != nil {
		log.WithError(err).Fatal("ethereum subscription init error")
	}

	return &ReportsAPI{controller: reports.NewController(db, ethClient, txLogs, ethSub)}
}

func (api *ReportsAPI) CreateReports(c context.Context) {
	// Listen for contract logs (such as withdrawals)
	for {
		select {
		case err := <-api.controller.EthSub.Err():
			log.WithError(err).Error("Ethereum subscription error")
		case vLog := <-api.controller.TxLogs:
			err := api.controller.CreateReports(c, vLog)
			if err != nil {
				log.WithError(err).WithField("tx log", vLog).Error("Create reports error")
			} else {
				log.WithFields(
					log.Fields{
						"tx": vLog.TxHash,
					}).Info("Created a transaction report")
			}
		}
	}
}

// PostSwapOrder godoc
// @Title        Get Transaction Report
// @Description  Get a transaction report
// @Tags         Transaction
// @Accept       json
// @Produce      json
// @Success      200  {object}  reports.Resp
// @Failure      400  {object}  http.Error
// @Failure      500  {object}  http.Error
// @Router       /api/v1/report/transaction [get]
func (api *ReportsAPI) GetTransactionReport(c *gin.Context) {
	response, err := api.controller.GetTransactionReport(c.Request.Context())
	if err != nil {
		c.JSON(err.GetStatusCode(), err)
		return
	}
	c.JSON(response.StatusCode, response)
}

// PostSwapOrder godoc
// @Title        Get Pause Report
// @Description  Get a pause report
// @Tags         Pause
// @Accept       json
// @Produce      json
// @Success      200  {object}  reports.Resp
// @Failure      400  {object}  http.Error
// @Failure      500  {object}  http.Error
// @Router       /api/v1/report/pause [get]
func (api *ReportsAPI) GetPauseReport(c *gin.Context) {
	response, err := api.controller.GetPauseReport(c.Request.Context())
	if err != nil {
		c.JSON(err.GetStatusCode(), err)
		return
	}
	c.JSON(response.StatusCode, response)
}

// PostSwapOrder godoc
// @Title        Get Slither Report
// @Description  Get a slither report
// @Tags         Slither
// @Accept       json
// @Produce      json
// @Success      200  {object}  reports.Resp
// @Failure      400  {object}  http.Error
// @Failure      500  {object}  http.Error
// @Router       /api/v1/report/slither [get]
func (api *ReportsAPI) GetSlitherReport(c *gin.Context) {
	response, err := api.controller.GetSlitherReport(c.Request.Context())
	if err != nil {
		c.JSON(err.GetStatusCode(), err)
		return
	}
	c.JSON(response.StatusCode, response)
}

// PostSwapOrder godoc
// @Title        Create Slither Report
// @Description  Creates or Updates a slither order
// @Tags         Slither
// @Accept       json
// @Produce      json
// @Param        contact_file  formData  file  true  "Contract file to upload"
// @Param        contact_name  formData  string  true  "Contract name"
// @Success      200  {object}  reports.Resp
// @Failure      400  {object}  http.Error
// @Failure      500  {object}  http.Error
// @Router       /api/v1/report/slither [post]
func (api *ReportsAPI) PostSlitherReport(c *gin.Context) {
	// Parse the form data with a maximum memory of 10MB
	err := c.Request.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		c.JSON(http.StatusBadRequest, httplib.ErrBadRequest)
		return
	}
	contractName := c.Request.FormValue("contract_name")
	if contractName == "" {
		c.JSON(http.StatusBadRequest, httplib.ErrBadRequest)
		return
	}
	// Retrieve the file from the form data
	file, handler, err := c.Request.FormFile("contract_file")
	if err != nil {
		c.JSON(http.StatusInternalServerError, httplib.ErrInternalServer)
		return
	}
	defer file.Close()
	// Create the destination file where the uploaded file will be saved
	targetSmartContract := filepath.Join(config.Default.App.AssetsPath, "contract", handler.Filename)
	dst, err := os.Create(targetSmartContract)
	if err != nil {
		c.JSON(http.StatusInternalServerError, httplib.ErrInternalServer)
		return
	}
	defer dst.Close()
	// Copy the uploaded file to the destination
	_, err = io.Copy(dst, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, httplib.ErrInternalServer)
		return
	}
	reportLocation := filepath.Join(config.Default.App.AssetsPath, "contract/report.json")
	openzeppelinLocation := fmt.Sprintf("@openzeppelin=%s", filepath.Join(config.Default.App.AssetsPath, "contract/@openzeppelin"))
	cmd := exec.Command("slither", dst.Name(),
		"--json", reportLocation, "--solc-remaps",
		openzeppelinLocation, "--filter-paths", "@openzeppelin")
	//cmd.Stderr = os.Stderr
	err = cmd.Run()
	_, reportLocationErr := os.Stat(reportLocation)
	if os.IsNotExist(reportLocationErr) {
		if err != nil {
			c.JSON(http.StatusInternalServerError, httplib.ErrInternalServer)
			return
		}
	}
	reportData, err := os.ReadFile(reportLocation)
	if err != nil {
		c.JSON(http.StatusInternalServerError, httplib.ErrInternalServer)
		return
	}
	response, httpErr := api.controller.PostSlitherReport(c.Request.Context(), contractName, reportData)
	if httpErr != nil {
		c.JSON(httpErr.GetStatusCode(), err)
		return
	}
	_ = os.Remove(targetSmartContract)
	_ = os.Remove(reportLocation)
	c.JSON(response.StatusCode, response)
}
