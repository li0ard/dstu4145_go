<p align="center">
    <b>dstu4145_go</b><br>
    <b>DSTU 4145-2002 curves and DSA</b>
    <br><br>
    <a href="https://pkg.go.dev/github.com/li0ard/dstu4145_go"><img src="https://pkg.go.dev/badge/github.com/li0ard/dstu4145_go.svg" /></a>
    <a href="https://github.com/li0ard/dstu4145_go/blob/main/LICENSE"><img src="https://img.shields.io/github/license/li0ard/dstu4145_go" /></a>
    <br>
    <hr>
</p>

## Examples
```go
package main

import (
	"encoding/hex"
	"fmt"

	dstu4145go "github.com/li0ard/dstu4145_go"
)

c := dstu4145go.NewCurve(dstu4145go.DSTU_163_TEST())
privateKey, _ := hex.DecodeString("00000000000000000000000183F60FDF7951FF47D67193F8D073790C1C9B5A3E")
publicKey := c.GetPublicKey(privateKey)
fmt.Println("x =", publicKey.X.Text(16))
fmt.Println("y =", publicKey.Y.Text(16))
hash, _ := hex.DecodeString("09C9C44277910C9AAEE486883A2EB95B7180166DDF73532EEB76EDAEF52247FF")
var signature = c.Sign(privateKey, hash, nil)
fmt.Println("r =", signature.R.Text(16))
fmt.Println("s =", signature.S.Text(16))
fmt.Println("check =", c.Verify(publicKey, hash, signature))
```