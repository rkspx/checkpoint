package checkpoint

import (
	"crypto/rand"
	"math"
	"math/big"
)

func randInt64(max int64) (int64, error) {
	if max == 0 {
		return 0, nil
	}
	if max <= 0 {
		max = int64(math.Abs(float64(max)))
	}
	bign, err := rand.Int(rand.Reader, big.NewInt(max))
	if err != nil {
		return 0, err
	}

	return bign.Int64(), nil
}
