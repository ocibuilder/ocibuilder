/*
Copyright 2019 BlackRock, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package crypto

import (
	"io"
	"os"
	"testing"

	"github.com/ocibuilder/ocibuilder/pkg/apis/ocibuilder/v1alpha1"
	"github.com/stretchr/testify/assert"
)

func TestDecodeKey(t *testing.T) {
	_, err := DecodeKey(privKey)
	assert.Equal(t, nil, err)

	_, err = DecodeKey("invalid-key")
	assert.Error(t, err, io.EOF)
}

func TestSignDigest(t *testing.T) {
	privKey, pubKey, err := ValidateKeysPacket(signKeyPlain)
	assert.Equal(t, nil, err)
	e := CreateEntityFromKeys(privKey, pubKey)

	_, err = SignDigest("SHAtestDigeste3413412", "", e)
	assert.Equal(t, nil, err)
}

func TestValidateKeys(t *testing.T) {
	err := os.Setenv("PRI_KEY", "this-is-a-private-key")
	assert.Equal(t, nil, err)
	err = os.Setenv("PUB_KEY", "this-is-a-public-key")
	assert.Equal(t, nil, err)

	privKey, pubKey, err := ValidateKeys(signKeyEnv)

	assert.Equal(t, nil, err)
	assert.Equal(t, "this-is-a-private-key", privKey)
	assert.Equal(t, "this-is-a-public-key", pubKey)

	signKeyEnv.EnvPublicKey = ""
	_, _, err = ValidateKeys(signKeyEnv)
	assert.Error(t, err, "no private and public keys found in specification")

	err = os.Unsetenv("PRI_KEY")
	assert.Equal(t, nil, err)
	err = os.Unsetenv("PUB_KEY")
	assert.Equal(t, nil, err)
}

func TestValidateKeysPacket(t *testing.T) {
	_, _, err := ValidateKeysPacket(signKeyPlain)
	assert.Equal(t, nil, err)
}

var signKeyEnv = &v1alpha1.SignKey{
	EnvPrivateKey: "PRI_KEY",
	EnvPublicKey:  "PUB_KEY",
}

var signKeyPlain = &v1alpha1.SignKey{
	PlainPrivateKey: privKey,
	PlainPublicKey:  pubKey,
}

// This key is NOT used anywhere other than this unit test
var privKey = `-----BEGIN PGP PRIVATE KEY BLOCK-----
Version: BCPG C# v1.6.1.0

lQOsBF4m0SQBCAC+VXNpTIglJlXIeEaseUL3aTqJmWnJE1Vpu7MYqT9rCKtnKlcN
BU2+WfcrG9ZJD/zUETxtw2m8nCZH/0K7XxjZPLo0qDfbM/giP9EcSJzmDeaAUMEo
buG/M5yxHReuysGZI/2X3Fw5swMr1mGOTSf6JLY6xecqlxgpI/N1IEWIKGSUmwlZ
fzhuBuV0EATov5zJ7XHdDUljrP3EdhJ3nurPwaUWkFjVeZEgqzj8QP0u8dIKQd4R
m/qWeR/cVa6btY8r2t4Cay4/ER8iYAHeqJHJlJmtCys+xEFzVwcemw27D+Mq4Ck6
wzc+3eypF2JJdrhicBqJ5JauG9NUCi8rynN7ABEBAAH/AwMCpW8FKB/oQuxgfeCG
X50b0fVmpRi+AESOGKDM/uUSISXQZUPxWojoKPG61E43oAD8utzpsI7TnnbjH1os
bEilVH/6QesYVXqMZIhMPdVyFPnTEr05MwkwaA4UeXLGHX5JsGd0l3nZtFGQlDp8
Tyxj7nUSUGEnQRmNZWzMnD6wfiMuaB2nfPAYjFPPeoVF09FGXpJcMPunbNenhinH
76I8M/OWFiUcBg6pEOy64ZoVG1sblKVcxC2Mv1g0koLQAANGJ4M6mjjQfJCcL0MW
Qd8C7bupd3m0Ph/S6LUPmH2ljmUhtaf42VNmMK9MAVcmxXYEyEefPCB3PKYBS7M1
hlzjwnpkscB0pzzfaq4AdveQMujHNG0rWIaKvvL4TzDauWnkYp0FnylFq8IewIxO
eTiflSjM/eWmAbGtUbonUMajokaTmR+DftJpeb0TBgFWKsj7bGf8xOHHDGl2T9fZ
zQOOdVO6ACNs60hj/hozv9sMFOTaJk1zX7kFXGJJiFeB+aCbyRrVKTSg4CKQFpXn
pAPDwc2KhEZZ1hAKJvL/mC565RJJvAveWcB0CB7w+0QOMkFrOTaqMVAwVvYop1xN
/qAzrzK3+5RZwb3ajbk7ShGE6JgdauYO2CgLPP0+jy36Iw+pgyTZYrLlH9U9k1ID
hsFYO8nhk+IkWc2RjIDs5ol3ym2DpHBjakWEaGwGiNfFWvxBH3Cp568J81CS6bSu
p7XLZLB97gBv8VAIP2t1UlJmnF/NKaBU76oacIccFRHfSIm+FE/crHYn+jyRZHnS
6y1tg5F4D1xkuKv1P3o1rxAmxkaifgHMCiQWGs3BsVFJe7OKXvIkWcuvnlFOevkc
dHUAs713JXWJYDQHDpKOo3c0iu1mseFCvdPjOton1bQAiQEcBBABAgAGBQJeJtEk
AAoJEE0wzYTUb94tRj8IALeVuBaZMb3HVRd589hxLYXNlzaK4WuMyOZavXzOLzji
bhPceiq3LRsXFY4U1xx9CtzyhTb8t0QrlLYgrTWXNvovezXlrDPuWH7J+5jPyy4o
3CAKGqTL+pVBRVM3MI+4D/wRatKM8uc08iCNJmuZI55sAmbZJR8IeQCgBzGf3cY/
0WxIKje8zQHHms+M3T3sQul07OoDD1qAVVWtWbbLPah36u18Gc77GiC2DtVoi8ux
m6LlB08sbpjUhjmxF+A34jPuKsLVP/gfyGktMQ4phtDL3T3cRstbNKDfX/IY9fAF
ZH7U4zGg6Wi5dl1oEMoLLavkiYK4Czwf9pRIpxhNIpE=
=fQz/
-----END PGP PRIVATE KEY BLOCK-----
`

var pubKey = `-----BEGIN PGP PUBLIC KEY BLOCK-----
Version: BCPG C# v1.6.1.0

mQENBF4m0SQBCAC+VXNpTIglJlXIeEaseUL3aTqJmWnJE1Vpu7MYqT9rCKtnKlcN
BU2+WfcrG9ZJD/zUETxtw2m8nCZH/0K7XxjZPLo0qDfbM/giP9EcSJzmDeaAUMEo
buG/M5yxHReuysGZI/2X3Fw5swMr1mGOTSf6JLY6xecqlxgpI/N1IEWIKGSUmwlZ
fzhuBuV0EATov5zJ7XHdDUljrP3EdhJ3nurPwaUWkFjVeZEgqzj8QP0u8dIKQd4R
m/qWeR/cVa6btY8r2t4Cay4/ER8iYAHeqJHJlJmtCys+xEFzVwcemw27D+Mq4Ck6
wzc+3eypF2JJdrhicBqJ5JauG9NUCi8rynN7ABEBAAG0AIkBHAQQAQIABgUCXibR
JAAKCRBNMM2E1G/eLUY/CAC3lbgWmTG9x1UXefPYcS2FzZc2iuFrjMjmWr18zi84
4m4T3Hoqty0bFxWOFNccfQrc8oU2/LdEK5S2IK01lzb6L3s15awz7lh+yfuYz8su
KNwgChqky/qVQUVTNzCPuA/8EWrSjPLnNPIgjSZrmSOebAJm2SUfCHkAoAcxn93G
P9FsSCo3vM0Bx5rPjN097ELpdOzqAw9agFVVrVm2yz2od+rtfBnO+xogtg7VaIvL
sZui5QdPLG6Y1IY5sRfgN+Iz7irC1T/4H8hpLTEOKYbQy9093EbLWzSg31/yGPXw
BWR+1OMxoOlouXZdaBDKCy2r5ImCuAs8H/aUSKcYTSKR
=oqhx
-----END PGP PUBLIC KEY BLOCK-----
`
