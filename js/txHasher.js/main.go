package main

import (
	"encoding/json"

	"github.com/MadBase/MadNet/application/objs"
	"github.com/gopherjs/gopherjs/js"
	"github.com/miratronix/jopher"
)

func main() {
	js.Module.Get("exports").Set("TxHasher", jopher.Promisify(TxHasher))
}

func TxHasher(s string) (string, error) {
	tx := &objs.Tx{}
	err := json.Unmarshal([]byte(s), tx)
	if err != nil {
		return "", err
	}

	err = tx.SetTxHash()
	if err != nil {
		return "", err
	}
	for _, v := range tx.Vin {
		preHash, err := v.TXInLinker.MarshalBinary()
		if err != nil {
			return "", err
		}
		v.Signature = preHash
	}
	for _, v := range tx.Vout {
		switch {
		case v.HasDataStore():
			ds, err := v.DataStore()
			if err != nil {
				return "", err
			}
			preHash, err := ds.DSLinker.MarshalBinary()
			if err != nil {
				return "", err
			}
			ds.Signature.Signature = preHash
		default:
			// do nothing
		}
	}
	sb, err := json.Marshal(tx)
	if err != nil {
		return "", err
	}
	return string(sb), nil

}
