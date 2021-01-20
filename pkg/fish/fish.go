package fish

import (
	"bytes"
	"fmt"
	"golang.org/x/crypto/blowfish"
	"strings"
)

func Pad(src []byte, mod int) []byte {
	remainder := len(src) % mod
	if remainder != 0 {
		return append(src, bytes.Repeat([]byte{0}, mod-remainder)...)
	}
	return src
}

func BlowFishEncrypt(key string, src []byte) ([]byte, error) {
	cipher, err := blowfish.NewCipher([]byte(key))
	if err != nil {
		return nil, err
	}
	const bs = blowfish.BlockSize
	dst := make([]byte, len(src))
	for i := 0; i < len(src); i += bs {
		cipher.Encrypt(dst[i:i+bs], src[i:i+bs])
	}

	return dst, nil

}

func BlowFishDecrypt(key string, src []byte) ([]byte, error) {
	cipher, err := blowfish.NewCipher([]byte(key))
	if err != nil {
		return nil, err
	}
	const bs = blowfish.BlockSize
	dst := make([]byte, len(src))
	for i := 0; i < len(src); i += bs {
		cipher.Decrypt(dst[i:i+bs], src[i:i+bs])
	}
	return bytes.TrimRight(dst, "\x00"), nil

}

func Base64Encode(src []byte) string {
	charset := "./0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	//确保 src 为8的倍数
	src = Pad(src, 8)
	buf := make([]byte, len(src)/2*3)
	left := 0
	right := 0

	for j, k := 0, 0; k < len(src); {
		for i := 24; i >= 0; i, k = i-8, k+1 {
			left += int(src[k]) << uint8(i)
		}
		for i := 24; i >= 0; i, k = i-8, k+1 {
			right += int(src[k]) << uint8(i)
		}
		for i := 0; i < 6; i, j = i+1, j+1 {
			buf[j] = charset[right&0x3F]
			right >>= 6
		}
		for i := 0; i < 6; i, j = i+1, j+1 {
			buf[j] = charset[left&0x3F]
			left >>= 6
		}
	}

	return string(buf)

}

func Base64Decode(src []byte) ([]byte, error) {

	if len(src) > 0 && len(src) < 12 {
		return nil, fmt.Errorf("invalid base64 input: %s", src)
	}

	charset := []byte("./0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	buf := make([]byte, len(src)/3*2)

	for j, k := 0, 0; k < len(src); {
		left := 0
		right := 0

		for i := uint8(0); i < 6; i, k = i+1, k+1 {
			v := bytes.IndexByte(charset, src[k])
			right |= v << (i * 6)
		}
		for i := uint8(0); i < 6; i, k = i+1, k+1 {
			v := bytes.IndexByte(charset, src[k])
			left |= v << (i * 6)
		}

		for i := uint8(0); i < 4; i, j = i+1, j+1 {
			w := left & (0xFF << ((3 - i) * 8))
			z := w >> ((3 - i) * 8)
			buf[j] = byte(z)
		}

		for i := uint8(0); i < 4; i, j = i+1, j+1 {
			w := right & (0xFF << ((3 - i) * 8))
			z := w >> ((3 - i) * 8)
			buf[j] = byte(z)
		}

		//for i := uint8(0); i < 6; i, k = i+1, k+1 {
		//	v := bytes.IndexByte(charset, src[k])
		//	right |= v << (i * 6)
		//}
		//for i := uint8(0); i < 6; i, k = i+1, k+1 {
		//	v := bytes.IndexByte(charset, src[k])
		//	left |= v << (i * 6)
		//}
		//for i := uint8(0); i < 4; i, j = i+1, j+1 {
		//	w := left & (0xFF << ((3 - i) * 8))
		//	z := w >> ((3 - i) * 8)
		//	buf[j] = byte(z)
		//}
		//for i := uint8(0); i < 4; i, j = i+1, j+1 {
		//	w := right & (0xFF << ((3 - i) * 8))
		//	z := w >> ((3 - i) * 8)
		//	buf[j] = byte(z)
		//}

		//for i := uint8(0); i < 6; i, k = i+1, k+1 {
		//	v := bytes.IndexByte(charset, src[k])
		//	left |= v << (i * 6)
		//}
		//for i := uint8(0); i < 4; i, j = i+1, j+1 {
		//	w := left & (0xFF << ((3 - i) * 8))
		//	z := w >> ((3 - i) * 8)
		//	buf[j] = byte(z)
		//}
		//for i := uint8(0); i < 4; i, j = i+1, j+1 {
		//	w := right & (0xFF << ((3 - i) * 8))
		//	z := w >> ((3 - 1) * 8)
		//	buf[j] = byte(z)
		//}
	}
	return buf, nil
}

func Encrypt(key string, message string) (string, error) {
	encr, err := BlowFishEncrypt(key, Pad([]byte(message), 8))
	if err != nil {
		return "", err
	}
	return "+OK " + Base64Encode(encr), nil
}

func IsEncrypted(s string) bool {
	return strings.HasPrefix(s, "+OK ") || strings.HasPrefix(s, "mcps ")
}

func Decrypt(key string, message string) (string, error) {

	if strings.HasPrefix(message, "+OK ") {
		message = message[4:]
	} else if strings.HasPrefix(message, "mcps ") {
		message = message[5:]
	} else {
		return message, nil
	}
	b, err := Base64Decode([]byte(message))
	if err != nil {
		return "", err
	}
	dec, err := BlowFishDecrypt(key, b)
	if err != nil {
		return "", err
	}
	return string(dec), nil
}
