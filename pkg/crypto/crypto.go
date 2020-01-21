/*
Copyright © 2019 BlackRock Inc.

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
	"crypto"
	"errors"
	"io"
	"os"
	"strings"

	"github.com/ocibuilder/ocibuilder/pkg/apis/ocibuilder/v1alpha1"
	"github.com/ocibuilder/ocibuilder/pkg/util"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
	"golang.org/x/crypto/openpgp/packet"
)

func CreateEntityFromKeys(privKey *packet.PrivateKey, pubKey *packet.PublicKey) *openpgp.Entity {
	config := packet.Config{
		DefaultHash:            crypto.SHA3_256,
		DefaultCipher:          packet.CipherAES256,
		DefaultCompressionAlgo: packet.CompressionZLIB,
	}

	e := openpgp.Entity{
		PrimaryKey: pubKey,
		PrivateKey: privKey,
		Identities: make(map[string]*openpgp.Identity),
	}

	currentTime := config.Now()
	uid := packet.NewUserId("", "", "")
	isPrimaryId := false

	e.Identities[uid.Id] = &openpgp.Identity{
		Name:   "",
		UserId: nil,
		SelfSignature: &packet.Signature{
			CreationTime: currentTime,
			SigType:      packet.SigTypePositiveCert,
			PubKeyAlgo:   packet.PubKeyAlgoRSA,
			Hash:         config.Hash(),
			IsPrimaryId:  &isPrimaryId,
			FlagsValid:   true,
			FlagSign:     true,
			FlagCertify:  true,
			IssuerKeyId:  &e.PrimaryKey.KeyId,
		},
	}

	return &e
}

func ValidateKeys(key *v1alpha1.SignKey) (privKey, pubKey string, err error) {
	log := util.Logger

	if key.EnvPrivateKey != "" && key.EnvPublicKey != "" {
		privKey = os.Getenv(key.EnvPrivateKey)
		if privKey == "" {
			log.Warn("environment variable empty for private key")
		}

		pubKey = os.Getenv(key.EnvPublicKey)
		if pubKey == "" {
			log.Warn("environment variable empy for public key")
		}

		return privKey, pubKey, nil
	}

	if key.PlainPrivateKey != "" && key.PlainPublicKey != "" {
		return key.PlainPrivateKey, key.PlainPublicKey, nil
	}

	return "", "", errors.New("no private and public keys found in specification")
}

func ValidateKeysPacket(key *v1alpha1.SignKey) (*packet.PrivateKey, *packet.PublicKey, error) {

	privKeyStr, pubKeyStr, err := ValidateKeys(key)
	if err != nil {
		return nil, nil, err
	}

	privPack, err := DecodeKey(privKeyStr)
	if err != nil {
		return nil, nil, err
	}
	privKey, ok := privPack.(*packet.PrivateKey)
	if !ok {
		return nil, nil, errors.New("invalid pgp private key when validating for signing image")
	}

	pubPack, err := DecodeKey(pubKeyStr)
	if err != nil {
		return nil, nil, err
	}
	pubKey, ok := pubPack.(*packet.PublicKey)
	if !ok {
		return nil, nil, errors.New("invalid pgp private key when validating for signing image")
	}

	return privKey, pubKey, nil
}

func DecodeKey(key string) (packet.Packet, error) {
	block, err := armor.Decode(strings.NewReader(key))
	log := util.Logger
	if err != nil {
		log.Error("error decoding key - no block found due to invalid PGP key")
		return nil, err
	}

	reader := packet.NewReader(block.Body)
	pkt, err := reader.Next()
	if err != nil {
		log.Error("error decoding key - unknown packet due to invalid PGP key")
		return nil, err
	}

	return pkt, nil
}

func SignDigest(digest string, signer *openpgp.Entity) (io.Writer, error) {

	if err := signer.PrivateKey.Decrypt([]byte{}); err != nil {
		return nil, err
	}

	if err := openpgp.ArmoredDetachSignText(os.Stdout, signer, strings.NewReader(digest), nil); err != nil {
		return nil, err
	}

	return nil, nil
}
