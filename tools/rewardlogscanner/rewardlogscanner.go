package main

import (
	"context"
	"math/big"
	"sort"
	"strconv"

	"github.com/golang/protobuf/proto"
	"github.com/tealeg/xlsx"
	"go.uber.org/zap"

	"github.com/iotexproject/iotex-core/action"
	"github.com/iotexproject/iotex-core/action/protocol/rewarding/rewardingpb"
	"github.com/iotexproject/iotex-core/action/protocol/rolldpos"
	"github.com/iotexproject/iotex-core/config"
	"github.com/iotexproject/iotex-core/pkg/log"
	"github.com/iotexproject/iotex-core/server/itx"
)

const endEpochNum = 300

var (
	aliasToOperatorAddr = map[string]string{
		"cobo":         "io1vlsmjs87jlk93624nppccfn24nk9nplu9uhu53",
		"iotexteam":    "io1fra0fx6akacny9asewt7vyvggqq4rtq3eznccr",
		"rosemary9":    "io1u5ff879gg2dw9vfpxr2tsmuaz07e2rea6gvl7s",
		"consensusnet": "io1kugyfdkdss7acy0y9x8fmfjd5qyxfcye3gxt99",
		"iotexmainnet": "io12wc9ra4la98yay4cqdav5mwxxuzwpt6hk23n3z",
		"metanyx":      "io10reczcaelglh5xmkay65h9vw3e5dp82e8vw0rz",
		"pubxpayments": "io1rfeltdr5wgmhm4rdx9eun3pp68ahm7fq00wcmm",
		"ratels":       "io1n3llm0zzlles6pvpzugzajyzyjehztztxf2rff",
		"rosemary3":    "io1xsx5n94kg2zv64r7tm8vyz9mh86amfak9ka9xx",
		"rosemary13":   "io1wv5m0xyermvr2n0wjx2cjsqwyk863drdl5qfyn",
		"stanfordcpc":  "io1qqaswtu7rcevahucpfxc0zxx088pylwtxnkrrl",
		"iotask":       "io1e3w03ulnrsxth2g0rgsq6v406fhdccsgfq3hz7",
		"laomao":       "io1gfq9el2gnguus64ex3hu8ajd6e4yzk3f9cz5vx",
		"rosemary0":    "io15fqav3tugm96ge7anckx0k4gukz5m4mqf0jpv3",
		"rosemary12":   "io127ftn4ry6wgxdrw4hcd6gdwqlq70ujk98dvtw5",
		"royalland":    "io108h7sa5sap44e244hz649zyk5y4rqzsvnpzxh5",
		"hashbuy":      "io1sd5t5dwxrk2t50z8yl86n6ht8c99umt4u6rknl",
		"infstones":    "io17cmrextyfeu4gddwd89g5qncedsnc553dhz7xa",
		"rosemary8":    "io1gf08snppu2a2wfd50pjas2j6q2kcxjzqph3pep",
		"gamefantasy":  "io1qnec80ark9shjc6uzk45dhm8s50dpc27sjuver",
		"rosemary6":    "io1c3r4th3zrk4uhv83a9gr4gvn3y6pzaj6mc84ea",
		"rosemary11":   "io12yxdwewry70gr9fs6fphyfaky9c7gurmzk8f4f",
		"capitalmu":    "io1zy9q4myp3jezxdsv82z3f865ruamhqm7a6mggh",
		"iosg":         "io14u5d66rt465ykm7t2847qllj0reml27q30kr75",
		"iotexlab":     "io1nf0rvzgq3tqym6n3trttsrt7d4gqqsmqfzy0da",
		"rosemary5":    "io159fv8mu9d5djk8u2t0flgw4yqmt6fg98uqjka8",
		"rosemary7":    "io14vmhs9c75r2ptxdaqrtk0dz7skct30pxmt69d9",
		"thebottoken":  "io1xj0u5n20tsqwxh5a3xdtmzuz9wasft0pqjrq8t",
		"wannodes":     "io13v3dhwds82hg0uc9l4puer00k93qagdh62j0mz",
		"whales":       "io1wl83n3up9w8nedf30lnyxzple0gu5pme0dyrds",
		"yvalidator":   "io1f0rh94m3ctkwep3rlsswwq5vxwlntx4s574l3q",
		"coingecko":    "io1we5gqn4xeak9ycnu4l9lv0qfq3euapymnzffx6",
		"droute":       "io1kr8c6krd7dhxaaqwdkr6erqgu4z0scug3drgja",
		"iotexunion":   "io17dm3tq9ezgs2uvwrqu8e004l5nqq33jm46r0mg",
		"keysiotex":    "io1aqf30kqz5rqh6zn82c00j684p2h2t5cg30wm8t",
		"longz":        "io1ddjluttkzljqfgdtcz9eu3r3s832umy7ylx0ts",
		"rosemary1":    "io1x9kjkr0qv2fa7j4t2as8lrj223xxsqt4tl7xp7",
		"rosemary4":    "io1vtm2zgn830pn6auc2cvnchgwdaefa9gr4z0s86",
		"rosemary10":   "io1du4eq4f88n4wyc026l3gamjwetlgsg4jz7j884",
		"iotexgeeks":   "io1lm85qfm24eqrc042spnplac9ran28thuh5048n",
		"iotxplorerio": "io1et7zkzc76m9twa4gn5xht3urt9mwj05qvdtj66",
		"preangle":     "io1z4sxtefurklkyrfmmdtjmw4h8csnxlv9747hyd",
		"rosemary":     "io1ar5l5s268rtgzshltnqv88mua06ucm58dx678y",
	}

	rewardAddrMap = map[string]string{
		"io1vlsmjs87jlk93624nppccfn24nk9nplu9uhu53": "cobo",
		"io1fra0fx6akacny9asewt7vyvggqq4rtq3eznccr": "iotexteam",
		"io1u5ff879gg2dw9vfpxr2tsmuaz07e2rea6gvl7s": "rosemary9",
		"io1kugyfdkdss7acy0y9x8fmfjd5qyxfcye3gxt99": "consensusnet",
		"io12wc9ra4la98yay4cqdav5mwxxuzwpt6hk23n3z": "iotexmainnet",
		"io1knfzqrzwlx9cx6t2trgcppn767vptlce0m2kje": "metanyx",
		"io1rfeltdr5wgmhm4rdx9eun3pp68ahm7fq00wcmm": "pubxpayments",
		"io1n3llm0zzlles6pvpzugzajyzyjehztztxf2rff": "ratels",
		"io1xsx5n94kg2zv64r7tm8vyz9mh86amfak9ka9xx": "rosemary3",
		"io1wv5m0xyermvr2n0wjx2cjsqwyk863drdl5qfyn": "rosemary13",
		"io1qqaswtu7rcevahucpfxc0zxx088pylwtxnkrrl": "stanfordcpc",
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
		"io1c6a0ugwqg9p96l69fzh00xk5l4ffjrqfcv54ue": "thebottoken",
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
	for alias := range aliasToOperatorAddr {
		epochMeta[alias] = make(map[int]*epochInfo)
	}

	for epochNum := 1; epochNum <= int(endEpochNum); epochNum++ {
		epochStartHeight := rp.GetEpochHeight(uint64(epochNum))
		epochEndHeight := rp.GetEpochLastBlockHeight(uint64(epochNum))
		for height := int(epochStartHeight); height <= int(epochEndHeight); height++ {
			blk, err := bc.GetBlockByHeight(uint64(height))
			if err != nil {
				log.L().Fatal("Failed to get block", zap.Error(err))
			}
			for _, selp := range blk.Actions {
				if _, ok := selp.Action().(*action.GrantReward); !ok {
					continue
				}
				receipt, err := bc.GetReceiptByActionHash(selp.Hash())
				if err != nil {
					log.L().Fatal("Failed to get receipt by action hash", zap.Error(err))
				}
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
