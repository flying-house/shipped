package main

import (
	"encoding/json"
	"fmt"
	"image/color"
	"io/ioutil"
	"log"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	screenWidth  = 1000
	screenHeight = 1000
)

// TrailPoint -
type TrailPoint struct {
	X, Y float64
}

// Game -
type Game struct {
	vessel    *Vessel
	trail     []TrailPoint
	trailChan chan TrailPoint
	lastX     float64
	lastY     float64
}

// BoatData -
type BoatData struct {
	Boats []map[string]interface{} `json:"boats"`
}

// EngineData -
type EngineData struct {
	Engines []map[string]interface{} `json:"engines"`
}

// LoadBoats -
func LoadBoats(filename string) ([]map[string]interface{}, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var boatData BoatData
	err = json.Unmarshal(data, &boatData)
	if err != nil {
		return nil, err
	}

	return boatData.Boats, nil
}

// LoadEngines -
func LoadEngines(filename string) ([]map[string]interface{}, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var engineData EngineData
	err = json.Unmarshal(data, &engineData)
	if err != nil {
		return nil, err
	}

	return engineData.Engines, nil
}

// FindBoatByName -
func FindBoatByName(boats []map[string]interface{}, name string) (map[string]interface{}, error) {
	for _, boat := range boats {
		if boat["name"].(string) == name {
			return boat, nil
		}
	}
	return nil, fmt.Errorf("boat '%s' not found", name)
}

// FindEngineByName -
func FindEngineByName(engines []map[string]interface{}, name string) (map[string]interface{}, error) {
	for _, engine := range engines {
		if engine["name"].(string) == name {
			return engine, nil
		}
	}
	return nil, fmt.Errorf("engine '%s' not found", name)
}

// Update -
func (g *Game) Update() error {
	g.vessel.HandleInput()
	UpdatePhysics(g.vessel)
	g.updateTrailFromChannel()
	return nil
}

func (g *Game) updateTrailFromChannel() {
	for {
		select {
		case point := <-g.trailChan:
			g.trail = append(g.trail, point)
			if len(g.trail) > 1000 {
				g.trail = g.trail[1:]
			}
		default:
			return
		}
	}
}

// DrawTrail -
func (g *Game) DrawTrail(screen *ebiten.Image) {
	if len(g.trail) < 2 {
		return
	}

	trailColor := color.RGBA{128, 128, 128, 255}
	for i := 0; i < len(g.trail)-1; i++ {
		ebitenutil.DrawLine(screen,
			g.trail[i].X, g.trail[i].Y,
			g.trail[i+1].X, g.trail[i+1].Y,
			trailColor)
	}
}

// Draw -
func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0, 120, 255, 255})
	g.DrawTrail(screen)
	DrawVessel(screen, g.vessel)
	DrawHUD(screen, g)
}

// Layout -
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	boats, err := LoadBoats("items/boats.json")
	if err != nil {
		log.Fatalf("Failed to load boats: %v", err)
	}

	engines, err := LoadEngines("items/engines.json")
	if err != nil {
		log.Fatalf("Failed to load engines: %v", err)
	}

	selectedBoat, err := FindBoatByName(boats, "15' Johnboat")
	if err != nil {
		log.Fatalf("Failed to find boat: %v", err)
	}

	selectedEngine, err := FindEngineByName(engines, "Jawnsen 10hp")
	if err != nil {
		log.Fatalf("Failed to find engine: %v", err)
	}

	engineSpecs := []map[string]interface{}{selectedEngine}
	vessel := NewVessel(screenWidth/2, screenHeight/2, selectedBoat, engineSpecs)

	game := &Game{
		vessel:    vessel,
		trail:     make([]TrailPoint, 0),
		trailChan: make(chan TrailPoint, 100),
		lastX:     vessel.X,
		lastY:     vessel.Y,
	}

	go func() {
		ticker := time.NewTicker(250 * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				currentX, currentY := vessel.X, vessel.Y
				distance := math.Hypot(currentX-game.lastX, currentY-game.lastY)
				if distance > 2.0 {
					select {
					case game.trailChan <- TrailPoint{X: currentX, Y: currentY}:
						game.lastX, game.lastY = currentX, currentY
					default:
					}
				}
			}
		}
	}()

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Shipped - Component-Based Maritime Simulation")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
