package types

import (
	"crypto/sha256"
	"math/rand"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const NFT_ID_PREFIX = "vconft"

// const NONCE = "000000000000000000000000000000000000000000000000"

// const nonce []byte

func NewNftId(class string, version int) (string, error) {
	prefix := NFT_ID_PREFIX + "-" + class + "-" + strconv.Itoa(version) + "-"
	data := []byte(NONCE)
	// data := append(nonce...)
	randbytes := make([]byte, 32)
	rand.Read(randbytes)
	data = append(data, randbytes...)
	hash := sha256.Sum256(data)

	nftId, err := sdk.Bech32ifyAddressBytes(prefix, hash[:])
	if err != nil {
		return "", err
	}
	return nftId, nil
}
