package util

func StringHashToUint32(s string) uint32 {
	var hash uint32
	bs := []byte(s)
	for _, v := range bs {
		hash = hash*131 + uint32(v)
	}

	return hash
}
