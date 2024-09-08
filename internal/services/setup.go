package services

import (
	"context"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"

	"github.com/rocky2015aaa/ethdefender/internal/config"
	"github.com/rocky2015aaa/ethdefender/internal/repository/postgres"
	"github.com/rocky2015aaa/ethdefender/pkg/viper"
)

const (
	defaultConfigPath = "config.yml"
	defaultEnvPath    = ".env"

	EnvReportDBUri        = "REPORTER_DB_URI"
	AddedContractFunction = "getOwner"
)

var ParsedABI abi.ABI

func Setup() {
	InitConfig()
	InitLogging()
	InitContractABI()
}

func SetupWithDB() {
	Setup()
	InitDatabase()
}

func InitConfig() {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = defaultConfigPath
	}

	viper.Load(configPath, &config.Default)

	envPath := os.Getenv("ENV_PATH")
	if envPath == "" {
		envPath = defaultEnvPath
	}
	err := godotenv.Load(envPath)
	if err != nil {
		log.WithError(err).Fatal("env file loading error")
	}
}

func InitLogging() {
	logLevel, err := log.ParseLevel(config.Default.LogLevel)
	if err != nil {
		log.WithError(err).Fatal("Log level parsing error")
	}

	log.SetLevel(logLevel)
}

func InitDatabase() {
	db, err := postgres.New(os.Getenv(EnvReportDBUri), config.Default.Database.Log)
	if err != nil {
		log.WithError(err).Fatal("Database init error")
	}

	if err := postgres.Setup(db); err != nil {
		log.WithError(err).Fatal("Database setup error")
	}
}

func InitContractABI() {
	ethClient, err := ethclient.Dial(config.Default.Ethereum.SubscriptionUrl)
	if err != nil {
		log.Fatal(err)
	}
	// Address of the contract
	contractAddress := common.HexToAddress(config.Default.Ethereum.ContractAddress)
	updatedAbiData, err := os.ReadFile(config.Default.Ethereum.UpdatedContractABILocation)
	if err != nil {
		log.WithError(err).Fatal("Contract ABI read error")
	}
	updatedParsedABI, err := abi.JSON(strings.NewReader(string(updatedAbiData)))
	if err != nil {
		log.WithError(err).Fatal("Contract ABI conversion error")
	}
	// Call the getOwner function
	callData, err := updatedParsedABI.Pack(AddedContractFunction)
	if err != nil {
		log.WithError(err).Fatal("Contract ABI parsing error")
	}
	// Prepare the call message
	msg := ethereum.CallMsg{
		To:   &contractAddress,
		Data: callData,
	}
	// Call the contract function
	result, err := ethClient.CallContract(context.Background(), msg, nil)
	if err != nil {
		log.WithError(err).Println("Contract ABI test function call error. try to use an original ABI")
	}
	// If calling added function is successful, then keep using updated abi or not use original abi
	if len(result) == 0 {
		abiData, err := os.ReadFile(config.Default.Ethereum.ContractABILocation)
		if err != nil {
			log.WithError(err).Fatal("Contract ABI read error")
		}
		ParsedABI, err = abi.JSON(strings.NewReader(string(abiData)))
		if err != nil {
			log.WithError(err).Fatal("Contract ABI conversion error")
		}
	} else {
		ParsedABI = updatedParsedABI
	}
}
