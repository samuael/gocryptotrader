package starkex

import "math/big"

// TransferParams represents a starkex asset transfer parameters. Type value: 4
type TransferParams struct {
	AssetID            *big.Int `json:"asset_id"`
	AssetIDFee         *big.Int `json:"asset_id_fee"`
	ReceiverPublicKey  *big.Int `json:"receiver_public_key"`
	SenderPositionID   uint64   `json:"sender_position_id"`
	ReceiverPositionID uint64   `json:"receiver_position_id"`
	SrcFeePositionID   uint64   `json:"src_fee_position_id"`
	Amount             uint64   `json:"amount"`
	MaxAmountFee       uint64   `json:"max_amount_fee"`
	Nonce              uint32   `json:"nonce"`
	ExpTimestampHrs    uint32   `json:"expiration_timestamp"`
}

// ConditionalTransferParams represents a conditional transfer parameters. Type value: 5
type ConditionalTransferParams struct {
	AssetID            *big.Int `json:"asset_id"`
	AssetIDFee         *big.Int `json:"asset_id_fee"`
	ReceiverPublicKey  *big.Int `json:"receiver_public_key"`
	Condition          *big.Int `json:"condition"`
	SenderPositionID   uint64   `json:"sender_position_id"`
	ReceiverPositionID uint64   `json:"receiver_position_id"`
	SrcFeePositionID   uint64   `json:"src_fee_position_id"`
	Amount             uint64   `json:"amount"`
	MaxAmountFee       uint64   `json:"max_amount_fee"`
	Nonce              uint32   `json:"nonce"`
	ExpTimestampHrs    uint32   `json:"expiration_timestamp"`
}

// CreateOrderWithFeeParams represents a starkex create order parameters. Order Prefix: 3
type CreateOrderWithFeeParams struct {
	OrderType               string
	AssetIdSynthetic        *big.Int
	AssetIdCollateral       *big.Int
	AssetIdFee              *big.Int
	QuantumAmountSynthetic  *big.Int
	QuantumAmountCollateral *big.Int
	QuantumAmountFee        *big.Int
	IsBuyingSynthetic       bool
	PositionId              *big.Int // Users Account ID
	Nonce                   *big.Int
	ExpirationEpochHours    *big.Int
}

// WithdrawalToAddressParams represents a starkex withdrawal to address parameters. Type value: 7.
type WithdrawalToAddressParams struct {
	NetworkID            int64    `json:"-"`
	AssetIDCollateral    *big.Int `json:"asset_id_collateral"`
	EthAddress           *big.Int `json:"eth_address"`
	PositionID           *big.Int `json:"position_id"`
	Amount               *big.Int `json:"amount"`
	Nonce                *big.Int `json:"nonce"`
	ExpirationEpochHours *big.Int `json:"expiration_timestamp"`
}

// WithdrawalParams represents a starkex withdrawal parameters. Type value: 6.
type WithdrawalParams struct {
	NetworkID            int64    `json:"-"`
	AssetIDCollateral    *big.Int `json:"asset_id_collateral"`
	PositionID           *big.Int `json:"position_id"`
	Amount               *big.Int `json:"amount"`
	Nonce                *big.Int `json:"nonce"`
	ExpirationEpochHours *big.Int `json:"expiration_timestamp"`
}
