package shopify

import (
	"math/rand"
	"time"
)

const (
	alphaNumerics = "abcdefghijklmnopqrstuvwxyz0123456789"
)

func randomAlphaNumerics(n int) string {
	s := ""
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < n; i++ {
		anidx := r.Intn(len(alphaNumerics))
		s += string(alphaNumerics[anidx])
	}
	return s
}

func (a *App) GenerateNonce() (string, error) {
	return a.nonceGenerate()
}

func DefaultNonceGenerate() (string, error) {
	return randomAlphaNumerics(24), nil
}
