package apexpro

import (
	"github.com/thrasher-corp/gocryptotrader/common/convert"
	"github.com/thrasher-corp/gocryptotrader/types"
)

// AllSymbolsConfigs represents all symbols configurations.
type AllSymbolsConfigs struct {
	SpotConfig struct {
		Assets []struct {
			TokenID       string `json:"tokenId"`
			Token         string `json:"token"`
			DisplayName   string `json:"displayName"`
			Decimals      int64  `json:"decimals"`
			ShowStep      string `json:"showStep"`
			IconURL       string `json:"iconUrl"`
			L2WithdrawFee string `json:"l2WithdrawFee"`
		} `json:"assets"`
		Global struct {
			DefaultRegisterTransferToken      string `json:"defaultRegisterTransferToken"`
			DefaultRegisterTransferTokenID    string `json:"defaultRegisterTransferTokenId"`
			DefaultRegisterSubAccountID       string `json:"defaultRegisterSubAccountId"`
			DefaultChangePubKeyZklinkChainID  string `json:"defaultChangePubKeyZklinkChainId"`
			DefaultChangePubKeyFeeTokenID     string `json:"defaultChangePubKeyFeeTokenId"`
			DefaultChangePubKeyFeeToken       string `json:"defaultChangePubKeyFeeToken"`
			DefaultChangePubKeyFee            string `json:"defaultChangePubKeyFee"`
			RegisterTransferLpAccountID       string `json:"registerTransferLpAccountId"`
			RegisterTransferLpSubAccount      string `json:"registerTransferLpSubAccount"`
			RegisterTransferLpSubAccountL2Key string `json:"registerTransferLpSubAccountL2Key"`
			PerpLpAccountID                   string `json:"perpLpAccountId"`
			PerpLpSubAccount                  string `json:"perpLpSubAccount"`
			PerpLpSubAccountL2Key             string `json:"perpLpSubAccountL2Key"`
			ContractAssetPoolAccountID        string `json:"contractAssetPoolAccountId"`
			ContractAssetPoolZkAccountID      string `json:"contractAssetPoolZkAccountId"`
			ContractAssetPoolSubAccount       string `json:"contractAssetPoolSubAccount"`
			ContractAssetPoolL2Key            string `json:"contractAssetPoolL2Key"`
			ContractAssetPoolEthAddress       string `json:"contractAssetPoolEthAddress"`
		} `json:"global"`
		Spot       []any `json:"spot"`
		MultiChain struct {
			Chains []struct {
				Chain              string `json:"chain"`
				ChainID            string `json:"chainId"`
				ChainType          string `json:"chainType"`
				L1ChainID          string `json:"l1ChainId"`
				ChainIconURL       string `json:"chainIconUrl"`
				ContractAddress    string `json:"contractAddress"`
				StopDeposit        bool   `json:"stopDeposit"`
				FeeLess            bool   `json:"feeLess"`
				GasLess            bool   `json:"gasLess"`
				GasToken           string `json:"gasToken"`
				DynamicFee         bool   `json:"dynamicFee"`
				FeeGasLimit        int64  `json:"feeGasLimit"`
				BlockTimeSeconds   int64  `json:"blockTimeSeconds"`
				RPCURL             string `json:"rpcUrl"`
				WebRPCURL          string `json:"webRpcUrl"`
				WebTxURL           string `json:"webTxUrl"`
				TxConfirm          int64  `json:"txConfirm"`
				WithdrawGasFeeLess bool   `json:"withdrawGasFeeLess"`
				Tokens             []struct {
					Decimals          int64  `json:"decimals"`
					IconURL           string `json:"iconUrl"`
					Token             string `json:"token"`
					TokenAddress      string `json:"tokenAddress"`
					PullOff           bool   `json:"pullOff"`
					WithdrawEnable    bool   `json:"withdrawEnable"`
					Slippage          string `json:"slippage"`
					IsDefaultToken    bool   `json:"isDefaultToken"`
					DisplayToken      string `json:"displayToken"`
					NeedResetApproval bool   `json:"needResetApproval"`
					MinFee            string `json:"minFee"`
					MaxFee            string `json:"maxFee"`
					FeeRate           string `json:"feeRate"`
				} `json:"tokens"`
			} `json:"chains"`
			MaxWithdraw string `json:"maxWithdraw"`
			MinDeposit  string `json:"minDeposit"`
			MinWithdraw string `json:"minWithdraw"`
		} `json:"multiChain"`
	} `json:"spotConfig"`
	ContractConfig struct {
		Assets []AssetInfo `json:"assets"`
		Tokens []struct {
			Token    string       `json:"token"`
			StepSize types.Number `json:"stepSize"`
			IconURL  string       `json:"iconUrl"`
		} `json:"tokens"`
		Global struct {
			FeeAccountID             string `json:"feeAccountId"`
			FeeAccountL2Key          string `json:"feeAccountL2Key"`
			ContractAssetLpAccountID string `json:"contractAssetLpAccountId"`
			ContractAssetLpL2Key     string `json:"contractAssetLpL2Key"`
			OperationAccountID       string `json:"operationAccountId"`
			OperationL2Key           string `json:"operationL2Key"`
			ExperienceMoneyAccountID string `json:"experienceMoneyAccountId"`
			ExperienceMoneyL2Key     string `json:"experienceMoneyL2Key"`
			AgentAccountID           string `json:"agentAccountId"`
			AgentL2Key               string `json:"agentL2Key"`
			FinxFeeAccountID         string `json:"finxFeeAccountId"`
			FinxFeeL2Key             string `json:"finxFeeL2Key"`
			NegativeRateAccountID    string `json:"negativeRateAccountId"`
			NegativeRateL2Key        string `json:"negativeRateL2Key"`
			BrokerAccountID          string `json:"brokerAccountId"`
			BrokerL2Key              string `json:"brokerL2Key"`
		} `json:"global"`
		PerpetualContract []struct {
			BaselinePositionValue            string               `json:"baselinePositionValue"`
			CrossID                          int64                `json:"crossId"`
			CrossSymbolID                    int64                `json:"crossSymbolId"`
			CrossSymbolName                  string               `json:"crossSymbolName"`
			DigitMerge                       string               `json:"digitMerge"`
			DisplayMaxLeverage               string               `json:"displayMaxLeverage"`
			DisplayMinLeverage               string               `json:"displayMinLeverage"`
			EnableDisplay                    bool                 `json:"enableDisplay"`
			EnableOpenPosition               bool                 `json:"enableOpenPosition"`
			EnableTrade                      bool                 `json:"enableTrade"`
			FundingImpactMarginNotional      string               `json:"fundingImpactMarginNotional"`
			FundingInterestRate              types.Number         `json:"fundingInterestRate"`
			IncrementalInitialMarginRate     types.Number         `json:"incrementalInitialMarginRate"`
			IncrementalMaintenanceMarginRate types.Number         `json:"incrementalMaintenanceMarginRate"`
			IncrementalPositionValue         string               `json:"incrementalPositionValue"`
			InitialMarginRate                types.Number         `json:"initialMarginRate"`
			MaintenanceMarginRate            types.Number         `json:"maintenanceMarginRate"`
			MaxOrderSize                     types.Number         `json:"maxOrderSize"`
			MaxPositionSize                  types.Number         `json:"maxPositionSize"`
			MinOrderSize                     types.Number         `json:"minOrderSize"`
			MaxMarketPriceRange              string               `json:"maxMarketPriceRange"`
			SettleAssetID                    string               `json:"settleAssetId"`
			BaseTokenID                      string               `json:"baseTokenId"`
			StepSize                         types.Number         `json:"stepSize"`
			Symbol                           string               `json:"symbol"`
			SymbolDisplayName                string               `json:"symbolDisplayName"`
			TickSize                         types.Number         `json:"tickSize"`
			MaxMaintenanceMarginRate         types.Number         `json:"maxMaintenanceMarginRate"`
			MaxPositionValue                 string               `json:"maxPositionValue"`
			TagIconURL                       string               `json:"tagIconUrl"`
			Tag                              string               `json:"tag"`
			RiskTip                          bool                 `json:"riskTip"`
			DefaultInitialMarginRate         types.Number         `json:"defaultInitialMarginRate"`
			KlineStartTime                   convert.ExchangeTime `json:"klineStartTime"`
			MaxMarketSizeBuffer              string               `json:"maxMarketSizeBuffer"`
			EnableFundingSettlement          bool                 `json:"enableFundingSettlement"`
			IndexPriceDecimals               int64                `json:"indexPriceDecimals"`
			IndexPriceVarRate                types.Number         `json:"indexPriceVarRate"`
			OpenPositionOiLimitRate          types.Number         `json:"openPositionOiLimitRate"`
			FundingMaxRate                   types.Number         `json:"fundingMaxRate"`
			FundingMinRate                   types.Number         `json:"fundingMinRate"`
			FundingMaxValue                  string               `json:"fundingMaxValue"`
			EnableFundingMxValue             bool                 `json:"enableFundingMxValue"`
			L2PairID                         string               `json:"l2PairId"`
			SettleTimeStamp                  convert.ExchangeTime `json:"settleTimeStamp"`
			IsPrelaunch                      bool                 `json:"isPrelaunch"`
			RiskLimitConfig                  struct {
				PositionSteps []string `json:"positionSteps"`
				ImrSteps      []string `json:"imrSteps"`
				MmrSteps      []string `json:"mmrSteps"`
			} `json:"riskLimitConfig"`
		} `json:"perpetualContract"`
		PrelaunchContract []struct {
			BaselinePositionValue            string               `json:"baselinePositionValue"`
			CrossID                          int64                `json:"crossId"`
			CrossSymbolID                    int64                `json:"crossSymbolId"`
			CrossSymbolName                  string               `json:"crossSymbolName"`
			DigitMerge                       string               `json:"digitMerge"`
			DisplayMaxLeverage               string               `json:"displayMaxLeverage"`
			DisplayMinLeverage               string               `json:"displayMinLeverage"`
			EnableDisplay                    bool                 `json:"enableDisplay"`
			EnableOpenPosition               bool                 `json:"enableOpenPosition"`
			EnableTrade                      bool                 `json:"enableTrade"`
			FundingImpactMarginNotional      string               `json:"fundingImpactMarginNotional"`
			FundingInterestRate              types.Number         `json:"fundingInterestRate"`
			IncrementalInitialMarginRate     types.Number         `json:"incrementalInitialMarginRate"`
			IncrementalMaintenanceMarginRate types.Number         `json:"incrementalMaintenanceMarginRate"`
			IncrementalPositionValue         string               `json:"incrementalPositionValue"`
			InitialMarginRate                types.Number         `json:"initialMarginRate"`
			MaintenanceMarginRate            types.Number         `json:"maintenanceMarginRate"`
			MaxOrderSize                     string               `json:"maxOrderSize"`
			MaxPositionSize                  types.Number         `json:"maxPositionSize"`
			MinOrderSize                     types.Number         `json:"minOrderSize"`
			MaxMarketPriceRange              string               `json:"maxMarketPriceRange"`
			SettleAssetID                    string               `json:"settleAssetId"`
			BaseTokenID                      string               `json:"baseTokenId"`
			StepSize                         types.Number         `json:"stepSize"`
			Symbol                           string               `json:"symbol"`
			SymbolDisplayName                string               `json:"symbolDisplayName"`
			TickSize                         types.Number         `json:"tickSize"`
			MaxMaintenanceMarginRate         types.Number         `json:"maxMaintenanceMarginRate"`
			MaxPositionValue                 string               `json:"maxPositionValue"`
			TagIconURL                       string               `json:"tagIconUrl"`
			Tag                              string               `json:"tag"`
			RiskTip                          bool                 `json:"riskTip"`
			DefaultLeverage                  types.Number         `json:"defaultLeverage"`
			KlineStartTime                   convert.ExchangeTime `json:"klineStartTime"`
			MaxMarketSizeBuffer              string               `json:"maxMarketSizeBuffer"`
			EnableFundingSettlement          bool                 `json:"enableFundingSettlement"`
			IndexPriceDecimals               int64                `json:"indexPriceDecimals"`
			IndexPriceVarRate                types.Number         `json:"indexPriceVarRate"`
			OpenPositionOiLimitRate          types.Number         `json:"openPositionOiLimitRate"`
			FundingMaxRate                   types.Number         `json:"fundingMaxRate"`
			FundingMinRate                   types.Number         `json:"fundingMinRate"`
			FundingMaxValue                  string               `json:"fundingMaxValue"`
			EnableFundingMxValue             bool                 `json:"enableFundingMxValue"`
			L2PairID                         string               `json:"l2PairId"`
			SettleTimeStamp                  convert.ExchangeTime `json:"settleTimeStamp"`
			IsPrelaunch                      bool                 `json:"isPrelaunch"`
			RiskLimitConfig                  struct {
				PositionSteps any `json:"positionSteps"`
				ImrSteps      any `json:"imrSteps"`
				MmrSteps      any `json:"mmrSteps"`
			} `json:"riskLimitConfig"`
		} `json:"prelaunchContract"`
		MaxMarketBalanceBuffer string `json:"maxMarketBalanceBuffer"`
	} `json:"contractConfig"`
}

// AssetInfo represents an asset detail information.
type AssetInfo struct {
	TokenID       string       `json:"tokenId"`
	Token         string       `json:"token"`
	DisplayName   string       `json:"displayName"`
	Decimals      types.Number `json:"decimals"`
	ShowStep      string       `json:"showStep"`
	IconURL       string       `json:"iconUrl"`
	L2WithdrawFee types.Number `json:"l2WithdrawFee"`
}

// NewTradingData represents a new trading data detail.
type NewTradingData struct {
	Side      string               `json:"S"`
	Volume    types.Number         `json:"v"`
	Price     types.Number         `json:"p"`
	Symbol    string               `json:"s"`
	TradeTime convert.ExchangeTime `json:"T"`
}

// MarketDepthV3 represents a market depth information.
type MarketDepthV3 []struct {
	Assks      [][]types.Number     `json:"a"` // Sell
	Bids       [][]types.Number     `json:"b"` // Buy
	Symbol     string               `json:"s"`
	UpdateTime convert.ExchangeTime `json:"u"`
}
