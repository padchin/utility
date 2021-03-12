package crypt

import (
	"crypto/md5"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"time"
)

// ComputeMD5 вычисляет контрольную сумму MD5 файла по указанному пути.
func ComputeMD5(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		log.Printf("невозможно открыть файл: %v", err)
		return "", err
	}
	defer f.Close()
	h := md5.New()
	if _, err = io.Copy(h, f); err != nil {
		log.Printf("невозможно вычислить MD5: %v", err)
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
	passPhrase := passPhraseGen(12)
	byteHash, errHash := bcrypt.GenerateFromPassword(passPhrase, 10)
	if errHash != nil {
		return nil, errHash
	}
	errWrite := secretHashWrite(&byteHash)
	if errWrite != nil {
		return nil, errWrite
	} else {
		return passPhrase, nil
	}
}

func passPhraseGen(iPhraseLength int) []byte {
	var bytePassPhrase []byte
	rand.Seed(time.Now().UnixNano())
	for i := 1; i <= iPhraseLength; i++ {
	Gen:
		uiRandNumber := uint8(rand.Intn(125-33) + 33)
		switch uiRandNumber {
		// убираются некоторые неудобные символы
		case 34, 39, 44, 45, 46, 94, 96, 124:
			//todo сделать защиту от бесконечного цикла
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
	errWrite := ioutil.WriteFile("secret", *secret, 0777)
	return errWrite
}
