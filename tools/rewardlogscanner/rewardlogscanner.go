package main

import (
	"context"
	"math/big"
	"sort"
	"strconv"

	"github.com/tealeg/xlsx"
	"go.uber.org/zap"
	"github.com/golang/protobuf/proto"

	"github.com/iotexproject/iotex-core/action/protocol/rolldpos"
	"github.com/iotexproject/iotex-core/config"
	"github.com/iotexproject/iotex-core/pkg/log"
	"github.com/iotexproject/iotex-core/server/itx"
	"github.com/iotexproject/iotex-core/action"
	"github.com/iotexproject/iotex-core/action/protocol/rewarding/rewardingpb"
	"github.com/iotexproject/iotex-core/pkg/enc"
	"github.com/iotexproject/iotex-core/protogen/iotextypes"
	"github.com/iotexproject/iotex-core/pkg/hash"
	"fmt"
)

const endEpochNum = 300

var (
	aliasMap = map[string]bool{
		"cobo":         true,
		"iotexteamandpubxpayments":    true,
		"consensusnet": true,
		"iotexmainnet": true,
		"metanyx":      true,
		"ratels":       true,
		"stanfordcpc":  true,
		"iotask":       true,
		"laomao":       true,
		"royalland":    true,
		"hashbuy":      true,
		"infstones":    true,
		"gamefantasy":  true,
		"capitalmu":    true,
		"iosg":         true,
		"iotexlab":     true,
		"thebottoken":  true,
		"wannodes":     true,
		"whales":       true,
		"yvalidator":   true,
		"coingecko":    true,
		"droute":       true,
		"iotexunion":   true,
		"keysiotex":    true,
		"longz":        true,
		"iotexgeeks":   true,
		"iotxplorerio": true,
		"preangle":     true,
	}

	rewardAddrMap = map[string]string{
		"io1vlsmjs87jlk93624nppccfn24nk9nplu9uhu53": "cobo",
		"io1ppq3s9fwnhqzd6vg0tt8y5qa3hyk8yy776xpwq": "iotexteamandpubxpayments",
		"io1p6xzyl5z4p5z43406l0rqypqdk2zgv56wf8wjp": "consensusnet",
		"io1mrm4gchhdwtp4ulyd8yqus9s6tgn5nygtswjxx": "iotexmainnet",
		"io1knfzqrzwlx9cx6t2trgcppn767vptlce0m2kje": "metanyx",
		"io1n3llm0zzlles6pvpzugzajyzyjehztztxf2rff": "ratels",
		"io1qqaswtu7rcevahucpfxc0zxx088pylwtxnkrrl": "stanfordcpc",
		"io1e3w03ulnrsxth2g0rgsq6v406fhdccsgfq3hz7": "iotask",
		"io1464yt9pmty64dlun95kvuvdtr2n3wf2cx4t9p2": "laomao",
		"io19ngaugnmgcshfsdzk56h6trp9gpsk9gkkl0h3j": "royalland",
		"io1sd5t5dwxrk2t50z8yl86n6ht8c99umt4u6rknl": "hashbuy",
		"io1w7enh0vv87uw9rt732z089nq6652dradk30vlu": "infstones",
		"io1yn6q5zs7942vfwjd6dmtexjysdfpcj2n27ljar": "gamefantasy",
		"io1yjjqtktt5yh54zcum4dfj2xh4r6xatq8p8nflv": "capitalmu",
		"io14u5d66rt465ykm7t2847qllj0reml27q30kr75": "iosg",
		"io1nf0rvzgq3tqym6n3trttsrt7d4gqqsmqfzy0da": "iotexlab",
		"io1c6a0ugwqg9p96l69fzh00xk5l4ffjrqfcv54ue": "thebottoken",
		"io1kjuludtup293j92kqhnnptge8hpg9lez40mych": "wannodes",
		"io1wl83n3up9w8nedf30lnyxzple0gu5pme0dyrds": "whales",
		"io1n23j989qrynzfg9h5vg9lhnsu0zuk7n6kqd0dm": "yvalidator",
		"io1dfv45vq9vx69uwz660nqaqt9q6pv5j2496dq27": "coingecko",
		"io1xwhjtczj3v3jx6fltd5h8y47zraq2vxkvltnrw": "droute",
		"io17dm3tq9ezgs2uvwrqu8e004l5nqq33jm46r0mg": "iotexunion",
		"io1f6vcjkudxnfdzv4v66sjljpych38h7387ynwvw": "keysiotex",
		"io1cuua32hupytntjvedh9egqprav349rldgjk2uq": "longz",
		"io1g2hsz6v68v4rqjjyqkywms605kepspejj3esxg": "iotexgeeks",
		"io1et7zkzc76m9twa4gn5xht3urt9mwj05qvdtj66": "iotxplorerio",
		"io13v4pnjhfhxdynpq0uhwv2lmu3a5uvwtclls63p": "preangle",
		"io1fsrkpzrt6juvkuermz6dd2ng3a044uda6ywu8l": "iotxplorerio",
	}
)

type epochInfo struct {
	blockReward      *big.Int
	epochReward      *big.Int
	foundationReward *big.Int
}

func main() {
	cfg := config.Default
	cfg.Genesis.NumSubEpochs = 15
	cfg.Consensus.Scheme = config.NOOPScheme
	cfg.DB.DbPath = cfg.Chain.ChainDBPath

	ctx := context.Background()
	svr, err := itx.NewServer(cfg)
	if err != nil {
		log.L().Fatal("Failed to create a new server.", zap.Error(err))
	}
	if err := svr.Start(ctx); err != nil {
		log.L().Fatal("Failed to start server.", zap.Error(err))
	}
	defer func() {
		if err := svr.Stop(ctx); err != nil {
			log.L().Fatal("Failed to stop server", zap.Error(err))
		}
	}()
	cs := svr.ChainService(cfg.Chain.ID)
	bc := cs.Blockchain()
	kvstore := bc.KVStore()

	registry := cs.Registry()
	p, ok := registry.Find(rolldpos.ProtocolID)
	if !ok {
		log.L().Fatal("Failed to find rolldpos protocol.", zap.Error(err))
	}
	rp, ok := p.(*rolldpos.Protocol)
	if !ok {
		log.L().Fatal("Failed to cast rolldpos protocol.")
	}

	endEpochNum := endEpochNum
	epochMeta := make(map[string]map[int]*epochInfo)
	for alias := range aliasMap {
		epochMeta[alias] = make(map[int]*epochInfo)
	}

	for epochNum := 1; epochNum <= int(endEpochNum); epochNum++ {
		fmt.Println(epochNum)

		epochStartHeight := rp.GetEpochHeight(uint64(epochNum))
		epochEndHeight := rp.GetEpochLastBlockHeight(uint64(epochNum))
		for height := int(epochStartHeight); height <= int(epochEndHeight); height++ {
			blk, err := bc.GetBlockByHeight(uint64(height))
			if err != nil {
				log.L().Fatal("Failed to get block", zap.Error(err))
			}
			var heightBytes [8]byte
			enc.MachineEndian.PutUint64(heightBytes[:], uint64(height))
			receiptsBytes, err := kvstore.Get("rpt", heightBytes[:])
			if err != nil {
				log.L().Fatal("Failed to get receipts from db")
			}
			receipts := iotextypes.Receipts{}
			if err := proto.Unmarshal(receiptsBytes, &receipts); err != nil {
				log.L().Fatal("Failed to unmarshal receipts")
			}
			receiptMap := make(map[hash.Hash256]*action.Receipt)
			for _, receipt := range receipts.Receipts {
				r := &action.Receipt{}
				r.ConvertFromReceiptPb(receipt)
				receiptMap[r.ActionHash] = r
			}
			for _, selp := range blk.Actions {
				if _, ok := selp.Action().(*action.GrantReward); !ok {
					continue
				}
				receipt := receiptMap[selp.Hash()]
				for _, l := range receipt.Logs {
					rewardLog := &rewardingpb.RewardLog{}
					if err := proto.Unmarshal(l.Data, rewardLog); err != nil {
						log.L().Fatal("Failed to unmarshal receipt data into rewardLog", zap.Error(err))
					}
					if _, ok := epochMeta[rewardAddrMap[rewardLog.Addr]]; !ok {
						continue
					}
					epochMap := epochMeta[rewardAddrMap[rewardLog.Addr]]
					if _, ok := epochMap[epochNum]; !ok {
						epochMap[epochNum] = &epochInfo{
							blockReward:      big.NewInt(0),
							epochReward:      big.NewInt(0),
							foundationReward: big.NewInt(0),
						}
					}
					epochInfo := epochMap[epochNum]
					amount, ok := big.NewInt(0).SetString(rewardLog.Amount, 10)
					if !ok {
						log.L().Fatal("Failed to convert reward amount from string to big int")
					}
					switch rewardLog.Type {
					case rewardingpb.RewardLog_BLOCK_REWARD:
						epochInfo.blockReward.Add(epochInfo.blockReward, amount)
					case rewardingpb.RewardLog_EPOCH_REWARD:
						epochInfo.epochReward.Add(epochInfo.epochReward, amount)
					case rewardingpb.RewardLog_FOUNDATION_BONUS:
						epochInfo.foundationReward.Add(epochInfo.foundationReward, amount)
					default:
						log.L().Fatal("Unknown type of reward")
					}
				}

			}
		}
	}

	for alias, epochMap := range epochMeta {
		if err := writeExcel(alias+"_log.xlsx", epochMap); err != nil {
			log.L().Fatal("Failed to write reward result into excel form.")
		}
	}
}

func writeExcel(
	fileName string,
	epochMap map[int]*epochInfo,
) error {
	file := xlsx.NewFile()
	sheet, err := file.AddSheet("sheet1")
	if err != nil {
		return err
	}
	row := sheet.AddRow()
	cell1 := row.AddCell()
	cell1.Value = "Epoch Number"
	cell2 := row.AddCell()
	cell2.Value = "Block Reward"
	cell3 := row.AddCell()
	cell3.Value = "Epoch Reward"
	cell4 := row.AddCell()
	cell4.Value = "Foundation Reward"

	keys := make([]int, 0)
	for epochNum := range epochMap {
		keys = append(keys, epochNum)
	}
	sort.Ints(keys)
	for _, epochNum := range keys {
		epochInfo := epochMap[epochNum]
		row = sheet.AddRow()
		cell1 = row.AddCell()
		cell1.Value = strconv.Itoa(epochNum)
		cell2 = row.AddCell()
		cell2.Value = epochInfo.blockReward.String()
		cell3 = row.AddCell()
		cell3.Value = epochInfo.epochReward.String()
		cell4 = row.AddCell()
		cell4.Value = epochInfo.foundationReward.String()
	}

	for _, col := range sheet.Cols {
		col.Width = 40
	}
	return file.Save(fileName)
}
