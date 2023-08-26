package keeper

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"vcoa/x/vnft/types"

	sdkctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/nft"

	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
)

func generatePubKeys(n int) []cryptotypes.PubKey {
	pks := make([]cryptotypes.PubKey, n)
	for i := 0; i < n; i++ {
		pks[i] = secp256k1.GenPrivKey().PubKey()
	}
	return pks
}

func (k msgServer) MintNft(goCtx context.Context, msg *types.MsgMintNft) (*types.MsgMintNftResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: Handling the message
	has := k.nftKeeper.HasClass(ctx, msg.ClassId)
	if !has {
		return nil, types.ErrFailedToMarshalData.Wrapf("%s, %s", "No such class ID", msg.ClassId)
	}

	// check if nft class transferable
	class, hasClass := k.nftKeeper.GetClass(ctx, msg.ClassId)
	if !hasClass {
		return nil, types.ErrNftClassNotFound.Wrapf("Class id %s not found", msg.ClassId)
	}
	// Unmarshal class data
	var classData types.ClassData
	if err := classData.Unmarshal(class.Data.Value); err != nil {
		return nil, types.ErrFailedToUnmarshalData.Wrapf(err.Error())
	}

	type format struct {
		CField    []string `protobuf:"bytes,1,opt,name=cField,proto3" json:"cField,omitempty"`
		CType     []string `protobuf:"bytes,2,opt,name=cType,proto3" json:"cType,omitempty"`
		COptional []bool   `protobuf:"bytes,3,opt,name=cOptional,proto3" json:"cOptional,omitempty"`
	}

	type metadata struct {
		Version      int             `protobuf:"bytes,1,opt,name=version,proto3" json:"version,omitempty"`
		Format       types.JsonInput `protobuf:"bytes,2,opt,name=format,proto3" json:"format,omitempty"`
		MultiAuthor  bool            `protobuf:"bytes,3,opt,name=multiAuthor,proto3" json:"multiAuthor,omitempty"`
		Burnable     bool            `protobuf:"bytes,4,opt,name=burnable,proto3" json:"burnable,omitempty"`
		Transferable bool            `protobuf:"bytes,5,opt,name=transferable,proto3" json:"transferable,omitempty"`
	}

	var classMeta metadata
	json.Unmarshal([]byte(classData.Metadata), &classMeta)

	var f format
	json.Unmarshal([]byte(classMeta.Format), &f)

	var token nft.NFT

	userAddress, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return nil, types.ErrFailedToGetUserAddress.Wrapf("User Address %s", err.Error())
	}
	address := userAddress

	fmt.Println("msg.Creator: ", msg.Creator)
	fmt.Println("userAddress: ", userAddress)

	token.ClassId = msg.ClassId

	token.Uri = msg.Uri
	token.UriHash = msg.UriHash

	databyte := []byte(msg.Data.Metadata)
	var meta types.JsonInput
	if json.Unmarshal(databyte, &meta) != nil {
		return nil, types.ErrFailedToMarshalData.Wrapf("Unmarshal %s", "error")
	}

	var nftMeta types.JsonInput

	fmt.Println("meta: ", meta)
	fmt.Println("class.Name: ", class.Name)
	fmt.Println("classMeta.Version: ", classMeta.Version)

	// totalSupply := k.nftKeeper.GetTotalSupply(ctx, class.Id)

	token.Id, err = types.NewNftId(class.Name, classMeta.Version, token.UriHash, msg.Creator)
	fmt.Println("token.Id: ", token.Id)
	if err != nil {
		return nil, types.ErrFailedToCreateNftId.Wrapf("%s", err.Error())
	}

	pubKey := "" // public key of the creator

	if class.Name == "ID" {
		token.Uri = ""
		token.UriHash = ""

		type nftID struct {
			PubKey string
			Record []string
		}

		var rd nftID
		json.Unmarshal([]byte(meta), &rd)

		fmt.Printf("rd: %s %d\n", rd, len(rd.Record))

		if len(rd.Record) < 1 {
			return nil, types.ErrFailedToMarshalData.Wrapf("Not enough record to process")
		}

		pubKey = rd.PubKey

		rd.Record[0] = strings.ReplaceAll(rd.Record[0], " ", "")
		rd.Record[0] = regexp.MustCompile(`[^a-zA-Z0-9 ]+`).ReplaceAllString(rd.Record[0], "")
		rd.Record[0] = strings.ToLower(rd.Record[0])

		if rd.Record[0] == "" {
			return nil, types.ErrFailedToMarshalData.Wrapf("Unmarshal %s", "error")
		}

		nj := "{\""
		if len(rd.Record) >= 2 {
			rd.Record[1] = regexp.MustCompile(`[^a-zA-Z0-9 ]+`).ReplaceAllString(rd.Record[1], "") // remove all unknown char
			nj = nj + "address" + "\": \"" + address.String() + "\", \"" + f.CField[0] + "\": \"" + rd.Record[0] + "\", \"" + f.CField[1] + "\": \"" + rd.Record[1]
		} else {
			nj = nj + "address" + "\": \"" + address.String() + "\", \"" + f.CField[0] + "\": \"" + rd.Record[0] + "\", \"" + f.CField[1] + "\": \""
		}
		nj = nj + "\"}"

		nftMeta = types.JsonInput(nj)
		fmt.Println("nftMeta: ", nftMeta)

		token.Id = "vconft-" + class.Name + "-" + rd.Record[0]

	} else {

		type nftMetaType struct {
			PubKey string
			Record []string
		}
		var rd nftMetaType
		json.Unmarshal([]byte(meta), &rd)
		fmt.Printf("rd: %s %d\n", rd, len(rd.Record))
		pubKey = rd.PubKey

		if rd.Record[0] == "" {
			return nil, types.ErrFailedToMarshalData.Wrapf("Unmarshal %s", "error")
		}

		nj := "{\""
		nj = nj + "address" + "\": \"" + address.String() + "\""
		for index, element := range f.CField {
			if element != "tag" {
				rd.Record[index] = regexp.MustCompile(`[^a-zA-Z0-9 ]+`).ReplaceAllString(rd.Record[index], "") // remove all unknown char
			} else {
				rd.Record[index] = regexp.MustCompile(`[^\w\s-]+`).ReplaceAllString(rd.Record[index], "") // remove all unknown char, numbers
				rd.Record[index] = strings.ToLower(rd.Record[index])
			}
			nj = nj + ",\"" + element + "\":\"" + rd.Record[index] + "\""
		}
		nj = nj + "}"

		nftMeta = types.JsonInput(nj)
		fmt.Println("nftMeta: ", nftMeta)
	}

	var pubKeyBytes []byte
	if pubKey != "" {
		pubKeyBytes, err = hex.DecodeString(pubKey)
		if err != nil {
			return nil, types.ErrFailedToMarshalData.Wrapf("%s", err.Error())
		}
		fmt.Println("pubKey: ", pubKey)
		fmt.Println("pubKeyBytes: ", pubKeyBytes)
		fmt.Println("&secp256k1.PubKey{Key: pubKeyBytes}: ", &secp256k1.PubKey{Key: pubKeyBytes})

	} else {

		// for no public key was given
		fmt.Println("msg.UriHash: ", msg.UriHash)
		pk, _ := base64.StdEncoding.DecodeString(msg.UriHash)
		creatorPubKey := &secp256k1.PubKey{Key: pk}

		var aa sdk.AccAddress
		aa.Unmarshal(creatorPubKey.Address())

		fmt.Println("pk: ", pk)
		fmt.Println("creatorPubKey: ", creatorPubKey)
		fmt.Println("aa: ", aa.String())
		pubKeyBytes = creatorPubKey.Key
		fmt.Println("pubKeyBytes: ", pubKeyBytes)
	}

	if !classMeta.Transferable {
		// scramble pubKey to make no one knows the private key
		var hashBytes []byte
		h := sha256.New()
		h.Write(pubKeyBytes[1:])
		bs := h.Sum(nil)
		hashBytes = append(hashBytes, pubKeyBytes[0])
		hashBytes = append(hashBytes, bs...)
		fmt.Println("hashBytes: ", hashBytes)
		hashedKey := &secp256k1.PubKey{Key: hashBytes}
		fmt.Println("&secp256k1.PubKey{Key: hashBytes}: ", hashedKey)

		var mAddr sdk.AccAddress
		mAddr.Unmarshal(hashedKey.Address())
		userAddress = mAddr

		// to make sure mAddr is never used for this class
		balance := k.nftKeeper.GetBalance(ctx, msg.ClassId, mAddr)
		if balance > 0 {
			return nil, types.ErrOwnerAddressUsedBefore.Wrapf("Address registered before")
		}

	}

	fmt.Println()
	fmt.Println()
	fmt.Println("-----------------------------------------------------------------------------")

	var nftmeta types.JsonInput
	if json.Unmarshal([]byte(nftMeta), &nftmeta) != nil {
		return nil, types.ErrFailedToMarshalData.Wrapf("Unmarshal %s", "error")
	}

	nftData := types.NftData{
		Metadata: nftmeta,
	}

	classDataInAny, err := sdkctypes.NewAnyWithValue(&nftData)
	if err != nil {
		return nil, types.ErrFailedToMarshalData.Wrapf("%s", err.Error())
	}

	token.Data = classDataInAny

	if err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("%s", err.Error())
	}

	err = k.nftKeeper.Mint(ctx, token, userAddress)
	if err != nil {
		return nil, types.ErrFailedToMintNFT.Wrapf("%s", err.Error())
	}

	return &types.MsgMintNftResponse{}, nil
}
