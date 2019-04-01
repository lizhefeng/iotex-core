package main

import (
	"context"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/iotexproject/iotex-core/pkg/log"
	"github.com/iotexproject/iotex-core/protogen/iotexapi"
)

func main() {
	grpcAddr := "api.iotex.one:80"

	grpcctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	conn, err := grpc.DialContext(grpcctx, grpcAddr, grpc.WithBlock(), grpc.WithInsecure())
	if err != nil {
		log.L().Error("Failed to connect to API server.")
	}

	client := iotexapi.NewAPIServiceClient(conn)

	chainMetaRes, err := client.GetChainMeta(context.Background(), &iotexapi.GetChainMetaRequest{})
	if err != nil {
		log.L().Fatal("Failed to get chain metadata.", zap.Error(err))
	}
	chainMeta := chainMetaRes.ChainMeta
	log.L().Info("Blockchain Metadata", zap.Uint64("current height", chainMeta.Height))

	production := make(map[string]uint64)
	interval := 10000
	start := 1
	for start+interval-1 <= int(chainMeta.Height) {
		getBlockMetasRequest := &iotexapi.GetBlockMetasRequest{
			Lookup: &iotexapi.GetBlockMetasRequest_ByIndex{
				ByIndex: &iotexapi.GetBlockMetasByIndexRequest{
					Start: uint64(start),
					Count: uint64(interval),
				},
			},
		}
		blockMetasRes, err := client.GetBlockMetas(context.Background(), getBlockMetasRequest)
		if err != nil {
			log.L().Fatal("Failed to get block metadata.", zap.Error(err))
		}
		blkMetas := blockMetasRes.BlkMetas
		for _, blk := range blkMetas {
			production[blk.ProducerAddress]++
		}
		start += interval
	}
	if start <= int(chainMeta.Height) {
		getBlockMetasRequest := &iotexapi.GetBlockMetasRequest{
			Lookup: &iotexapi.GetBlockMetasRequest_ByIndex{
				ByIndex: &iotexapi.GetBlockMetasByIndexRequest{
					Start: uint64(start),
					Count: chainMeta.Height - uint64(start) + 1,
				},
			},
		}
		blockMetasRes, err := client.GetBlockMetas(context.Background(), getBlockMetasRequest)
		if err != nil {
			log.L().Fatal("Failed to get block metadata.", zap.Error(err))
		}
		blkMetas := blockMetasRes.BlkMetas
		for _, blk := range blkMetas {
			production[blk.ProducerAddress]++
		}
	}
	var producerCount uint64
	for bp, produce := range production {
		producerCount++
		log.L().Info(bp, zap.Uint64("produce", produce))
	}
	log.L().Info("Block Production Summary", zap.Uint64("number of producers", producerCount), zap.Uint64("Average Production", chainMeta.Height/producerCount))
}
