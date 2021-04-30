package crypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	crypto_rand "crypto/rand"
	"encoding/hex"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"io"
	"io/ioutil"
	math_rand "math/rand"
	"os"
	"time"
)

// ComputeMD5 вычисляет контрольную сумму MD5 файла по указанному пути.
func ComputeMD5(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := md5.New()
	if _, err = io.Copy(h, f); err != nil {
		return "", err
	}
	return fmt.Sprintf("%X", h.Sum(nil)), nil
}

// CheckPassword сравнивает пароль и хэш. Если хэш явно не указан, то он читается из файла secret в текущем каталоге.
// Возвращает true, если хэш соответствует паролю.
func CheckPassword(password string, hash ...string) (bool, error) {
	var _hash []byte
	var err error

	if len(hash) == 0 {
		_hash, err = secretHashRead()
		if err != nil {
			return false, err
		}
	} else {
		_hash = []byte(hash[0])
	}
	// Comparing the password with the hash
	err = bcrypt.CompareHashAndPassword(_hash, []byte(password))

	if err != nil {
		return false, err
	}

	return true, nil
}

// GeneratePasswordAndStoreNewHash вызывает passPhraseGen и записывает новый хэш пароля на диск и возвращает сгенерированный пароль
// из простых символов ascii с некоторым исключением для удобства набора на клавиатуре.
func GeneratePasswordAndStoreNewHash() ([]byte, error) {
	phrase := passPhraseGen(12)
	hash, err := bcrypt.GenerateFromPassword(phrase, 10)

	if err != nil {
		return nil, err
	}

	err = secretHashWrite(&hash)

	if err != nil {
		return nil, err
	}

	return phrase, nil
}

func passPhraseGen(iPhraseLength int) []byte {
	var bytePassPhrase []byte

	math_rand.Seed(time.Now().UnixNano())

	for i := 1; i <= iPhraseLength; i++ {
	Gen:
		uiRandNumber := uint8(math_rand.Intn(125-33) + 33)
		switch uiRandNumber {
		// убираются некоторые неудобные символы
		case 34, 39, 44, 45, 46, 94, 96, 124:
			// todo сделать защиту от бесконечного цикла
			goto Gen
		default:
			bytePassPhrase = append(bytePassPhrase, uiRandNumber)
		}
	}

	return bytePassPhrase
}

func secretHashRead() ([]byte, error) {
	secret, errRead := ioutil.ReadFile("secret")

	if errRead != nil {
		return nil, errRead
	}

	return secret, nil
}

// secretHashWrite записывает хэш в файл secret текущего каталога.
func secretHashWrite(secret *[]byte) error {
	errWrite := ioutil.WriteFile("secret", *secret, 0600)

	return errWrite
}

func Encrypt(stringToEncrypt string, keyString string) (encryptedString string) {

	//Since the key is in string, we need to convert decode it to bytes
	key, _ := hex.DecodeString(keyString)
	plaintext := []byte(stringToEncrypt)

	//Create a new Cipher Block from the key
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}

	//Create a new GCM - https://en.wikipedia.org/wiki/Galois/Counter_Mode
	//https://golang.org/pkg/crypto/cipher/#NewGCM
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}

	//Create a nonce. Nonce should be from GCM
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(crypto_rand.Reader, nonce); err != nil {
		panic(err.Error())
	}

	//Encrypt the data using aesGCM.Seal
	//Since we don't want to save the nonce somewhere else in this case, we add it as a prefix to the encrypted data. The first nonce argument in Seal is the prefix.
	ciphertext := aesGCM.Seal(nonce, nonce, plaintext, nil)
	return fmt.Sprintf("%x", ciphertext)
}

func Decrypt(encryptedString string, keyString string) (decryptedString string) {

	key, _ := hex.DecodeString(keyString)
	enc, _ := hex.DecodeString(encryptedString)

	//Create a new Cipher Block from the key
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}

	//Create a new GCM
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}

	//Get the nonce size
	nonceSize := aesGCM.NonceSize()

	//Extract the nonce from the encrypted data
	nonce, ciphertext := enc[:nonceSize], enc[nonceSize:]

	//Decrypt the data
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		panic(err.Error())
	}

	return fmt.Sprintf("%s", plaintext)
}

func main() {

	bytes := make([]byte, 32) //generate a random 32 byte key for AES-256
	if _, err := crypto_rand.Read(bytes); err != nil {
		panic(err.Error())
	}

	key := hex.EncodeToString(bytes) //encode key in bytes to string and keep as secret, put in a vault
	fmt.Printf("key to Encrypt/Decrypt : %s\n", key)

	encrypted := Encrypt("Hello Encrypt", key)
	fmt.Printf("encrypted : %s\n", encrypted)

	decrypted := Decrypt(encrypted, key)
	fmt.Printf("decrypted : %s\n", decrypted)
}