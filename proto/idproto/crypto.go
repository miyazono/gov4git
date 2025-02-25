package idproto

import (
	"context"
	"crypto/ed25519"

	"github.com/gov4git/gov4git/lib/form"
)

type Ed25519PublicKey = form.Bytes

type Ed25519PrivateKey = form.Bytes

func GenerateCredentials(publicURL, privateURL string) (*PrivateCredentials, error) {
	pubKey, privKey, err := ed25519.GenerateKey(nil)
	if err != nil {
		return nil, err
	}
	return &PrivateCredentials{
		PrivateURL:        privateURL,
		PrivateKeyEd25519: Ed25519PrivateKey(privKey),
		PublicCredentials: PublicCredentials{
			ID:               GenerateUniqueID(),
			PublicURL:        publicURL,
			PublicKeyEd25519: Ed25519PublicKey(pubKey),
		},
	}, nil
}

type SignedPlaintext struct {
	Plaintext        form.Bytes       `json:"plaintext"`
	Signature        form.Bytes       `json:"signature"`
	PublicKeyEd25519 Ed25519PublicKey `json:"ed25519_public_key"`
}

func ParseSignedPlaintext(ctx context.Context, data []byte) (*SignedPlaintext, error) {
	var signed SignedPlaintext
	if err := form.DecodeForm(ctx, data, &signed); err != nil {
		return nil, err
	}
	return &signed, nil
}

func SignPlaintext(ctx context.Context, priv *PrivateCredentials, plaintext []byte) (*SignedPlaintext, error) {
	signature := ed25519.Sign(ed25519.PrivateKey(priv.PrivateKeyEd25519), plaintext)
	return &SignedPlaintext{
		Plaintext:        plaintext,
		Signature:        signature,
		PublicKeyEd25519: priv.PublicCredentials.PublicKeyEd25519,
	}, nil
}

func (signed *SignedPlaintext) Verify() bool {
	return ed25519.Verify(ed25519.PublicKey(signed.PublicKeyEd25519), signed.Plaintext, signed.Signature)
}
