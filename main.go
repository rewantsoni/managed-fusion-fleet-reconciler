package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/red-hat-storage/managed-fusion-fleet-reconciler/pkg/cloud"
	"github.com/red-hat-storage/managed-fusion-fleet-reconciler/pkg/types"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/red-hat-storage/managed-fusion-fleet-reconciler/pkg/db"
	"github.com/red-hat-storage/managed-fusion-fleet-reconciler/pkg/forman"
	"github.com/red-hat-storage/managed-fusion-fleet-reconciler/pkg/reconciler"

	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

type config struct {
	DB struct {
		Host     string            `yaml:"host"`
		Port     int               `yaml:"port"`
		User     string            `yaml:"user"`
		Password string            `yaml:"password"`
		Name     string            `yaml:"name"`
		Tables   map[string]string `yaml:"tables"`
	} `yaml:"db"`
	Reconcile struct {
		Concurrency int `yaml:"concurrency"`
	} `yaml:"reconcile"`
	AWS struct {
		AccessKey           string
		SecretKey           string
		Region              string
		ProviderTemplateURL string
	}
}

func loadAndValidateConfig(filePath string) (*config, error) {
	// parse configuration from yaml file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	cfg := &config{}
	if err := yaml.Unmarshal([]byte(data), cfg); err != nil {
		return nil, err
	}
	if cfg.DB.Host == "" {
		return nil, fmt.Errorf("config.db.host is not set")
	}
	if cfg.DB.Port == 0 {
		return nil, fmt.Errorf("config.db.port is not set")
	}
	if cfg.DB.User == "" {
		return nil, fmt.Errorf("config.db.user is not set")
	}
	if cfg.DB.Password == "" {
		return nil, fmt.Errorf("config.db.password is not set")
	}
	if cfg.DB.Name == "" {
		return nil, fmt.Errorf("config.db.name is not set")
	}
	if len(cfg.DB.Tables) == 0 {
		return nil, fmt.Errorf("config.db.tables is not set")
	}
	if cfg.Reconcile.Concurrency == 0 {
		return nil, fmt.Errorf("config.reconcile.concurrency is not set")
	}
	return cfg, nil
}

const configFileEnvVarName = "FLEET_RECONCILER_CONFIG"

func main() {
	//logger := logr.Logger{}
	ctx := context.Background()
	configFilePath := os.Getenv(configFileEnvVarName)
	if configFilePath == "" {
		log.Fatalf("%q environment variable not set", configFileEnvVarName)
	}
	// parse configuration from yaml file
	conf, err := loadAndValidateConfig(configFilePath)
	if err != nil {
		log.Fatal(err)
	}

	zapConfig := zap.Config{
		Level:            zap.NewAtomicLevel(),
		Encoding:         "json",
		EncoderConfig:    zap.NewProductionEncoderConfig(),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}
	logger, err := zapConfig.Build()
	if err != nil {
		log.Fatal(err)
	}

	connString := db.GetConnectionString(conf.DB.Host, conf.DB.User, conf.DB.Password, conf.DB.Name, conf.DB.Port)

	dbClient, err := db.NewClient(ctx, connString, conf.DB.Tables)
	if err != nil {
		logger.Fatal(err.Error())
	}
	defer dbClient.Close(ctx)

	awsconf, err := awsConfig.LoadDefaultConfig(context.TODO(),
		awsConfig.WithRegion(conf.AWS.Region),
		awsConfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(conf.AWS.AccessKey, conf.AWS.SecretKey, "")),
	)
	if err != nil {
		log.Fatal(err)
	}
	awsProvider, err := cloud.NewCloudProvider(types.AWSCloudProvider, logger)
	if err != nil {
		logger.Fatal(err.Error())
	}

	reqChan := forman.GoForman(logger, conf.Reconcile.Concurrency,
		func(req forman.Request) forman.Result {
			return reconciler.Reconcile(logger, dbClient, req)
		},
	)

	if err := dbClient.OnProvider(ctx, logger, true, func(provideName string) {
		req := forman.Request{}
		req.Name = provideName
		reqChan <- req
	}); err != nil {
		logger.Fatal(err.Error())
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// wait for signal to shutdown
	<-sigChan
	logger.Info("Received signal, shutting down")
}
