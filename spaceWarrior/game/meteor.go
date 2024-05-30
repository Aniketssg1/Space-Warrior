package game

import (
	"math"
	"math/rand"
	"spaceWarrior/assets"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	rotationSpeedMin = -0.02
	rotationSpeedMax = 0.02
)

type Meteor struct {
	position      Vector
	rotation      float64
	movement      Vector
	rotationSpeed float64
	sprite        *ebiten.Image
}

func NewMeteor(baseVelocity float64, p *Player) *Meteor {
	target := Vector{
		X: p.position.X,
		Y: p.position.Y,
	}

	angle := rand.Float64() * 2 * math.Pi
	r := screenWidth / 2.0
	pos := Vector{
		X: target.X + math.Cos(angle)*r,
		Y: target.Y + math.Sin(angle)*r,
	}

	velocity := baseVelocity + rand.Float64()*1.5
	direction := Vector{
		X: target.X - pos.X,
		Y: target.Y - pos.Y,
	}

	normalizedDirection := direction.Normalize()
	movement := Vector{
		X: normalizedDirection.X * velocity,
		Y: normalizedDirection.Y * velocity,
	}
	sprite := assets.MeteorSprites[rand.Intn(len(assets.MeteorSprites))]

	m := &Meteor{
		position:      pos,
		movement:      movement,
		rotationSpeed: rotationSpeedMin + rand.Float64()*(rotationSpeedMax-rotationSpeedMin),
		sprite:        sprite,
	}
	return m
}

func (m *Meteor) Update() {
	m.position.X += m.movement.X
	m.position.Y += m.movement.Y
	m.rotation += m.rotationSpeed
}

func (m *Meteor) Draw(screen *ebiten.Image) {
	bounds := m.sprite.Bounds()
	halfW := float64(bounds.Dx()) / 2
	halfH := float64(bounds.Dy()) / 2

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-halfW, -halfH)
	op.GeoM.Rotate(m.rotation)
	op.GeoM.Translate(halfW, halfH)

	op.GeoM.Translate(m.position.X, m.position.Y)

	screen.DrawImage(m.sprite, op)
}

func (m *Meteor) Collider() Circle {
	bounds := m.sprite.Bounds()

	return NewCircle(
		m.position.X,
		m.position.Y,
		math.Sqrt(math.Pow(float64(bounds.Dx()), 2)+math.Pow(float64(bounds.Dy()), 2))/2,
	)
}
