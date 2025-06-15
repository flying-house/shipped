package main

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// DrawVessel -
func DrawVessel(screen *ebiten.Image, v *Vessel) {
	headingRad := v.Heading * math.Pi / 180

	points := []struct{ x, y float64 }{
		{0, -v.Length / 2},
		{-v.Width / 2, -v.Length / 4},
		{v.Width / 2, -v.Length / 4},
		{-v.Width / 2, -v.Length / 4},
		{-v.Width / 2, v.Length / 2},
		{v.Width / 2, v.Length / 2},
		{v.Width / 2, -v.Length / 4},
	}

	var rotatedPoints []struct{ x, y float64 }
	sin, cos := math.Sin(headingRad), math.Cos(headingRad)

	for _, p := range points {
		rotX := p.x*cos - p.y*sin
		rotY := p.x*sin + p.y*cos

		rotatedPoints = append(rotatedPoints, struct{ x, y float64 }{
			x: rotX + v.X,
			y: rotY + v.Y,
		})
	}

	vesselColor := color.RGBA{100, 150, 255, 255}

	for i := 3; i < 7; i++ {
		j := i + 1
		if j >= 7 {
			j = 3
		}
		vector.StrokeLine(screen,
			float32(rotatedPoints[i].x), float32(rotatedPoints[i].y),
			float32(rotatedPoints[j].x), float32(rotatedPoints[j].y),
			2, vesselColor, false)
	}

	for i := 0; i < 3; i++ {
		j := (i + 1) % 3
		vector.StrokeLine(screen,
			float32(rotatedPoints[i].x), float32(rotatedPoints[i].y),
			float32(rotatedPoints[j].x), float32(rotatedPoints[j].y),
			2, vesselColor, false)
	}

	vector.StrokeLine(screen,
		float32(rotatedPoints[1].x), float32(rotatedPoints[1].y),
		float32(rotatedPoints[3].x), float32(rotatedPoints[3].y),
		2, vesselColor, false)
	vector.StrokeLine(screen,
		float32(rotatedPoints[2].x), float32(rotatedPoints[2].y),
		float32(rotatedPoints[6].x), float32(rotatedPoints[6].y),
		2, vesselColor, false)

	ebitenutil.DrawRect(screen, v.X-2, v.Y-2, 4, 4, color.RGBA{255, 0, 0, 255})

	if v.VelocityX != 0 || v.VelocityY != 0 {
		speed := math.Hypot(v.VelocityX, v.VelocityY)
		scale := 50.0
		endX := v.X + (v.VelocityX/speed)*speed*scale
		endY := v.Y + (v.VelocityY/speed)*speed*scale
		vector.StrokeLine(screen, float32(v.X), float32(v.Y), float32(endX), float32(endY), 2, color.RGBA{0, 255, 0, 255}, false)
	}

	if v.Gear != GearNeutral && v.Throttle > 0 {
		var thrustPointX, thrustPointY float64
		if len(v.EngineMountPoints) > 0 {
			thrustPointX = v.EngineMountPoints[0][0]
			thrustPointY = v.EngineMountPoints[0][1]
		}

		vesselAngleRad := v.Heading * math.Pi / 180
		engineWorldX := v.X + (thrustPointX*math.Cos(vesselAngleRad) - thrustPointY*math.Sin(vesselAngleRad))
		engineWorldY := v.Y + (thrustPointX*math.Sin(vesselAngleRad) + thrustPointY*math.Cos(vesselAngleRad))

		thrustDirection := v.Heading + v.RudderAngle
		thrustAngleRad := (thrustDirection - 90) * math.Pi / 180

		displayAngleRad := thrustAngleRad + math.Pi

		thrust := v.GetCurrentThrust()
		scale := 20.0
		thrustMagnitude := math.Abs(thrust)

		endX := engineWorldX + math.Cos(displayAngleRad)*thrustMagnitude*scale/100
		endY := engineWorldY + math.Sin(displayAngleRad)*thrustMagnitude*scale/100

		vector.StrokeLine(screen,
			float32(engineWorldX), float32(engineWorldY),
			float32(endX), float32(endY),
			2, color.RGBA{0, 150, 255, 255}, false)
	}
}

// DrawHUD -
func DrawHUD(screen *ebiten.Image, g *Game) {
	gearStr := ""
	switch g.vessel.Gear {
	case GearReverse:
		gearStr = "R"
	case GearForward:
		gearStr = "F"
	case GearNeutral:
		gearStr = "N"
	}

	speed := math.Hypot(g.vessel.VelocityX, g.vessel.VelocityY)

	hudText := fmt.Sprintf(
		"Gear: %s\nPos: (%.0f, %.0f)\nThrottle: %.0f%%\nRudder: %.1f\nHeading: %.0f°\nSpeed: %.1fkts",
		gearStr,
		g.vessel.X, g.vessel.Y,
		g.vessel.Throttle*100,
		g.vessel.RudderAngle,
		g.vessel.Heading,
		(speed * 12),
	)

	ebitenutil.DebugPrintAt(screen, hudText, 10, 10)

	controlsText := "Controls:\nW/S: Gear Up/Down\nA/D: Rudder Left/Right\nG: Reset/Neutral"
	ebitenutil.DebugPrintAt(screen, controlsText, 10, screenHeight-80)
}
