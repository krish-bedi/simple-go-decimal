package tinydecimal

import (
	"errors"
	"math"
	"math/big"
	"strconv"
	"strings"
)

// precision of 4 lets us to store up to 4 digits after the decimal (0.0001)
const precision = 4

// scale is our multiplier that lets us store fractional values as integers
var scale int64 = int64(math.Pow10(precision)) // 10^4 = 10,000

var (
	ErrOverflow 	   = errors.New("tinyDecimal: value is too large")
	ErrPrecisionLoss  = errors.New("tinyDecimal: value has too many decimal places")
	ErrDivisionByZero = errors.New("tinyDecimal: division by zero")
)

// max decimal value with a precision of 4 is: 922,337,203,685,477.5807
// (fixed: int64 max, 9,223,372,036,854,775,807, divided by scale, 10,000)
type Decimal struct {
	fixed int64
}

// New builds a Decimal from (value * 10^exp)
// Example: New(123, -2) -> 123 * 10^-2 = 1.23. 
// int8 is used for the exponent as its enough to cover the decimal shift for int64 (19 digits)
func New(value int64, exponent int8) (Decimal, error) {
	shift := int32(exponent) + int32(precision)
	fixed := value

	// convert: value * 10^exp to fixed / scale representation
	if shift > 0 {
		for range shift {
			// Check for positive and negative overflow before multiplying
			if fixed > math.MaxInt64 / 10 || fixed < math.MinInt64 / 10 {
				return Decimal{}, ErrOverflow
			}

			fixed *= 10
		}
	} else if shift < 0 {
		for range -shift {
			// Check for precision loss, ex: New(1, -5) = 0.00001
			// which is less than our precision of 0.0001
			if fixed % 10 != 0 {
				return Decimal{}, ErrPrecisionLoss
			}

			fixed /= 10
		}
	}

	return Decimal{fixed: fixed}, nil
}

func (d Decimal) String() string {
	if d.fixed == 0 {
		return "0"
	}

	negative := d.fixed < 0

	str := strconv.FormatInt(d.fixed, 10)
	if negative {
		str = strings.TrimPrefix(str, "-")
	}

	padLength := (1 + precision) - len(str)
	if padLength > 0 {
		str = strings.Repeat("0", padLength) + str
	}

	decimalPlace := len(str) - precision
	decimalStr := str[:decimalPlace] + "." + str[decimalPlace:]

	resultStr := strings.TrimRight(decimalStr, "0")
	resultStr = strings.TrimSuffix(resultStr, ".")

	if negative {
		resultStr = "-" + resultStr
	}

	return resultStr
}

func (d Decimal) Add(other Decimal) (Decimal, error) {
	// +overflow
	if other.fixed > 0 && d.fixed > 0 && d.fixed > math.MaxInt64 - other.fixed {
		return Decimal{}, ErrOverflow
	}
	// -overflow
	if other.fixed < 0 && d.fixed < 0 && d.fixed < math.MinInt64 - other.fixed {
		return Decimal{}, ErrOverflow
	}
	return Decimal{fixed: d.fixed + other.fixed}, nil
}

func (d Decimal) Sub(other Decimal) (Decimal, error) {
	// +overflow
	if d.fixed >= 0 && other.fixed < 0 && d.fixed > math.MaxInt64 + other.fixed {
		return Decimal{}, ErrOverflow
	}
	// -overflow
	if d.fixed < 0 && other.fixed > 0 && d.fixed < math.MinInt64 + other.fixed {
		return Decimal{}, ErrOverflow
	}
	return Decimal{fixed: d.fixed - other.fixed}, nil
}

func (d Decimal) Multiply(other Decimal) (Decimal, error) {
	// Use BigInt to store transient product that would otherwise overflow int64
	product := new(big.Int).Mul(
		big.NewInt(d.fixed), 
		big.NewInt(other.fixed),
	)

	// Scale back down
	product.Quo(product, big.NewInt(scale))

	// Check if it fits back in int64
	if !product.IsInt64() {
		return Decimal{}, ErrOverflow
	}

	return Decimal{fixed: product.Int64()}, nil
}

func (d Decimal) Divide(other Decimal) (Decimal, error) {
	// Division by Zero error
	if other.fixed == 0 {
		return Decimal{}, ErrDivisionByZero
	}

	// Store transient product in BigInt to prevent overflow
	// Scale up first to prevent integer division precision loss
	product := new(big.Int).Mul(
		big.NewInt(d.fixed),
		big.NewInt(scale),
	)

	product.Quo(product, big.NewInt(other.fixed))

	// Check if it fits back in int64
	if !product.IsInt64() {
		return Decimal{}, ErrOverflow
	}

	return Decimal{fixed: product.Int64()}, nil
}