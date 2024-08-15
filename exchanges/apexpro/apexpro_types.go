package apexpro

import (
	"time"

	"github.com/thrasher-corp/gocryptotrader/common/convert"
	"github.com/thrasher-corp/gocryptotrader/currency"
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
			IncrementalPositionValue         types.Number         `json:"incrementalPositionValue"`
			InitialMarginRate                types.Number         `json:"initialMarginRate"`
			MaintenanceMarginRate            types.Number         `json:"maintenanceMarginRate"`
			MaxOrderSize                     types.Number         `json:"maxOrderSize"`
			MaxPositionSize                  types.Number         `json:"maxPositionSize"`
			MinOrderSize                     types.Number         `json:"minOrderSize"`
			MaxMarketPriceRange              types.Number         `json:"maxMarketPriceRange"`
			SettleAssetID                    string               `json:"settleAssetId"` // Collateral asset ID.
			BaseTokenID                      string               `json:"baseTokenId"`
			StepSize                         types.Number         `json:"stepSize"`
			Symbol                           string               `json:"symbol"`
			SymbolDisplayName                string               `json:"symbolDisplayName"`
			TickSize                         types.Number         `json:"tickSize"`
			MaxMaintenanceMarginRate         types.Number         `json:"maxMaintenanceMarginRate"`
			MaxPositionValue                 types.Number         `json:"maxPositionValue"`
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
			FundingMaxValue                  types.Number         `json:"fundingMaxValue"`
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
			DisplayMaxLeverage               types.Number         `json:"displayMaxLeverage"`
			DisplayMinLeverage               types.Number         `json:"displayMinLeverage"`
			EnableDisplay                    bool                 `json:"enableDisplay"`
			EnableOpenPosition               bool                 `json:"enableOpenPosition"`
			EnableTrade                      bool                 `json:"enableTrade"`
			FundingImpactMarginNotional      types.Number         `json:"fundingImpactMarginNotional"`
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
			SettleAssetID                    string               `json:"settleAssetId"`
			BaseTokenID                      string               `json:"baseTokenId"`
			StepSize                         types.Number         `json:"stepSize"`
			Symbol                           string               `json:"symbol"`
			SymbolDisplayName                string               `json:"symbolDisplayName"`
			TickSize                         types.Number         `json:"tickSize"`
			MaxMaintenanceMarginRate         types.Number         `json:"maxMaintenanceMarginRate"`
			MaxPositionValue                 types.Number         `json:"maxPositionValue"`
			TagIconURL                       string               `json:"tagIconUrl"`
			Tag                              string               `json:"tag"`
			RiskTip                          bool                 `json:"riskTip"`
			DefaultLeverage                  types.Number         `json:"defaultLeverage"`
			KlineStartTime                   convert.ExchangeTime `json:"klineStartTime"`
			MaxMarketSizeBuffer              string               `json:"maxMarketSizeBuffer"`
			EnableFundingSettlement          bool                 `json:"enableFundingSettlement"`
			IndexPriceDecimals               float64              `json:"indexPriceDecimals"`
			IndexPriceVarRate                types.Number         `json:"indexPriceVarRate"`
			OpenPositionOiLimitRate          types.Number         `json:"openPositionOiLimitRate"`
			FundingMaxRate                   types.Number         `json:"fundingMaxRate"`
			FundingMinRate                   types.Number         `json:"fundingMinRate"`
			FundingMaxValue                  types.Number         `json:"fundingMaxValue"`
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
	Symbol               string               `json:"symbol"`
	Price24HPcnt         types.Number         `json:"price24hPcnt"`
	LastPrice            types.Number         `json:"lastPrice"`
	HighPrice24H         types.Number         `json:"highPrice24h"`
	LowPrice24H          types.Number         `json:"lowPrice24h"`
	MarkPrice            types.Number         `json:"markPrice"`
	IndexPrice           types.Number         `json:"indexPrice"`
	OpenInterest         types.Number         `json:"openInterest"`
	Turnover24H          types.Number         `json:"turnover24h"`
	Volume24H            types.Number         `json:"volume24h"`
	FundingRate          types.Number         `json:"fundingRate"`
	PredictedFundingRate types.Number         `json:"predictedFundingRate"`
	NextFundingTime      convert.ExchangeTime `json:"nextFundingTime"`
	TradeCount           types.Number         `json:"tradeCount"`
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

// CurrencyInfo represents a currency detail.
type CurrencyInfo struct {
	ID                string       `json:"id"` // Settlement Currency ID.
	StarkExAssetID    string       `json:"starkExAssetId"`
	StarkExResolution string       `json:"starkExResolution"`
	StepSize          types.Number `json:"stepSize"`
	ShowStep          string       `json:"showStep"`
	IconURL           string       `json:"iconUrl"`
}

// V2ConfigData v2 assets and symbols configuration response.
type V2ConfigData struct {
	Data struct {
		USDCConfig struct {
			Currency []CurrencyInfo `json:"currency"`
			Global   struct {
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
				Symbol                           string               `json:"symbol"`
				BaselinePositionValue            string               `json:"baselinePositionValue"`
				CrossID                          int64                `json:"crossId"`
				CrossSymbolID                    int64                `json:"crossSymbolId"`
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
				StepSize                         types.Number         `json:"stepSize"`
				SymbolDisplayName                string               `json:"symbolDisplayName"`
				SymbolDisplayName2               string               `json:"symbolDisplayName2"`
				TickSize                         types.Number         `json:"tickSize"`
				UnderlyingCurrencyID             string               `json:"underlyingCurrencyId"`
				MaxMaintenanceMarginRate         types.Number         `json:"maxMaintenanceMarginRate"`
				MaxPositionValue                 types.Number         `json:"maxPositionValue"`
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
		USDTConfig struct {
			Currency []CurrencyInfo `json:"currency"`
			Global   struct {
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

// V1CurrencyConfig represents a V1 currency configuration.
type V1CurrencyConfig struct {
	ID                string       `json:"id"`
	StarkExAssetID    string       `json:"starkExAssetId"`
	StarkExResolution string       `json:"starkExResolution"`
	StepSize          types.Number `json:"stepSize"`
	ShowStep          string       `json:"showStep"`
	IconURL           string       `json:"iconUrl"`
}

// AllSymbolsV1Config represents a configuration information
type AllSymbolsV1Config struct {
	Data struct {
		Currency []V1CurrencyConfig `json:"currency"`
		Global   struct {
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
		PerpetualContract []PerpetualContractDetail `json:"perpetualContract"`
		MultiChain        struct {
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

// PerpetualContractDetail represents a perpetual contract detail.
type PerpetualContractDetail struct {
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
	Cs        int64                `json:"cs"`
	Timestamp convert.ExchangeTime `json:"ts"`
}

// UserData represents an account user information.
type UserData struct {
	EthereumAddress string `json:"ethereumAddress"`
	IsRegistered    bool   `json:"isRegistered"`
	Email           string `json:"email"`
	Username        string `json:"username"`
	UserData        struct {
	} `json:"userData"`
	IsEmailVerified          bool `json:"isEmailVerified"`
	EmailNotifyGeneralEnable bool `json:"emailNotifyGeneralEnable"`
	EmailNotifyTradingEnable bool `json:"emailNotifyTradingEnable"`
	EmailNotifyAccountEnable bool `json:"emailNotifyAccountEnable"`
	PopupNotifyTradingEnable bool `json:"popupNotifyTradingEnable"`
}

// UserResponse represents a user account detail response.
type UserResponse struct {
	Data    interface{} `json:"data"`
	Code    int64       `json:"code"`
	Message string      `json:"msg"`
}

// WsCandlesticks represents a list of candlestick data.
type WsCandlesticks struct {
	Topic     string               `json:"topic"`
	Data      []CandlestickData    `json:"data"`
	Timestamp convert.ExchangeTime `json:"ts"`
	Type      string               `json:"type"`
}

// WsSymbolsTickerInformaton represents a ticker information for assets.
type WsSymbolsTickerInformaton struct {
	Topic string `json:"topic"`
	Data  []struct {
		Symbol                    string       `json:"s"`
		LastPrice                 types.Number `json:"p"`
		Price24HrChangePercentage types.Number `json:"pr"`
		Highest24Hr               types.Number `json:"h"`
		Lowest24Hr                types.Number `json:"l"`
		OpeningPrice              types.Number `json:"op,omitempty"`
		IndexPrice                types.Number `json:"xp"`
		Turnover24Hr              types.Number `json:"to"`
		Volume24Hr                types.Number `json:"v"`
		FundingRate               types.Number `json:"fr"`
		OpenInterest              types.Number `json:"o"`
		TradeCount24Hr            types.Number `json:"tc"`
		MarkPrice                 types.Number `json:"mp,omitempty"`
	} `json:"data"`
	Type      string               `json:"type"`
	Timestamp convert.ExchangeTime `json:"ts"`
}

// RegistrationAndOnboardingResponse represents a registration and onboarding response.
type RegistrationAndOnboardingResponse struct {
	APIKey struct {
		APIKey string   `json:"apiKey"`
		Key    string   `json:"key"`
		Secret string   `json:"secret"`
		Remark string   `json:"remark"`
		Ips    []string `json:"ips"`
	} `json:"apiKey"`
	User struct {
		EthereumAddress          string       `json:"ethereumAddress"`
		IsRegistered             bool         `json:"isRegistered"`
		Email                    string       `json:"email"`
		Username                 string       `json:"username"`
		ReferredByAffiliateLink  string       `json:"referredByAffiliateLink"`
		AffiliateLink            string       `json:"affiliateLink"`
		ApexTokenBalance         types.Number `json:"apexTokenBalance"`
		StakedApexTokenBalance   types.Number `json:"stakedApexTokenBalance"`
		IsEmailVerified          bool         `json:"isEmailVerified"`
		IsSharingUsername        bool         `json:"isSharingUsername"`
		IsSharingAddress         bool         `json:"isSharingAddress"`
		Country                  string       `json:"country"`
		ID                       string       `json:"id"`
		AvatarURL                string       `json:"avatarUrl"`
		AvatarBorderURL          string       `json:"avatarBorderUrl"`
		EmailNotifyGeneralEnable bool         `json:"emailNotifyGeneralEnable"`
		EmailNotifyTradingEnable bool         `json:"emailNotifyTradingEnable"`
		EmailNotifyAccountEnable bool         `json:"emailNotifyAccountEnable"`
		PopupNotifyTradingEnable bool         `json:"popupNotifyTradingEnable"`
		AppNotifyTradingEnable   bool         `json:"appNotifyTradingEnable"`
	} `json:"user"`
	Account struct {
		EthereumAddress string `json:"ethereumAddress"`
		L2Key           string `json:"l2Key"`
		ID              string `json:"id"`
		Version         string `json:"version"`
		SpotAccount     struct {
			CreatedAt            convert.ExchangeTime `json:"createdAt"`
			UpdatedAt            convert.ExchangeTime `json:"updatedAt"`
			ZkAccountID          string               `json:"zkAccountId"`
			IsMultiSigEthAddress bool                 `json:"isMultiSigEthAddress"`
			DefaultSubAccountID  string               `json:"defaultSubAccountId"`
			Nonce                int                  `json:"nonce"`
			Status               string               `json:"status"`
			SubAccounts          []struct {
				SubAccountID       string `json:"subAccountId"`
				L2Key              string `json:"l2Key"`
				Nonce              int    `json:"nonce"`
				NonceVersion       int    `json:"nonceVersion"`
				ChangePubKeyStatus string `json:"changePubKeyStatus"`
			} `json:"subAccounts"`
		} `json:"spotAccount"`
		SpotWallets []struct {
			UserID                   string               `json:"userId"`
			AccountID                string               `json:"accountId"`
			SubAccountID             string               `json:"subAccountId"`
			Balance                  types.Number         `json:"balance"`
			TokenID                  string               `json:"tokenId"`
			PendingDepositAmount     types.Number         `json:"pendingDepositAmount"`
			PendingWithdrawAmount    types.Number         `json:"pendingWithdrawAmount"`
			PendingTransferOutAmount types.Number         `json:"pendingTransferOutAmount"`
			PendingTransferInAmount  types.Number         `json:"pendingTransferInAmount"`
			CreatedAt                convert.ExchangeTime `json:"createdAt"`
			UpdatedAt                convert.ExchangeTime `json:"updatedAt"`
		} `json:"spotWallets"`
		ExperienceMoney []struct {
			AvailableAmount types.Number `json:"availableAmount"`
			TotalNumber     types.Number `json:"totalNumber"`
			TotalAmount     types.Number `json:"totalAmount"`
			RecycledAmount  types.Number `json:"recycledAmount"`
			Token           string       `json:"token"`
		} `json:"experienceMoney"`
		ContractAccount struct {
			CreatedAt             convert.ExchangeTime `json:"createdAt"`
			TakerFeeRate          types.Number         `json:"takerFeeRate"`
			MakerFeeRate          types.Number         `json:"makerFeeRate"`
			MinInitialMarginRate  types.Number         `json:"minInitialMarginRate"`
			Status                string               `json:"status"`
			UnrealizePnlPriceType string               `json:"unrealizePnlPriceType"`
		} `json:"contractAccount"`
		ContractWallets []struct {
			UserID                   string       `json:"userId"`
			AccountID                string       `json:"accountId"`
			Balance                  types.Number `json:"balance"`
			Asset                    string       `json:"asset"`
			PendingDepositAmount     types.Number `json:"pendingDepositAmount"`
			PendingWithdrawAmount    types.Number `json:"pendingWithdrawAmount"`
			PendingTransferOutAmount types.Number `json:"pendingTransferOutAmount"`
			PendingTransferInAmount  types.Number `json:"pendingTransferInAmount"`
		} `json:"contractWallets"`
		Positions []struct {
			IsPrelaunch             bool                 `json:"isPrelaunch"`
			Symbol                  string               `json:"symbol"`
			Status                  string               `json:"status"`
			Side                    string               `json:"side"`
			Size                    types.Number         `json:"size"`
			EntryPrice              types.Number         `json:"entryPrice"`
			ExitPrice               types.Number         `json:"exitPrice"`
			CreatedAt               convert.ExchangeTime `json:"createdAt"`
			UpdatedTime             convert.ExchangeTime `json:"updatedTime"`
			Fee                     types.Number         `json:"fee"`
			FundingFee              types.Number         `json:"fundingFee"`
			LightNumbers            types.Number         `json:"lightNumbers"`
			CustomInitialMarginRate string               `json:"customInitialMarginRate"`
		} `json:"positions"`
		IsNewUser bool `json:"isNewUser"`
	} `json:"account"`
}

// EditUserDataParams represents a request parameter to edit user data.
type EditUserDataParams struct {
	Email                    string `json:"email"`
	UserData                 string `json:"userData"`
	Username                 string `json:"username"`
	IsSharingUsername        bool   `json:"isSharingUsername"`
	IsSharingAddress         bool   `json:"isSharingAddress"`
	Country                  string `json:"country"`
	EmailNotifyGeneralEnable bool   `json:"emailNotifyGeneralEnable"`
	EmailNotifyTradingEnable bool   `json:"emailNotifyTradingEnable"`
	EmailNotifyAccountEnable bool   `json:"emailNotifyAccountEnable"`
	PopupNotifyTradingEnable bool   `json:"popupNotifyTradingEnable"`
}

// UserDataResponse represents a user data response.
type UserDataResponse struct {
	EthereumAddress string `json:"ethereumAddress"`
	IsRegistered    bool   `json:"isRegistered"`
	Email           string `json:"email"`
	Username        string `json:"username"`
	UserData        struct {
	} `json:"userData"`
	IsEmailVerified          bool `json:"isEmailVerified"`
	EmailNotifyGeneralEnable bool `json:"emailNotifyGeneralEnable"`
	EmailNotifyTradingEnable bool `json:"emailNotifyTradingEnable"`
	EmailNotifyAccountEnable bool `json:"emailNotifyAccountEnable"`
	PopupNotifyTradingEnable bool `json:"popupNotifyTradingEnable"`
}

// UserAccountV2 represents a V2 user account detail.
type UserAccountV2 struct {
	StarkKey        string `json:"starkKey"`
	PositionID      string `json:"positionId"`
	EthereumAddress string `json:"ethereumAddress"`
	ID              string `json:"id"`
	ExperienceMoney []struct {
		AvailableAmount types.Number `json:"availableAmount"`
		TotalNumber     types.Number `json:"totalNumber"`
		TotalAmount     types.Number `json:"totalAmount"`
		RecycledAmount  types.Number `json:"recycledAmount"`
		Token           string       `json:"token"`
	} `json:"experienceMoney"`
	Accounts []struct {
		CreatedAt             convert.ExchangeTime `json:"createdAt"`
		TakerFeeRate          types.Number         `json:"takerFeeRate"`
		MakerFeeRate          types.Number         `json:"makerFeeRate"`
		MinInitialMarginRate  types.Number         `json:"minInitialMarginRate"`
		Status                string               `json:"status"`
		Token                 string               `json:"token"`
		UnrealizePnlPriceType string               `json:"unrealizePnlPriceType"`
	} `json:"accounts"`
	Wallets   any `json:"wallets"`
	Positions []struct {
		Token                   string               `json:"token"`
		Symbol                  string               `json:"symbol"`
		Status                  string               `json:"status"`
		Side                    string               `json:"side"`
		Size                    types.Number         `json:"size"`
		EntryPrice              types.Number         `json:"entryPrice"`
		ExitPrice               types.Number         `json:"exitPrice"`
		CreatedAt               convert.ExchangeTime `json:"createdAt"`
		UpdatedTime             convert.ExchangeTime `json:"updatedTime"`
		Fee                     types.Number         `json:"fee"`
		FundingFee              types.Number         `json:"fundingFee"`
		LightNumbers            string               `json:"lightNumbers"`
		CustomInitialMarginRate types.Number         `json:"customInitialMarginRate"`
	} `json:"positions"`
}

// UserAccountDetail represents a user account detail.
type UserAccountDetail struct {
	EthereumAddress string `json:"ethereumAddress"`
	L2Key           string `json:"l2Key"`
	ID              string `json:"id"` // position ID or account ID
	Version         string `json:"version"`
	SpotAccount     struct {
		CreatedAt            convert.ExchangeTime `json:"createdAt"`
		UpdatedAt            convert.ExchangeTime `json:"updatedAt"`
		ZkAccountID          string               `json:"zkAccountId"`
		IsMultiSigEthAddress bool                 `json:"isMultiSigEthAddress"`
		DefaultSubAccountID  string               `json:"defaultSubAccountId"`
		Nonce                int64                `json:"nonce"`
		Status               string               `json:"status"`
		SubAccounts          []struct {
			SubAccountID       string `json:"subAccountId"`
			L2Key              string `json:"l2Key"`
			Nonce              int64  `json:"nonce"`
			NonceVersion       int64  `json:"nonceVersion"`
			ChangePubKeyStatus string `json:"changePubKeyStatus"`
		} `json:"subAccounts"`
	} `json:"spotAccount"`
	SpotWallets []struct {
		UserID                   string               `json:"userId"`
		AccountID                string               `json:"accountId"`
		SubAccountID             string               `json:"subAccountId"`
		Balance                  types.Number         `json:"balance"`
		TokenID                  string               `json:"tokenId"`
		PendingDepositAmount     types.Number         `json:"pendingDepositAmount"`
		PendingWithdrawAmount    types.Number         `json:"pendingWithdrawAmount"`
		PendingTransferOutAmount types.Number         `json:"pendingTransferOutAmount"`
		PendingTransferInAmount  types.Number         `json:"pendingTransferInAmount"`
		CreatedAt                convert.ExchangeTime `json:"createdAt"`
		UpdatedAt                convert.ExchangeTime `json:"updatedAt"`
	} `json:"spotWallets"`
	ExperienceMoney []struct {
		AvailableAmount types.Number `json:"availableAmount"`
		TotalNumber     types.Number `json:"totalNumber"`
		TotalAmount     types.Number `json:"totalAmount"`
		RecycledAmount  types.Number `json:"recycledAmount"`
		Token           string       `json:"token"`
	} `json:"experienceMoney"`
	ContractAccount struct {
		CreatedAt             convert.ExchangeTime `json:"createdAt"`
		TakerFeeRate          types.Number         `json:"takerFeeRate"`
		MakerFeeRate          types.Number         `json:"makerFeeRate"`
		MinInitialMarginRate  types.Number         `json:"minInitialMarginRate"`
		Status                string               `json:"status"`
		UnrealizePnlPriceType string               `json:"unrealizePnlPriceType"`
		Token                 string               `json:"token"`
	} `json:"contractAccount"`
	ContractWallets []struct {
		UserID                   string       `json:"userId"`
		AccountID                string       `json:"accountId"`
		Asset                    string       `json:"asset"`
		Balance                  types.Number `json:"balance"`
		PendingDepositAmount     types.Number `json:"pendingDepositAmount"`
		PendingWithdrawAmount    types.Number `json:"pendingWithdrawAmount"`
		PendingTransferOutAmount types.Number `json:"pendingTransferOutAmount"`
		PendingTransferInAmount  types.Number `json:"pendingTransferInAmount"`
	} `json:"contractWallets"`
	Positions []struct {
		IsPrelaunch             bool                 `json:"isPrelaunch"`
		Symbol                  string               `json:"symbol"`
		Status                  string               `json:"status"`
		Side                    string               `json:"side"`
		Size                    types.Number         `json:"size"`
		EntryPrice              types.Number         `json:"entryPrice"`
		ExitPrice               string               `json:"exitPrice"`
		CreatedAt               convert.ExchangeTime `json:"createdAt"`
		UpdatedTime             convert.ExchangeTime `json:"updatedTime"`
		Fee                     types.Number         `json:"fee"`
		FundingFee              types.Number         `json:"fundingFee"`
		LightNumbers            string               `json:"lightNumbers"`
		CustomInitialMarginRate string               `json:"customInitialMarginRate"`
	} `json:"positions"`
	IsNewUser bool `json:"isNewUser"`
}

// UserAccountDetailV1 represents a user account detail throught the v1 API endpoint.
type UserAccountDetailV1 struct {
	StarkKey     string               `json:"starkKey"`
	PositionID   string               `json:"positionId"`
	TakerFeeRate types.Number         `json:"takerFeeRate"`
	MakerFeeRate types.Number         `json:"makerFeeRate"`
	CreatedAt    convert.ExchangeTime `json:"createdAt"`
	Wallets      []struct {
		UserID                   string       `json:"userId"`
		AccountID                string       `json:"accountId"`
		Asset                    string       `json:"asset"`
		Balance                  types.Number `json:"balance"`
		PendingDepositAmount     types.Number `json:"pendingDepositAmount"`
		PendingWithdrawAmount    types.Number `json:"pendingWithdrawAmount"`
		PendingTransferOutAmount types.Number `json:"pendingTransferOutAmount"`
		PendingTransferInAmount  types.Number `json:"pendingTransferInAmount"`
	} `json:"wallets"`
	OpenPositions []struct {
		Symbol       string               `json:"symbol"`
		Side         string               `json:"side"`
		Size         types.Number         `json:"size"`
		EntryPrice   types.Number         `json:"entryPrice"`
		Fee          types.Number         `json:"fee"`
		FundingFee   types.Number         `json:"fundingFee"`
		CreatedAt    convert.ExchangeTime `json:"createdAt"`
		UpdatedTime  convert.ExchangeTime `json:"updatedTime"`
		LightNumbers string               `json:"lightNumbers"`
	} `json:"openPositions"`
	ID string `json:"id"`
}

// UserAccountBalanceResponse represents a user account balance.
type UserAccountBalanceResponse struct {
	TotalEquityValue    types.Number `json:"totalEquityValue"`
	AvailableBalance    types.Number `json:"availableBalance"`
	InitialMargin       types.Number `json:"initialMargin"`
	MaintenanceMargin   types.Number `json:"maintenanceMargin"`
	SymbolToOraclePrice map[string]struct {
		OraclePrice types.Number         `json:"oraclePrice"`
		CreatedTime convert.ExchangeTime `json:"createdTime"`
	} `json:"symbolToOraclePrice"`
}

type AutoGenerated struct {
	TotalEquityValue    string `json:"totalEquityValue"`
	AvailableBalance    string `json:"availableBalance"`
	InitialMargin       string `json:"initialMargin"`
	MaintenanceMargin   string `json:"maintenanceMargin"`
	SymbolToOraclePrice struct {
		BTCUSDC struct {
			OraclePrice string `json:"oraclePrice"`
			CreatedTime int    `json:"createdTime"`
		} `json:"BTC-USDC"`
	} `json:"symbolToOraclePrice"`
}

// UserAccountBalanceV2Response represents a V2 user account balance information.
type UserAccountBalanceV2Response struct {
	USDTBalance *UserAccountBalanceResponse `json:"usdtBalance"`
	USDCBalance *UserAccountBalanceResponse `json:"usdcBalance"`
}

// UserWithdrawals represents users withdrawals list.
type UserWithdrawals struct {
	Transfers []UserWithdrawal `json:"transfers"`
}

// UserWithdrawalsV2 represents users withdrawals list.
type UserWithdrawalsV2 struct {
	Transfers []UserWithdrawalV2 `json:"transfers"`
}

// UserWithdrawalV2 represents a user asset withdrawal info
type UserWithdrawalV2 struct {
	ID              string               `json:"id"`
	Type            string               `json:"type"`
	CurrencyID      string               `json:"currencyId"`
	Amount          types.Number         `json:"amount"`
	TransactionHash string               `json:"transactionHash"`
	Status          string               `json:"status"`
	CreatedAt       convert.ExchangeTime `json:"createdAt"`
	UpdatedTime     convert.ExchangeTime `json:"updatedTime"`
	ConfirmedAt     convert.ExchangeTime `json:"confirmedAt"`
	ClientID        string               `json:"clientId"`
	ConfirmedCount  int64                `json:"confirmedCount"`
	RequiredCount   int64                `json:"requiredCount"`
	OrderID         string               `json:"orderId"`
	ChainID         string               `json:"chainId"`
	Fee             types.Number         `json:"fee"`
}

// UserWithdrawal represents a user withdrawal information.
type UserWithdrawal struct {
	ID              string               `json:"id"`
	Type            string               `json:"type"`
	Amount          types.Number         `json:"amount"`
	TransactionHash string               `json:"transactionHash"`
	Status          string               `json:"status"`
	CreatedAt       convert.ExchangeTime `json:"createdAt"`
	UpdatedAt       convert.ExchangeTime `json:"updatedAt"`
	ConfirmedAt     convert.ExchangeTime `json:"confirmedAt"`
	FromTokenID     string               `json:"fromTokenId"`
	ToTokenID       string               `json:"toTokenId"`
	ChainID         string               `json:"chainId"`
	OrderID         string               `json:"orderId"`
	EthAddress      string               `json:"ethAddress"`
	FromEthAddress  string               `json:"fromEthAddress"`
	ToEthAddress    string               `json:"toEthAddress"`
	Fee             types.Number         `json:"fee"`
	ClientID        string               `json:"clientId"`
}

// WithdrawalToAddressParams represents a withdrawal parameter to an address through the V2 API
type WithdrawalToAddressParams struct {
	Amount          float64
	ClientID        string
	ExpirationTime  time.Time
	Asset           currency.Code
	EthereumAddress string
}

// AssetWithdrawalParams represents a user asset withdrawal parameter.
type AssetWithdrawalParams struct {
	Amount           float64
	ClientWithdrawID string
	Timestamp        time.Time
	EthereumAddress  string
	Signature        string
	ZKAccountID      string
	SubAccountID     string
	L2Key            string
	ToChainID        string
	L2SourceTokenID  currency.Code // L2 currency(Token ID). Eg. 'USDT' or 'USDC'
	L1TargetTokenID  currency.Code // L1 currency(Token ID). Eg. 'USDT' or 'USDC'
	Fee              float64
	Nonce            string
	IsFastWithdraw   bool
}

// WithdrawalResponse represents a withdrawal placing response.
type WithdrawalResponse struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

// WithdrawalFeeInfos represents an asset withdrawal fee information
type WithdrawalFeeInfos struct {
	WithdrawFeeAndPoolBalances []struct {
		ChainID                 string       `json:"chainId"`
		TokenID                 string       `json:"tokenId"`
		Fee                     types.Number `json:"fee"`
		ZkAvailableAmount       types.Number `json:"zkAvailableAmount"`
		FastpoolAvailableAmount types.Number `json:"fastpoolAvailableAmount"`
	} `json:"withdrawFeeAndPoolBalances"`
}

// ContractTransferLimit represents a contract transfer limit detail.
type ContractTransferLimit struct {
	WithdrawAvailableAmount          types.Number `json:"withdrawAvailableAmount"`
	TransferAvailableAmount          types.Number `json:"transferAvailableAmount"`
	ExperienceMoneyAvailableAmount   types.Number `json:"experienceMoneyAvailableAmount"`
	ExperienceMoneyRecycledAmount    types.Number `json:"experienceMoneyRecycledAmount"`
	WithdrawAvailableOriginAmount    types.Number `json:"withdrawAvailableOriginAmount"`
	ExperienceMoneyNeedRecycleAmount types.Number `json:"experienceMoneyNeedRecycleAmount"`
}

// TradeHistory represents a trade history
type TradeHistory struct {
	Orders    []TradeFill `json:"orders"`
	TotalSize int64       `json:"totalSize"`
}

// TradeFill  represents a trade fill information.
type TradeFill struct {
	ID                   string               `json:"id"`
	ClientID             string               `json:"clientId"`
	AccountID            string               `json:"accountId"`
	Symbol               string               `json:"symbol"`
	Side                 string               `json:"side"`
	Price                types.Number         `json:"price"`
	LimitFee             types.Number         `json:"limitFee"`
	Fee                  types.Number         `json:"fee"`
	TriggerPrice         types.Number         `json:"triggerPrice"`
	TrailingPercent      types.Number         `json:"trailingPercent"`
	Size                 types.Number         `json:"size"`
	Type                 string               `json:"type"`
	CreatedAt            convert.ExchangeTime `json:"createdAt"`
	UpdatedTime          convert.ExchangeTime `json:"updatedTime"`
	ExpiresAt            convert.ExchangeTime `json:"expiresAt"`
	Status               string               `json:"status"`
	TimeInForce          string               `json:"timeInForce"`
	PostOnly             bool                 `json:"postOnly"`
	ReduceOnly           bool                 `json:"reduceOnly"`
	LatestMatchFillPrice string               `json:"latestMatchFillPrice"`
	CumMatchFillSize     types.Number         `json:"cumMatchFillSize"`
	CumMatchFillValue    types.Number         `json:"cumMatchFillValue"`
	CumMatchFillFee      types.Number         `json:"cumMatchFillFee"`
	CumSuccessFillSize   types.Number         `json:"cumSuccessFillSize"`
	CumSuccessFillValue  types.Number         `json:"cumSuccessFillValue"`
	CumSuccessFillFee    types.Number         `json:"cumSuccessFillFee"`
}

// SymbolWorstPrice represents a worst price of a contract.
type SymbolWorstPrice struct {
	WorstPrice  types.Number `json:"worstPrice"`
	BidOnePrice types.Number `json:"bidOnePrice"`
	AskOnePrice types.Number `json:"askOnePrice"`
}

// OrderDetail represents an order detail.
type OrderDetail struct {
	ID              string               `json:"id"`
	ClientOrderID   string               `json:"clientOrderId"`
	AccountID       string               `json:"accountId"`
	Symbol          string               `json:"symbol"`
	Side            string               `json:"side"`
	Price           types.Number         `json:"price"`
	TriggerPrice    types.Number         `json:"triggerPrice"`
	TrailingPercent string               `json:"trailingPercent"`
	Size            types.Number         `json:"size"`
	OrderType       string               `json:"type"`
	CreatedAt       convert.ExchangeTime `json:"createdAt"`
	ExpiresAt       convert.ExchangeTime `json:"expiresAt"`
	Status          string               `json:"status"`
	TimeInForce     string               `json:"timeInForce"`
	PostOnly        bool                 `json:"postOnly"`

	// Included in the V3 API response.
	LimitFee             types.Number         `json:"limitFee"`
	Fee                  types.Number         `json:"fee"`
	UpdatedTime          convert.ExchangeTime `json:"updatedTime"`
	ReduceOnly           bool                 `json:"reduceOnly"`
	LatestMatchFillPrice types.Number         `json:"latestMatchFillPrice"`
	CumMatchFillSize     types.Number         `json:"cumMatchFillSize"`
	CumMatchFillValue    types.Number         `json:"cumMatchFillValue"`
	CumMatchFillFee      types.Number         `json:"cumMatchFillFee"`
	CumSuccessFillSize   types.Number         `json:"cumSuccessFillSize"`
	CumSuccessFillValue  types.Number         `json:"cumSuccessFillValue"`
	CumSuccessFillFee    types.Number         `json:"cumSuccessFillFee"`

	// used by the V1 API endpoint response.
	UnfillableAt convert.ExchangeTime `json:"unfillableAt"`
	CancelReason string               `json:"cancelReason"`
}

// OrderHistoryResponse represents list of order.
type OrderHistoryResponse struct {
	Orders    []OrderDetail `json:"orders"`
	TotalSize int64         `json:"totalSize"`
}

// FundingRateResponse represents a list of funding rates.
type FundingRateResponse struct {
	FundingValues []struct {
		ID            string               `json:"id"`
		Symbol        string               `json:"symbol"`
		FundingValue  string               `json:"fundingValue"`
		Rate          types.Number         `json:"rate"`
		PositionSize  types.Number         `json:"positionSize"`
		Price         types.Number         `json:"price"`
		Side          string               `json:"side"`
		Status        string               `json:"status"`
		FundingTime   convert.ExchangeTime `json:"fundingTime"`
		TransactionID string               `json:"transactionId"`
	} `json:"fundingValues"`
	TotalSize int64 `json:"totalSize"`
}

// PNLHistory represents positions profit and loss(PNL) history
type PNLHistory struct {
	HistoricalPnl []PNLDetail `json:"historicalPnl"`
	TotalSize     int64       `json:"totalSize"`
}

// PNLDetail represents a profit and loss information of a symbol
type PNLDetail struct {
	Symbol       string               `json:"symbol"`
	Size         types.Number         `json:"size"`
	TotalPnl     types.Number         `json:"totalPnl"`
	Price        types.Number         `json:"price"`
	CreatedAt    convert.ExchangeTime `json:"createdAt"`
	OrderType    string               `json:"type"`
	IsLiquidate  bool                 `json:"isLiquidate"`
	IsDeleverage bool                 `json:"isDeleverage"`
}

// AssetValueHistory represents a historical value of an asset.
type AssetValueHistory struct {
	HistoryValues []struct {
		AccountTotalValue types.Number         `json:"accountTotalValue"`
		DateTime          convert.ExchangeTime `json:"dateTime"`
	} `json:"historyValues"`
}

// WithdrawalsV2 represents an asset withdrawal details
type WithdrawalsV2 struct {
	Transfers []struct {
		ID              string               `json:"id"`
		Type            string               `json:"type"`
		CurrencyID      string               `json:"currencyId"`
		Amount          types.Number         `json:"amount"`
		TransactionHash string               `json:"transactionHash"`
		Status          string               `json:"status"`
		CreatedAt       convert.ExchangeTime `json:"createdAt"`
		UpdatedTime     convert.ExchangeTime `json:"updatedTime"`
		ConfirmedAt     convert.ExchangeTime `json:"confirmedAt"`
		ClientID        string               `json:"clientId"`
		OrderID         string               `json:"orderId"`
		ChainID         string               `json:"chainId"`
		Fee             types.Number         `json:"fee"`
	} `json:"transfers"`
	TotalSize int64 `json:"totalSize"`
}

// FastAndCrossChainWithdrawalFees represents a fast and cross-chain uncommon withdrawal fees
type FastAndCrossChainWithdrawalFees struct {
	Fee                 types.Number `json:"fee"`
	PoolAvailableAmount types.Number `json:"poolAvailableAmount"`
}

// TransferAndWithdrawalLimit represents an asset transfer and withdrawal limit detail.
type TransferAndWithdrawalLimit struct {
	WithdrawAvailableAmount types.Number `json:"withdrawAvailableAmount"`
	TransferAvailableAmount types.Number `json:"transferAvailableAmount"`
}

// CreateOrderParams represents a request parameter for creating order.
type CreateOrderParams struct {
	Symbol          currency.Pair
	Side            string
	OrderType       string
	Size            float64
	Price           float64
	LimitFee        float64
	Nonce           string
	ExpirationTime  time.Time
	TimeInForce     string
	TriggerPrice    float64
	TrailingPercent float64
	ClientOrderID   int64
	ReduceOnly      bool
}

// WithdrawalParams represents an asset withdrawal parameters
type WithdrawalParams struct {
	Amount         float64
	ClientID       string
	ExpirationTime time.Time
	Asset          currency.Code
}

// FastWithdrawalParams represents a cross-chain withdrawal parameters
type FastWithdrawalParams struct {
	Amount       float64
	ClientID     string
	Expiration   time.Time
	Asset        currency.Code
	ERC20Address string
	ChainID      string
	Fees         float64
	IPAccountID  string
}

// Withdraw to Address : 7 | OK
// /v1/create-withdrawal-to-address
// /v2/create-withdrawal-to-address

// Cross-Chain Transfer : 4 | OK
// /v1/cross-chain-withdraw
// /v2/cross-chain-withdraw

// Create Order : 3 | OK
// /v1/create-order
// /v2/create-order

// Conditional Transfer : 5 | OK
// /v1/fast-withdraw
// /v2/fast-withdraw
