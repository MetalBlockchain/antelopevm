package resource

import (
	"fmt"

	"github.com/MetalBlockchain/antelopevm/config"
	"github.com/MetalBlockchain/antelopevm/math"
	"github.com/MetalBlockchain/antelopevm/utils"
)

//go:generate msgp
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
	if num%den > 0 {
		return num/den + 1
	} else {
		return num / den
	}
}

type ExponentialMovingAverageAccumulator struct {
	LastOrdinal uint32 `json:"last_ordinal"`
	ValueEx     uint64 `json:"value_ex"`
	Consumed    uint64 `json:"consumed"`
}

func makeRatio(numerator uint64, denominator uint64) Ratio {
	return Ratio{numerator, denominator}
}

func MultiWithRatio(value uint64, ratio Ratio) uint64 {
	if ratio.Denominator == 0 {
		panic("usage exceeds maximum value representable after extending for precision")
	}

	return value * ratio.Numerator / ratio.Denominator
}

func DowngradeCast(val math.Uint128) int64 {
	max := uint64(math.MaxInt64)

	if val.High != 0 && val.Low > max {
		panic("usage exceeds maximum value representable after extending for precision")
	}

	return int64(val.Low)
}

func (ema *ExponentialMovingAverageAccumulator) Average() uint64 {
	return IntegerDivideCeil(ema.ValueEx, uint64(config.RateLimitingPrecision))
}

func (ema *ExponentialMovingAverageAccumulator) Add(units uint64, ordinal uint32, windowSize uint32) {
	valueExContrib := IntegerDivideCeil(units*uint64(config.RateLimitingPrecision), uint64(windowSize))

	if ema.LastOrdinal != ordinal {
		if ema.LastOrdinal+windowSize > ordinal {
			delta := ordinal - ema.LastOrdinal
			decay := makeRatio(uint64(windowSize-delta), uint64(windowSize))
			ema.ValueEx = MultiWithRatio(ema.ValueEx, decay)
		} else {
			ema.ValueEx = 0
		}

		ema.LastOrdinal = ordinal
		ema.Consumed = ema.Average()
	}

	ema.Consumed += units
	ema.ValueEx += valueExContrib
}

type UsageAccumulator struct {
	ExponentialMovingAverageAccumulator
}
