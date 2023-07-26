package types

import (
	"crypto/sha256"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const CLASS_ID_PREFIX = "vconftclass"
const NONCE = "000000000000000000000000000000000000000000000000"

// const nonce []byte

func NewClassId(prefix []byte, serial int) (string, error) {
	nonce := []byte(NONCE)
	data := append(prefix, nonce...)
	data = append(data, []byte(strconv.Itoa(serial))...)
	hash := sha256.Sum256(data)

	classId, err := sdk.Bech32ifyAddressBytes(CLASS_ID_PREFIX, hash[:])
	if err != nil {
		return "", err
	}
	return classId, nil
}
