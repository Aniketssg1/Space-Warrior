package game

import (
	"math"
	"spaceWarrior/assets"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	bulletSpeedPerSecond = 300.0
	fireDisplayDuration  = 200 * time.Millisecond // Display fire for 200ms
)

type Bullet struct {
	position Vector
	rotation float64
	sprite   *ebiten.Image
}

type Fire struct {
	position  Vector
	rotation  float64
	sprite    *ebiten.Image
	timer     *Timer
	isVisible bool
}

func NewFire(pos Vector, rotation float64) *Fire {
	sprite := assets.FireSprite

	bounds := sprite.Bounds()
	halfW := float64(bounds.Dx()) / 2
	halfH := float64(bounds.Dy()) / 2

	pos.X -= halfW
	pos.Y -= halfH

	f := &Fire{
		position:  pos,
		rotation:  rotation,
		sprite:    sprite,
		timer:     NewTimer(fireDisplayDuration),
		isVisible: true,
	}

	return f
}

func NewBullet(pos Vector, rotation float64) *Bullet {
	sprite := assets.LaserSprite

	bounds := sprite.Bounds()
	halfW := float64(bounds.Dx()) / 2
	halfH := float64(bounds.Dy()) / 2

	pos.X -= halfW
	pos.Y -= halfH

	b := &Bullet{
		position: pos,
		rotation: rotation,
		sprite:   sprite,
	}

	// Play gunshot sound
	assets.GunShot.Rewind()
	assets.GunShot.Play()
	return b
}

func (b *Bullet) Draw(screen *ebiten.Image) {
	bounds := b.sprite.Bounds()
	halfW := float64(bounds.Dx()) / 2
	halfH := float64(bounds.Dy()) / 2

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-halfW, -halfH)
	op.GeoM.Rotate(b.rotation)
	op.GeoM.Translate(halfW, halfH)

	op.GeoM.Translate(b.position.X, b.position.Y)

	screen.DrawImage(b.sprite, op)
}

func (f *Fire) Draw(screen *ebiten.Image) {
	if !f.isVisible {
		return
	}

	bounds := f.sprite.Bounds()
	halfW := float64(bounds.Dx()) / 2
	halfH := float64(bounds.Dy()) / 2

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-halfW, -halfH)
	op.GeoM.Rotate(f.rotation)
	op.GeoM.Translate(halfW, halfH)

	op.GeoM.Translate(f.position.X, f.position.Y)

	screen.DrawImage(f.sprite, op)
}

func (b *Bullet) Update() {
	speed := bulletSpeedPerSecond / float64(ebiten.TPS())
	b.position.X += math.Sin(b.rotation) * speed
	b.position.Y += math.Cos(b.rotation) * -speed
}

func (f *Fire) Update() {
	f.timer.Update()
	if f.timer.IsReady() {
		f.isVisible = false
	}
}

func (b *Bullet) Collider() Rect {
	bounds := b.sprite.Bounds()

	return NewRect(
		b.position.X,
		b.position.Y,
		float64(bounds.Dx()),
		float64(bounds.Dy()),
	)
}
