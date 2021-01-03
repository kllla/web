package id

import (
	"crypto/md5"
	"encoding/hex"
	"time"
)

func GetID(entropy string) string {
	hash := md5.Sum([]byte(entropy + time.Now().String()))
	shrtnID := hex.EncodeToString(hash[:3])
	return shrtnID
}
