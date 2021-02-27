package utility

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"os"
)

// ComputeMD5 вычисляет контрольную сумму MD5 файла
func ComputeMD5(s_name string) (string, error) {
	f, err := os.Open(s_name)
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
