package kucoin

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/thrasher-corp/gocryptotrader/currency"
	"github.com/thrasher-corp/gocryptotrader/exchanges/orderbook"
)

var (
	validPeriods = []string{
		"1min", "3min", "5min", "15min", "30min", "1hour", "2hour", "4hour", "6hour", "8hour", "12hour", "1day", "1week",
	}

	errInvalidResponseReciever = errors.New("invalid response receiver")
	errInvalidPrice            = errors.New("invalid price")
	errInvalidSize             = errors.New("invalid size")
	errMalformedData           = errors.New("malformed data")
)

var offlineTradeFee = map[currency.Code]float64{
	currency.BTC:   0.0005,
	currency.ETH:   0.005,
	currency.BNB:   0.01,
	currency.USDT:  25,
	currency.SOL:   0.01,
	currency.ADA:   1.000,
	currency.XRP:   0.5,
	currency.DOT:   0.1,
	currency.USDC:  20,
	currency.DOGE:  20,
	currency.AVAX:  0.01,
	currency.SHIB:  600000,
	currency.LUNA:  0.15,
	currency.LTC:   0.001,
	currency.CRO:   50,
	currency.UNI:   1.2,
	currency.BUSD:  1,
	currency.LINK:  1,
	currency.MATIC: 10,
	currency.ALGO:  0.1,
	currency.BCH:   0.01,
	currency.VET:   30,
	currency.XLM:   0.02,
	currency.ICP:   0.0005,
	currency.AXS:   0.6,
	currency.EGLD:  0.005,
	currency.TRX:   1.5,
	currency.FTT:   0.35,
	currency.UST:   4,
	currency.MANA:  10,
	currency.THETA: 0.2,
	currency.ETC:   0.01,
	currency.FIL:   0.01,
	currency.ATOM:  0.01,
	currency.DAI:   6,
	currency.APE:   1.5,
	currency.HBAR:  3,
	currency.NEAR:  0.01,
	currency.FTM:   30,
	currency.XTZ:   0.2,
	currency.XCN:   100,
	currency.HNT:   0.05,
	currency.XMR:   0.001,
	currency.GRT:   60,
	currency.EOS:   0.2,
	currency.FLOW:  0.05,
	currency.KLAY:  0.5,
	currency.SAND:  10,
	currency.CAKE:  0.05,
	currency.AAVE:  0.2,
	currency.LRC:   20,
	currency.XEC:   5000,
	currency.KSM:   0.01,
	currency.ONE:   100,
	currency.MKR:   0.0075,
	currency.KDA:   0.5,
	currency.BSV:   0.01,
	currency.BTT:   300000,
	currency.NEO:   0,
	currency.RUNE:  0.05,
	currency.USDD:  1,
	currency.QNT:   0.04,
	currency.CHZ:   4,
	currency.STX:   1.5,
	currency.ZEC:   0.005,
	currency.WAVES: 0.002,
	currency.AR:    0.02,
	currency.AMP:   2100,
	currency.DASH:  0.002,
	currency.KCS:   0.75,
	currency.CELO:  0.1,
	currency.COMP:  0.15,
	currency.TFUEL: 5,
	currency.CRV:   8.5,
	currency.XEM:   4,
	currency.BAT:   25,
	currency.HT:    0.1,
	currency.IMX:   10,
	currency.QTUM:  0.01,
	currency.DCR:   0.01,
	currency.ICX:   1,
	currency.OMG:   3,
	currency.TUSD:  15,
	currency.RVN:   2,
	currency.ROSE:  0.1,
	currency.ZEN:   0.002,
	currency.ZIL:   10,
	currency.SUSHI: 5,
	currency.AUDIO: 24,
	currency.LPT:   0.85,
	currency.XDC:   2,
	currency.SCRT:  0.25,
	currency.UMA:   3.5,
	currency.VLX:   10,
	currency.ANKR:  275,
	currency.GMT:   0.5,
	currency.PERP:  7.5,
	currency.TEL:   5500,
	currency.SNX:   6,
}

// UnmarshalTo acts as interface to exchange API response
type UnmarshalTo interface {
	GetError() error
}

// Error defines all error information for each request
type Error struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
}

// GetError checks and returns an error if it is supplied.
func (e Error) GetError() error {
	code, err := strconv.ParseInt(e.Code, 10, 64)
	if err != nil {
		return err
	}
	switch code {
	case 200000, 200:
		return nil
	default:
		return fmt.Errorf("Code: %s Message: %s", e.Code, e.Msg)
	}
}

// kucoinTimeMilliSec provides an internal conversion helper
type kucoinTimeMilliSec time.Time

// UnmarshalJSON is custom type json unmarshaller for kucoinTimeMilliSec
func (k *kucoinTimeMilliSec) UnmarshalJSON(data []byte) error {
	var timestamp int64
	err := json.Unmarshal(data, &timestamp)
	if err != nil {
		return err
	}
	*k = kucoinTimeMilliSec(time.UnixMilli(timestamp))
	return nil
}

// Time returns a time.Time object
func (k kucoinTimeMilliSec) Time() time.Time {
	return time.Time(k)
}

// kucoinTimeMilliSecStr provides an internal conversion helper
type kucoinTimeMilliSecStr time.Time

// UnmarshalJSON is custom type json unmarshaller for kucoinTimeMilliSecStr
func (k *kucoinTimeMilliSecStr) UnmarshalJSON(data []byte) error {
	var timestamp string
	err := json.Unmarshal(data, &timestamp)
	if err != nil {
		return err
	}

	t, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		return err
	}
	*k = kucoinTimeMilliSecStr(time.UnixMilli(t))
	return nil
}

// Time returns a time.Time object
func (k kucoinTimeMilliSecStr) Time() time.Time {
	return time.Time(k)
}

// kucoinTimeNanoSec provides an internal conversion helper
type kucoinTimeNanoSec time.Time

// UnmarshalJSON is custom type json unmarshaller for kucoinTimeNanoSec
func (k *kucoinTimeNanoSec) UnmarshalJSON(data []byte) error {
	var timestamp int64
	err := json.Unmarshal(data, &timestamp)
	if err != nil {
		return err
	}
	*k = kucoinTimeNanoSec(time.Unix(0, timestamp))
	return nil
}

// Time returns a time.Time object
func (k kucoinTimeNanoSec) Time() time.Time {
	return time.Time(k)
}

// SymbolInfo stores symbol information
type SymbolInfo struct {
	Symbol          string  `json:"symbol"`
	Name            string  `json:"name"`
	BaseCurrency    string  `json:"baseCurrency"`
	QuoteCurrency   string  `json:"quoteCurrency"`
	FeeCurrency     string  `json:"feeCurrency"`
	Market          string  `json:"market"`
	BaseMinSize     float64 `json:"baseMinSize,string"`
	QuoteMinSize    float64 `json:"quoteMinSize,string"`
	BaseMaxSize     float64 `json:"baseMaxSize,string"`
	QuoteMaxSize    float64 `json:"quoteMaxSize,string"`
	BaseIncrement   float64 `json:"baseIncrement,string"`
	QuoteIncrement  float64 `json:"quoteIncrement,string"`
	PriceIncrement  float64 `json:"priceIncrement,string"`
	PriceLimitRate  float64 `json:"priceLimitRate,string"`
	MinFunds        float64 `json:"minFunds,string"`
	IsMarginEnabled bool    `json:"isMarginEnabled"`
	EnableTrading   bool    `json:"enableTrading"`
}

// Ticker stores ticker data
type Ticker struct {
	Sequence    string  `json:"sequence"`
	BestAsk     float64 `json:"bestAsk,string"`
	Size        float64 `json:"size,string"`
	Price       float64 `json:"price,string"`
	BestBidSize float64 `json:"bestBidSize,string"`
	BestBid     float64 `json:"bestBid,string"`
	BestAskSize float64 `json:"bestAskSize,string"`
	Time        uint64  `json:"time"`
}

type tickerInfoBase struct {
	Symbol           string  `json:"symbol"`
	Buy              float64 `json:"buy,string"`
	Sell             float64 `json:"sell,string"`
	ChangeRate       float64 `json:"changeRate,string"`
	ChangePrice      float64 `json:"changePrice,string"`
	High             float64 `json:"high,string"`
	Low              float64 `json:"low,string"`
	Volume           float64 `json:"vol,string"`
	VolumeValue      float64 `json:"volValue,string"`
	Last             float64 `json:"last,string"`
	AveragePrice     float64 `json:"averagePrice,string"`
	TakerFeeRate     float64 `json:"takerFeeRate,string"`
	MakerFeeRate     float64 `json:"makerFeeRate,string"`
	TakerCoefficient float64 `json:"takerCoefficient,string"`
	MakerCoefficient float64 `json:"makerCoefficient,string"`
}

// TickerInfo stores ticker information
type TickerInfo struct {
	tickerInfoBase
	SymbolName string `json:"symbolName"`
}

// Stats24hrs stores 24 hrs statistics
type Stats24hrs struct {
	tickerInfoBase
	Time uint64 `json:"time"`
}

// Orderbook stores the orderbook data
type Orderbook struct {
	Bids []orderbook.Item
	Asks []orderbook.Item
	Time time.Time
}

type orderbookResponse struct {
	Data struct {
		Asks     [][2]string        `json:"asks"`
		Bids     [][2]string        `json:"bids"`
		Time     kucoinTimeMilliSec `json:"time"`
		Sequence string             `json:"sequence"`
	} `json:"data"`
	Error
}

// Trade stores trade data
type Trade struct {
	Sequence string            `json:"sequence"`
	Price    float64           `json:"price,string"`
	Size     float64           `json:"size,string"`
	Side     string            `json:"side"`
	Time     kucoinTimeNanoSec `json:"time"`
}

// Kline stores kline data
type Kline struct {
	StartTime time.Time
	Open      float64
	Close     float64
	High      float64
	Low       float64
	Volume    float64 // Transaction volume
	Amount    float64 // Transaction amount
}

type currencyBase struct {
	Currency        string `json:"currency"` // a unique currency code that will never change
	Name            string `json:"name"`     // will change after renaming
	Fullname        string `json:"fullName"`
	Precision       int64  `json:"precision"`
	Confirms        int64  `json:"confirms"`
	ContractAddress string `json:"contractAddress"`
	IsMarginEnabled bool   `json:"isMarginEnabled"`
	IsDebitEnabled  bool   `json:"isDebitEnabled"`
}

// Currency stores currency data
type Currency struct {
	currencyBase
	WithdrawalMinSize float64 `json:"withdrawalMinSize,string"`
	WithdrawalMinFee  float64 `json:"withdrawalMinFee,string"`
	IsWithdrawEnabled bool    `json:"isWithdrawEnabled"`
	IsDepositEnabled  bool    `json:"isDepositEnabled"`
}

// Chain stores blockchain data
type Chain struct {
	Name              string  `json:"chainName"`
	Confirms          int64   `json:"confirms"`
	ContractAddress   string  `json:"contractAddress"`
	WithdrawalMinSize float64 `json:"withdrawalMinSize,string"`
	WithdrawalMinFee  float64 `json:"withdrawalMinFee,string"`
	IsWithdrawEnabled bool    `json:"isWithdrawEnabled"`
	IsDepositEnabled  bool    `json:"isDepositEnabled"`
}

// CurrencyDetail stores currency details
type CurrencyDetail struct {
	currencyBase
	Chains []Chain `json:"chains"`
}

// MarkPrice stores mark price data
type MarkPrice struct {
	Symbol      string             `json:"symbol"`
	Granularity int64              `json:"granularity"`
	TimePoint   kucoinTimeMilliSec `json:"timePoint"`
	Value       float64            `json:"value"`
}

// MarginConfiguration stores margin configuration
type MarginConfiguration struct {
	CurrencyList     []string `json:"currencyList"`
	WarningDebtRatio float64  `json:"warningDebtRatio,string"`
	LiqDebtRatio     float64  `json:"liqDebtRatio,string"`
	MaxLeverage      float64  `json:"maxLeverage"`
}

// MarginAccount stores margin account data
type MarginAccount struct {
	CurrencyList  float64 `json:"availableBalance,string"`
	Currency      string  `json:"currency"`
	HoldBalance   float64 `json:"holdBalance,string"`
	Liability     float64 `json:"liability,string"`
	MaxBorrowSize float64 `json:"maxBorrowSize,string"`
	TotalBalance  float64 `json:"totalBalance,string"`
}

// MarginAccounts stores margin accounts data
type MarginAccounts struct {
	Accounts  []MarginAccount `json:"accounts"`
	DebtRatio float64         `json:"debtRatio,string"`
}

// MarginRiskLimit stores margin risk limit
type MarginRiskLimit struct {
	Currency        string  `json:"currency"`
	BorrowMaxAmount float64 `json:"borrowMaxAmount,string"`
	BuyMaxAmount    float64 `json:"buyMaxAmount,string"`
	Precision       int64   `json:"precision"`
}

// PostBorrowOrderResp stores borrow order resposne
type PostBorrowOrderResp struct {
	OrderID  string `json:"orderId"`
	Currency string `json:"currency"`
}

// BorrowOrder stores borrow order
type BorrowOrder struct {
	OrderID   string  `json:"orderId"`
	Currency  string  `json:"currency"`
	Size      float64 `json:"size,string"`
	Filled    float64 `json:"filled"`
	MatchList []struct {
		Currency     string                `json:"currency"`
		DailyIntRate float64               `json:"dailyIntRate,string"`
		Size         float64               `json:"size,string"`
		Term         int64                 `json:"term"`
		Timestamp    kucoinTimeMilliSecStr `json:"timestamp"`
		TradeID      string                `json:"tradeId"`
	} `json:"matchList"`
	Status string `json:"status"`
}

type baseRecord struct {
	TradeID      string  `json:"tradeId"`
	Currency     string  `json:"currency"`
	DailyIntRate float64 `json:"dailyIntRate,string"`
	Principal    float64 `json:"principal,string"`
	RepaidSize   float64 `json:"repaidSize,string"`
	Term         int64   `json:"term"`
}

// OutstandingRecord stores outstanding record
type OutstandingRecord struct {
	baseRecord
	AccruedInterest float64               `json:"accruedInterest,string"`
	Liability       float64               `json:"liability,string"`
	MaturityTime    kucoinTimeMilliSecStr `json:"maturityTime"`
	CreatedAt       kucoinTimeMilliSecStr `json:"createdAt"`
}

// RepaidRecord stores repaid record
type RepaidRecord struct {
	baseRecord
	Interest  float64               `json:"interest,string"`
	RepayTime kucoinTimeMilliSecStr `json:"repayTime"`
}

// LendOrder stores lend order
type LendOrder struct {
	OrderID      string                `json:"orderId"`
	Currency     string                `json:"currency"`
	Size         float64               `json:"size,string"`
	FilledSize   float64               `json:"filledSize,string"`
	DailyIntRate float64               `json:"dailyIntRate,string"`
	Term         int64                 `json:"term"`
	CreatedAt    kucoinTimeMilliSecStr `json:"createdAt"`
}

// LendOrderHistory stores lend order history
type LendOrderHistory struct {
	LendOrder
	Status string `json:"status"`
}

// UnsettleLendOrder stores unsettle lend order
type UnsettleLendOrder struct {
	TradeID         string                `json:"tradeId"`
	Currency        string                `json:"currency"`
	Size            float64               `json:"size,string"`
	AccruedInterest float64               `json:"accruedInterest,string"`
	Repaid          float64               `json:"repaid,string"`
	DailyIntRate    float64               `json:"dailyIntRate,string"`
	Term            int64                 `json:"term"`
	MaturityTime    kucoinTimeMilliSecStr `json:"maturityTime"`
}

// SettleLendOrder stores  settled lend order
type SettleLendOrder struct {
	TradeID      string             `json:"tradeId"`
	Currency     string             `json:"currency"`
	Size         float64            `json:"size,string"`
	Interest     float64            `json:"interest,string"`
	Repaid       float64            `json:"repaid,string"`
	DailyIntRate float64            `json:"dailyIntRate,string"`
	Term         int64              `json:"term"`
	SettledAt    kucoinTimeMilliSec `json:"settledAt"`
	Note         string             `json:"note"`
}

// LendRecord stores lend record
type LendRecord struct {
	Currency        string  `json:"currency"`
	Outstanding     float64 `json:"outstanding,string"`
	FilledSize      float64 `json:"filledSize,string"`
	AccruedInterest float64 `json:"accruedInterest,string"`
	RealizedProfit  float64 `json:"realizedProfit,string"`
	IsAutoLend      bool    `json:"isAutoLend"`
}

// LendMarketData stores lend market data
type LendMarketData struct {
	DailyIntRate float64 `json:"dailyIntRate,string"`
	Term         int64   `json:"term"`
	Size         float64 `json:"size,string"`
}

// MarginTradeData stores margin trade data
type MarginTradeData struct {
	TradeID      string            `json:"tradeId"`
	Currency     string            `json:"currency"`
	Size         float64           `json:"size,string"`
	DailyIntRate float64           `json:"dailyIntRate,string"`
	Term         int64             `json:"term"`
	Timestamp    kucoinTimeNanoSec `json:"timestamp"`
}

type IsolatedMarginPairConfig struct {
	Symbol                string  `json:"symbol"`
	SymbolName            string  `json:"symbolName"`
	BaseCurrency          string  `json:"baseCurrency"`
	QuoteCurrency         string  `json:"quoteCurrency"`
	MaxLeverage           int64   `json:"maxLeverage"`
	LiquidationDebtRatio  float64 `json:"flDebtRatio,string"`
	TradeEnable           bool    `json:"tradeEnable"`
	AutoRenewMaxDebtRatio float64 `json:"autoRenewMaxDebtRatio,string"`
	BaseBorrowEnable      bool    `json:"baseBorrowEnable"`
	QuoteBorrowEnable     bool    `json:"quoteBorrowEnable"`
	BaseTransferInEnable  bool    `json:"baseTransferInEnable"`
	QuoteTransferInEnable bool    `json:"quoteTransferInEnable"`
}

type baseAsset struct {
	Currency         string  `json:"currency"`
	TotalBalance     float64 `json:"totalBalance,string"`
	HoldBalance      float64 `json:"holdBalance,string"`
	AvailableBalance float64 `json:"availableBalance,string"`
	Liability        float64 `json:"liability,string"`
	Interest         float64 `json:"interest,string"`
	BorrowableAmount float64 `json:"borrowableAmount,string"`
}

type AssetInfo struct {
	Symbol     string    `json:"symbol"`
	Status     string    `json:"status"`
	DebtRatio  float64   `json:"debtRatio,string"`
	BaseAsset  baseAsset `json:"baseAsset"`
	QuoteAsset baseAsset `json:"quoteAsset"`
}

type IsolatedMarginAccountInfo struct {
	TotalConversionBalance     float64     `json:"totalConversionBalance,string"`
	LiabilityConversionBalance float64     `json:"liabilityConversionBalance,string"`
	Assets                     []AssetInfo `json:"assets"`
}

type baseRepaymentRecord struct {
	LoanID            string  `json:"loanId"`
	Symbol            string  `json:"symbol"`
	Currency          string  `json:"currency"`
	PrincipalTotal    float64 `json:"principalTotal,string"`
	InterestBalance   float64 `json:"interestBalance,string"`
	CreatedAt         int64   `json:"createdAt"`
	Period            int64   `json:"period"`
	RepaidSize        float64 `json:"repaidSize,string"`
	DailyInterestRate float64 `json:"dailyInterestRate,string"`
}

type OutstandingRepaymentRecord struct {
	baseRepaymentRecord
	LiabilityBalance float64 `json:"liabilityBalance,string"`
	MaturityTime     int64   `json:"maturityTime"`
}

type CompletedRepaymentRecord struct {
	baseRepaymentRecord
	RepayFinishAt int64 `json:"repayFinishAt"`
}

type PostMarginOrderResp struct {
	OrderID     string  `json:"orderId"`
	BorrowSize  float64 `json:"borrowSize"`
	LoanApplyID string  `json:"loanApplyId"`
}

type OrderRequest struct {
	ClientOID   string  `json:"clientOid"`
	Symbol      string  `json:"symbol"`
	Side        string  `json:"side"`
	Type        string  `json:"type,omitempty"`      // optional
	Remark      string  `json:"remark,omitempty"`    // optional
	Stop        string  `json:"stop,omitempty"`      // optional
	StopPrice   string  `json:"stopPrice,omitempty"` // optional
	STP         string  `json:"stp,omitempty"`       // optional
	Price       float64 `json:"price,string,omitempty"`
	Size        float64 `json:"size,string,omitempty"`
	TimeInForce string  `json:"timeInForce,omitempty"` // optional
	CancelAfter int64   `json:"cancelAfter,omitempty"` // optional
	PostOnly    bool    `json:"postOnly,omitempty"`    // optional
	Hidden      bool    `json:"hidden,omitempty"`      // optional
	Iceberg     bool    `json:"iceberg,omitempty"`     // optional
	VisibleSize string  `json:"visibleSize,omitempty"` // optional
}

type PostBulkOrderResp struct {
	OrderRequest
	Channel string `json:"channel"`
	ID      string `json:"id"`
	Status  string `json:"status"`
	FailMsg string `json:"failMsg"`
}

type OrderDetail struct {
	OrderRequest
	Channel       string             `json:"channel"`
	ID            string             `json:"id"`
	OpType        string             `json:"opType"` // operation type: DEAL
	Funds         string             `json:"funds"`
	DealFunds     string             `json:"dealFunds"`
	DealSize      float64            `json:"dealSize,string"`
	Fee           float64            `json:"fee,string"`
	FeeCurrency   string             `json:"feeCurrency"`
	StopTriggered bool               `json:"stopTriggered"`
	Tags          string             `json:"tags"`
	IsActive      bool               `json:"isActive"`
	CancelExist   bool               `json:"cancelExist"`
	CreatedAt     kucoinTimeMilliSec `json:"createdAt"`
	TradeType     string             `json:"tradeType"`
}

type Fill struct {
	Symbol         string             `json:"symbol"`
	TradeID        string             `json:"tradeId"`
	OrderID        string             `json:"orderId"`
	CounterOrderId string             `json:"counterOrderId"`
	Side           string             `json:"side"`
	Liquidity      string             `json:"liquidity"`
	ForceTaker     bool               `json:"forceTaker"`
	Price          float64            `json:"price,string"`
	Size           float64            `json:"size,string"`
	Funds          float64            `json:"funds,string"`
	Fee            float64            `json:"fee,string"`
	FeeRate        float64            `json:"feeRate,string"`
	FeeCurrency    string             `json:"feeCurrency"`
	Stop           string             `json:"stop"`
	OrderType      string             `json:"type"`
	CreatedAt      kucoinTimeMilliSec `json:"createdAt"`
	TradeType      string             `json:"tradeType"`
}

type StopOrder struct {
	OrderRequest
	ID              string             `json:"id"`
	UserID          string             `json:"userId"`
	Status          string             `json:"status"`
	Funds           float64            `json:"funds,string"`
	Channel         string             `json:"channel"`
	Tags            string             `json:"tags"`
	DomainId        string             `json:"domainId"`
	TradeSource     string             `json:"tradeSource"`
	TradeType       string             `json:"tradeType"`
	FeeCurrency     string             `json:"feeCurrency"`
	TakerFeeRate    string             `json:"takerFeeRate"`
	MakerFeeRate    string             `json:"makerFeeRate"`
	CreatedAt       kucoinTimeMilliSec `json:"createdAt"`
	OrderTime       kucoinTimeNanoSec  `json:"orderTime"`
	StopTriggerTime kucoinTimeMilliSec `json:"stopTriggerTime"`
}

type baseAccount struct {
	Currency  string  `json:"currency"`
	Balance   float64 `json:"balance,string"`
	Available float64 `json:"available,string"`
	Holds     float64 `json:"holds,string"`
}

type AccountInfo struct {
	baseAccount
	ID   string `json:"id"`
	Type string `json:"type"`
}

type LedgerInfo struct {
	ID          string             `json:"id"`
	Currency    string             `json:"currency"`
	Amount      float64            `json:"amount,string"`
	Fee         float64            `json:"fee,string"`
	Balance     float64            `json:"balance,string"`
	AccountType string             `json:"accountType"`
	BizType     string             `json:"bizType"`
	Direction   string             `json:"direction"`
	CreatedAt   kucoinTimeMilliSec `json:"createdAt"`
	Context     string             `json:"context"`
}

type MainAccountInfo struct {
	baseAccount
	BaseCurrency      string  `json:"baseCurrency"`
	BaseCurrencyPrice float64 `json:"baseCurrencyPrice,string"`
	BaseAmount        float64 `json:"baseAmount,string"`
}

type SubAccountInfo struct {
	SubUserID      string            `json:"subUserId"`
	SubName        string            `json:"subName"`
	MainAccounts   []MainAccountInfo `json:"mainAccounts"`
	TradeAccounts  []MainAccountInfo `json:"tradeAccounts"`
	MarginAccounts []MainAccountInfo `json:"marginAccounts"`
}

type TransferableBalanceInfo struct {
	baseAccount
	Transferable float64 `json:"transferable,string"`
}

type DepositAddress struct {
	Address         string `json:"address"`
	Memo            string `json:"memo"`
	Chain           string `json:"chain"`
	ContractAddress string `json:"contractAddress"` // missing in case of futures
}

type baseDeposit struct {
	Currency   string  `json:"currency"`
	Amount     float64 `json:"amount"`
	WalletTxID string  `json:"walletTxId"`
	IsInner    bool    `json:"isInner"`
	Status     string  `json:"status"`
}

type Deposit struct {
	baseDeposit
	Address   string  `json:"address"`
	Memo      string  `json:"memo"`
	Fee       float64 `json:"fee"`
	Remark    string  `json:"remark"`
	CreatedAt kucoinTimeMilliSec
	UpdatedAt kucoinTimeMilliSec
}

type HistoricalDepositWithdrawal struct {
	baseDeposit
	CreatedAt kucoinTimeMilliSec `json:"createAt"`
}

type Withdrawal struct {
	Deposit
	ID string `json:"id"`
}

type WithdrawalQuota struct {
	Currency            string  `json:"currency"`
	LimitBTCAmount      float64 `json:"limitBTCAmount,string"`
	UsedBTCAmount       float64 `json:"usedBTCAmount,string"`
	RemainAmount        float64 `json:"remainAmount,string"`
	AvailableAmount     float64 `json:"availableAmount,string"`
	WithdrawMinFee      float64 `json:"withdrawMinFee,string"`
	InnerWithdrawMinFee float64 `json:"innerWithdrawMinFee,string"`
	WithdrawMinSize     float64 `json:"withdrawMinSize,string"`
	IsWithdrawEnabled   bool    `json:"isWithdrawEnabled"`
	Precision           int64   `json:"precision"`
	Chain               string  `json:"chain"`
}

type Fees struct {
	Symbol       string  `json:"symbol"`
	TakerFeeRate float64 `json:"takerFeeRate,string"`
	MakerFeeRate float64 `json:"makerFeeRate,string"`
}

// WSInstanceServers response
type WSInstanceServers struct {
	Token           string           `json:"token"`
	InstanceServers []InstanceServer `json:"instanceServers"`
}

// InstanceServer represents a a single websocket instance server information.
type InstanceServer struct {
	Endpoint     string `json:"endpoint"`
	Encrypt      bool   `json:"encrypt"`
	Protocol     string `json:"protocol"`
	PingInterval int64  `json:"pingInterval"`
	PingTimeout  int64  `json:"pingTimeout"`
}

// WSConnMessages represents response messages ping, pong, and welcome message structures.
type WSConnMessages struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

// WsSubscriptionInput represents a subscription information structure.
type WsSubscriptionInput struct {
	ID             int64  `json:"id"`
	Type           string `json:"type"`
	Topic          string `json:"topic"`
	PrivateChannel bool   `json:"privateChannel"`
	Response       bool   `json:"response,omitempty"`
}

// WSSubscriptionResponse represents a subscription response.
type WSSubscriptionResponse struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

// WsPushData represents a push data from a server.
type WsPushData struct {
	ID          string          `json:"id"`
	Type        string          `json:"type"`
	Topic       string          `json:"topic"`
	Subject     string          `json:"subject"`
	ChannelType string          `json:"channelType,omitempty"`
	Data        json.RawMessage `json:"data"`
}

// WsTicker represents a ticker push data from server.
type WsTicker struct {
	Sequence    string  `json:"sequence"`
	BestAsk     float64 `json:"bestAsk,string"`
	Size        float64 `json:"size,string"`
	BestBidSize float64 `json:"bestBidSize,string"`
	Price       float64 `json:"price,string"`
	BestAskSize float64 `json:"bestAskSize,string"`
	BestBid     float64 `json:"bestBid,string"`
}

// WsTickerDetail represents a ticker snapshot data from server.
type WsTickerDetail struct {
	Sequence string `json:"sequence"`
	Data     []struct {
		Trading         bool    `json:"trading"`
		Symbol          string  `json:"symbol"`
		Buy             float64 `json:"buy"`
		Sell            float64 `json:"sell"`
		Sort            int     `json:"sort"`
		VolValue        float64 `json:"volValue"`
		BaseCurrency    string  `json:"baseCurrency"`
		Market          string  `json:"market"`
		QuoteCurrency   string  `json:"quoteCurrency"`
		SymbolCode      string  `json:"symbolCode"`
		Datetime        int64   `json:"datetime"`
		High            float64 `json:"high"`
		Vol             float64 `json:"vol"`
		Low             float64 `json:"low"`
		ChangePrice     float64 `json:"changePrice"`
		ChangeRate      float64 `json:"changeRate"`
		LastTradedPrice float64 `json:"lastTradedPrice"`
		Board           int     `json:"board"`
		Mark            int     `json:"mark"`
	} `json:"data"`
}

// WsOrderbook represents orderbook information.
type WsOrderbook struct {
	Changes struct {
		Asks [][3]string `json:"asks"`
		Bids [][3]string `json:"bids"`
	} `json:"changes"`
	SequenceEnd   int64  `json:"sequenceEnd"`
	SequenceStart int64  `json:"sequenceStart"`
	Symbol        string `json:"symbol"`
	TimeMS        int64  `json:"time"`
}

// WsCandlestickData represents candlestick information push data for a symbol.
type WsCandlestickData struct {
	Symbol  string    `json:"symbol"`
	Candles [7]string `json:"candles"`
	Time    int64     `json:"time"`
}

// WsCandlestick represents candlestick information push data for a symbol.
type WsCandlestick struct {
	Symbol  string `json:"symbol"`
	Candles struct {
		StartTime         time.Time
		OpenPrice         float64
		ClosePrice        float64
		HighPrice         float64
		LowPrice          float64
		TransactionVolume float64
		TransactionAmount float64
	} `json:"candles"`
	Time time.Time `json:"time"`
}

func (a *WsCandlestickData) getCandlestickData() (*WsCandlestick, error) {
	cand := &WsCandlestick{
		Symbol: a.Symbol,
		Time:   time.UnixMilli(a.Time),
	}
	timeStamp, err := strconv.ParseInt(a.Candles[0], 10, 64)
	if err != nil {
		return nil, err
	}
	cand.Candles.StartTime = time.UnixMilli(timeStamp)
	cand.Candles.OpenPrice, err = strconv.ParseFloat(a.Candles[1], 64)
	if err != nil {
		return nil, err
	}
	cand.Candles.ClosePrice, err = strconv.ParseFloat(a.Candles[2], 64)
	if err != nil {
		return nil, err
	}
	cand.Candles.HighPrice, err = strconv.ParseFloat(a.Candles[3], 64)
	if err != nil {
		return nil, err
	}
	cand.Candles.LowPrice, err = strconv.ParseFloat(a.Candles[4], 64)
	if err != nil {
		return nil, err
	}
	cand.Candles.TransactionVolume, err = strconv.ParseFloat(a.Candles[5], 64)
	if err != nil {
		return nil, err
	}
	cand.Candles.TransactionAmount, err = strconv.ParseFloat(a.Candles[6], 64)
	if err != nil {
		return nil, err
	}
	return cand, nil
}

// WsTrade represents a trade push data.
type WsTrade struct {
	Sequence     string  `json:"sequence"`
	Type         string  `json:"type"`
	Symbol       string  `json:"symbol"`
	Side         string  `json:"side"`
	Price        float64 `json:"price,string"`
	Size         float64 `json:"size,string"`
	TradeID      string  `json:"tradeId"`
	TakerOrderID string  `json:"takerOrderId"`
	MakerOrderID string  `json:"makerOrderId"`
	Time         int64   `json:"time,string"`
}

// WsPriceIndicator represents index price or mark price indicator push data.
type WsPriceIndicator struct {
	Symbol      string  `json:"symbol"`
	Granularity float64 `json:"granularity"`
	Timestamp   int64   `json:"timestamp"`
	Value       float64 `json:"value"`
}

// WsMarginFundingBook represents order book changes on margin.
type WsMarginFundingBook struct {
	Sequence           int64     `json:"sequence"`
	Currency           string    `json:"currency"`
	DailyInterestRate  float64   `json:"dailyIntRate,string"`
	AnnualInterestRate float64   `json:"annualIntRate,string"`
	Term               int64     `json:"term"`
	Size               float64   `json:"size,string"`
	Side               string    `json:"side"`
	Timestamp          time.Time `json:"ts"` // In Nanosecond

}

// WsTradeOrder represents a private trade order push data.
type WsTradeOrder struct {
	Symbol     string    `json:"symbol"`
	OrderType  string    `json:"orderType"`
	Side       string    `json:"side"`
	OrderID    string    `json:"orderId"`
	Type       string    `json:"type"`
	OrderTime  time.Time `json:"orderTime"`
	Size       float64   `json:"size,string"`
	FilledSize float64   `json:"filledSize,string"`
	Price      float64   `json:"price,string"`
	ClientOid  string    `json:"clientOid"`
	RemainSize float64   `json:"remainSize,string"`
	Status     string    `json:"status"`
	Timestamp  time.Time `json:"ts"`
	Liquidity  string    `json:"liquidity,omitempty"`
	MatchPrice string    `json:"matchPrice,omitempty"`
	MatchSize  string    `json:"matchSize,omitempty"`
	TradeID    string    `json:"tradeId,omitempty"`
	OldSize    string    `json:"oldSize,omitempty"`
}

// WsAccountBalance represents a Account Balance push data.
type WsAccountBalance struct {
	Total           float64 `json:"total,string"`
	Available       float64 `json:"available,string"`
	AvailableChange float64 `json:"availableChange,string"`
	Currency        string  `json:"currency"`
	Hold            float64 `json:"hold,string"`
	HoldChange      float64 `json:"holdChange,string"`
	RelationEvent   string  `json:"relationEvent"`
	RelationEventID string  `json:"relationEventId"`
	RelationContext struct {
		Symbol  string `json:"symbol"`
		TradeID string `json:"tradeId"`
		OrderID string `json:"orderId"`
	} `json:"relationContext"`
	Time time.Time `json:"time,string"`
}

// WsDebtRatioChange represents a push data
type WsDebtRatioChange struct {
	DebtRatio float64           `json:"debtRatio"`
	TotalDebt string            `json:"totalDebt"`
	DebtList  map[string]string `json:"debtList"`
	Timestamp time.Time         `json:"timestamp"`
}

// WsPositionStatus represents a position status push data.
type WsPositionStatus struct {
	Type        string    `json:"type"`
	TimestampMS time.Time `json:"timestamp"`
}

// WsMarginTradeOrderEntersEvent represents a push data to the lenders
// when the order enters the order book or when the order is executed.
type WsMarginTradeOrderEntersEvent struct {
	Currency     string    `json:"currency"`
	OrderID      string    `json:"orderId"`      //Trade ID
	DailyIntRate float64   `json:"dailyIntRate"` //Daily interest rate.
	Term         int64     `json:"term"`         //Term (Unit: Day)
	Size         float64   `json:"size"`         //Size
	LentSize     float64   `json:"lentSize"`     //Size executed -- filled when the subject is order.update
	Side         string    `json:"side"`         //Lend or borrow. Currently, only "Lend" is available
	Timestamp    time.Time `json:"ts"`           //Timestamp (nanosecond)
}

// WsMarginTradeOrderDoneEvent represents a push message to the lenders when the order is completed.
type WsMarginTradeOrderDoneEvent struct {
	Currency  string    `json:"currency"`
	OrderID   string    `json:"orderId"`
	Reason    string    `json:"reason"`
	Side      string    `json:"side"`
	Timestamp time.Time `json:"ts"`
}

// WsStopOrder represents a stop order.
// When a stop order is received by the system, you will receive a message with "open" type.
// It means that this order entered the system and waited to be triggered.
type WsStopOrder struct {
	CreatedAt      time.Time `json:"createdAt"`
	OrderID        string    `json:"orderId"`
	OrderPrice     float64   `json:"orderPrice,string"`
	OrderType      string    `json:"orderType"`
	Side           string    `json:"side"`
	Size           float64   `json:"size,string"`
	Stop           string    `json:"stop"`
	StopPrice      float64   `json:"stopPrice,string"`
	Symbol         string    `json:"symbol"`
	TradeType      string    `json:"tradeType"`
	TriggerSuccess bool      `json:"triggerSuccess"`
	Timestamp      time.Time `json:"ts"`
	Type           string    `json:"type"`
}