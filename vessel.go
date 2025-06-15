package main

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	// GearReverse -
	GearReverse = -1
	// GearNeutral -
	GearNeutral = 0
	// GearForward -
	GearForward = 1
)

const (
	fuelDiesel   = 6.7
	fuelGasoline = 6.0
)

// Vessel -
type Vessel struct {
	// Position and movement
	X, Y                 float64
	Heading              float64
	VelocityX, VelocityY float64
	AngularVelocity      float64

	// Control systems
	Gear           int
	Throttle       float64
	RudderRate     float64
	RudderAngle    float64
	RudderLimit    float64
	LastGearChange time.Time
	GearCooldown   time.Duration

	// Hull specifications (from boat component)
	Length              float64
	Width               float64
	EmptyWeight         float64
	MaxSpeed            float64
	MaxPower            float64 // Maximum power the hull can handle
	Drag                float64
	FuelCapacityGallons float64
	FuelTankLocationX   float64
	FuelTankLocationY   float64
	EngineMount         string
	EngineMountPoints   [][2]float64

	// Physics properties (calculated from components)
	Mass            float64 // Total mass including fuel, cargo, etc.
	CenterOfMassX   float64 // Calculated center of mass
	CenterOfMassY   float64
	MomentOfInertia float64 // Calculated from mass distribution
	RotationalDrag  float64

	// Engine system
	Engines             []*Engine
	TotalInstalledPower float64 // Sum of all engine power
	MaxThrust           float64 // Calculated from engines and hull
	CurrentSpeed        float64
}

// NewVessel creates a vessel from boat and engine specifications
func NewVessel(x, y float64, boatSpec map[string]interface{}, engineSpecs []map[string]interface{}) *Vessel {
	v := &Vessel{
		X:              x,
		Y:              y,
		Heading:        0,
		Gear:           GearNeutral,
		Throttle:       0,
		RudderAngle:    0,
		RudderRate:     0.85, // degrees/frame
		GearCooldown:   time.Second,
		RotationalDrag: 0.985,
	}

	// Load hull specifications
	v.loadHullSpecs(boatSpec)

	// Load and install engines
	v.loadEngines(engineSpecs)

	// Calculate derived properties
	v.calculatePhysicsProperties()

	return v
}

// loadHullSpecs loads boat specifications from JSON data
func (v *Vessel) loadHullSpecs(spec map[string]interface{}) {
	v.Length = spec["length"].(float64)
	v.Width = spec["width"].(float64)
	v.EmptyWeight = spec["emptyWeight"].(float64)
	v.MaxSpeed = spec["maxSpeed"].(float64)
	v.MaxPower = spec["maxPower"].(float64)
	v.Drag = spec["drag"].(float64)
	v.FuelCapacityGallons = spec["fuelCapacityGallons"].(float64)
	v.EngineMount = spec["engineMount"].(string)

	// Load fuel tank location
	fuelLoc := spec["fuelTankLocation"].([]interface{})
	v.FuelTankLocationX = fuelLoc[0].(float64)
	v.FuelTankLocationY = fuelLoc[1].(float64)

	// Load engine mount points
	mountPoints := spec["engineMountPoints"].([]interface{})
	v.EngineMountPoints = make([][2]float64, len(mountPoints))
	for i, point := range mountPoints {
		pointArray := point.([]interface{})
		v.EngineMountPoints[i] = [2]float64{
			pointArray[0].(float64),
			pointArray[1].(float64),
		}
	}

	// Set rudder limit based on boat size (larger boats turn slower)
	v.RudderLimit = 75.0 - (v.Length-15.0)*2.0 // Smaller boats turn sharper
	if v.RudderLimit < 30.0 {
		v.RudderLimit = 30.0
	}
}

// loadEngines creates engine instances from specifications
func (v *Vessel) loadEngines(engineSpecs []map[string]interface{}) {
	v.Engines = make([]*Engine, len(engineSpecs))
	v.TotalInstalledPower = 0

	for i, spec := range engineSpecs {
		v.Engines[i] = NewEngineFromSpec(spec)
		v.TotalInstalledPower += v.Engines[i].RatedPower
	}
}

func (v *Vessel) calculatePhysicsProperties() {
	fuelWeight := v.FuelCapacityGallons * fuelDiesel
	v.Mass = v.EmptyWeight + fuelWeight

	// TODO: calculate center of mass accurately
	v.CenterOfMassX = 0
	v.CenterOfMassY = 0

	// simplified rectangular approximation: I = m(L²+W²)/12
	v.MomentOfInertia = v.Mass * (v.Length*v.Length + v.Width*v.Width) / 12.0

	thrustPerHP := 3.5
	maxEngineThrust := v.TotalInstalledPower * thrustPerHP

	hullLimitedThrust := v.MaxPower * thrustPerHP

	if maxEngineThrust < hullLimitedThrust {
		v.MaxThrust = maxEngineThrust
	} else {
		v.MaxThrust = hullLimitedThrust
	}

	v.MaxThrust = v.MaxThrust / 10.0
}

// GetCurrentThrust -
func (v *Vessel) GetCurrentThrust() float64 {
	if v.Gear == GearNeutral {
		return 0
	}

	// TODO: calculate thrust based on power and RPM
	totalThrust := 0.0
	for _, engine := range v.Engines {
		enginePower := engine.GetPower(v.Throttle * 100)
		totalThrust += enginePower * engine.GetRPM(v.Throttle*100, v.CurrentSpeed)
		v.CurrentSpeed = v.CurrentSpeed + (totalThrust / v.Mass)
	}

	return totalThrust * float64(v.Gear)
}

// HandleInput -
func (v *Vessel) HandleInput() {
	if ebiten.IsKeyPressed(ebiten.KeyY) {
		v.X = screenWidth / 2
		v.Y = screenHeight / 2
	}

	if ebiten.IsKeyPressed(ebiten.KeyA) {
		v.RudderAngle += v.RudderRate
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		v.RudderAngle -= v.RudderRate
	}
	if v.RudderAngle < -v.RudderLimit {
		v.RudderAngle = -v.RudderLimit
	}
	if v.RudderAngle > v.RudderLimit {
		v.RudderAngle = v.RudderLimit
	}
	if ebiten.IsKeyPressed(ebiten.KeyG) {
		v.RudderAngle = 0
	}

	now := time.Now()
	if now.Sub(v.LastGearChange) >= v.GearCooldown {
		if ebiten.IsKeyPressed(ebiten.KeyW) && v.Gear < GearForward {
			v.Gear++
			v.LastGearChange = now
			v.Throttle = 0
		}
		if ebiten.IsKeyPressed(ebiten.KeyS) && v.Gear > GearReverse {
			v.Gear--
			v.LastGearChange = now
			v.Throttle = 0
		}
	}

	if v.Gear != GearNeutral {
		if v.Gear == GearForward {
			if ebiten.IsKeyPressed(ebiten.KeyW) && v.Throttle < 1.00 {
				v.Throttle += 0.01
			}
			if ebiten.IsKeyPressed(ebiten.KeyS) && v.Throttle > 0.00 {
				v.Throttle -= 0.01
			}
		} else if v.Gear == GearReverse {
			if ebiten.IsKeyPressed(ebiten.KeyS) && v.Throttle < 1.00 {
				v.Throttle += 0.01
			}
			if ebiten.IsKeyPressed(ebiten.KeyW) && v.Throttle > 0.00 {
				v.Throttle -= 0.01
			}
		}
	}
}
