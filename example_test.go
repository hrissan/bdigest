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

	"pgregory.net/bdigest"
)

func ExampleNewDigest() {
	r := rand.New(rand.NewSource(0))
	d := bdigest.NewDigest(0.01)

	for i := 0; i < 100000; i++ {
		v := math.Exp(r.NormFloat64())
		d.Add(v)
	}

	fmt.Printf("%v buckets\n", d.Size())
	for _, q := range []float64{0, 0.1, 0.25, 0.5, 0.75, 0.9, 0.95, 0.99, 0.999, 0.9999, 1} {
		fmt.Printf("%v\tq%v\n", d.Quantile(q), q)
	}

	// Output:
	// 486 buckets
	// 0.016404706383449263	q0
	// 0.2808056914403475	q0.1
	// 0.5116715635730887	q0.25
	// 1.01	q0.5
	// 1.993661701417351	q0.75
	// 3.6327611266266073	q0.9
	// 5.207005866310038	q0.95
	// 10.278225915562174	q0.99
	// 21.978242872649467	q0.999
	// 42.5242714134434	q0.9999
	// 267.7721258158505	q1
}
