package trustmanager

import (
	"github.com/docker/docker-credential-helpers/credentials"
	"github.com/docker/docker-credential-helpers/client"
	"encoding/base64"
	"errors"
	"github.com/docker/notary/tuf/data"
	"strings"
	"fmt"
	"github.com/docker/notary"
	//"os/exec"
	//"regexp"
)

// KeyNativeStore is an implementation of Storage that keeps
// the contents in the keychain access.
type KeyNativeStore struct {
	notary.PassRetriever
	newProgFunc client.ProgramFunc
}

// NewKeyNativeStore creates a KeyNativeStore
func NewKeyNativeStore(passphraseRetriever notary.PassRetriever) (*KeyNativeStore, error) {
	if defaultCredentialsStore=="" {
		return nil, errors.New("Native storage on your operating system is not yet supported")
	}
	credCommand:="docker-credential-"+defaultCredentialsStore
	x:=client.NewShellProgramFunc(credCommand)
	return &KeyNativeStore{
		PassRetriever: passphraseRetriever,
		newProgFunc:x,
	}, nil
}

//AddKey writes data new KeyChain in the native keychain store
func (k *KeyNativeStore) AddKey(keyInfo KeyInfo, privKey data.PrivateKey) error {
	pemPrivKey, err := KeyToPEM(privKey, keyInfo.Role)
	if err!=nil{
		return err
	}
	secretByte:=base64.StdEncoding.EncodeToString(pemPrivKey)
	keyCredentials:=credentials.Credentials{
		ServerURL:privKey.ID(),
		Username:keyInfo.Gun+"<notary_key>"+keyInfo.Role,
		Secret:secretByte,
	}
	err=client.Store(k.newProgFunc,&(keyCredentials))
	return err
}

// GetKey returns the credentials from the native keychain store given a server name
func (k *KeyNativeStore) GetKey(keyID string) (data.PrivateKey, string, error) {
	serverName:=keyID
	gotCredentials,err:=client.Get(k.newProgFunc,serverName)
	if err!=nil {
		return nil, "", err
	}
	gotSecret:=gotCredentials.Secret
	gotSecretByte,err:=base64.StdEncoding.DecodeString(gotSecret)
	privKey, err := ParsePEMPrivateKey(gotSecretByte, "")
	role:=strings.SplitAfter(gotCredentials.Username, "<notary_key>")[1]
	return privKey, role, err
}

// GetKeyInfo returns the corresponding gun and role key info for a keyID
func (k *KeyNativeStore) GetKeyInfo(keyID string) (KeyInfo, error) {
	serverName:=keyID
	gotCredentials,err:=client.Get(k.newProgFunc,serverName)
	if err!=nil {
		return KeyInfo{}, err
	}
	if strings.Contains(gotCredentials.Username, "<notary_key>") {
		keyinfo := strings.SplitAfter(gotCredentials.Username, "<notary_key>")
		gun := keyinfo[0][:(len(keyinfo[0])-12)]
		return KeyInfo{
		Gun: gun,
		Role: keyinfo[1],
		}, err
	}
	return KeyInfo{}, fmt.Errorf("The keyID doesn't belong to a Notary key")
}

// ListKeys lists all the Keys inside of a native store
// Just a placeholder for now- returns an empty slice
func (k *KeyNativeStore) ListKeys() map[string]KeyInfo {
	m := make(map[string]KeyInfo)
	//working list functionality on osxkeychain. if uncommenting, also uncomment the imports at the file start
	//if defaultCredentialsStore=="osxkeychain" {
	//	out, _ := exec.Command("security","dump").Output()
	//	output:=string(out)
	//	re := regexp.MustCompile("path\"<blob>=\"(\\d|[a-z]){64}\"")
	//	keys:=re.FindAllString(output,-1)
	//	for _,key := range keys {
	//		keyId:=key[13:(len(key)-1)]
	//		keyInfo, _:=k.GetKeyInfo(keyId)
	//		m[keyId]=keyInfo
	//	}
	//}
	return m
}

//RemoveKey removes a KeyChain (identified by server name- a string) from the keychain access store
//Currently when we try to remove a key that doesn't exist: in osx, RemoveKey throws an error. in linux, RemoveKey doesn't throw an error
//This is due to inconsistent behaviour from the credentials-helper which can be corrected if necessary
//The inconsistency can be seen clearly if we read through TestRemoveFromNativeStoreNoPanic in keynativestore_test.go
func (k *KeyNativeStore) RemoveKey(keyID string) error {
	err:=client.Erase(k.newProgFunc,keyID)
	return err
}

//ExportKey removes a KeyChain from the keychain access store as an encrypted byte string
func (k *KeyNativeStore) ExportKey(keyID string) ([]byte, error) {
	serverName:=keyID
	gotKey, role, err:= k.GetKey(serverName)
	if err!=nil {
		return nil, err
	}
	// take in a passphrase with the given retriever
	var (
		chosenPassphrase string
		giveup           bool
	)

	for attempts := 0; ; attempts++ {
		chosenPassphrase, giveup, err = k.PassRetriever(keyID, role, true, attempts)
		if giveup {
			return nil, errors.New("Given up")
		}
		if attempts > 3 {
			return nil, errors.New("Exceeded attempts, please select a secure passphrase and type it with care")
		}
		if err != nil {
			continue
		}
		break
	}
	// encrypt the byte string
	encSecret, err:=EncryptPrivateKey(gotKey, role, chosenPassphrase)
	return encSecret, err
}

// Name returns a user friendly name for the location this store
// keeps its data, here it is the name of the native store on this operating system
func (k *KeyNativeStore) Name() string {
	return fmt.Sprintf("Native keychain store: %s", defaultCredentialsStore)
}

