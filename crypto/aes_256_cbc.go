package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/binary"
	"fmt"
	"hash"

	"golang.org/x/crypto/pbkdf2"
)

const (
	SHA256 = "SHA256"
	SHA384 = "SHA384"
	SHA512 = "SHA512"
)

type Aes256CBC struct {
	Iterations int32
	KeySize    int
	version    int16
	saltSize   int16
	hashAlgo   string
}

func NewAes256CBC() *Aes256CBC {
	return &Aes256CBC{
		Iterations: 10000,
		KeySize:    256,
		version:    2,
		saltSize:   64,
		hashAlgo:   SHA384,
	}
}

func (a *Aes256CBC) SetHashAlgo(hashAlgo string) error {
	switch hashAlgo {
	case SHA256:
		a.hashAlgo = SHA256
	case SHA384:
		a.hashAlgo = SHA384
	case SHA512:
		a.hashAlgo = SHA512
	default:
		return fmt.Errorf("invalid hash algo")
	}

	return nil
}

func (a *Aes256CBC) Encrypt(key []byte, data []byte) (encryptedData []byte, err error) {
	return a.EncryptWithMetadata(key, data, nil)
}

func (a *Aes256CBC) EncryptWithMetadata(key []byte, data []byte, metadata []byte) (encryptedData []byte, err error) {

	// 1. version   2
	// 2. metadataSize  4
	// 3. iterations  4
	// 4. symmetricSaltSize 2
	// 5. signingSaltSize 2
	// 6. symmetricSalt 16
	// 7. signingSalt 16
	// 8. iv 16

	saltSize := a.saltSize / 8
	buf := new(bytes.Buffer)
	err = binary.Write(buf, binary.LittleEndian, a.version)
	if err != nil {
		return nil, err
	}

	metadataSize := int32(len(metadata))
	err = binary.Write(buf, binary.LittleEndian, metadataSize)
	if err != nil {
		return nil, err
	}

	err = binary.Write(buf, binary.LittleEndian, a.Iterations)
	if err != nil {
		return nil, err
	}

	err = binary.Write(buf, binary.LittleEndian, saltSize)
	if err != nil {
		return nil, err
	}

	err = binary.Write(buf, binary.LittleEndian, saltSize)
	if err != nil {
		return nil, err
	}

	symetricSalt, err := RandBytes(int(saltSize))
	if err != nil {
		return nil, err
	}
	signingSalt, err := RandBytes(int(saltSize))
	if err != nil {
		return nil, err
	}

	iv, err := RandBytes(16)
	if err != nil {
		return nil, err
	}

	err = binary.Write(buf, binary.LittleEndian, symetricSalt)
	if err != nil {
		return nil, err
	}

	err = binary.Write(buf, binary.LittleEndian, signingSalt)
	if err != nil {
		return nil, err
	}

	err = binary.Write(buf, binary.LittleEndian, iv)
	if err != nil {
		return nil, err
	}

	if metadataSize > 0 {
		buf.Write(metadata)
	}

	keySize := a.KeySize / 8
	cdr := pbkdf2.Key(key, symetricSalt, int(a.Iterations), keySize, sha256.New)
	paddedData := pad(data)
	ciphertext := make([]byte, len(paddedData))
	c, _ := aes.NewCipher(cdr)
	ctr := cipher.NewCBCEncrypter(c, iv)
	ctr.CryptBlocks(ciphertext, paddedData)

	hdr := pbkdf2.Key(key, signingSalt, int(a.Iterations), keySize, sha256.New)
	h := a.NewHmac(hdr)
	h.Write(ciphertext)
	hash := h.Sum(nil)

	bufLen := buf.Len()
	hashLen := len(hash)
	ciphertextLen := len(ciphertext)

	result := make([]byte, bufLen+hashLen+ciphertextLen)
	copy(result, buf.Bytes())
	copy(result[bufLen:], hash)
	copy(result[bufLen+hashLen:], ciphertext)

	return result, nil
}

func (a *Aes256CBC) Decrypt(key []byte, encryptedData []byte) (data []byte, err error) {
	decryptedData, _, err := a.DecryptWithMetadata(key, encryptedData)
	return decryptedData, err
}

func (a *Aes256CBC) DecryptWithMetadata(key []byte, encryptedData []byte) (data []byte, metadata []byte, err error) {
	keySize := a.KeySize / 8

	var version int16
	err = binary.Read(bytes.NewReader(encryptedData), binary.LittleEndian, &version)
	if err != nil {
		return nil, nil, err
	}

	if version != a.version {
		return nil, nil, fmt.Errorf("invalid version for Aes256CBC")
	}

	var metadataSize int32
	err = binary.Read(bytes.NewReader(encryptedData[2:6]), binary.LittleEndian, &metadataSize)
	if err != nil {
		return nil, nil, err
	}

	var iterations int32
	err = binary.Read(bytes.NewReader(encryptedData[6:10]), binary.LittleEndian, &iterations)
	if err != nil {
		return nil, nil, err
	}

	var symmetricSaltSize int16
	err = binary.Read(bytes.NewReader(encryptedData[10:12]), binary.LittleEndian, &symmetricSaltSize)
	if err != nil {
		return nil, nil, err
	}

	var signingSaltSize int16
	err = binary.Read(bytes.NewReader(encryptedData[12:14]), binary.LittleEndian, &signingSaltSize)
	if err != nil {
		return nil, nil, err
	}

	symmetricSalt := encryptedData[14 : 14+symmetricSaltSize]
	signingSalt := encryptedData[14+symmetricSaltSize : 14+symmetricSaltSize+signingSaltSize]
	iv := encryptedData[14+symmetricSaltSize+signingSaltSize : 14+symmetricSaltSize+signingSaltSize+16]

	metadata = encryptedData[14+symmetricSaltSize+signingSaltSize+16 : 14+int(symmetricSaltSize)+int(signingSaltSize)+16+int(metadataSize)]

	hashStart := 14 + int(symmetricSaltSize) + int(signingSaltSize) + 16 + int(metadataSize)

	hash := encryptedData[hashStart : hashStart+a.GetHashSize()]
	ciphertext := encryptedData[hashStart+a.GetHashSize():]

	hdr := pbkdf2.Key(key, signingSalt, int(iterations), keySize, sha256.New)
	h := a.NewHmac(hdr)
	h.Write(ciphertext)
	expectedHash := h.Sum(nil)

	if !hmac.Equal(hash, expectedHash) {
		return nil, nil, fmt.Errorf("hash mismatch")
	}

	cdr := pbkdf2.Key(key, symmetricSalt, int(iterations), keySize, sha256.New)
	c, _ := aes.NewCipher(cdr)
	ctr := cipher.NewCBCDecrypter(c, iv)
	plaintext := make([]byte, len(ciphertext))
	ctr.CryptBlocks(plaintext, ciphertext)

	return unpad(plaintext), metadata, nil
}

func (a *Aes256CBC) GetHashSize() int {
	switch a.hashAlgo {
	case SHA256:
		return sha256.Size
	case SHA384:
		return sha512.Size384
	case SHA512:
		return sha512.Size
	default:
		return sha512.Size384
	}
}

func (a *Aes256CBC) NewHmac(key []byte) hash.Hash {
	switch a.hashAlgo {
	case SHA256:
		return hmac.New(sha256.New, key)
	case SHA384:
		return hmac.New(sha512.New384, key)
	case SHA512:
		return hmac.New(sha512.New, key)
	default:
		return hmac.New(sha512.New384, key)
	}
}
