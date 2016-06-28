// +build pkcs11

package client

import (
	"fmt"
	"net/http"
	"github.com/docker/notary/passphrase"
	"github.com/docker/notary"
	"github.com/docker/notary/trustmanager"
	"github.com/docker/notary/trustmanager/yubikey"
	"github.com/docker/notary/trustpinning"
)

// NewNotaryRepository is a helper method that returns a new notary repository.
// It takes the base directory under where all the trust files will be stored
// (usually ~/.docker/trust/).
func NewNotaryRepository(baseDir, gun, baseURL string, rt http.RoundTripper,
	retriever notary.PassRetriever, trustPinning trustpinning.TrustPinConfig, UseNative bool) (
	*NotaryRepository, error) {

	fileKeyStore, err := trustmanager.NewKeyFileStore(baseDir, retriever)
	if err != nil {
		return nil, fmt.Errorf("failed to create private key store in directory: %s", baseDir)
	}

	keyStores := []trustmanager.KeyStore{fileKeyStore}
	if UseNative {
		nativeKeyStore, err := trustmanager.NewKeyNativeStore(passphrase.ConstantRetriever("password"))
		if err == nil {
			// Note that the order is important, since we want to prioritize
			// the native key store
			keyStores = append([]trustmanager.KeyStore{nativeKeyStore}, keyStores...)
		}
	}
	yubiKeyStore, _ := yubikey.NewYubiStore(fileKeyStore, retriever)
	if yubiKeyStore != nil {
		// Note that the order is important, since we want to prioritize
		// the yubi key store
		keyStores=append([]trustmanager.KeyStore{yubiKeyStore},keyStores...)
	}

	return repositoryFromKeystores(baseDir, gun, baseURL, rt, keyStores, trustPinning)
}
