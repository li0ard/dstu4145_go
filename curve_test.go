package dstu4145go_test

import (
	"encoding/hex"
	"testing"

	dstu4145go "github.com/li0ard/dstu4145_go"
)

func TestPubKey(t *testing.T) {
	c := dstu4145go.NewCurve(dstu4145go.DSTU_163_TEST())
	privateKey, _ := hex.DecodeString("0183F60FDF7951FF47D67193F8D073790C1C9B5A3E")
	publicKeyPoint := c.GetPublicKey(privateKey)

	expectedX := dstu4145go.HI("0x57de7fde023ff929cb6ac785ce4b79cf64abdc2da")
	expectedY := dstu4145go.HI("0x3e85444324bcf06ad85abf6ad7b5f34770532b9aa")

	if publicKeyPoint.X.Cmp(expectedX) != 0 || publicKeyPoint.Y.Cmp(expectedY) != 0 {
		t.Errorf("expected X = %s\nreceived X = %s\n\nexpected Y = %s\nreceived Y = %s",
			expectedX.Text(16),
			publicKeyPoint.X.Text(16),
			expectedY.Text(16),
			publicKeyPoint.Y.Text(16),
		)
	}
}

func TestSign(t *testing.T) {
	c := dstu4145go.NewCurve(dstu4145go.DSTU_163_TEST())
	privateKey, _ := hex.DecodeString("0183F60FDF7951FF47D67193F8D073790C1C9B5A3E")
	hash, _ := hex.DecodeString("09C9C44277910C9AAEE486883A2EB95B7180166DDF73532EEB76EDAEF52247FF")
	opts := &dstu4145go.SignOptions{
		Rand: dstu4145go.HI("0x1025E40BD97DB012B7A1D79DE8E12932D247F61C6"),
	}
	signature := c.Sign(privateKey, hash, opts)

	expectedR := dstu4145go.HI("0x274ea2c0caa014a0d80a424f59ade7a93068d08a7")
	expectedS := dstu4145go.HI("0x2100d86957331832b8e8c230f5bd6a332b3615aca")

	if signature.R.Cmp(expectedR) != 0 || signature.S.Cmp(expectedS) != 0 {
		t.Errorf("expected X = %s\nreceived X = %s\n\nexpected Y = %s\nreceived Y = %s",
			expectedR.Text(16),
			signature.R.Text(16),
			expectedS.Text(16),
			signature.S.Text(16),
		)
	}
}

func TestVerify(t *testing.T) {
	c := dstu4145go.NewCurve(dstu4145go.DSTU_163_TEST())

	publicKey := c.Point(
		dstu4145go.HI("0x57de7fde023ff929cb6ac785ce4b79cf64abdc2da"),
		dstu4145go.HI("0x3e85444324bcf06ad85abf6ad7b5f34770532b9aa"),
	)
	hash, _ := hex.DecodeString("09C9C44277910C9AAEE486883A2EB95B7180166DDF73532EEB76EDAEF52247FF")

	result := c.Verify(publicKey, hash, &dstu4145go.Signature{
		R: dstu4145go.HI("0x274ea2c0caa014a0d80a424f59ade7a93068d08a7"),
		S: dstu4145go.HI("0x2100d86957331832b8e8c230f5bd6a332b3615aca"),
	})

	if result != true {
		t.Error("")
	}
}
