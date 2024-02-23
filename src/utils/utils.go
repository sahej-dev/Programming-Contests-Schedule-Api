package utils

import (
	"io"
	"math/rand"
	"os"
	"time"
)

func ExecuteEvery(d time.Duration, f func(tick time.Time)) *time.Ticker {
	ticker := time.NewTicker(d)

	go func() {
		for tick := range ticker.C {
			f(tick)
		}
	}()

	return ticker
}

func RandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	result := make([]byte, length)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}

	return string(result)
}

func DoesFileExists(filepath string) bool {
	_, err := os.Stat(filepath)

	return err == nil
}

func EnsureDirExists(dirPath string) error {
	_, err := os.Stat(dirPath)

	if err == nil {
		return nil
	}

	if os.IsNotExist(err) {
		err = os.MkdirAll(dirPath, 0755)
	}

	return err
}

func CopyFile(sourcePath string, destinationPath string) error {
	sFile, err := os.Open(sourcePath)
	if err != nil {
		return err
	}
	defer sFile.Close()

	dFile, err := os.Create(destinationPath)
	if err != nil {
		return err
	}
	defer dFile.Close()

	_, err = io.Copy(dFile, sFile)

	return err
}
