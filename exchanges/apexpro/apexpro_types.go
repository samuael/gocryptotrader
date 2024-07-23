package apexpro

import (
	"time"

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
type MarketDepthV3 struct {
	Asks       [][2]types.Number    `json:"a"` // Sell
	Bids       [][2]types.Number    `json:"b"` // Buy
	Symbol     string               `json:"s"`
	UpdateTime convert.ExchangeTime `json:"u"`
}

// CandlestickData represents a candlestick chart data.
type CandlestickData struct {
	Start    convert.ExchangeTime `json:"start"`
	Symbol   string               `json:"symbol"`
	Interval string               `json:"interval"`
	Low      types.Number         `json:"low"`
	High     types.Number         `json:"high"`
	Open     types.Number         `json:"open"`
	Close    types.Number         `json:"close"`
	Volume   types.Number         `json:"volume"`
	Turnover string               `json:"turnover"`
}

// TickerData represents a price ticker data.
type TickerData struct {
	Symbol               string       `json:"symbol"`
	Price24HPcnt         types.Number `json:"price24hPcnt"`
	LastPrice            types.Number `json:"lastPrice"`
	HighPrice24H         types.Number `json:"highPrice24h"`
	LowPrice24H          types.Number `json:"lowPrice24h"`
	MarkPrice            types.Number `json:"markPrice"`
	IndexPrice           types.Number `json:"indexPrice"`
	OpenInterest         types.Number `json:"openInterest"`
	Turnover24H          types.Number `json:"turnover24h"`
	Volume24H            types.Number `json:"volume24h"`
	FundingRate          types.Number `json:"fundingRate"`
	PredictedFundingRate types.Number `json:"predictedFundingRate"`
	NextFundingTime      string       `json:"nextFundingTime"`
	TradeCount           types.Number `json:"tradeCount"`
}

// FundingRateHistory represents a funding rate history response.
type FundingRateHistory struct {
	HistoryFunds []struct {
		Symbol           string               `json:"symbol"`
		Rate             types.Number         `json:"rate"`
		Price            types.Number         `json:"price"`
		FundingTime      convert.ExchangeTime `json:"fundingTime"`
		FundingTimestamp convert.ExchangeTime `json:"fundingTimestamp"`
	} `json:"historyFunds"`
	TotalSize int64 `json:"totalSize"`
}

// V2ConfigData v2 assets and symbols configuration response.
type V2ConfigData struct {
	Data struct {
		UsdcConfig struct {
			Currency []struct {
				ID                string `json:"id"`
				StarkExAssetID    string `json:"starkExAssetId"`
				StarkExResolution string `json:"starkExResolution"`
				StepSize          string `json:"stepSize"`
				ShowStep          string `json:"showStep"`
				IconURL           string `json:"iconUrl"`
			} `json:"currency"`
			Global struct {
				FeeAccountID                    string `json:"feeAccountId"`
				FeeAccountL2Key                 string `json:"feeAccountL2Key"`
				StarkExCollateralCurrencyID     string `json:"starkExCollateralCurrencyId"`
				StarkExFundingValidityPeriod    int    `json:"starkExFundingValidityPeriod"`
				StarkExMaxFundingRate           string `json:"starkExMaxFundingRate"`
				StarkExOrdersTreeHeight         int    `json:"starkExOrdersTreeHeight"`
				StarkExPositionsTreeHeight      int    `json:"starkExPositionsTreeHeight"`
				StarkExPriceValidityPeriod      int    `json:"starkExPriceValidityPeriod"`
				StarkExContractAddress          string `json:"starkExContractAddress"`
				RegisterEnvID                   int    `json:"registerEnvId"`
				CrossChainAccountID             string `json:"crossChainAccountId"`
				CrossChainL2Key                 string `json:"crossChainL2Key"`
				FastWithdrawAccountID           string `json:"fastWithdrawAccountId"`
				FastWithdrawFactRegisterAddress string `json:"fastWithdrawFactRegisterAddress"`
				FastWithdrawL2Key               string `json:"fastWithdrawL2Key"`
				FastWithdrawMaxAmount           string `json:"fastWithdrawMaxAmount"`
				BybitWithdrawAccountID          string `json:"bybitWithdrawAccountId"`
				BybitWithdrawL2Key              string `json:"bybitWithdrawL2Key"`
				ExperienceMonenyAccountID       string `json:"experienceMonenyAccountId"`
				ExperienceMonenyL2Key           string `json:"experienceMonenyL2Key"`
				ExperienceMoneyAccountID        string `json:"experienceMoneyAccountId"`
				ExperienceMoneyL2Key            string `json:"experienceMoneyL2Key"`
			} `json:"global"`
			PerpetualContract []struct {
				BaselinePositionValue            string               `json:"baselinePositionValue"`
				CrossID                          int                  `json:"crossId"`
				CrossSymbolID                    int                  `json:"crossSymbolId"`
				CrossSymbolName                  string               `json:"crossSymbolName"`
				DigitMerge                       string               `json:"digitMerge"`
				DisplayMaxLeverage               types.Number         `json:"displayMaxLeverage"`
				DisplayMinLeverage               types.Number         `json:"displayMinLeverage"`
				EnableDisplay                    bool                 `json:"enableDisplay"`
				EnableOpenPosition               bool                 `json:"enableOpenPosition"`
				EnableTrade                      bool                 `json:"enableTrade"`
				FundingImpactMarginNotional      string               `json:"fundingImpactMarginNotional"`
				FundingInterestRate              types.Number         `json:"fundingInterestRate"`
				IncrementalInitialMarginRate     types.Number         `json:"incrementalInitialMarginRate"`
				IncrementalMaintenanceMarginRate types.Number         `json:"incrementalMaintenanceMarginRate"`
				IncrementalPositionValue         types.Number         `json:"incrementalPositionValue"`
				InitialMarginRate                types.Number         `json:"initialMarginRate"`
				MaintenanceMarginRate            types.Number         `json:"maintenanceMarginRate"`
				MaxOrderSize                     types.Number         `json:"maxOrderSize"`
				MaxPositionSize                  types.Number         `json:"maxPositionSize"`
				MinOrderSize                     types.Number         `json:"minOrderSize"`
				MaxMarketPriceRange              types.Number         `json:"maxMarketPriceRange"`
				SettleCurrencyID                 string               `json:"settleCurrencyId"`
				StarkExOraclePriceQuorum         string               `json:"starkExOraclePriceQuorum"`
				StarkExResolution                string               `json:"starkExResolution"`
				StarkExRiskFactor                string               `json:"starkExRiskFactor"`
				StarkExSyntheticAssetID          string               `json:"starkExSyntheticAssetId"`
				StepSize                         string               `json:"stepSize"`
				Symbol                           string               `json:"symbol"`
				SymbolDisplayName                string               `json:"symbolDisplayName"`
				SymbolDisplayName2               string               `json:"symbolDisplayName2"`
				TickSize                         string               `json:"tickSize"`
				UnderlyingCurrencyID             string               `json:"underlyingCurrencyId"`
				MaxMaintenanceMarginRate         types.Number         `json:"maxMaintenanceMarginRate"`
				MaxPositionValue                 string               `json:"maxPositionValue"`
				TagIconURL                       string               `json:"tagIconUrl"`
				Tag                              string               `json:"tag"`
				RiskTip                          bool                 `json:"riskTip"`
				DefaultLeverage                  string               `json:"defaultLeverage"`
				KlineStartTime                   convert.ExchangeTime `json:"klineStartTime"`
			} `json:"perpetualContract"`
			MultiChain struct {
				Chains []struct {
					Chain             string       `json:"chain"`
					ChainID           int64        `json:"chainId"`
					ChainIconURL      string       `json:"chainIconUrl"`
					ContractAddress   string       `json:"contractAddress"`
					DepositGasFeeLess bool         `json:"depositGasFeeLess"`
					StopDeposit       bool         `json:"stopDeposit"`
					FeeLess           bool         `json:"feeLess"`
					FeeRate           string       `json:"feeRate"`
					GasLess           bool         `json:"gasLess"`
					GasToken          string       `json:"gasToken"`
					MinFee            types.Number `json:"minFee"`
					DynamicFee        bool         `json:"dynamicFee"`
					RPCURL            string       `json:"rpcUrl"`
					WebRPCURL         string       `json:"webRpcUrl"`
					WebTxURL          string       `json:"webTxUrl"`
					BlockTime         string       `json:"blockTime"`
					TxConfirm         int          `json:"txConfirm"`
					Tokens            []struct {
						Decimals       int64  `json:"decimals"`
						IconURL        string `json:"iconUrl"`
						Token          string `json:"token"`
						TokenAddress   string `json:"tokenAddress"`
						PullOff        bool   `json:"pullOff"`
						WithdrawEnable bool   `json:"withdrawEnable"`
						Slippage       string `json:"slippage"`
						IsDefaultToken bool   `json:"isDefaultToken"`
						DisplayToken   string `json:"displayToken"`
					} `json:"tokens"`
					WithdrawGasFeeLess bool `json:"withdrawGasFeeLess"`
					IsGray             bool `json:"isGray"`
				} `json:"chains"`
				Currency    string       `json:"currency"`
				MaxWithdraw types.Number `json:"maxWithdraw"`
				MinDeposit  types.Number `json:"minDeposit"`
				MinWithdraw types.Number `json:"minWithdraw"`
			} `json:"multiChain"`
			DepositFromBybit bool `json:"depositFromBybit"`
		} `json:"usdcConfig"`
		UsdtConfig struct {
			Currency []struct {
				ID                string       `json:"id"`
				StarkExAssetID    string       `json:"starkExAssetId"`
				StarkExResolution string       `json:"starkExResolution"`
				StepSize          types.Number `json:"stepSize"`
				ShowStep          string       `json:"showStep"`
				IconURL           string       `json:"iconUrl"`
			} `json:"currency"`
			Global struct {
				FeeAccountID                    string `json:"feeAccountId"`
				FeeAccountL2Key                 string `json:"feeAccountL2Key"`
				StarkExCollateralCurrencyID     string `json:"starkExCollateralCurrencyId"`
				StarkExFundingValidityPeriod    int64  `json:"starkExFundingValidityPeriod"`
				StarkExMaxFundingRate           string `json:"starkExMaxFundingRate"`
				StarkExOrdersTreeHeight         int64  `json:"starkExOrdersTreeHeight"`
				StarkExPositionsTreeHeight      int64  `json:"starkExPositionsTreeHeight"`
				StarkExPriceValidityPeriod      int64  `json:"starkExPriceValidityPeriod"`
				StarkExContractAddress          string `json:"starkExContractAddress"`
				RegisterEnvID                   int64  `json:"registerEnvId"`
				CrossChainAccountID             string `json:"crossChainAccountId"`
				CrossChainL2Key                 string `json:"crossChainL2Key"`
				FastWithdrawAccountID           string `json:"fastWithdrawAccountId"`
				FastWithdrawFactRegisterAddress string `json:"fastWithdrawFactRegisterAddress"`
				FastWithdrawL2Key               string `json:"fastWithdrawL2Key"`
				FastWithdrawMaxAmount           string `json:"fastWithdrawMaxAmount"`
				BybitWithdrawAccountID          string `json:"bybitWithdrawAccountId"`
				BybitWithdrawL2Key              string `json:"bybitWithdrawL2Key"`
				ExperienceMonenyAccountID       string `json:"experienceMonenyAccountId"`
				ExperienceMonenyL2Key           string `json:"experienceMonenyL2Key"`
				ExperienceMoneyAccountID        string `json:"experienceMoneyAccountId"`
				ExperienceMoneyL2Key            string `json:"experienceMoneyL2Key"`
			} `json:"global"`
			PerpetualContract []struct {
				BaselinePositionValue            string               `json:"baselinePositionValue"`
				CrossID                          int                  `json:"crossId"`
				CrossSymbolID                    int                  `json:"crossSymbolId"`
				CrossSymbolName                  string               `json:"crossSymbolName"`
				DigitMerge                       string               `json:"digitMerge"`
				DisplayMaxLeverage               string               `json:"displayMaxLeverage"`
				DisplayMinLeverage               string               `json:"displayMinLeverage"`
				EnableDisplay                    bool                 `json:"enableDisplay"`
				EnableOpenPosition               bool                 `json:"enableOpenPosition"`
				EnableTrade                      bool                 `json:"enableTrade"`
				FundingImpactMarginNotional      string               `json:"fundingImpactMarginNotional"`
				FundingInterestRate              string               `json:"fundingInterestRate"`
				IncrementalInitialMarginRate     string               `json:"incrementalInitialMarginRate"`
				IncrementalMaintenanceMarginRate string               `json:"incrementalMaintenanceMarginRate"`
				IncrementalPositionValue         string               `json:"incrementalPositionValue"`
				InitialMarginRate                types.Number         `json:"initialMarginRate"`
				MaintenanceMarginRate            types.Number         `json:"maintenanceMarginRate"`
				MaxOrderSize                     types.Number         `json:"maxOrderSize"`
				MaxPositionSize                  types.Number         `json:"maxPositionSize"`
				MinOrderSize                     types.Number         `json:"minOrderSize"`
				MaxMarketPriceRange              types.Number         `json:"maxMarketPriceRange"`
				SettleCurrencyID                 string               `json:"settleCurrencyId"`
				StarkExOraclePriceQuorum         string               `json:"starkExOraclePriceQuorum"`
				StarkExResolution                string               `json:"starkExResolution"`
				StarkExRiskFactor                string               `json:"starkExRiskFactor"`
				StarkExSyntheticAssetID          string               `json:"starkExSyntheticAssetId"`
				StepSize                         types.Number         `json:"stepSize"`
				Symbol                           string               `json:"symbol"`
				SymbolDisplayName                string               `json:"symbolDisplayName"`
				SymbolDisplayName2               string               `json:"symbolDisplayName2"`
				TickSize                         types.Number         `json:"tickSize"`
				UnderlyingCurrencyID             string               `json:"underlyingCurrencyId"`
				MaxMaintenanceMarginRate         string               `json:"maxMaintenanceMarginRate"`
				MaxPositionValue                 string               `json:"maxPositionValue"`
				TagIconURL                       string               `json:"tagIconUrl"`
				Tag                              string               `json:"tag"`
				RiskTip                          bool                 `json:"riskTip"`
				DefaultLeverage                  string               `json:"defaultLeverage"`
				KlineStartTime                   convert.ExchangeTime `json:"klineStartTime"`
			} `json:"perpetualContract"`
			MultiChain struct {
				Chains []struct {
					Chain             string `json:"chain"`
					ChainID           int64  `json:"chainId"`
					ChainIconURL      string `json:"chainIconUrl"`
					ContractAddress   string `json:"contractAddress"`
					DepositGasFeeLess bool   `json:"depositGasFeeLess"`
					StopDeposit       bool   `json:"stopDeposit"`
					FeeLess           bool   `json:"feeLess"`
					FeeRate           string `json:"feeRate"`
					GasLess           bool   `json:"gasLess"`
					GasToken          string `json:"gasToken"`
					MinFee            string `json:"minFee"`
					DynamicFee        bool   `json:"dynamicFee"`
					RPCURL            string `json:"rpcUrl"`
					WebRPCURL         string `json:"webRpcUrl"`
					WebTxURL          string `json:"webTxUrl"`
					BlockTime         string `json:"blockTime"`
					TxConfirm         int64  `json:"txConfirm"`
					Tokens            []struct {
						Decimals       int    `json:"decimals"`
						IconURL        string `json:"iconUrl"`
						Token          string `json:"token"`
						TokenAddress   string `json:"tokenAddress"`
						PullOff        bool   `json:"pullOff"`
						WithdrawEnable bool   `json:"withdrawEnable"`
						Slippage       string `json:"slippage"`
						IsDefaultToken bool   `json:"isDefaultToken"`
						DisplayToken   string `json:"displayToken"`
					} `json:"tokens"`
					WithdrawGasFeeLess bool `json:"withdrawGasFeeLess"`
					IsGray             bool `json:"isGray"`
				} `json:"chains"`
				Currency    string `json:"currency"`
				MaxWithdraw string `json:"maxWithdraw"`
				MinDeposit  string `json:"minDeposit"`
				MinWithdraw string `json:"minWithdraw"`
			} `json:"multiChain"`
		} `json:"usdtConfig"`
	} `json:"data"`
	TimeCost int64 `json:"timeCost"`
}

// AllSymbolsV1Config represents a configuration information
type AllSymbolsV1Config struct {
	Data struct {
		Currency []struct {
			ID                string       `json:"id"`
			StarkExAssetID    string       `json:"starkExAssetId"`
			StarkExResolution string       `json:"starkExResolution"`
			StepSize          types.Number `json:"stepSize"`
			ShowStep          string       `json:"showStep"`
			IconURL           string       `json:"iconUrl"`
		} `json:"currency"`
		Global struct {
			FeeAccountID                    string       `json:"feeAccountId"`
			FeeAccountL2Key                 string       `json:"feeAccountL2Key"`
			StarkExCollateralCurrencyID     string       `json:"starkExCollateralCurrencyId"`
			StarkExFundingValidityPeriod    int64        `json:"starkExFundingValidityPeriod"`
			StarkExMaxFundingRate           types.Number `json:"starkExMaxFundingRate"`
			StarkExOrdersTreeHeight         int64        `json:"starkExOrdersTreeHeight"`
			StarkExPositionsTreeHeight      int64        `json:"starkExPositionsTreeHeight"`
			StarkExPriceValidityPeriod      int64        `json:"starkExPriceValidityPeriod"`
			StarkExContractAddress          string       `json:"starkExContractAddress"`
			RegisterEnvID                   int64        `json:"registerEnvId"`
			CrossChainAccountID             string       `json:"crossChainAccountId"`
			CrossChainL2Key                 string       `json:"crossChainL2Key"`
			FastWithdrawAccountID           string       `json:"fastWithdrawAccountId"`
			FastWithdrawFactRegisterAddress string       `json:"fastWithdrawFactRegisterAddress"`
			FastWithdrawL2Key               string       `json:"fastWithdrawL2Key"`
			FastWithdrawMaxAmount           string       `json:"fastWithdrawMaxAmount"`
		} `json:"global"`
		PerpetualContract []struct {
			BaselinePositionValue            string       `json:"baselinePositionValue"`
			CrossID                          int64        `json:"crossId"`
			CrossSymbolID                    int64        `json:"crossSymbolId"`
			CrossSymbolName                  string       `json:"crossSymbolName"`
			DigitMerge                       string       `json:"digitMerge"`
			DisplayMaxLeverage               types.Number `json:"displayMaxLeverage"`
			DisplayMinLeverage               types.Number `json:"displayMinLeverage"`
			EnableDisplay                    bool         `json:"enableDisplay"`
			EnableOpenPosition               bool         `json:"enableOpenPosition"`
			EnableTrade                      bool         `json:"enableTrade"`
			FundingImpactMarginNotional      string       `json:"fundingImpactMarginNotional"`
			FundingInterestRate              types.Number `json:"fundingInterestRate"`
			IncrementalInitialMarginRate     types.Number `json:"incrementalInitialMarginRate"`
			IncrementalMaintenanceMarginRate string       `json:"incrementalMaintenanceMarginRate"`
			IncrementalPositionValue         string       `json:"incrementalPositionValue"`
			InitialMarginRate                string       `json:"initialMarginRate"`
			MaintenanceMarginRate            string       `json:"maintenanceMarginRate"`
			MaxOrderSize                     string       `json:"maxOrderSize"`
			MaxPositionSize                  string       `json:"maxPositionSize"`
			MinOrderSize                     string       `json:"minOrderSize"`
			MaxMarketPriceRange              string       `json:"maxMarketPriceRange"`
			SettleCurrencyID                 string       `json:"settleCurrencyId"`
			StarkExOraclePriceQuorum         string       `json:"starkExOraclePriceQuorum"`
			StarkExResolution                string       `json:"starkExResolution"`
			StarkExRiskFactor                string       `json:"starkExRiskFactor"`
			StarkExSyntheticAssetID          string       `json:"starkExSyntheticAssetId"`
			StepSize                         string       `json:"stepSize"`
			Symbol                           string       `json:"symbol"`
			SymbolDisplayName                string       `json:"symbolDisplayName"`
			TickSize                         types.Number `json:"tickSize"`
			UnderlyingCurrencyID             string       `json:"underlyingCurrencyId"`
			MaxMaintenanceMarginRate         types.Number `json:"maxMaintenanceMarginRate"`
			MaxPositionValue                 types.Number `json:"maxPositionValue"`
		} `json:"perpetualContract"`
		MultiChain struct {
			Chains []struct {
				Chain             string       `json:"chain"`
				ChainID           int64        `json:"chainId"`
				ChainIconURL      string       `json:"chainIconUrl"`
				ContractAddress   string       `json:"contractAddress"`
				DepositGasFeeLess bool         `json:"depositGasFeeLess"`
				FeeLess           bool         `json:"feeLess"`
				FeeRate           types.Number `json:"feeRate"`
				GasLess           bool         `json:"gasLess"`
				GasToken          string       `json:"gasToken"`
				MinFee            string       `json:"minFee"`
				RPCURL            string       `json:"rpcUrl"`
				WebTxURL          string       `json:"webTxUrl"`
				TxConfirm         int64        `json:"txConfirm"`
				Tokens            []struct {
					Decimals     int64  `json:"decimals"`
					IconURL      string `json:"iconUrl"`
					Token        string `json:"token"`
					TokenAddress string `json:"tokenAddress"`
					PullOff      bool   `json:"pullOff"`
				} `json:"tokens"`
				WithdrawGasFeeLess bool `json:"withdrawGasFeeLess"`
			} `json:"chains"`
			Currency    string       `json:"currency"`
			MaxWithdraw types.Number `json:"maxWithdraw"`
			MinDeposit  types.Number `json:"minDeposit"`
			MinWithdraw types.Number `json:"minWithdraw"`
		} `json:"multiChain"`
	} `json:"data"`
	TimeCost int64 `json:"timeCost"`
}

// NonceResponse represents a nonce response.
type NonceResponse struct {
	Nonce        string               `json:"nonce"`
	NonceExpired convert.ExchangeTime `json:"nonceExpired"`
}

// WsMessage represents a websocket input message.
type WsMessage struct {
	Operation string   `json:"op"`
	Args      []string `json:"args"`
}

// WsDepth represents a websocket orderbook data.
type WsDepth struct {
	Topic string `json:"topic"`
	Type  string `json:"type"`
	Data  struct {
		Symbol   string            `json:"s"`
		Bids     [][2]types.Number `json:"b"`
		Asks     [][2]types.Number `json:"a"`
		UpdateID int64             `json:"u"`
	} `json:"data"`
	Cs        int64                `json:"cs"`
	Timestamp convert.ExchangeTime `json:"ts"`
}

// WsTrade represents a trade data pushed through the websocket stream.
type WsTrade struct {
	Topic string `json:"topic"`
	Type  string `json:"type"`
	Data  []struct {
		Timestamp       convert.ExchangeTime `json:"T"`
		Symbol          string               `json:"s"`
		Side            string               `json:"S"`
		Volume          types.Number         `json:"v"`
		Price           types.Number         `json:"p"`
		TickerDirection string               `json:"L"`
		OrderID         string               `json:"i"`
	} `json:"data"`
	Cs        int64                `json:"cs"`
	Timestamp convert.ExchangeTime `json:"ts"`
}

// WsTicker represents a ticker item data.
type WsTicker struct {
	Topic string `json:"topic"`
	Type  string `json:"type"`
	Data  struct {
		Symbol               string       `json:"symbol"`
		LastPrice            types.Number `json:"lastPrice"`
		Price24HPcnt         types.Number `json:"price24hPcnt"`
		HighPrice24H         types.Number `json:"highPrice24h"`
		LowPrice24H          types.Number `json:"lowPrice24h"`
		Turnover24H          types.Number `json:"turnover24h"`
		Volume24H            types.Number `json:"volume24h"`
		NextFundingTime      time.Time    `json:"nextFundingTime"`
		OraclePrice          types.Number `json:"oraclePrice"`
		IndexPrice           types.Number `json:"indexPrice"`
		OpenInterest         types.Number `json:"openInterest"`
		TradeCount           types.Number `json:"tradeCount"`
		FundingRate          types.Number `json:"fundingRate"`
		PredictedFundingRate types.Number `json:"predictedFundingRate"`
	} `json:"data"`
	Cs int   `json:"cs"`
	Ts int64 `json:"ts"`
}
