package utils

import (
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"google.golang.org/protobuf/proto"
)

func GetServerStoragePathByNames(packageNames ...string) string {
	if len(packageNames) == 0 {
		return ""
	}
	return filepath.Join(packageNames...)
}

func IsFileExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

func WriteFile(filePath, content string) error {
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("failed to open or create file: %w", err)
	}
	defer file.Close()
	if _, err := io.WriteString(file, content); err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}
	return nil
}

func Contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// ReadLastNLines reads the last n lines from a file.
func ReadLastNLines(file *os.File, n int) (string, error) {
	if n <= 0 {
		return "", fmt.Errorf("invalid number of lines: %d", n)
	}

	stat, err := file.Stat()
	if err != nil {
		return "", err
	}

	fileSize := stat.Size()
	if fileSize == 0 {
		return "", nil
	}

	bufferSize := 1024
	buf := make([]byte, bufferSize)
	lines := make([]string, 0, n)
	offset := int64(0)
	lineCount := 0

	for offset < fileSize && lineCount < n {
		readSize := min(bufferSize, int(fileSize-offset))
		offset += int64(readSize)

		_, err := file.Seek(-offset, io.SeekEnd)
		if err != nil {
			return "", err
		}

		_, err = file.Read(buf[:readSize])
		if err != nil {
			return "", err
		}

		// Reverse the buffer to process lines from end to start
		for i := readSize - 1; i >= 0 && lineCount < n; i-- {
			if buf[i] == '\n' || i == 0 {
				start := i
				if buf[i] == '\n' {
					start++
				}
				line := string(buf[start:readSize])
				if line != "" || i == 0 {
					lines = append([]string{line}, lines...)
					lineCount++
					readSize = i
				}
			}
		}
	}
	return strings.Join(lines, "\n"), nil
}

func MergePath(paths ...string) string {
	pathArr := make([]string, 0)
	for _, path := range paths {
		pathArr = append(pathArr, strings.Split(path, "/")...)
	}
	return strings.Join(pathArr, "/")
}

func SerializeToBase64(msg proto.Message) (string, error) {
	data, err := proto.Marshal(msg)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(data), nil
}

func DeserializeFromBase64(data string, msg proto.Message) error {
	decoded, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return err
	}
	return proto.Unmarshal(decoded, msg)
}
