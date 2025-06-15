package main

// Engine -
type Engine struct {
	Name       string
	Mount      string
	IdlePower  float64
	IdleRPM    float64
	RatedPower float64
	RatedRPM   float64
	SFC        float64 // Specific Fuel Consumption (lbs/hp/min)
	FuelType   string
	EngineType string
	Cost       float64
	Weight     float64
}

// NewEngine -
func NewEngine(
	name string,
	mount string,
	ratedPower float64,
	idlePower float64,
	idleRPM float64,
	ratedRPM float64,
	sfc float64,
	fuelType string,
	engineType string,
) *Engine {
	return &Engine{
		Name:       name,
		Mount:      mount,
		IdlePower:  idlePower,
		IdleRPM:    idleRPM,
		RatedPower: ratedPower,
		RatedRPM:   ratedRPM,
		SFC:        sfc,
		FuelType:   fuelType,
		EngineType: engineType,
		Cost:       0,
		Weight:     0,
	}
}

// NewEngineFromSpec -
func NewEngineFromSpec(spec map[string]interface{}) *Engine {
	return &Engine{
		Name:       spec["name"].(string),
		Mount:      spec["mount"].(string),
		IdlePower:  spec["idlePower"].(float64),
		IdleRPM:    spec["idleRPM"].(float64),
		RatedPower: spec["power"].(float64),
		RatedRPM:   spec["ratedRPM"].(float64),
		SFC:        spec["sfc"].(float64),
		FuelType:   spec["fuelType"].(string),
		EngineType: spec["engineType"].(string),
		Cost:       spec["cost"].(float64),
		Weight:     spec["weight"].(float64),
	}
}

// GetPower -
func (e *Engine) GetPower(throttlePercent float64) float64 {
	if throttlePercent <= 0 {
		return 0
	}
	if throttlePercent >= 100 {
		return e.RatedPower
	}

	idleThrottlePercent := 2.0

	if throttlePercent <= idleThrottlePercent {
		return e.IdlePower * (throttlePercent / idleThrottlePercent)
	}

	powerRange := e.RatedPower - e.IdlePower
	throttleRange := 100.0 - idleThrottlePercent
	throttleAboveIdle := throttlePercent - idleThrottlePercent

	return e.IdlePower + (powerRange * throttleAboveIdle / throttleRange)
}

// GetFuelFlow -
func (e *Engine) GetFuelFlow(power float64) float64 {
	return power * e.SFC
}

// GetRPM -
func (e *Engine) GetRPM(throttlePercent float64, speed float64) float64 {
	if throttlePercent <= 0 {
		return e.IdleRPM
	}
	if throttlePercent >= 100 {
		return e.RatedRPM
	}

	rpmRange := e.RatedRPM - e.IdleRPM
	throttleRange := 100.0
	throttleAboveIdle := throttlePercent - 2.0

	return e.IdleRPM + (rpmRange * throttleAboveIdle / throttleRange)
	// TODO: calculate "loaded"RPM based on vessel speed
	// return e.IdleRPM + (rpmRange * speed / e.RatedRPM)
}
