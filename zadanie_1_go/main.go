package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

// flipRandomBit zmienia jeden losowy bit w danych, ale nie w IV
func flipRandomBit(data []byte, blockSize int) ([]byte, error) {
	// Upewnij się, że dane są wystarczająco długie, aby można było zmieniać bity poza IV
	if len(data) <= blockSize {
		return nil, fmt.Errorf("data is too short to flip a bit without affecting the IV")
	}

	// Obliczamy maksymalny możliwy indeks bitu (poza IV)
	maxBitIndex := (len(data) - blockSize) * 8

	// Generujemy losowy bit w zakresie od 0 do maxBitIndex
	randomByte := make([]byte, 1)
	if _, err := rand.Read(randomByte); err != nil {
		return nil, fmt.Errorf("error generating random byte: %v", err)
	}

	// Losujemy indeks bitu w danych (poza IV)
	bitIndex := int(randomByte[0])%maxBitIndex + blockSize*8

	// Obliczamy indeks bajtu i pozycję bitu w tym bajcie
	byteIndex := bitIndex / 8
	bitPos := bitIndex % 8

	// Zmieniamy wybrany bit
	data[byteIndex] ^= 1 << bitPos

	return data, nil
}

func encryptECB(plainText, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	cipherText := make([]byte, len(plainText))
	for bs, be := 0, block.BlockSize(); bs < len(plainText); bs, be = bs+block.BlockSize(), be+block.BlockSize() {
		block.Encrypt(cipherText[bs:be], plainText[bs:be])
	}

	return cipherText, nil
}

func decryptECB(cipherText, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	plainText := make([]byte, len(cipherText))
	for bs, be := 0, block.BlockSize(); bs < len(cipherText); bs, be = bs+block.BlockSize(), be+block.BlockSize() {
		block.Decrypt(plainText[bs:be], cipherText[bs:be])
	}

	return plainText, nil
}

func encryptCBC(plainText, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	padding := block.BlockSize() - len(plainText)%block.BlockSize()
	padText := append(plainText, make([]byte, padding)...)

	cipherText := make([]byte, aes.BlockSize+len(padText))
	iv := cipherText[:aes.BlockSize]

	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(cipherText[aes.BlockSize:], padText) // skipping first block because it's IV

	return cipherText, nil
}

func decryptCBC(cipherText, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	if len(cipherText) < aes.BlockSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	if len(cipherText)%aes.BlockSize != 0 {
		return nil, fmt.Errorf("ciphertext is not a multiple of the block size")
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(cipherText, cipherText)

	return cipherText, nil
}

func encryptCTR(plainText, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	cipherText := make([]byte, aes.BlockSize+len(plainText))
	iv := cipherText[:aes.BlockSize]

	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], plainText)

	return cipherText, nil
}

func decryptCTR(cipherText, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	if len(cipherText) < aes.BlockSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(cipherText, cipherText)

	return cipherText, nil
}

func processFile(fileName string, key []byte) {
	plainText, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	// Encrypt and decrypt using ECB
	ecbCipherText, err := encryptECB(plainText, key)
	if err != nil {
		fmt.Println("Error encrypting ECB:", err)
		return
	}
	os.WriteFile(fileName+"_ecb_enc.txt", ecbCipherText, 0644)

	ecbPlainText, err := decryptECB(ecbCipherText, key)
	if err != nil {
		fmt.Println("Error decrypting ECB:", err)
		return
	}
	os.WriteFile(fileName+"_ecb_dec.txt", ecbPlainText, 0644)

	// Flip a random bit in the ciphertext
	ecbCipherText, err = flipRandomBit(ecbCipherText, aes.BlockSize)
	if err != nil {
		fmt.Println("Error flipping random bit:", err)
		return
	}

	ecbPlainText, err = decryptECB(ecbCipherText, key)
	if err != nil {
		fmt.Println("Error decrypting ECB after flipping a bit:", err)
		return
	}

	// Encrypt and decrypt using CBC
	cbcCipherText, err := encryptCBC(plainText, key)
	if err != nil {
		fmt.Println("Error encrypting CBC:", err)
		return
	}
	os.WriteFile(fileName+"_cbc_enc.txt", cbcCipherText, 0644)

	cbcPlainText, err := decryptCBC(cbcCipherText, key)
	if err != nil {
		fmt.Println("Error decrypting CBC:", err)
		return
	}
	os.WriteFile(fileName+"_cbc_dec.txt", cbcPlainText, 0644)

	// Flip a random bit in the ciphertext
	cbcCipherText, err = flipRandomBit(cbcCipherText, aes.BlockSize)
	if err != nil {
		fmt.Println("Error flipping random bit:", err)
		return
	}

	cbcPlainText, err = decryptCBC(cbcCipherText, key)
	if err != nil {
		fmt.Println("Error decrypting CBC after flipping a bit:", err)
		return
	}
	os.WriteFile(fileName+"_cbc_dec_flipped.txt", cbcPlainText, 0644)

	// Encrypt and decrypt using CTR
	ctrCipherText, err := encryptCTR(plainText, key)
	if err != nil {
		fmt.Println("Error encrypting CTR:", err)
		return
	}
	os.WriteFile(fileName+"_ctr_enc.txt", ctrCipherText, 0644)

	ctrPlainText, err := decryptCTR(ctrCipherText, key)
	if err != nil {
		fmt.Println("Error decrypting CTR:", err)
		return
	}
	os.WriteFile(fileName+"_ctr_dec.txt", ctrPlainText, 0644)

	// Flip a random bit in the ciphertext
	ctrCipherText, err = flipRandomBit(ctrCipherText, aes.BlockSize)
	if err != nil {
		fmt.Println("Error flipping random bit:", err)
		return
	}

	ctrPlainText, err = decryptCTR(ctrCipherText, key)
	if err != nil {
		fmt.Println("Error decrypting CTR after flipping a bit:", err)
		return
	}
	os.WriteFile(fileName+"_ctr_dec_flipped.txt", ctrPlainText, 0644)

}

func main() {
	key := []byte("examplekey123456") // 16 bytes key for AES-128

	files := []string{"0_5.txt", "1.txt", "5.txt"}
	for _, file := range files {
		processFile(file, key)
	}
}
