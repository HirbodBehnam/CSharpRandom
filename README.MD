# C# Random
![Go](https://github.com/HirbodBehnam/CSharpRandom/workflows/Go/badge.svg)

A library to use C# random number generator in Golang and create exact same number with a given seed.

## Documents
https://pkg.go.dev/github.com/HirbodBehnam/CSharpRandom

## Example
```go
package main

import (
  "fmt"
  "github.com/HirbodBehnam/CSharpRandom"
  "math"
  "time"
)

func main() {
  rng := CSharpRandom.NewRandom(int32(time.Now().Unix() % (math.MinInt32 + 1)))
  fmt.Println(rng.NextInt())
  fmt.Println(rng.NextDouble())
  number, _ := rng.NextCeiling(1000)
  fmt.Println(number)
  number, _ = rng.NextRange(-1000,1000)
  fmt.Println(number)
}
```