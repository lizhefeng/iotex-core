package main

import (
	"context"
	"math/big"
	"sort"
	"strconv"

	"github.com/pkg/errors"
	"github.com/tealeg/xlsx"
	"go.uber.org/zap"

	"github.com/iotexproject/iotex-core/action/protocol/rewarding"
	"github.com/iotexproject/iotex-core/action/protocol/rolldpos"
	"github.com/iotexproject/iotex-core/address"
	"github.com/iotexproject/iotex-core/blockchain"
	"github.com/iotexproject/iotex-core/config"
	"github.com/iotexproject/iotex-core/crypto"
	"github.com/iotexproject/iotex-core/pkg/log"
	"github.com/iotexproject/iotex-core/server/itx"
	"github.com/iotexproject/iotex-core/state"
)

var (
	aliasToAddr = map[string]string{
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

	addrToAlias = map[string]string{
		"io1vlsmjs87jlk93624nppccfn24nk9nplu9uhu53": "cobo",
		"io1fra0fx6akacny9asewt7vyvggqq4rtq3eznccr": "iotexteam",
		"io1u5ff879gg2dw9vfpxr2tsmuaz07e2rea6gvl7s": "rosemary9",
		"io1kugyfdkdss7acy0y9x8fmfjd5qyxfcye3gxt99": "consensusnet",
		"io12wc9ra4la98yay4cqdav5mwxxuzwpt6hk23n3z": "iotexmainnet",
		"io10reczcaelglh5xmkay65h9vw3e5dp82e8vw0rz": "metanyx",
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

	exemptAddrs = map[string]bool{
		"io15fqav3tugm96ge7anckx0k4gukz5m4mqf0jpv3": true,
		"io1x9kjkr0qv2fa7j4t2as8lrj223xxsqt4tl7xp7": true,
		"io1ar5l5s268rtgzshltnqv88mua06ucm58dx678y": true,
		"io1xsx5n94kg2zv64r7tm8vyz9mh86amfak9ka9xx": true,
		"io1vtm2zgn830pn6auc2cvnchgwdaefa9gr4z0s86": true,
		"io159fv8mu9d5djk8u2t0flgw4yqmt6fg98uqjka8": true,
		"io1c3r4th3zrk4uhv83a9gr4gvn3y6pzaj6mc84ea": true,
		"io14vmhs9c75r2ptxdaqrtk0dz7skct30pxmt69d9": true,
		"io1gf08snppu2a2wfd50pjas2j6q2kcxjzqph3pep": true,
		"io1u5ff879gg2dw9vfpxr2tsmuaz07e2rea6gvl7s": true,
		"io1du4eq4f88n4wyc026l3gamjwetlgsg4jz7j884": true,
		"io12yxdwewry70gr9fs6fphyfaky9c7gurmzk8f4f": true,
		"io1lx53nfgq2dnzrz5ecz2ecs7vvl6qll0mkn970w": true,
		"io127ftn4ry6wgxdrw4hcd6gdwqlq70ujk98dvtw5": true,
		"io1wv5m0xyermvr2n0wjx2cjsqwyk863drdl5qfyn": true,
		"io1v0q5g2f8z6l3v25krl677chdx7g5pwt9kvqfpc": true,
		"io1xsdegzr2hdj5sv5ad4nr0yfgpsd98e40u6svem": true,
		"io1fks575kklxafq4fwjccmz5d3pmq5ynxk5h6h0v": true,
		"io15npzu93ug8r3zdeysppnyrcdu2xssz0lcam9l9": true,
		"io1gh7xfrsnj6p5uqgjpk9xq6jg9na28aewgp7a9v": true,
		"io1nyjs526mnqcsx4twa7nptkg08eclsw5c2dywp4": true,
		"io1jafqlvntcxgyp6e0uxctt3tljzc3vyv5hg4ukh": true,
		"io1z7mjef7w528nasnsafan0rp6yuvkvq405l6r8j": true,
		"io1cup9k8hl8fp40vrj29ex8djc346780dk223end": true,
		"io1scs89jur7qklzh5vfrmha3c40u8yajjx6kvzg9": true,
	}
)

type epochInfo struct {
	isBlockProducer       bool
	isActiveBlockProducer bool
	produce               int
	rewardAddr            string
	votes                 string
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

	endEpochNum := rp.GetEpochNum(bc.TipHeight())
	epochMeta := make(map[string]map[int]*epochInfo)
	for alias := range aliasToAddr {
		epochMeta[alias] = make(map[int]*epochInfo)
	}

	totalVotes := make(map[int]*big.Int)
	robotVotes := make(map[int]*big.Int)
	for epochNum := 1; epochNum <= int(endEpochNum); epochNum++ {
		candidates, err := bc.CandidatesByHeight(rp.GetEpochHeight(uint64(epochNum)))
		if err != nil {
			log.L().Fatal("Failed to get candidates by height", zap.Error(err))
		}
		totalVotes[epochNum] = big.NewInt(0)
		robotVotes[epochNum] = big.NewInt(0)
		for _, cand := range candidates {
			totalVotes[epochNum].Add(totalVotes[epochNum], cand.Votes)
			if _, ok := exemptAddrs[cand.Address]; ok {
				robotVotes[epochNum].Add(robotVotes[epochNum], cand.Votes)
			}
			if epochMap, ok := epochMeta[addrToAlias[cand.Address]]; ok {
				epochMap[epochNum] = &epochInfo{
					rewardAddr: cand.RewardAddress,
					votes:      cand.Votes.String(),
				}
			}
		}

		numDelegates := int(rp.NumCandidateDelegates())
		if len(candidates) < numDelegates {
			numDelegates = len(candidates)
		}
		delegateList := candidates[:numDelegates]
		for _, delegate := range delegateList {
			if epochMap, ok := epochMeta[addrToAlias[delegate.Address]]; ok {
				epochMap[epochNum].isBlockProducer = true
			}
		}

		activeBPs, err := readActiveBlockProducersByEpoch(rp, bc, uint64(epochNum))
		if err != nil {
			log.L().Fatal("Failed to read active block producers")
		}

		for _, activeBP := range activeBPs {
			if epochMap, ok := epochMeta[addrToAlias[activeBP.Address]]; ok {
				epochMap[epochNum].isActiveBlockProducer = true
			}
		}

		epochStartHeight := rp.GetEpochHeight(uint64(epochNum))
		epochEndHeight := rp.GetEpochLastBlockHeight(uint64(epochNum))
		for height := int(epochStartHeight); height <= int(epochEndHeight); height++ {
			blk, err := bc.GetBlockByHeight(uint64(height))
			if err != nil {
				log.L().Fatal("Failed to get block", zap.Error(err))
			}
			if epochMap, ok := epochMeta[addrToAlias[blk.ProducerAddress()]]; ok {
				epochMap[epochNum].produce++
			}
		}
	}
	ws, err := bc.GetFactory().NewWorkingSet()
	if err != nil {
		log.L().Fatal("Failed to create a new working set", zap.Error(err))
	}
	p, ok = registry.Find(rewarding.ProtocolID)
	if !ok {
		log.L().Fatal("Failed to find rewarding protocol.", zap.Error(err))
	}
	rewardProtocol, ok := p.(*rewarding.Protocol)
	if !ok {
		log.L().Fatal("Failed to cast rewarding protocol.")
	}

	rewardAddrMap := make(map[string][]string)
	for alias, epochMap := range epochMeta {
		rewardAddrs := make([]string, 0)
		for _, epochInfo := range epochMap {
			rewardAddrs = append(rewardAddrs, epochInfo.rewardAddr)
		}
		rewardAddrMap[alias] = uniqueAddress(rewardAddrs)
	}

	rewards := make(map[string]map[string]*big.Int)
	for alias, rewardAddrs := range rewardAddrMap {
		rewardMap := make(map[string]*big.Int)
		for _, rewardAddr := range rewardAddrs {
			addr, err := address.FromString(rewardAddr)
			if err != nil {
				log.L().Fatal("Failed to get address from string", zap.Error(err))
			}
			reward, err := rewardProtocol.UnclaimedBalance(context.Background(), ws, addr)
			if err != nil {
				log.L().Fatal("Failed to get unclaimed balance", zap.Error(err))
			}
			rewardMap[rewardAddr] = reward
		}
		rewards[alias] = rewardMap
	}

	for alias, epochMap := range epochMeta {
		if err := writeExcel(alias+"_reward.xlsx", int(cfg.Genesis.NumSubEpochs), epochMap, totalVotes, robotVotes, rewards[alias]); err != nil {
			log.L().Fatal("Failed to write reward result into excel form.")
		}
	}
}

func readBlockProducersByEpoch(p *rolldpos.Protocol, bc blockchain.Blockchain, epochNum uint64) (state.CandidateList, error) {
	epochHeight := p.GetEpochHeight(epochNum)
	delegates, err := bc.CandidatesByHeight(epochHeight)
	if err != nil {
		return nil, err
	}
	var blockProducers state.CandidateList
	for i, delegate := range delegates {
		if uint64(i) >= p.NumCandidateDelegates() {
			break
		}
		blockProducers = append(blockProducers, delegate)
	}
	return blockProducers, nil
}

func readActiveBlockProducersByEpoch(p *rolldpos.Protocol, bc blockchain.Blockchain, epochNum uint64) (state.CandidateList, error) {
	blockProducers, err := readBlockProducersByEpoch(p, bc, epochNum)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get active block producers in epoch %d", epochNum)
	}

	var blockProducerList []string
	blockProducerMap := make(map[string]*state.Candidate)
	for _, bp := range blockProducers {
		blockProducerList = append(blockProducerList, bp.Address)
		blockProducerMap[bp.Address] = bp
	}

	epochHeight := p.GetEpochHeight(epochNum)
	crypto.SortCandidates(blockProducerList, epochHeight, crypto.CryptoSeed)

	length := int(p.NumDelegates())
	if len(blockProducerList) < int(p.NumDelegates()) {
		length = len(blockProducerList)
	}

	var activeBlockProducers state.CandidateList
	for i := 0; i < length; i++ {
		activeBlockProducers = append(activeBlockProducers, blockProducerMap[blockProducerList[i]])
	}
	return activeBlockProducers, nil
}

func writeExcel(
	fileName string,
	NumSubEpochs int,
	epochMap map[int]*epochInfo,
	totalVotes map[int]*big.Int,
	robotVotes map[int]*big.Int,
	rewards map[string]*big.Int,
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
	cell2.Value = "Is Block Producer"
	cell3 := row.AddCell()
	cell3.Value = "Is Active Block Producer"
	cell4 := row.AddCell()
	cell4.Value = "Production"
	cell5 := row.AddCell()
	cell5.Value = "Expected Production"
	cell6 := row.AddCell()
	cell6.Value = "Votes"
	cell7 := row.AddCell()
	cell7.Value = "Total Votes"
	cell8 := row.AddCell()
	cell8.Value = "Robot Votes"
	cell9 := row.AddCell()
	cell9.Value = "Reward Address"

	keys := make([]int, 0)
	for epochNum := range totalVotes {
		keys = append(keys, epochNum)
	}
	sort.Ints(keys)
	for _, epochNum := range keys {
		epochInfo := epochMap[epochNum]
		row = sheet.AddRow()
		cell1 = row.AddCell()
		cell1.Value = strconv.Itoa(epochNum)
		cell2 = row.AddCell()
		cell2.Value = "No"
		if epochInfo != nil && epochInfo.isBlockProducer {
			cell2.Value = "Yes"
		}
		cell3 = row.AddCell()
		cell3.Value = "No"
		if epochInfo != nil && epochInfo.isActiveBlockProducer {
			cell3.Value = "Yes"
		}
		cell4 = row.AddCell()
		cell4.Value = "0"
		if epochInfo != nil {
			cell4.Value = strconv.Itoa(epochInfo.produce)
		}
		cell5 = row.AddCell()
		cell5.Value = "0"
		if epochInfo != nil && epochInfo.isActiveBlockProducer {
			cell5.Value = strconv.Itoa(NumSubEpochs)
		}
		cell6 = row.AddCell()
		cell6.Value = "0"
		if epochInfo != nil {
			cell6.Value = epochInfo.votes
		}
		cell7 = row.AddCell()
		cell7.Value = totalVotes[epochNum].String()
		cell8 := row.AddCell()
		cell8.Value = robotVotes[epochNum].String()
		cell9 = row.AddCell()
		if epochInfo != nil {
			cell9.Value = epochInfo.rewardAddr
		}
	}

	sheet2, err := file.AddSheet("sheet2")
	if err != nil {
		return err
	}

	row = sheet2.AddRow()
	cell1 = row.AddCell()
	cell1.Value = "Reward Address"
	cell2 = row.AddCell()
	cell2.Value = "Unclaimed Reward Balance"
	for rewardAddr, reward := range rewards {
		row = sheet2.AddRow()
		cell1 = row.AddCell()
		cell1.Value = rewardAddr
		cell2 = row.AddCell()
		cell2.Value = reward.String()
	}

	for _, col := range sheet.Cols {
		col.Width = 40
	}

	for _, col := range sheet2.Cols {
		col.Width = 40
	}
	return file.Save(fileName)
}

func uniqueAddress(rewardAddrs []string) []string {
	check := make(map[string]bool)
	unique := make([]string, 0)
	for _, rewardAddr := range rewardAddrs {
		if _, ok := check[rewardAddr]; ok {
			continue
		}
		check[rewardAddr] = true
		unique = append(unique, rewardAddr)
	}
	return unique
}
