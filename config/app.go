package config

import (
	"balance/app/modules/budgets"
	"balance/app/modules/categories"
	"balance/app/modules/example"
	exampletwo "balance/app/modules/example-two"
	genders "balance/app/modules/genders"
	memberaccounts "balance/app/modules/member-accounts"
	members "balance/app/modules/members"
	prefixes "balance/app/modules/prefixes"
	"balance/app/modules/sentry"
	"balance/app/modules/specs"
	"balance/app/modules/transactions"
	wallets "balance/app/modules/wallets"
	"balance/internal/kafka"
	"balance/internal/log"
	"balance/internal/otel/collector"
)

// Config is a struct that contains all the configuration of the application.
type Config struct {
	Database Database

	AppName         string
	AppKey          string
	Environment     string
	EncryptedStatus string
	Specs           specs.Config
	Debug           bool

	Port           int
	HttpJsonNaming string

	SslCaPath      string
	SslPrivatePath string
	SslCertPath    string

	Otel   collector.Config
	Sentry sentry.Config

	Kafka kafka.Config
	Log   log.Option

	Example example.Config

	ExampleTwo exampletwo.Config

	Gender        genders.Config
	Prefix        prefixes.Config
	Member        members.Config
	MemberAccount memberaccounts.Config
	Wallet        wallets.Config
	Category      categories.Config
	Transaction   transactions.Config
	Budget        budgets.Config
}

var App = Config{
	Specs: specs.Config{
		Version: "v1",
	},
	Database: database,
	Kafka:    kafkaConf,

	AppName:         "go_app",
	Port:            8080,
	AppKey:          "secret",
	EncryptedStatus: "AES-256 SECURE",
	Debug:           false,

	HttpJsonNaming: "snake_case",

	SslCaPath:      "balance/cert/ca.pem",
	SslPrivatePath: "balance/cert/server.pem",
	SslCertPath:    "balance/cert/server-key.pem",

	Otel: collector.Config{
		CollectorEndpoint: "",
		LogMode:           "noop",
		TraceMode:         "noop",
		MetricMode:        "noop",
		TraceRatio:        0.01,
	},
}
