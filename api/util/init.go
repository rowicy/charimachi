package util

var (
	violationRates []ViolationRate
)

func init() {

	violationRates, _ = LoadViolationRates("data/violation_rates.json")
}
