package main

import (
	"context"
	"time"
	"flag"
	"strconv"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"github.com/tealeg/xlsx"

	"github.com/iotexproject/iotex-core/pkg/log"
	"github.com/iotexproject/iotex-core/protogen/iotexapi"
)

func main() {
	var startHeight int
	var endHeight int
	var scanWindow int

	flag.IntVar(&startHeight, "start-height", 1, "start height")
	flag.IntVar(&endHeight, "end-height", 1, "end height")
	flag.IntVar(&scanWindow, "window", 10000, "scan window for getting block metadata")
	flag.Parse()

	if startHeight > endHeight {
		log.L().Fatal("start height cannot be greater than end height")
	}

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

	if endHeight == 1 || uint64(endHeight) > chainMeta.Height {
		endHeight = int(chainMeta.Height)
	}

	production := make(map[string]uint64)
	start := startHeight
	for start+scanWindow-1 <= endHeight {
		getBlockMetasRequest := &iotexapi.GetBlockMetasRequest{
			Lookup: &iotexapi.GetBlockMetasRequest_ByIndex{
				ByIndex: &iotexapi.GetBlockMetasByIndexRequest{
					Start: uint64(start),
					Count: uint64(scanWindow),
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
		start += scanWindow
	}
	if start <= endHeight {
		getBlockMetasRequest := &iotexapi.GetBlockMetasRequest{
			Lookup: &iotexapi.GetBlockMetasRequest_ByIndex{
				ByIndex: &iotexapi.GetBlockMetasByIndexRequest{
					Start: uint64(start),
					Count: uint64(endHeight) - uint64(start) + 1,
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
	producerCount := len(production)
	totalNumBlks := endHeight - startHeight + 1
	log.L().Info("Block Production Summary", zap.Int("number of producers", producerCount), zap.Int("Average Production", totalNumBlks /producerCount))

	if err := writeExcel("rehearsalbpstat.xlsx", production, totalNumBlks); err != nil {
		log.L().Fatal("Failed to write block producer status to excel form.", zap.Error(err))
	}
}

func writeExcel(fileName string, production map[string]uint64, totalNumBlks int) error {
	file := xlsx.NewFile()
	sheet, err := file.AddSheet("sheet1")
	if err != nil {
		return err
	}
	row := sheet.AddRow()
	cell1 := row.AddCell()
	cell1.Value = "Block Producer"
	cell2 := row.AddCell()
	cell2.Value = "Number of Productions"
	cell3 := row.AddCell()
	cell3.Value = "Total Number of Blocks"
	for bp, count := range production {
		row := sheet.AddRow()
		cell1 = row.AddCell()
		cell1.Value = bp
		cell2 = row.AddCell()
		cell2.Value = strconv.Itoa(int(count))
		cell3 = row.AddCell()
		cell3.Value = strconv.Itoa(totalNumBlks)
	}

	for _, col := range sheet.Cols {
		col.Width = 40
	}
	return file.Save(fileName)
}
