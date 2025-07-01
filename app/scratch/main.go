package main

import (
	"encoding/json"
	"fmt"
	"log"

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
	tx := Tx{
		FromID: "Bill",
		ToID:   "Aaron",
		Value:  1000,
	}

	privateKey, err := crypto.LoadECDSA("zblock/accounts/kennedy.ecdsa")
	if err != nil {
		return fmt.Errorf("unable to load private key for node: %w", err)
	}

	data, err := json.Marshal(tx)
	if err != nil {
		return fmt.Errorf("unable to marshal %w", err)
	}

	v := crypto.Keccak256(data)

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

	return nil
}
