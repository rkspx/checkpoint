package checkpoint

import (
	"math"
	"testing"
)

func TestRandomInt64(t *testing.T) {
	for i := -100; i <= 100; i++ {
		for j := 0; j <= 100; j++ {
			n, err := randInt64(int64(i))
			if err != nil {
				t.Errorf("failed to generate random number for max %d: %s", i, err.Error())
				continue
			}

			if n < 0 {
				t.Errorf("result number should not less than 0, got : %d", n)
				continue
			}

			if n > int64(math.Abs(float64(i))) {
				t.Errorf("result number should not more than limit %d, got %d", i, n)
				continue
			}

			if i == 0 && n != 0 {
				t.Errorf("result of max 0 should be 0, got %d", n)
				continue
			}
		}
	}
}
