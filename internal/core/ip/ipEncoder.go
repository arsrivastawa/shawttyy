package ip

import (
	"crypto/md5"
	"encoding/base64"
)

func IPEncoder(ip string) string {
	md5hash := md5.Sum([]byte(ip))
	md5ToBase64 := base64.RawURLEncoding.EncodeToString(md5hash[:])
	return md5ToBase64
}
