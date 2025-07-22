package utils

import (
	"3cognito/coderunner/types"
	"crypto/sha256"
	"fmt"
)

func HashData(input types.FileData) string {
	data := []byte(input.Language + ":" + input.Content)
	hash := sha256.Sum256(data)
	return fmt.Sprintf("%x", hash)
}
