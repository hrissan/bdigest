// Copyright 2020 Gregory Petrosyan <gregory.petrosyan@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package bdigest_test

import (
	"fmt"
	"math"
	"math/rand"
	"testing"

	"pgregory.net/bdigest"
)

const (
	benchElemCount = 100 * 1000
)

var (
	errors    = []float64{0.001, 0.01, 0.05}
	quantiles = []float64{0, 0.001, 0.01, 0.1, 0.25, 0.5, 0.75, 0.9, 0.95, 0.99, 0.999, 0.9999, 1}
)

func BenchmarkNewDigest(b *testing.B) {
	for _, err := range errors {
		b.Run(fmt.Sprintf("%v", err), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				bdigest.NewDigest(err)
			}
		})
	}
}

func BenchmarkDigest_Add(b *testing.B) {
	for _, err := range errors {
		b.Run(fmt.Sprintf("%v", err), func(b *testing.B) {
			r := rand.New(rand.NewSource(0))
			values := make([]float64, b.N)
			for i := 0; i < b.N; i++ {
				values[i] = math.Exp(r.NormFloat64())
			}
			d := bdigest.NewDigest(err)
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				d.Add(values[i])
			}
		})
	}
}

func BenchmarkDigest_Quantile(b *testing.B) {
	for _, err := range errors {
		b.Run(fmt.Sprintf("%v", err), func(b *testing.B) {
			d := logNormalDigest(err, 0, benchElemCount, 0)

			for _, q := range quantiles {
				b.Run(fmt.Sprintf("q%v", q), func(b *testing.B) {
					for i := 0; i < b.N; i++ {
						d.Quantile(q)
					}
				})
			}
		})
	}
}

func BenchmarkDigest_Merge(b *testing.B) {
	for _, err := range errors {
		b.Run(fmt.Sprintf("%v", err), func(b *testing.B) {
			d1 := logNormalDigest(err, 0, benchElemCount, 0)
			d2 := logNormalDigest(err, 1, benchElemCount, 0)
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				_ = d1.Merge(d2)
			}
		})
	}
}

func BenchmarkDigest_MarshalBinary(b *testing.B) {
	for _, err := range errors {
		b.Run(fmt.Sprintf("%v", err), func(b *testing.B) {
			d := logNormalDigest(err, 0, benchElemCount, 0)
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				_, err := d.MarshalBinary()
				if err != nil {
					b.Fatalf("unexpected error during marshaling: %v", err)
				}
			}
		})
	}
}

func BenchmarkDigest_UnmarshalBinary(b *testing.B) {
	for _, err := range errors {
		b.Run(fmt.Sprintf("%v", err), func(b *testing.B) {
			d1 := logNormalDigest(err, 0, benchElemCount, 0)
			d2 := &bdigest.Digest{}
			buf, err := d1.MarshalBinary()
			if err != nil {
				b.Fatalf("unexpected error during marshaling: %v", err)
			}
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				err := d2.UnmarshalBinary(buf)
				if err != nil {
					b.Fatalf("unexpected error during unmarshaling: %v", err)
				}
			}
		})
	}
}

func logNormalDigest(err float64, seed int64, count int, zeroChance int32) *bdigest.Digest {
	d := bdigest.NewDigest(err)

	r := rand.New(rand.NewSource(seed))
	for i := 0; i < count; i++ {
		if zeroChance > 0 && r.Int31n(zeroChance) == 0 {
			d.Add(0)
		} else {
			d.Add(math.Exp(r.NormFloat64()))
		}
	}

	return d
}
