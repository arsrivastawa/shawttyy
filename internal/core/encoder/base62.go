package encoder

import "strings"

const base62Alphabet = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// Base62Encode converts a non-negative int64 id into a Base62 short code.
func Base62Encode(id int64) string {
	if id == 0 {
		return "0"
	}

	base := int64(len(base62Alphabet))
	var sb strings.Builder
	for id > 0 {
		sb.WriteByte(base62Alphabet[id%base])
		id /= base
	}

	return reverse(sb.String())
}

func reverse(s string) string {
	b := []byte(s)
	for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}
	return string(b)
}
