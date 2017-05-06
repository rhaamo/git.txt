package tool

import (
	"encoding/hex"
	"crypto/sha1"
	"crypto/md5"
	"math/big"
	"crypto/rand"
)

const alphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

// MD5Bytes encodes string to MD5 bytes.
func MD5Bytes(str string) []byte {
	m := md5.New()
	m.Write([]byte(str))
	return m.Sum(nil)
}

// MD5 encodes string to MD5 hex value.
func MD5(str string) string {
	return hex.EncodeToString(MD5Bytes(str))
}

// SHA1 encodes string to SHA1 hex value.
func SHA1(str string) string {
	h := sha1.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

// ShortSHA1 truncates SHA1 string length to at most 10.
func ShortSHA1(sha1 string) string {
	if len(sha1) > 10 {
		return sha1[:10]
	}
	return sha1
}

// RandomString returns generated random string in given length of characters.
// It also returns possible error during generation.
func RandomString(n int) (string, error) {
	buffer := make([]byte, n)
	max := big.NewInt(int64(len(alphanum)))

	for i := 0; i < n; i++ {
		index, err := randomInt(max)
		if err != nil {
			return "", err
		}

		buffer[i] = alphanum[index]
	}

	return string(buffer), nil
}

func randomInt(max *big.Int) (int, error) {
	rand, err := rand.Int(rand.Reader, max)
	if err != nil {
		return 0, err
	}

	return int(rand.Int64()), nil
}