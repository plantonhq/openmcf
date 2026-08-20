package capacityestimator

import (
	"fmt"
	"math/big"
	"strings"
)

// The Kubernetes quantity suffixes, exact. Binary suffixes are powers of
// 1024, decimal suffixes powers of 1000; CPU's "m" is a millicore. The
// engine never floats a quantity: every parse multiplies exactly and
// refuses anything that does not land on a whole millicore or byte.
var (
	binarySuffixes = []struct {
		suffix     string
		multiplier *big.Int
	}{
		{"Ei", exp2(60)}, {"Pi", exp2(50)}, {"Ti", exp2(40)},
		{"Gi", exp2(30)}, {"Mi", exp2(20)}, {"Ki", exp2(10)},
	}
	decimalSuffixes = []struct {
		suffix     string
		multiplier *big.Int
	}{
		{"E", exp10(18)}, {"P", exp10(15)}, {"T", exp10(12)},
		{"G", exp10(9)}, {"M", exp10(6)}, {"k", exp10(3)},
	}
)

func exp2(n int64) *big.Int  { return new(big.Int).Lsh(big.NewInt(1), uint(n)) }
func exp10(n int64) *big.Int { return new(big.Int).Exp(big.NewInt(10), big.NewInt(n), nil) }

// parseQuantity parses one Kubernetes quantity exactly. With milli true
// (CPU) the result is millicores and the only legal suffix is "m"; with
// milli false (memory, storage) the result is bytes and the binary and
// decimal suffixes apply. A quantity that does not resolve to a whole
// millicore or byte is refused -- rounding a reservation would misstate
// what the scheduler reserves.
func parseQuantity(text string, milli bool) (*big.Int, error) {
	number := text
	multiplier := big.NewInt(1)
	if milli {
		multiplier = big.NewInt(1000)
		if trimmed, found := strings.CutSuffix(text, "m"); found {
			number, multiplier = trimmed, big.NewInt(1)
		}
	} else {
		for _, candidate := range binarySuffixes {
			if trimmed, found := strings.CutSuffix(text, candidate.suffix); found {
				number, multiplier = trimmed, candidate.multiplier
				break
			}
		}
		if multiplier.Cmp(big.NewInt(1)) == 0 {
			for _, candidate := range decimalSuffixes {
				if trimmed, found := strings.CutSuffix(text, candidate.suffix); found {
					number, multiplier = trimmed, candidate.multiplier
					break
				}
			}
		}
	}

	value, ok := new(big.Rat).SetString(number)
	if !ok || value.Sign() < 0 {
		return nil, fmt.Errorf("quantity %q is not a non-negative Kubernetes quantity", text)
	}
	value.Mul(value, new(big.Rat).SetInt(multiplier))
	if !value.IsInt() {
		return nil, fmt.Errorf("quantity %q does not resolve to a whole unit -- refusing to round a reservation", text)
	}
	return new(big.Int).Set(value.Num()), nil
}

// renderCpu renders a millicore total the way a person writes it: whole
// cores plain ("3"), anything finer in millicores ("1500m"). Zero renders
// empty -- the footprint omits what nothing reserves.
func renderCpu(milli *big.Int) string {
	if milli.Sign() == 0 {
		return ""
	}
	quotient, remainder := new(big.Int).QuoRem(milli, big.NewInt(1000), new(big.Int))
	if remainder.Sign() == 0 {
		return quotient.String()
	}
	return milli.String() + "m"
}

// renderBytes renders a byte total in the LARGEST suffix (binary or
// decimal) that divides it evenly -- binary inputs come back in binary
// units, decimal inputs in decimal units, and a total no suffix divides
// renders as plain bytes. Zero renders empty.
func renderBytes(bytes *big.Int) string {
	if bytes.Sign() == 0 {
		return ""
	}
	// Interleave the two families by multiplier size (Ei > E > Pi > P >
	// ... > Ki > k) so the first even division is the largest unit.
	for i := 0; i < len(binarySuffixes); i++ {
		for _, candidate := range []struct {
			suffix     string
			multiplier *big.Int
		}{binarySuffixes[i], decimalSuffixes[i]} {
			quotient, remainder := new(big.Int).QuoRem(bytes, candidate.multiplier, new(big.Int))
			if remainder.Sign() == 0 {
				return quotient.String() + candidate.suffix
			}
		}
	}
	return bytes.String()
}
