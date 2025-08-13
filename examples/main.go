package main

import (
	"encoding/hex"
	"fmt"

	dstu4145go "github.com/li0ard/dstu4145_go"
)

func DSTU_163() {
	fmt.Println("============== DSTU 163 ==============")
	c := dstu4145go.NewCurve(dstu4145go.DSTU_163_TEST())
	privateKey, _ := hex.DecodeString("00000000000000000000000183F60FDF7951FF47D67193F8D073790C1C9B5A3E")
	publicKey := c.GetPublicKey(privateKey)
	fmt.Println("Public key:")
	fmt.Println("x =", publicKey.X.Text(16))
	fmt.Println("y =", publicKey.Y.Text(16))
	hash, _ := hex.DecodeString("09C9C44277910C9AAEE486883A2EB95B7180166DDF73532EEB76EDAEF52247FF")
	opts := &dstu4145go.SignOptions{
		Rand: dstu4145go.HI("0x1025E40BD97DB012B7A1D79DE8E12932D247F61C6"),
	}

	var signature = c.Sign(privateKey, hash, opts)
	fmt.Println("\nSignature:")
	fmt.Println("r =", signature.R.Text(16))
	fmt.Println("s =", signature.S.Text(16))
	fmt.Println("check =", c.Verify(publicKey, hash, signature))
}

func DSTU_257() {
	fmt.Println("\n============== DSTU 257 ==============")
	c := dstu4145go.NewCurve(dstu4145go.DSTU_257())
	privateKey, _ := hex.DecodeString("77fd46e42b36b76f551426cafcecdd2f8f4e0df00ea62886e4343d59da35fb0fbf")
	publicKey := c.GetPublicKey(privateKey)
	fmt.Println("Public key:")
	fmt.Println("x =", publicKey.X.Text(16))
	fmt.Println("y =", publicKey.Y.Text(16))
	hash, _ := hex.DecodeString("09C9C44277910C9AAEE486883A2EB95B7180166DDF73532EEB76EDAEF52247FF")
	var signature = c.Sign(privateKey, hash, nil)
	fmt.Println("\nSignature:")
	fmt.Println("r =", signature.R.Text(16))
	fmt.Println("s =", signature.S.Text(16))
	fmt.Println("check =", c.Verify(publicKey, hash, signature))
}

func DSTU_431() {
	fmt.Println("\n============== DSTU 431 ==============")
	c := dstu4145go.NewCurve(dstu4145go.DSTU_431())
	privateKey, _ := hex.DecodeString("b307af47e9d37d53a9fcbf6d6e9e09340f580f903fe48a84bbf77b0bf7ab379507e04ce4a8ea29124910f1ab161a0e690b3ac1329609")
	publicKey := c.GetPublicKey(privateKey)
	fmt.Println("Public key:")
	fmt.Println("x =", publicKey.X.Text(16))
	fmt.Println("y =", publicKey.Y.Text(16))
	hash, _ := hex.DecodeString("09C9C44277910C9AAEE486883A2EB95B7180166DDF73532EEB76EDAEF52247FF")
	var signature = c.Sign(privateKey, hash, nil)
	fmt.Println("\nSignature:")
	fmt.Println("r =", signature.R.Text(16))
	fmt.Println("s =", signature.S.Text(16))
	fmt.Println("check =", c.Verify(publicKey, hash, signature))
}

func main() {
	DSTU_163()
	DSTU_257()
	DSTU_431()
	fmt.Println("======================================")
}
