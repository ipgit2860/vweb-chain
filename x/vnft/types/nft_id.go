package types

import (
	"crypto/sha256"
	// "math/rand"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const NFT_ID_PREFIX = "vconft"

// const NONCE = "000000000000000000000000000000000000000000000000"

// const nonce []byte

func NewNftId(class string, version int, totalSupply uint64) (string, error) {
	prefix := NFT_ID_PREFIX + "-" + class + "-" + strconv.Itoa(version) + "-"

	data := []byte(class)
	data = append(data, []byte(NONCE)...)
	data = append(data, []byte(strconv.Itoa(version))...)
	data = append(data, []byte(strconv.FormatUint(totalSupply+1, 10))...)
	hash := sha256.Sum256(data)

	nftId, err := sdk.Bech32ifyAddressBytes(prefix, hash[:])
	if err != nil {
		return "", err
	}
	return nftId, nil
}
