package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"io"
	"syscall/js"
)

func arrayBufferToBase64JS() js.Func {

	wasmFunc := js.FuncOf(func(this js.Value, args []js.Value) any {
		buffer := make([]byte, 0)
		js.CopyBytesToGo(buffer, args[0])
		return base64.StdEncoding.EncodeToString(buffer)
	})

	return wasmFunc
}

func arrayBufferToBase64(buffer []byte) string {
	return base64.StdEncoding.EncodeToString(buffer)
}

func base64toBuffer(base64Key string) []byte {
	decoded, err := base64.StdEncoding.DecodeString(base64Key)
	if err != nil {
		panic(err.Error())
	}
	return decoded
}

func generateAESKey() []byte {
	aesKey := make([]byte, 32)
	_, err := rand.Read(aesKey)
	if err != nil {
		panic(err.Error())
	}
	return aesKey
}

func exportAESKey(aesKey []byte) []byte {
	// if exported ASE key is in bytes --> this function is not necessary
	return aesKey
}

var aesKey []byte

func getAesKeyJS() js.Func {
	wasmFunc := js.FuncOf(func(this js.Value, args []js.Value) any {
		return arrayBufferToBase64(aesKey)
	})
	return wasmFunc
}

func generateAndEncryptAesKeyJS() js.Func {
	wasmFunc := js.FuncOf(func(this js.Value, args []js.Value) any {
		rsaPublicKey := args[0].String()
		aesKey = generateAESKey()
		exportedAesKey := exportAESKey(aesKey)

		// Convert the exported AES key to a Base64 string for logging
		// aesKeyBase64 := arrayBufferToBase64(aesKey)
		// log here

		encryptedAESKey := encryptAESKeyWithRSA(exportedAesKey, rsaPublicKey)
		encryptedAESKeyBase64 := arrayBufferToBase64(encryptedAESKey)
		return encryptedAESKeyBase64
	})
	return wasmFunc
}

func encryptAESKeyWithRSA(exportedAESKey []byte, rsaPublicKey string) []byte {
	rsaPublicKeyBytes := base64toBuffer(rsaPublicKey)
	rsaPublicKeyImported := parseRsaPublicKey(rsaPublicKeyBytes)
	encryptedAesKey, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, rsaPublicKeyImported, exportedAESKey, nil)

	if err != nil {
		panic(err.Error())
	}
	return encryptedAesKey
}

func parseRsaPublicKey(publicKeyBytes []byte) *rsa.PublicKey {
	publicKey, err := x509.ParsePKIXPublicKey(publicKeyBytes)
	if err != nil {
		panic(err.Error())
	}
	rsaPublicKey, ok := publicKey.(*rsa.PublicKey)
	if !ok {
		panic("Invalid RSA public key format")
	}
	return rsaPublicKey
}

func encryptMessageJS() js.Func {
	wasmFunc := js.FuncOf(func(this js.Value, args []js.Value) any {
		message := args[0].String()
		key := args[1].String()

		block, err := aes.NewCipher(base64toBuffer(key))
		if err != nil {
			panic(err.Error())
		}

		// Convert the message to a byte slice
		messageBytes := []byte(message)

		// Create a new AES-GCM cipher with the block
		aesGcm, err := cipher.NewGCM(block)
		if err != nil {
			panic(err.Error())
		}

		//Create a nonce. Nonce should be from GCM
		nonce := make([]byte, aesGcm.NonceSize())
		if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
			panic(err.Error())
		}

		encryptedMessage := aesGcm.Seal(nonce, nonce, messageBytes, nil)
		return arrayBufferToBase64(encryptedMessage)
	})
	return wasmFunc
}

func decryptMessageJS() js.Func {
	wasmFunc := js.FuncOf(func(this js.Value, args []js.Value) any {
		encryptedMessageBase64 := args[0].String()
		key := args[1].String()

		// Convert the Base64 string back to a byte slice
		combined := base64toBuffer(encryptedMessageBase64)

		// Create a new Cipher Block from the key
		block, err := aes.NewCipher(base64toBuffer(key))
		if err != nil {
			panic(err.Error())
		}

		// Create a new AES-GCM cipher with the block
		aesGcm, err := cipher.NewGCM(block)
		if err != nil {
			panic(err.Error())
		}

		// Extract nonce
		nonce := combined[:aesGcm.NonceSize()]

		// Extract the encrypted message
		encryptedMessage := combined[aesGcm.NonceSize():]

		// Decrypt the message
		decryptedMessage, err := aesGcm.Open(nil, nonce, encryptedMessage, nil)
		if err != nil {
			panic(err.Error())
		}
		return string(decryptedMessage)
	})
	return wasmFunc
}

func main() {
	fmt.Println("Go Web Assembly")
	js.Global().Set("generateAndEncryptAesKeyJS", generateAndEncryptAesKeyJS())
	js.Global().Set("encryptMessageJS", encryptMessageJS())
	js.Global().Set("decryptMessageJS", decryptMessageJS())
	js.Global().Set("arrayBufferToBase64JS", arrayBufferToBase64JS())
	js.Global().Set("getAesKeyJS", getAesKeyJS())

	<-make(chan bool)
}
