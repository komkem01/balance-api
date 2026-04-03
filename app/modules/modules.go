package modules

import (
	"log/slog"
	"sync"

	"balance/app/modules/budgets"
	"balance/app/modules/categories"
	"balance/app/modules/entities"
	"balance/app/modules/example"
	"balance/app/modules/genders"
	memberaccounts "balance/app/modules/member-accounts"
	"balance/app/modules/members"
	"balance/app/modules/prefixes"
	"balance/app/modules/sentry"
	"balance/app/modules/specs"
	"balance/app/modules/storage"
	"balance/app/modules/transactions"
	"balance/app/modules/wallets"
	"balance/internal/config"
	"balance/internal/database"
	"balance/internal/log"
	"balance/internal/otel/collector"

	exampletwo "balance/app/modules/example-two"

	appConf "balance/config"
	// "balance/app/modules/kafka"
)

type Modules struct {
	Conf    *config.Module[appConf.Config]
	Specs   *specs.Module
	Log     *log.Module
	OTEL    *collector.Module
	Sentry  *sentry.Module
	DB      *database.DatabaseModule
	ENT     *entities.Module
	Storage *storage.Module
	// Kafka *kafka.Module
	Example       *example.Module
	Example2      *exampletwo.Module
	Gender        *genders.Module
	Prefix        *prefixes.Module
	Member        *members.Module
	MemberAccount *memberaccounts.Module
	Wallet        *wallets.Module
	Category      *categories.Module
	Transaction   *transactions.Module
	Budget        *budgets.Module
}

func modulesInit() {
	confMod := config.New(&appConf.App)
	specsMod := specs.New(config.Conf[specs.Config](confMod.Svc))
	conf := confMod.Svc.Config()

	logMod := log.New(config.Conf[log.Option](confMod.Svc))
	otel := collector.New(config.Conf[collector.Config](confMod.Svc))
	log := log.With(slog.String("module", "modules"))

	sentryMod := sentry.New(config.Conf[sentry.Config](confMod.Svc))

	db := database.New(conf.Database.Sql)
	entitiesMod := entities.New(db.Svc.DB())
	exampleMod := example.New(config.Conf[example.Config](confMod.Svc), entitiesMod.Svc)
	exampleMod2 := exampletwo.New(config.Conf[exampletwo.Config](confMod.Svc), entitiesMod.Svc)
	genderMod := genders.New(config.Conf[genders.Config](confMod.Svc), entitiesMod.Svc)
	prefixMod := prefixes.New(config.Conf[prefixes.Config](confMod.Svc), entitiesMod.Svc)
	memberMod := members.New(config.Conf[members.Config](confMod.Svc), entitiesMod.Svc)
	memberAccountMod := memberaccounts.New(config.Conf[memberaccounts.Config](confMod.Svc), entitiesMod.Svc)
	walletMod := wallets.New(config.Conf[wallets.Config](confMod.Svc), entitiesMod.Svc)
	categoryMod := categories.New(config.Conf[categories.Config](confMod.Svc), entitiesMod.Svc)
	storageMod := storage.New(config.Conf[storage.Config](confMod.Svc))
	transactionMod := transactions.New(config.Conf[transactions.Config](confMod.Svc), entitiesMod.Svc, storageMod.Svc)
	budgetMod := budgets.New(config.Conf[budgets.Config](confMod.Svc), entitiesMod.Svc)
	// kafka := kafka.New(&conf.Kafka)
	mod = &Modules{
		Conf:          confMod,
		Specs:         specsMod,
		Log:           logMod,
		OTEL:          otel,
		Sentry:        sentryMod,
		DB:            db,
		ENT:           entitiesMod,
		Storage:       storageMod,
		Example:       exampleMod,
		Example2:      exampleMod2,
		Gender:        genderMod,
		Prefix:        prefixMod,
		Member:        memberMod,
		MemberAccount: memberAccountMod,
		Wallet:        walletMod,
		Category:      categoryMod,
		Transaction:   transactionMod,
		Budget:        budgetMod,
	}

	log.Infof("all modules initialized")
}

var (
	once sync.Once
	mod  *Modules
)

func Get() *Modules {
	once.Do(modulesInit)

	return mod
}
