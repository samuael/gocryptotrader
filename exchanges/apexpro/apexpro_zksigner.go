package apexpro

// ZKLinkSigner represents a ZK link signing information
type ZKLinkSigner struct {}

type ZKKeyInfo struct {
	Seeds         []byte
	L2Key         string
	PublicKeyHash []byte
}

func DriveZKKey(ethAddress string) (*ZKKeyInfo, error) {
	msgHeader := "ApeX Omni Mainnet"
	// signature =
	println(msgHeader)
	return nil, nil
}
