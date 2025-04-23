package controllers

// import (
// 	"math"
// 	"sort"

// 	"github.com/goccmack/godsp"
// )

// type Transform struct {
// 	st       []float64
// 	level    int
// 	sections []*transformSection
// }

// type transformSection struct {
// 	start int
// 	size  int
// }

// // Daubechies4 returns the DWT with Daubechies 4 coeficients to level.
// func Daubechies4(s []float64, level int) *Transform {
// 	t := &Transform{
// 		st:       make([]float64, len(s)),
// 		level:    level,
// 		sections: getTransformSections(len(s), level),
// 	}
// 	copy(t.st, s)
// 	for _, section := range t.sections {
// 		scaleSize := section.size
// 		for l := level; l > 0; l-- {
// 			max := section.start + scaleSize
// 			split(t.st[section.start:max])
// 			daubechies4(t.st[section.start:max])
// 			scaleSize /= 2
// 		}
// 	}

// 	return t
// }

// /*
// Return the series of length 2^k stages of the DWT
// */
// func getTransformSections(N, level int) (sections []*transformSection) {
// 	for start := 0; N-start >= 64*godsp.Pow2(level); {
// 		size := int(godsp.Pow2(godsp.Log2(N - start)))
// 		section := &transformSection{
// 			start: start,
// 			size:  size,
// 		}
// 		sections = append(sections, section)
// 		start += size
// 	}
// 	return
// }

// /*
// GetFrameSize returns the size of DWT frame required for the transform
// */
// // func GetFrameSize(s []float64) int {
// // 	logLen := math.Log2(float64(len(s)))
// // 	logLenInt := int(math.Ceil(logLen))

// // 	return godsp.Pow2(logLenInt)
// // }

// /*
// Split s into even and odd elements,
// where the even elements are in the first half
// of the vector and the odd elements are in the
// second half.
// */
// func split(s []float64) {
// 	half := len(s) / 2
// 	odd := make([]float64, half)
// 	for i := 1; i < len(s); i += 2 {
// 		odd[i/2] = s[i]
// 	}
// 	for i := 2; i < len(s); i += 2 {
// 		s[i/2] = s[i]
// 	}
// 	for i, v := range odd {
// 		s[half+i] = v
// 	}
// }

// /*
// After: Ripples section 3.4
// */
// func daubechies4(s []float64) {
// 	half := len(s) / 2

// 	// Update 1:
// 	for n := 0; n < half; n++ {
// 		s[n] = s[n] + math.Sqrt(3)*s[half+n]
// 	}

// 	// Predict:
// 	s[half] = s[half] -
// 		(math.Sqrt(3)/4)*s[0] -
// 		((math.Sqrt(3)-2)/4)*s[half-1]
// 	for n := 1; n < half; n++ {
// 		s[half+n] = s[half+n] -
// 			(math.Sqrt(3)/4)*s[n] -
// 			((math.Sqrt(3)-2)/4)*s[n-1]
// 	}

// 	// Update 2:
// 	for n := 0; n < half-1; n++ {
// 		s[n] = s[n] - s[half+n+1]
// 	}
// 	s[half-1] = s[half-1] - s[half]

// 	// Normalise:
// 	for n := 0; n < half; n++ {
// 		s[n] = ((math.Sqrt(3) - 1) / math.Sqrt(2)) * s[n]
// 		s[n+half] = ((math.Sqrt(3) + 1) / math.Sqrt(2)) * s[n+half]
// 	}
// }

// // GetCoefficients returns the coefficients of all transform levels
// func (t *Transform) GetCoefficients() [][]float64 {
// 	cfs := make([][]float64, t.level)
// 	for _, s := range t.sections {
// 		scfs := t.getSectionCoefficients(s)
// 		for i, c := range scfs {
// 			cfs[i] = append(cfs[i], c...)
// 		}
// 	}
// 	return cfs
// }

// // GetDownSampledCoefficients returns the coefficients of all the levels downsampled to
// // the length of the deepest level of the transform.
// func (t *Transform) GetDownSampledCoefficients() [][]float64 {
// 	dscfs := make([][]float64, t.level)
// 	for _, s := range t.sections {
// 		cfs := t.getSectionCoefficients(s)
// 		for i, cf := range cfs {
// 			if i < t.level-1 {
// 				dscfs[i] = append(dscfs[i],
// 					godsp.DownSample(cf, godsp.Pow2(t.level-(i+1)))...)
// 			}
// 		}
// 	}
// 	return dscfs
// }

// /*
// GetDecomposition returns the vector containing the DWT decomposion
// */
// func (t *Transform) GetDecomposition() []float64 {
// 	return t.st
// }

// // GetCoefficients returns the coefficients of all transform levels
// func (t *Transform) getSectionCoefficients(s *transformSection) [][]float64 {
// 	cfs := make([][]float64, t.level)
// 	half := s.size / 2
// 	for l := 1; l <= t.level; l++ {
// 		cfs[l-1] = t.st[s.start+half : s.start+2*half]
// 		half /= 2
// 	}
// 	return cfs
// }
// func Denoise(s []float64, level int) []float64 {
//     // 1. Perform DWT decomposition
//     t := Daubechies4(s, level)

//     // 2. Access coefficients
//     coeffs := t.GetCoefficients()

//     // 3. Calculate threshold from finest detail coefficients (cD1)
//     if len(coeffs) == 0 || len(coeffs[0]) == 0 {
//         return s
//     }
//     cD1 := coeffs[0]
//     threshold := calculateThreshold(cD1)

//     // 4. Apply threshold to detail coefficients (keep last 2 levels)
//     for i := 0; i < len(coeffs)-2; i++ {
//         softThreshold(coeffs[i], threshold)
//     }

//     // 5. Reconstruct the signal
//     return reconstruct(t, level)
// }

// // Helper functions
// func calculateThreshold(cD1 []float64) float64 {
//     mad := medianAbs(cD1)
//     sigma := mad / 0.6745
//     return sigma * math.Sqrt(2*math.Log(float64(len(cD1))))
// }

// func medianAbs(data []float64) float64 {
//     abs := make([]float64, len(data))
//     for i, v := range data {
//         abs[i] = math.Abs(v)
//     }
//     sort.Float64s(abs)

//     n := len(abs)
//     if n%2 == 0 {
//         return (abs[n/2-1] + abs[n/2]) / 2
//     }
//     return abs[n/2]
// }

// func softThreshold(coeffs []float64, threshold float64) {
//     for i := range coeffs {
//         val := coeffs[i]
//         absVal := math.Abs(val)
//         if absVal <= threshold {
//             coeffs[i] = 0.0
//         } else {
//             coeffs[i] = math.Copysign(absVal-threshold, val)
//         }
//     }
// }

// // Reconstruction implementation
// func reconstruct(t *Transform, level int) []float64 {
//     // Create working copy of the transformed data
//     reconstructed := make([]float64, len(t.st))
//     copy(reconstructed, t.st)

//     // Reverse the transformation process
//     for _, section := range t.sections {
//         scaleSize := section.size
//         for l := 0; l < level; l++ {
//             max := section.start + scaleSize
//             inverseDaubechies4(reconstructed[section.start:max])
//             unSplit(reconstructed[section.start:max])
//             scaleSize *= 2
//         }
//     }

//     return reconstructed
// }

// func inverseDaubechies4(s []float64) {
//     half := len(s) / 2

//     // Reverse normalization
//     a := (math.Sqrt(3) - 1) / math.Sqrt(2)
//     b := (math.Sqrt(3) + 1) / math.Sqrt(2)
//     for n := 0; n < half; n++ {
//         s[n] /= a
//         s[half+n] /= b
//     }

//     // Reverse update 2
//     s[half-1] += s[half]
//     for n := half - 2; n >= 0; n-- {
//         s[n] += s[half+n+1]
//     }

//     // Reverse predict step
//     s[half] += (math.Sqrt(3)/4)*s[0] + ((math.Sqrt(3)-2)/4)*s[half-1]
//     for n := 1; n < half; n++ {
//         s[half+n] += (math.Sqrt(3)/4)*s[n] + ((math.Sqrt(3)-2)/4)*s[n-1]
//     }

//     // Reverse update 1
//     for n := 0; n < half; n++ {
//         s[n] -= math.Sqrt(3) * s[half+n]
//     }
// }

// func unSplit(s []float64) {
//     half := len(s) / 2
//     even := make([]float64, half)
//     odd := make([]float64, half)

//     copy(even, s[:half])
//     copy(odd, s[half:])

//     for i := 0; i < half; i++ {
//         s[2*i] = even[i]
//         s[2*i+1] = odd[i]
//     }
// }
