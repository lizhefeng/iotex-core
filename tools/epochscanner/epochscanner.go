package main

import (
	"context"
	"flag"
	"strconv"
	"time"

	"github.com/tealeg/xlsx"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/iotexproject/iotex-core/pkg/log"
	"github.com/iotexproject/iotex-core/pkg/util/byteutil"
	"github.com/iotexproject/iotex-core/protogen/iotexapi"
	"github.com/iotexproject/iotex-core/state"
)

func main() {
	var startEpoch int
	var endEpoch int
	flag.IntVar(&startEpoch, "start-epoch", 30, "start epoch number")
	flag.IntVar(&endEpoch, "end-epoch", 140, "end epoch number")
	flag.Parse()

	grpcAddr := "api.iotex.one:80"

	grpcctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	conn, err := grpc.DialContext(grpcctx, grpcAddr, grpc.WithBlock(), grpc.WithInsecure())
	if err != nil {
		log.L().Error("Failed to connect to API server.")
	}

	client := iotexapi.NewAPIServiceClient(conn)

	participation := make(map[string]uint64)
	topDelegates := make(map[string]uint64)
	for epochNum := startEpoch; epochNum <= endEpoch; epochNum++ {
		epochHeight := (epochNum-1)*24*15 + 1
		res, err := client.ReadState(context.Background(), &iotexapi.ReadStateRequest{
			ProtocolID: []byte("poll"),
			MethodName: []byte("ActiveConsensusBlockProducersByHeight"),
			Arguments:  [][]byte{byteutil.Uint64ToBytes(uint64(epochHeight))},
		})
		if err != nil {
			log.L().Fatal("Failed to read active block producers", zap.Error(err))
		}
		var activeBlockProducers state.CandidateList
		if err := activeBlockProducers.Deserialize(res.Data); err != nil {
			log.L().Fatal("Failed to deserialize active block producers")
		}
		for _, bp := range activeBlockProducers {
			participation[bp.Address]++
		}

		res, err = client.ReadState(context.Background(), &iotexapi.ReadStateRequest{
			ProtocolID: []byte("poll"),
			MethodName: []byte("ConsensusBlockProducersByHeight"),
			Arguments:  [][]byte{byteutil.Uint64ToBytes(uint64(epochHeight))},
		})
		if err != nil {
			log.L().Fatal("Failed to read block producers", zap.Error(err))
		}
		var blockProducers state.CandidateList
		if err := blockProducers.Deserialize(res.Data); err != nil {
			log.L().Fatal("Failed to deserialize block producers")
		}
		for _, bp := range blockProducers {
			topDelegates[bp.Address]++
		}
	}

	file := xlsx.NewFile()
	sheet, err := file.AddSheet("sheet1")
	if err != nil {
		log.L().Fatal("Failed to add sheet", zap.Error(err))
	}
	row := sheet.AddRow()
	cell1 := row.AddCell()
	cell1.Value = "Block Producer"
	cell2 := row.AddCell()
	cell2.Value = "Number of Top-36 Epochs"
	cell3 := row.AddCell()
	cell3.Value = "Number of Active Epochs"
	cell4 := row.AddCell()
	cell4.Value = "Total Number of Epochs"
	for bp, topCount := range topDelegates {
		row := sheet.AddRow()
		cell1 = row.AddCell()
		cell1.Value = bp
		cell2 = row.AddCell()
		cell2.Value = strconv.Itoa(int(topCount))
		cell3 = row.AddCell()
		cell3.Value = "0"
		if activeCount, ok := participation[bp]; ok {
			cell3.Value = strconv.Itoa(int(activeCount))
		}
		cell4 = row.AddCell()
		cell4.Value = strconv.Itoa(endEpoch - startEpoch + 1)
	}

	for _, col := range sheet.Cols {
		col.Width = 40
	}
	if err := file.Save("epochscanner.xlsx"); err != nil {
		log.L().Fatal("Failed to save excel sheet")
	}
}
