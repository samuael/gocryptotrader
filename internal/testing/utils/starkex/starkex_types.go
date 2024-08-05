package starkex

import "math/big"

// WithdrawalToAddressParams represents a starkex withdrawal to address parameters. Type value: 7.
type WithdrawalToAddressParams struct {
	NetworkID           int      `json:"-"`
	AssetIDCollateral   *big.Int `json:"asset_id_collateral"`
	EthAddress          *big.Int `json:"eth_address"`
	PositionID          uint64   `json:"position_id"`
	Amount              float64  `json:"amount"`
	Nonce               uint32   `json:"nonce"`
	ExpirationTimestamp int64    `json:"expiration_timestamp"`
}

// CreateOrderWithFeeParams represents a starkex create order parameters. Type value: 3
type CreateOrderWithFeeParams struct {
	IsBuyingSynthetic uint64   `json:"is_buying_synthetic"`
	AssetIDSynthetic  *big.Int `json:"asset_id_synthetic"`
	AssetIDCollateral *big.Int `json:"asset_id_collateral"`
	AssetIDFee        *big.Int `json:"asset_id_fee"`
	AmountSynthetic   uint64   `json:"amount_synthetic"`
	AmountCollateral  uint64   `json:"amount_collateral"`
	MaxAmountFee      uint64   `json:"max_amount_fee"`
	PositionID        uint64   `json:"position_id"`
	Nonce             uint32   `json:"nonce"`
	ExpTimestampHrs   uint32   `json:"expiration_timestamp"`
}

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
