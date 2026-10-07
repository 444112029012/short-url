package domain

import (
	"crypto/rand"
	"io"
)

// base62 alphabet (ADR-003). 62^8 ≈ 47.6 bits. crypto/rand is the only source.
const (
	shortCodeLength = 8
	alphabet        = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	// 62 * 4 = 248. Reject bytes in [248, 255] so modulo does not bias the alphabet.
	alphabetUnbiasedLimit = 248
)

// ShortCodeGenerator implements ShortCodeGenerator.generate.
type ShortCodeGenerator struct {
	reader io.Reader
}

func NewShortCodeGenerator() *ShortCodeGenerator {
	return NewShortCodeGeneratorWithReader(rand.Reader)
}

// NewShortCodeGeneratorWithReader lets tests simulate CSPRNG failure.
// Production code uses NewShortCodeGenerator (crypto/rand).
func NewShortCodeGeneratorWithReader(reader io.Reader) *ShortCodeGenerator {
	if reader == nil {
		reader = rand.Reader
	}
	return &ShortCodeGenerator{reader: reader}
}

// Generate returns a length-8 base62 code or internal_error if the reader fails.
func (g *ShortCodeGenerator) Generate() (string, *AppError) {
	if g == nil || g.reader == nil {
		return "", Internal()
	}
	out := make([]byte, shortCodeLength)
	var buf [1]byte
	filled := 0
	for spins := 0; filled < shortCodeLength && spins < shortCodeLength*64; spins++ {
		n, err := g.reader.Read(buf[:])
		if err != nil || n != 1 {
			return "", InternalWrap(err)
		}
		if int(buf[0]) >= alphabetUnbiasedLimit {
			continue
		}
		out[filled] = alphabet[int(buf[0])%len(alphabet)]
		filled++
	}
	if filled != shortCodeLength {
		return "", Internal()
	}
	return string(out), nil
}
