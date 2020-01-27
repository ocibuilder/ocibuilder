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
	"github.com/ocibuilder/ocibuilder/testing/dummy"
	"github.com/stretchr/testify/assert"
)

func TestDecodeKey(t *testing.T) {
	_, err := DecodeKey(dummy.TestPrivKey)
	assert.Equal(t, nil, err)

	_, err = DecodeKey("invalid-key")
	assert.Error(t, err, io.EOF)
}

func TestSignDigest(t *testing.T) {
	privKey, pubKey, err := ValidateKeysPacket(signKeyPlain)
	assert.Equal(t, nil, err)
	e := CreateEntityFromKeys(privKey, pubKey)

	_, _, err = SignDigest("SHAtestDigeste3413412", "", e)
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
	PlainPrivateKey: dummy.TestPrivKey,
	PlainPublicKey:  dummy.TestPubKey,
}
