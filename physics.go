package main

import (
	"math"
)

// UpdatePhysics -
func UpdatePhysics(v *Vessel) {
	doThrust(v)
	doDrag(v)
	doRotate(v)

	v.X += v.VelocityX
	v.Y += v.VelocityY

	v.Heading = math.Mod(v.Heading+360, 360)
}

func doThrust(v *Vessel) {
	if v.Gear == GearNeutral {
		return
	}

	thrust := v.GetCurrentThrust()
	thrustDirection := v.Heading + v.RudderAngle
	thrustAngleRad := (thrustDirection - 90) * math.Pi / 180

	thrustX := thrust * math.Cos(thrustAngleRad)
	thrustY := thrust * math.Sin(thrustAngleRad)

	v.VelocityX += thrustX / v.Mass
	v.VelocityY += thrustY / v.Mass

	var thrustPointX, thrustPointY float64
	if len(v.EngineMountPoints) > 0 {
		thrustPointX = v.EngineMountPoints[0][0]
		thrustPointY = v.EngineMountPoints[0][1]
	}

	vesselAngleRad := v.Heading * math.Pi / 180
	thrustWorldX := thrustPointX*math.Cos(vesselAngleRad) - thrustPointY*math.Sin(vesselAngleRad)
	thrustWorldY := thrustPointX*math.Sin(vesselAngleRad) + thrustPointY*math.Cos(vesselAngleRad)

	leverArmX := thrustWorldX - v.CenterOfMassX
	leverArmY := thrustWorldY - v.CenterOfMassY

	torque := leverArmX*thrustY - leverArmY*thrustX

	angularAcceleration := torque / v.MomentOfInertia
	v.AngularVelocity += angularAcceleration
}

func doRotate(v *Vessel) {
	v.AngularVelocity *= v.RotationalDrag

	if math.Abs(v.AngularVelocity) < 0.001 {
		v.AngularVelocity = 0
	}

	v.Heading += v.AngularVelocity
}

func doDrag(v *Vessel) {
	v.VelocityX *= v.Drag
	v.VelocityY *= v.Drag

	vesselAngleRad := (v.Heading - 90) * math.Pi / 180
	forwardX := math.Cos(vesselAngleRad)
	forwardY := math.Sin(vesselAngleRad)

	forwardVel := v.VelocityX*forwardX + v.VelocityY*forwardY
	sidewaysVelX := v.VelocityX - forwardVel*forwardX
	sidewaysVelY := v.VelocityY - forwardVel*forwardY

	// Hull-specific sideways drag (keel effect)
	// This could be made vessel-specific in the future
	sidewaysDrag := 0.82
	v.VelocityX = forwardVel*forwardX + sidewaysVelX*sidewaysDrag
	v.VelocityY = forwardVel*forwardY + sidewaysVelY*sidewaysDrag

	speed := math.Hypot(v.VelocityX, v.VelocityY)
	if speed > 0.1 {
		ratioThreshold := 0.001
		if math.Abs(v.VelocityX)/speed < ratioThreshold {
			v.VelocityX = 0
		}
		if math.Abs(v.VelocityY)/speed < ratioThreshold {
			v.VelocityY = 0
		}
	}
}
