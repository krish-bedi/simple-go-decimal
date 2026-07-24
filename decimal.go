package tinydecimal

import (
	"errors"
	"math"
	"strconv"
	"strings"
)

// precision of 4 lets us to store up to 4 digits after the decimal (0.0001)
// This is enough to store monetary values (2 decimals) with headroom for sub-cent calculations
const precision = 4

// scale is our multiplier that lets us store fractional values as integers
const scale int64 = 10000 // 10^4 = 10,000

var (
	ErrOverflow 	  = errors.New("tinyDecimal: value is too large")
	ErrPrecisionLoss = errors.New("tinyDecimal: value has too many decimal places")
)

// max decimal value with a precision of 4 is: 922,337,203,685,477.5807
// (fixed: int64 max, 9,223,372,036,854,775,807, divided by scale, 10,000)
type Decimal struct {
	fixed int64
}

// New builds a Decimal from (value * 10^exp)
// Example: New(123, -2) -> 123 * 10^-2 = 1.23
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

	num := d.fixed

	negative := num < 0
	if negative {
		num = -num
	}

	str := strconv.FormatInt(num, 10)

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

func (d Decimal) Add(other Decimal) Decimal {
	return Decimal{fixed: d.fixed + other.fixed}
}

func (d Decimal) Sub(other Decimal) Decimal {
	return Decimal{fixed: d.fixed - other.fixed}
}