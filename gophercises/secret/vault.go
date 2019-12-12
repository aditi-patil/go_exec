package secret

import (
	"encoding/json"
	"errors"
	"fmt"
	"gophercises/secret/cipher"
	"io"
	"os"
	"sync"
)

// File is initialisation method for vault
func File(encodingKey, filepath string) *Vault {
	return &Vault{
		encodingKey: encodingKey,
		filepath:    filepath,
		keyValues:   make(map[string]string),
	}
}

// Vault is a struct which defines secret key parameters
type Vault struct {
	encodingKey string
	filepath    string
	mutex       sync.Mutex
	keyValues   map[string]string
}

//Load : test
func (v *Vault) Load() error {
	f, err := os.Open(v.filepath)
	if err != nil {
		v.keyValues = make(map[string]string)
		return nil
	}
	defer f.Close()
	r, err := cipher.NewDecryptReader(v.encodingKey, f)
	if err != nil {
		return err
	}
	return v.readKeyValues(r)
}

func (v *Vault) readKeyValues(r io.Reader) error {
	dec := json.NewDecoder(r)
	err := dec.Decode(&v.keyValues)
	if err == io.EOF {
		v.keyValues = make(map[string]string)
		return nil
	}
	return err
}

//Save Save
func (v *Vault) Save() error {
	f, err := os.OpenFile(v.filepath, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	defer f.Close()
	w, err := cipher.EncryptWriter(v.encodingKey, f)
	if err != nil {
		return err
	}
	return v.writeKeyValues(w)
}

func (v *Vault) writeKeyValues(w io.Writer) error {
	enc := json.NewEncoder(w)
	return enc.Encode(v.keyValues)
}

// Get will give the value of given key from the secret
func (v *Vault) Get(key string) (string, error) {
	v.mutex.Lock()
	defer v.mutex.Unlock()
	err := v.Load()
	if err != nil {
		return "", err
	}
	value, ok := v.keyValues[key]
	if !ok {
		return "", errors.New("secret: no value for that key")
	}
	return value, nil
}

// Set will add key value pair in the secret
func (v *Vault) Set(key, value string) error {
	v.mutex.Lock()
	defer v.mutex.Unlock()
	err := v.Load()
	fmt.Println(err)
	if err != nil {
		return err
	}
	v.keyValues[key] = value
	err = v.Save()
	return err
}

// Remove will remove given key from secret
func (v *Vault) Remove(key string) error {
	v.mutex.Lock()
	defer v.mutex.Unlock()
	err := v.Load()
	if err != nil {
		return err
	}
	if _, ok := v.keyValues[key]; ok {
		delete(v.keyValues, key)
	} else {
		return errors.New("secret: no value for that key")
	}
	err = v.Save()
	return err
}
