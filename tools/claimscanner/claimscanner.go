package main

import (
	"context"

	"go.uber.org/zap"

	"github.com/iotexproject/iotex-core/pkg/log"
	"github.com/iotexproject/iotex-core/config"
	"github.com/iotexproject/iotex-core/server/itx"
	"github.com/iotexproject/iotex-core/action"
	"github.com/iotexproject/iotex-core/address"
)

var (
	delegateMap = map[string]string{
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

	robotPublicKey = map[string]bool {
		"04e707e24ff89349260b5b34c7fdcf7078804dcd6c97c9097b41e020ada9c269e0b9b60bb79d5c77263290634f30952bcfce8cc9340162a568d86ccee5119ff866": true,
		"048b895e188f84ad8ec7448e09793bee4489fa64e659a803a34f27cb2fab9b71451d3fbaa58b41a6c283f90ee1a1e11693df7a05d2bee4797c3f30e2906dfaba08": true,
		"0476ca3357f7128eeb06e5ff84ec1c58ab363a40ebf9c691e0e2336cc7949f634b03348f2afdfeb548306b9756b75da8d3713cf914cd97ef1c3d697bbb94cb1905": true,
		"0459ce74a763c249359606d25dd3da2cf353d89a9b4dc284a6854b80fd66d3d6c736977d3805b666416b7fa56d2bdd45e2e18fcb5b68fee3e6efd8497bd132736f": true,
		"0464791bd4ce7d4510128089fb1dbad86a7ef30bffd3f2bbdce64a97f3fab4910ef8256d84137dca21722d915e8c2796a32cb958acb5e01fa25007468cf8258508": true,
		"046eac1eeff49eea4f2269bda1b5c8275fd74cb9bb388b25df72bf72ce42d2a916b45751a42277e80b3670450a8fefc6ee786f3122eaec6389d8ced8a5d32ae67e": true,
		"04ca65d5d295b86fe2daa2afb902d8dff09b39ebdcb12e140aa704e30a47e48d297a8c1064973a399dd385f98695279744a7ffd066d787ba7d3d1f1b784ae5cb5c": true,
		"049e7711cbbbaf967dff854fb8aad60bf2e161148eea4caf0ef788d637ba2fb46551f434aa5e44cde243ec71a628490eaf212281ea6d61bca8c08ca60ea1df957c": true,
		"04e9c3c6ceaebd947fe183871080d4e99f953d6f916fc48b574a4372202531e30e259eb631e08af5014f4e205136613e6c9b75bd355f821eabe226ae059dc88de8": true,
		"0433ee7530068315479cf6d59bfb6fa81af77c434e20cd38b5352e33edab28bd4ab59aa0f5c3039f60d23c3f1798a9515d2b1e54274747b3f163923aa04028b07d": true,
		"04367f33b5ecce97e1eeafa288a5f6d9127309f83e94e8065ae2d735f26505747da0a0ad187a2c343dbca7a12f6723769d97e970a8899f7becc731a43f3e57dbd4": true,
		"04ff2ef9d15496ee0b3d015d2f7ae3f0d6c53c3e1ee5b819a29e2d7d4e9b646599db583e962dc29c62f54d0870e89202c200e7b2c6a4c271336e3944b81c731f34": true,
		"047892302dafa1f0381c99eed0f4523c7bb04685e190845c5230c688b57838fe6f081f1ad02257de5109d65cb86d55038fb5dd13728335b15aabbb4b2ae978b138": true,
		"0408a95f8b0d1e0208b2bc1482f84e282309ea976b373720f36de889d9a9af32d2f8bd9bdf376bf845fc6ea8c46af846ed0ee664308bdc0d6ce021b9314b1e4d88": true,
		"0453d2446603d339a0e91032df30a8d0ff70cbe3cf91036369dfb6792cb5fb296cd3683d642a32ad0502cb1bf3d52fac2863c5afc88167e9e22382277f3c18b285": true,
		"04dc0a472eae0ded9d24f3246861488d6211b650f524d6a3d6f3d474c48dea711396680cf0682562fac32ef4cfe56f5435d7d3759839ddde36c94d2d9fa4c4b9b0": true,
		"040a2727d560f9970ade38ca9f15bf445e898d88433b253618c03a7bc64f158e9b70d8a99cdd16652ffdc4fce23148653bd92e27bfcdd16e16619f2791783503a4": true,
		"0438b73234fbd0c5c01709aa9809233154c5527988ee80957fb8ae40e850c2fd7c5d91c5367a8a05eb07f95a64edbdfff8860c1f6ace34e252598f424bd0afdfea": true,
		"04e1077403181df5bbd278dfbcab63ad9de514a61c1901ebeabd3181fee6572cf395b0444af67728048d5a6f173b74170bd0e76875dcf4fad8b322c1328e50225e": true,
		"048c9b3fc60ab78209ae36a2e12f3d306f7d279e9aa11e9d7501cad144725bfa4098a00c380b50b17946524e6f9158f85c06f8e9fa4390364dd839224eb3d796f1": true,
		"04c7c18d874b04147e6edb9380fd9a9de91115e5ff08a3b2c36a2904ceeaebfda3b5ee98c45edd34f11bc46eec2563aa0a2e8aa51a756e21692b1f9eaadcaa79e6": true,
		"04d0369b83797de18bfbc113b5766fb07434269fc676be3a6199af974172e8871bce7a6a5590af18a34f702bc53ef60f01f6e37132d23a177979100e36e06e1e63": true,
		"041033d158d4995415c7900168eb20ace8ab53270f56ee13a797588ca1648531ba6635f460cbf107397eeca03b7ce8fc79b0efd225f1e1a61409c228dae7d08b1d": true,
		"0417898d8a18883e98c1ad46604130042f6f0b4eff9cf93d8e094a655c1b94b0328216cee3cdda4ab475728753e6da6c70ea24e58a50526a9928d4c57f63befa6e": true,
	}
)

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

	for i := 1; i <= int(bc.TipHeight()); i++ {
		blk, err := bc.GetBlockByHeight(uint64(i))
		if err != nil {
			log.L().Fatal("Failed to get block by height")
		}
		for _, selp := range blk.Actions {
			claim, ok := selp.Action().(*action.ClaimFromRewardingFund)
			if !ok {
				continue
			}
			srcPubKey := claim.SrcPubkey()
			if _, ok := robotPublicKey[srcPubKey.HexString()]; !ok {
				srcAddr, err := address.FromBytes(srcPubKey.Hash())
				if err != nil {
					log.L().Fatal("Failed to get address from public key hash")
				}
				claimer := "NoOne"
				if sender, ok := delegateMap[srcAddr.String()]; ok {
					claimer = sender
				}
				actHash := claim.Hash()
				receipt, err := bc.GetReceiptByActionHash(actHash)
				if err != nil {
					log.L().Fatal("Failed to get receipt by action hash")
				}
				if receipt.Status == uint64(1) {
					log.L().Info("Claim Reward", zap.Uint64("Height", uint64(i)), zap.String("Claimer", claimer), zap.String("Amount", claim.Amount().String()))
				}
			}

		}
	}
}
