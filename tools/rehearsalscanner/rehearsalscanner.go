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

	delegateMap := map[string]string{
		"io1vlsmjs87jlk93624nppccfn24nk9nplu9uhu53": "cobo",
		"io1fra0fx6akacny9asewt7vyvggqq4rtq3eznccr": "iotexteam",
		"io1u5ff879gg2dw9vfpxr2tsmuaz07e2rea6gvl7s": "rosemary9",
		"io1kugyfdkdss7acy0y9x8fmfjd5qyxfcye3gxt99": "consensusnet",
		"io12wc9ra4la98yay4cqdav5mwxxuzwpt6hk23n3z": "iotexmainet",
		"io10reczcaelglh5xmkay65h9vw3e5dp82e8vw0rz": "metanyx",
		"io1rfeltdr5wgmhm4rdx9eun3pp68ahm7fq00wcmm": "pubxpayments",
		"io1n3llm0zzlles6pvpzugzajyzyjehztztxf2rff": "ratels",
		"io1xsx5n94kg2zv64r7tm8vyz9mh86amfak9ka9xx": "rosemary3",
		"io1wv5m0xyermvr2n0wjx2cjsqwyk863drdl5qfyn": "rosemary13",
		"io1qqaswtu7rcevahucpfxc0zxx088pylwtxnkrrl": "stanford cpc",
		"io1e3w03ulnrsxth2g0rgsq6v406fhdccsgfq3hz7": "iotask",
		"io1gfq9el2gnguus64ex3hu8ajd6e4yzk3f9cz5vx": "laomao",
		"io15fqav3tugm96ge7anckx0k4gukz5m4mqf0jpv3": "rosemary0",
		"io127ftn4ry6wgxdrw4hcd6gdwqlq70ujk98dvtw5": "rosemary12",
		"io108h7sa5sap44e244hz649zyk5y4rqzsvnpzxh5": "royalland",
		"io1sd5t5dwxrk2t50z8yl86n6ht8c99umt4u6rknl": "hashbuy",
		"io17cmrextyfeu4gddwd89g5qncedsnc553dhz7xa": "infstones",
		"io1gf08snppu2a2wfd50pjas2j6q2kcxjzqph3pep": "rosemary8",
		"io1qnec80ark9shjc6uzk45dhm8s50dpc27sjuver": "gamefantasy",
		"io1c3r4th3zrk4uhv83a9gr4gvn3y6pzaj6mc84ea": "rosemary6",
		"io12yxdwewry70gr9fs6fphyfaky9c7gurmzk8f4f": "rosemary11",
		"io1zy9q4myp3jezxdsv82z3f865ruamhqm7a6mggh": "capitalmu",
		"io14u5d66rt465ykm7t2847qllj0reml27q30kr75": "iosg",
		"io1nf0rvzgq3tqym6n3trttsrt7d4gqqsmqfzy0da": "iotexlab",
		"io159fv8mu9d5djk8u2t0flgw4yqmt6fg98uqjka8": "rosemary5",
		"io14vmhs9c75r2ptxdaqrtk0dz7skct30pxmt69d9": "rosemary7",
		"io1xj0u5n20tsqwxh5a3xdtmzuz9wasft0pqjrq8t": "thebottoken",
		"io13v3dhwds82hg0uc9l4puer00k93qagdh62j0mz": "wannodes",
		"io1wl83n3up9w8nedf30lnyxzple0gu5pme0dyrds": "whales",
		"io1f0rh94m3ctkwep3rlsswwq5vxwlntx4s574l3q": "yvalidator",
		"io1we5gqn4xeak9ycnu4l9lv0qfq3euapymnzffx6": "coingecko",
		"io1kr8c6krd7dhxaaqwdkr6erqgu4z0scug3drgja": "droute",
		"io17dm3tq9ezgs2uvwrqu8e004l5nqq33jm46r0mg": "iotexunion",
		"io1aqf30kqz5rqh6zn82c00j684p2h2t5cg30wm8t": "keysiotex",
		"io1ddjluttkzljqfgdtcz9eu3r3s832umy7ylx0ts": "longz",
		"io1x9kjkr0qv2fa7j4t2as8lrj223xxsqt4tl7xp7": "rosemary1",
		"io1vtm2zgn830pn6auc2cvnchgwdaefa9gr4z0s86": "rosemary4",
		"io1du4eq4f88n4wyc026l3gamjwetlgsg4jz7j884": "rosemary10",
		"io1lm85qfm24eqrc042spnplac9ran28thuh5048n": "iotexgeeks",
		"io1et7zkzc76m9twa4gn5xht3urt9mwj05qvdtj66": "iotxplorerio",
		"io1z4sxtefurklkyrfmmdtjmw4h8csnxlv9747hyd": "preangle",
		"io1ar5l5s268rtgzshltnqv88mua06ucm58dx678y": "rosemary",
	}

	if endHeight == 1 || uint64(endHeight) > chainMeta.Height {
		endHeight = int(chainMeta.Height)
	}

	production := make(map[string]uint64)
	start := startHeight
	for start <= endHeight {
		count := scanWindow
		if scanWindow > endHeight-start+1 {
			count = endHeight - start + 1
		}
		getBlockMetasRequest := &iotexapi.GetBlockMetasRequest{
			Lookup: &iotexapi.GetBlockMetasRequest_ByIndex{
				ByIndex: &iotexapi.GetBlockMetasByIndexRequest{
					Start: uint64(start),
					Count: uint64(count),
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
		start += count
	}
	producerCount := len(production)
	totalNumBlks := endHeight - startHeight + 1
	log.L().Info("Block Production Summary", zap.Int("number of producers", producerCount), zap.Int("Average Production", totalNumBlks/producerCount))

	if err := writeExcel("rehearsalbpstat.xlsx", production, delegateMap, totalNumBlks); err != nil {
		log.L().Fatal("Failed to write block producer status to excel form.", zap.Error(err))
	}
}

func writeExcel(fileName string, production map[string]uint64, delegateMap map[string]string, totalNumBlks int) error {
	file := xlsx.NewFile()
	sheet, err := file.AddSheet("sheet1")
	if err != nil {
		return err
	}
	row := sheet.AddRow()
	cell1 := row.AddCell()
	cell1.Value = "Block Producer"
	cell2 := row.AddCell()
	cell2.Value = "Alias"
	cell3 := row.AddCell()
	cell3.Value = "Number of Productions"
	cell4 := row.AddCell()
	cell4.Value = "Total Number of Blocks"
	for bp, count := range production {
		row := sheet.AddRow()
		cell1 = row.AddCell()
		cell1.Value = bp
		cell2 = row.AddCell()
		if alias, ok := delegateMap[bp]; ok {
			cell2.Value = alias
		}
		cell3 = row.AddCell()
		cell3.Value = strconv.Itoa(int(count))
		cell4 = row.AddCell()
		cell4.Value = strconv.Itoa(totalNumBlks)
	}

	for _, col := range sheet.Cols {
		col.Width = 40
	}
	return file.Save(fileName)
}
