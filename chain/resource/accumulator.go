package resource

import (
	"fmt"

	"github.com/MetalBlockchain/antelopevm/config"
	"github.com/MetalBlockchain/antelopevm/math"
	"github.com/MetalBlockchain/antelopevm/utils"
)

type Ratio struct {
	Numerator   uint64 `json:"numerator"`
	Denominator uint64 `json:"denominator"`
}

type ElasticLimitParameters struct {
	Target        uint64 `json:"target"`
	Max           uint64 `json:"max"`
	Periods       uint32 `json:"periods"`
	MaxMultiplier uint32 `json:"max_multiplier"`
	ContractRate  Ratio  `json:"contract_rate"`
	ExpandRate    Ratio  `json:"expand_rate"`
}

func (e ElasticLimitParameters) Validate() error {
	if e.Periods <= 0 {
		return fmt.Errorf("elastic limit parameter 'periods' cannot be zero")
	}

	if e.ContractRate.Denominator <= 0 {
		return fmt.Errorf("elastic limit parameter 'contract_rate' is not a well-defined ratio")
	}

	if e.ExpandRate.Denominator <= 0 {
		return fmt.Errorf("elastic limit parameter 'expand_rate' is not a well-defined ratio")
	}

	return nil
}

func UpdateElasticLimit(currentLimit uint64, averageUsage uint64, params ElasticLimitParameters) uint64 {
	result := currentLimit
	if averageUsage > params.Target {
		result = result * params.ContractRate.Numerator / params.ContractRate.Denominator
	} else {
		result = result * params.ExpandRate.Numerator / params.ExpandRate.Denominator
	}
	return utils.Min(utils.Max(result, params.Max), uint64(params.Max*uint64(params.MaxMultiplier)))
}

func IntegerDivideCeil(num uint64, den uint64) uint64 {
	return (num / den) + func() uint64 {
		if num%den > 0 {
			return 1
		}
		return 0
	}()
}

type exponentialMovingAverageAccumulator struct {
	LastOrdinal uint32 `json:"last_ordinal"`
	ValueEx     uint64 `json:"value_ex"`
	Consumed    uint64 `json:"consumed"`

	precision   uint64
	maxRawValue uint64
}

func NewExponentialMovingAverageAccumulator() exponentialMovingAverageAccumulator {
	return exponentialMovingAverageAccumulator{
		LastOrdinal: 0,
		ValueEx:     0,
		Consumed:    0,

		precision:   config.RateLimitingPrecision,
		maxRawValue: math.MaxUint64 / config.RateLimitingPrecision,
	}
}

func (ema *exponentialMovingAverageAccumulator) Average() uint64 {
	return IntegerDivideCeil(ema.ValueEx, ema.precision)
}

func (ema *exponentialMovingAverageAccumulator) Add(units uint64, ordinal uint32, windowSize uint32) error {
	if units > ema.maxRawValue {
		return fmt.Errorf("usage exceeds maximum value representable after extending for precision")
	} else if math.MaxUint64-ema.Consumed < units {
		return fmt.Errorf("overflow in tracked usage when adding usage!")
	}

	valueExContrib := IntegerDivideCeil(units*config.RateLimitingPrecision, uint64(windowSize))
	if math.MaxUint64-ema.ValueEx < valueExContrib {
		return fmt.Errorf("overflow in accumulated value when adding usage!")
	}

	if ema.LastOrdinal != ordinal {
		if ordinal <= ema.LastOrdinal {
			return fmt.Errorf("new ordinal cannot be less than the previous ordinal")
		}

		if ema.LastOrdinal+windowSize > ordinal {
			delta := ordinal - ema.LastOrdinal
			decay := makeRatio(uint64(windowSize-delta), uint64(windowSize))

			if value, err := decay.Mul(ema.ValueEx); err != nil {
				return err
			} else {
				ema.ValueEx = value
			}
		} else {
			ema.ValueEx = 0
		}

		ema.LastOrdinal = ordinal
		ema.Consumed = ema.Average()
	}

	ema.Consumed += units
	ema.ValueEx += valueExContrib

	return nil
}

func makeRatio(numerator uint64, denominator uint64) Ratio {
	return Ratio{numerator, denominator}
}

func (r *Ratio) Mul(value uint64) (uint64, error) {
	if r.Numerator != 0 && math.MaxUint64/r.Numerator < value {
		return 0, fmt.Errorf("usage exceeds maximum value representable after extending for precision")
	}

	return (value * r.Numerator) / r.Denominator, nil
}

func DowngradeCast(val math.Uint128) int64 {
	max := uint64(math.MaxInt64)

	if val.High != 0 && val.Low > max {
		panic("usage exceeds maximum value representable after extending for precision")
	}

	return int64(val.Low)
}

type UsageAccumulator struct {
	exponentialMovingAverageAccumulator
}
