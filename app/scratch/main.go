package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"math/big"

	"github.com/ardanlabs/blockchain/foundation/blockchain/database"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

// basically money value from an account to another account
type Tx struct {
	FromID string `json:"from"`
	ToID   string `json:"to"`
	Value  uint64 `json:"value"`
}

func main() {
	if err := run(); err != nil {
		log.Fatalln(err)
	}
}

func run() error {
	privateKey, err := crypto.LoadECDSA("zblock/accounts/kennedy.ecdsa")
	if err != nil {
		return fmt.Errorf("unable to load private key for node: %w", err)
	}

	tx := Tx{
		FromID: "Bill",
		ToID:   "Aaron",
		Value:  1000,
	}

	data, err := json.Marshal(tx)
	if err != nil {
		return fmt.Errorf("unable to marshal %w", err)
	}
	// salt with our stamp
	stamp := []byte(fmt.Sprintf("\x19Ardan Signed Message:\n%d", len(data)))

	v := crypto.Keccak256(stamp, data)

	sig, err := crypto.Sign(v, privateKey)
	if err != nil {
		return fmt.Errorf("unable to sign transaction: %w", err)
	}

	fmt.Println("SIG: ", hexutil.Encode(sig))

	// =================================================================================
	// OVER THE WIRE

	publicKey, err := crypto.SigToPub(v, sig)
	if err != nil {
		return fmt.Errorf("unable to recover public key from signature: %w", err)
	}

	fmt.Println("PUB: ", crypto.PubkeyToAddress(*publicKey).String())

	// =================================================================================

	tx = Tx{
		FromID: "0xF01813E4B85e178A83e29B8E7bF26BD830a25f32",
		ToID:   "Frank",
		Value:  250,
	}

	data, err = json.Marshal(tx)
	if err != nil {
		return fmt.Errorf("unable to marshal %w", err)
	}

	stamp = []byte(fmt.Sprintf("\x19Ardan Signed Message:\n%d", len(data)))

	v2 := crypto.Keccak256(stamp, data)

	sig2, err := crypto.Sign(v2, privateKey)
	if err != nil {
		return fmt.Errorf("unable to sign transaction: %w", err)
	}

	fmt.Println("SIG: ", hexutil.Encode(sig2))

	// =================================================================================
	// had hard time understand what are we trying to accomplish:
	// basically we are using the previous sginature to acquire a new transactions public address which is wrong!!!
	// imagine you misspell it, this is scary! your public address changes!
	// but you can have from as the public address so you can validate!!!
	tx2 := Tx{
		FromID: "0xF01813E4B85e178A83e29B8E7bF26BD830a25f32",
		ToID:   "Frank",
		Value:  250,
	}

	data, err = json.Marshal(tx2)
	if err != nil {
		return fmt.Errorf("unable to marshal %w", err)
	}

	stamp = []byte(fmt.Sprintf("\x19Ardan Signed Message:\n%d", len(data)))

	v2 = crypto.Keccak256(stamp, data)

	publicKey, err = crypto.SigToPub(v2, sig2)
	if err != nil {
		return fmt.Errorf("unable to recover public key from signature: %w", err)
	}

	fmt.Println("PUB: ", crypto.PubkeyToAddress(*publicKey).String())

	vv, r, s, err := ToVRSFromHexSignature(hexutil.Encode(sig2))
	if err != nil {
		return fmt.Errorf("unable to VRS: %w", err)
	}

	fmt.Println("V|R|S: ", vv, r, s)
	// =================================================================================

	fmt.Println("==========================TX==============================")

	billTx, err := database.NewTx(1, 1,
		"0xF01813E4B85e178A83e29B8E7bF26BD830a25f32",
		"0xbEE6ACE826eC3DE1B6349888B9151B92522F7F76",
		1000,
		0,
		nil)
	if err != nil {
		return fmt.Errorf("unable to BillTx: %w", err)
	}
	signedTx, err := billTx.Sign(privateKey)
	if err != nil {
		return fmt.Errorf("unable to sign BillTx: %w", err)
	}
	fmt.Println(signedTx)
	return nil
}

// ToVRSFromHexSignature converts a hex representation of the signature into
// its R, S and V parts.
func ToVRSFromHexSignature(sigStr string) (v, r, s *big.Int, err error) {
	sig, err := hex.DecodeString(sigStr[2:])
	if err != nil {
		return nil, nil, nil, err
	}

	r = big.NewInt(0).SetBytes(sig[:32])
	s = big.NewInt(0).SetBytes(sig[32:64])
	v = big.NewInt(0).SetBytes([]byte{sig[64]})

	return v, r, s, nil
}
