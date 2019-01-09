package crypto4go

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
)

func packageData(originalData []byte, packageSize int) (r [][]byte) {
	var src = make([]byte, len(originalData))
	copy(src, originalData)

	r = make([][]byte, 0)
	if len(src) <= packageSize {
		return append(r, src)
	}
	for len(src) > 0 {
		var p = src[:packageSize]
		r = append(r, p)
		src = src[packageSize:]
		if len(src) <= packageSize {
			r = append(r, src)
			break
		}
	}
	return r
}

func RSAEncrypt(plaintext, key []byte) ([]byte, error) {
	pub, err := ParsePKCS1PublicKey(key)
	if err != nil {
		return nil, err
	}

	var data = packageData(plaintext, pub.N.BitLen()/8-11)
	var cipherData = make([]byte, 0, 0)

	for _, d := range data {
		var c, e = rsa.EncryptPKCS1v15(rand.Reader, pub, d)
		if e != nil {
			return nil, e
		}
		cipherData = append(cipherData, c...)
	}

	return cipherData, nil
}

func RSADecrypt(ciphertext, key []byte) ([]byte, error) {
	pri, err := ParsePKCS1PrivateKey(key)
	if err != nil {
		return nil, err
	}

	var data = packageData(ciphertext, pri.PublicKey.N.BitLen()/8)
	var plainData = make([]byte, 0, 0)

	for _, d := range data {
		var p, e = rsa.DecryptPKCS1v15(rand.Reader, pri, d)
		if e != nil {
			return nil, e
		}
		plainData = append(plainData, p...)
	}
	return plainData, nil
}

func SignPKCS1v15(src, key []byte, hash crypto.Hash) ([]byte, error) {
	pri, err := ParsePKCS1PrivateKey(key)
	if err != nil {
		return nil, err
	}

	var h = hash.New()
	h.Write(src)
	var hashed = h.Sum(nil)

	return rsa.SignPKCS1v15(rand.Reader, pri, hash, hashed)
}

func VerifyPKCS1v15(src, sig, key []byte, hash crypto.Hash) error {
	pub, err := ParsePKCS1PublicKey(key)
	if err != nil {
		return err
	}

	var h = hash.New()
	h.Write(src)
	var hashed = h.Sum(nil)

	return rsa.VerifyPKCS1v15(pub, hash, hashed, sig)
}
