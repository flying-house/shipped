# Shipped
*A 2D maritime trading and exploration game*

## Overview
**Shipped** is a top-down 2D game where players captain vessels through an expansive archipelago, transporting passengers and cargo for hire. The game combines realistic boat physics with sandbox-style gameplay, offering both structured missions and free exploration.

**Core Philosophy:** "Sandbox on rails" - plenty of guided opportunities alongside complete freedom to explore and trade.

## 🎮 Gameplay Vision
- **Perspective:** Always north-up oriented
- **Scale:** Dynamic zoom based on vessel speed underway, fixed in port
- **Focus:** Realistic maritime physics with economic trading gameplay
- **Scope:** Inter-island routes designed for 2-5 minute sailing sessions

---

## 🚢 Development Roadmap

### Phase 1: Core Physics Engine
**Goal:** Establish realistic water-based vessel physics

#### ✅ Basic Vessel Control
- [ ] Single-engine vessel
- [ ] Gear system: Reverse → Neutral → Forward
- [ ] Simple controls (WASD):
    - `A` to steer left, `D` to steer right, `G` to recenter
    - Gear/throttle: `W` for (R)->(N)->(F), `S` for (F)->(N)->(R)
    - Cooldown before throttle can increase from in-gear (~1.0s)
    - In reverse, throttle increases with `S`/decreases with `W`

- [ ] Advanced controls (two engine):
    - `A` to steer left, `F` to steer right, `G` to recenter
    - Engine 1: `W` increases gear, `S` decreases gear
    - Engine 2: `E` increases gear, `D` decreases gear

#### Modular Component System
- [ ] **Vessel objects:** Max engine HP, weight, drag, fuel capacity, engine mounts
- [ ] **Engine objects:** Rated HP, RPM range, fuel consumption, mount type
- [ ] **Propeller objects:** Thrust efficiency, reverse scalar, drag coefficients
- [ ] **JSON configuration:** Loadable vessel/engine/prop specifications

#### ✅ Physics Architecture
- [ ] Separate `physics.go` for movement calculations
- [ ] Separate `render.go` for display logic  
- [ ] Separate `vessel.go` for vessel management
- [ ] Realistic thrust, drag, and inertia modeling

### Phase 2: Environmental Systems
**Goal:** Add dynamic environmental factors and basic interactions

#### ✅ Port Interactions
- [ ] Mooring system with proximity detection
- [ ] Mooring cleat positioning and safe docking
- [ ] Basic port entry/exit mechanics

### Phase 3: Inter-Island Gameplay
**Goal:** Implement core trading and navigation gameplay

#### ✅ World Navigation
- [ ] Accelerated map mode for inter-island travel
- [ ] Multiple ports with distinct characteristics
- [ ] Route planning and navigation aids

#### ✅ Economic System
- [ ] Dynamic port markets with price fluctuations
- [ ] Seasonal and time-based market changes
- [ ] Port specializations (tourism, fishing, industrial)
- [ ] Cargo and passenger transport missions

#### ✅ Maritime Hazards
- [ ] Shallow water detection and grounding
- [ ] Obstacle navigation (rocks, reefs)
- [ ] Traffic system (fishing vessels, other boats)
- [ ] Weather impact on navigation

### Phase 4: Shipyard & Progression
**Goal:** Vessel customization and economic progression

#### ✅ Shipyard System
- [ ] Vessel upgrades and modifications
- [ ] Engine swapping and tuning
- [ ] Repair and maintenance mechanics
- [ ] New vessel purchases

#### ✅ Port Variation
- [ ] Location-specific pricing and availability
- [ ] Reputation system affecting prices
- [ ] Specialized shipyard capabilities per port

### Phase 5: Exploration & Discovery
**Goal:** Extended gameplay through exploration mechanics

#### ✅ Discovery System
- [ ] Hidden/unmapped port locations
- [ ] Abandoned vessel recovery and repair
- [ ] Treasure diving and salvage operations
- [ ] Easter eggs and secret locations

#### ✅ Specialized Activities
- [ ] Sightseeing tour missions
- [ ] Fishing charter operations
- [ ] Diving expedition support
- [ ] Emergency rescue missions

---

## 🛠️ Technical Stack
- **Engine:** Go + Ebiten (2D game framework)
- **Assets:** Bitmap graphics, JSON configuration files
- **Architecture:** Component-based entity system
- **Platform:** Cross-platform desktop (Windows, macOS, Linux)

## 🎯 Current Focus
Working on **Phase 1** - establishing the foundation of realistic boat physics and intuitive control systems that will make the maritime experience authentic and engaging.