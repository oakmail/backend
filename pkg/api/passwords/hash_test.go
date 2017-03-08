package passwords_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/dchest/uniuri"
	. "github.com/smartystreets/goconvey/convey"

	"github.com/oakmail/backend/pkg/api/passwords"
)

const count = 25

func TestHash(t *testing.T) {
	Convey("Hash should succeed in hashing them", t, func() {
		timings := []time.Duration{}

		for i := 0; i < count; i++ {
			start := time.Now()

			hash := passwords.Hash(uniuri.NewLen(128))
			timings = append(timings, time.Now().Sub(start))

			So(len(hash), ShouldBeGreaterThan, 0)
		}

		var sum time.Duration = 0
		for _, v := range timings {
			sum += v
		}

		avg := float64(sum) / float64(len(timings))

		fmt.Printf("avg %.2fms", avg/float64(time.Millisecond))
	})

	Convey("Verify should succeed in verifying them", t, func() {
		timings := []time.Duration{}

		for i := 0; i < count; i++ {
			plain := uniuri.NewLen(128)
			hash := passwords.Hash(plain)

			start := time.Now()

			res := passwords.Verify(hash, plain)
			timings = append(timings, time.Now().Sub(start))

			So(len(hash), ShouldBeGreaterThan, 0)
			So(res, ShouldBeTrue)
		}

		var sum time.Duration = 0
		for _, v := range timings {
			sum += v
		}

		avg := float64(sum) / float64(len(timings))

		fmt.Printf("avg %.2fms", avg/float64(time.Millisecond))
	})
}
