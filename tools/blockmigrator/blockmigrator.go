package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"go.uber.org/zap"

	"github.com/iotexproject/iotex-core/blockchain/blockdao"
	"github.com/iotexproject/iotex-core/config"
	"github.com/iotexproject/iotex-core/db"
	"github.com/iotexproject/iotex-core/pkg/log"
)

// migrateHeight is the blockchain height being migrated to
var migrateHeight int

func init() {
	flag.IntVar(&migrateHeight, "migrate-height", 0, "blockchain migration height")
	flag.Usage = func() {
		_, _ = fmt.Fprintf(os.Stderr,
			"usage: migrate -current-path=[string]\n -new-path=[string]\n -recovery-height=[int]\n")
		flag.PrintDefaults()
		os.Exit(2)
	}
	flag.Parse()
}

func main() {
	cfg, err := config.New()
	if err != nil {
		log.L().Fatal("Failed to new config", zap.Error(err))
	}
	if err := os.Rename(cfg.Chain.ChainDBPath, cfg.Chain.ChainDBPath+".old"); err != nil {
		log.L().Fatal("Failed to rename old chain db", zap.Error(err))
	}

	currentPath := cfg.Chain.ChainDBPath + ".old"
	cfg.DB.DbPath = currentPath
	_, gateway := cfg.Plugins[config.GatewayPlugin]
	currentDAO := blockdao.NewBlockDAO(
		db.NewBoltDB(cfg.DB),
		gateway && !cfg.Chain.EnableAsyncIndexWrite,
		cfg.Chain.CompressBlock,
		cfg.Chain.MaxCacheSize,
		cfg.DB,
	)

	cfg.DB.DbPath = cfg.Chain.ChainDBPath
	newDAO := blockdao.NewBlockDAO(db.NewBoltDB(cfg.DB),
		gateway && !cfg.Chain.EnableAsyncIndexWrite,
		cfg.Chain.CompressBlock,
		cfg.Chain.MaxCacheSize,
		cfg.DB,
	)

	ctx := context.Background()
	if err := currentDAO.Start(ctx); err != nil {
		log.L().Fatal("Failed to start the current block DAO")
	}
	if err := newDAO.Start(ctx); err != nil {
		log.L().Fatal("Failed to start the new block DAO")
	}
	defer func() {
		if err := currentDAO.Stop(ctx); err != nil {
			log.L().Fatal("Failed to stop the current block DAO")
		}
		if err := newDAO.Stop(ctx); err != nil {
			log.L().Fatal("Failed to stop the new block DAO")
		}
	}()

	tipHeight, err := currentDAO.GetBlockchainHeight()
	if err != nil {
		log.L().Fatal("Failed to get blockchain tip height")
	}
	if migrateHeight < 0 || uint64(migrateHeight) > tipHeight {
		log.L().Fatal("Invalid block migration height")
	}

	for i := uint64(1); i <= uint64(migrateHeight); i++ {
		hash, err := currentDAO.GetBlockHash(i)
		if err != nil {
			log.S().Fatalf("Failed to get block hash on height %d", i, zap.Error(err))
		}
		blk, err := currentDAO.GetBlock(hash)
		if err != nil {
			log.S().Fatalf("Failed to get block on height %d", i, zap.Error(err))
		}
		if err := newDAO.PutBlock(blk); err != nil {
			log.S().Fatalf("Failed to migrate block on height %d", i, zap.Error(err))
		}
	}
}
