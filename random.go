// A library to mimic C# random number generator. The code is scrapped from https://referencesource.microsoft.com/#mscorlib/system/random.cs
package CSharpRandom

import (
	"errors"
	"math"
)

// The random struct. Do not manually edit these values.
type CSRandom struct{
	inext     int32
	inextp    int32
	seedArray [56]int32
}

// The constant used in C# source code
const MSEED = 161803398
const MBIG = math.MaxInt32

// Creates a new random struct with the given seed
func NewRandom(seed int32) *CSRandom {
	rng := &CSRandom{}
	var ii,mj,mk,subtraction int32
	//Initialize our Seed array.
	//This algorithm comes from Numerical Recipes in C (2nd Ed.)
	if seed == math.MinInt32{
		subtraction = math.MaxInt32
	}else{
		subtraction = seed
		if subtraction < 0{
			subtraction = -subtraction
		}
	}
	mj = MSEED - subtraction
	rng.seedArray[55] = mj
	mk=1
	for i := int32(1); i < 55; i++ {  //Apparently the range [1..55] is special (Knuth) and so we're wasting the 0'th position.
		ii = (21 * i) % 55
		rng.seedArray[ii] = mk
		mk = mj - mk
		if mk < 0 {
			mk += MBIG
		}
		mj = rng.seedArray[ii]
	}
	for k := 1; k < 5; k++ {
		for i := 1; i < 56; i++ {
			rng.seedArray[i] -= rng.seedArray[1 + (i + 30) % 55]
			if rng.seedArray[i] < 0{
				rng.seedArray[i] += MBIG
			}
		}
	}
	rng.inext = 0
	rng.inextp = 21
	//seed = 1;
	return rng
}

// Return a new random number [0..1) and reSeed the Seed array.
func (rng *CSRandom) sample() float64 {
	//Including this division at the end gives us significantly improved
	//random number distribution.
	return float64(rng.internalSample()) * (1.0/MBIG)
}

func (rng *CSRandom) internalSample() int32 {
	var retVal int32
	locINext := rng.inext
	locINextp := rng.inextp
	locINext++
	if locINext >= 56{
		locINext=1
	}
	locINextp++
	if locINextp >= 56{
		locINextp = 1
	}
	retVal = rng.seedArray[locINext] - rng.seedArray[locINextp]
	if retVal == MBIG {
		retVal--
	}
	if retVal<0 {
		retVal+=MBIG
	}
	rng.seedArray[locINext] = retVal
	rng.inext = locINext
	rng.inextp = locINextp
	return retVal
}

// Returns an int [0..math.MaxInt32)
func (rng *CSRandom) NextInt() int32 {
	return rng.internalSample()
}

func (rng *CSRandom) getSampleForLargeRange() float64 {
	// The distribution of double value returned by Sample
	// is not distributed well enough for a large range.
	// If we use Sample for a range [Int32.MinValue..Int32.MaxValue)
	// We will end up getting even numbers only.
	result := rng.internalSample()
	// Note we can't use addition here. The distribution will be bad if we do that.
	negative := rng.internalSample() % 2 == 0 // decide the sign based on second sample
	if negative {
		result = -result
	}
	d := float64(result)
	d += math.MaxInt32 - 1 // get a number in range [0 .. 2 * Int32MaxValue - 1)
	d /= float64(2 * uint32(math.MaxInt32) - 1)
	return d
}

// Returns an int [minvalue..maxvalue)
func (rng *CSRandom) NextRange(minValue,maxValue int32) (int32, error) {
	if minValue > maxValue{
		return 0, errors.New("minValue cannot be bigger than maxValue")
	}
	r := int64(maxValue) - int64(minValue)
	if r <= math.MaxInt32 {
		return int32(rng.sample() * float64(r)) + minValue, nil
	} else {
		return int32(int64(rng.getSampleForLargeRange() * float64(r)) + int64(minValue)), nil
	}
}

// Returns an int [0..maxValue)
func (rng *CSRandom) NextCeiling(maxValue int32) (int32, error) {
	if maxValue < 0{
		return 0, errors.New("max value must be more than 0")
	}
	return int32(rng.sample() * float64(maxValue)), nil
}

// Returns a double [0..1)
func (rng *CSRandom) NextDouble() float64 {
	return rng.sample()
}

// Fills the byte array with random bytes [0..255]. The entire array is filled.
func (rng *CSRandom) NextBytes(buffer []byte) error {
	if buffer == nil{
		return errors.New("buffer is nil")
	}
	for i := range buffer{
		buffer[i] = byte(rng.internalSample() % (math.MaxUint8 + 1))
	}
	return nil
}