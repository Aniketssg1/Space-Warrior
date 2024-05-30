package game

import (
	"math"
	"time"

	"github.com/ThreeDotsLabs/meteors/assets"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	rotationPerSecond = math.Pi / 30
	moveSpeed         = 3
	bulletSpawnOffset = 50.0
	bulletSpawnTime   = 1 * time.Second
)

var shootCooldown = time.Millisecond * 400

type Player struct {
	game *Game

	name     string
	score    int
	position Vector
	rotation float64
	forward  float64
	backward float64
	sprite   *ebiten.Image

	bulletSpawnTimer *Timer
	shootCooldown    *Timer
}

func NewPlayer(game *Game) *Player {
	sprite := assets.PlayerSprite

	bounds := sprite.Bounds()
	halfW := float64(bounds.Dx()) / 2
	halfH := float64(bounds.Dy()) / 2

	pos := Vector{
		X: screenWidth/2 - halfW,
		Y: screenHeight/2 - halfH,
	}

	return &Player{
		game:             game,
		position:         pos,
		rotation:         0,
		forward:          0,
		backward:         0,
		sprite:           sprite,
		bulletSpawnTimer: NewTimer(bulletSpawnTime),
		shootCooldown:    NewTimer(shootCooldown),
	}
}

func (p *Player) Draw(screen *ebiten.Image) {
	bounds := p.sprite.Bounds()
	halfW := float64(bounds.Dx()) / 2
	halfH := float64(bounds.Dy()) / 2

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-halfW, -halfH)
	op.GeoM.Rotate(p.rotation)
	op.GeoM.Translate(halfW, halfH)

	op.GeoM.Translate(p.position.X, p.position.Y)

	screen.DrawImage(p.sprite, op)
}

func (p *Player) Update() {
	/* Space ship moving logic */
	speed := rotationPerSecond / float64(math.Pi)

	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		p.rotation -= speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		p.rotation += speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		p.position.X += moveSpeed * math.Sin(p.rotation)
		p.position.Y -= moveSpeed * math.Cos(p.rotation)
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		p.position.X -= moveSpeed * math.Sin(p.rotation)
		p.position.Y += moveSpeed * math.Cos(p.rotation)
	}

	/* Bullets Speedup */
	p.bulletSpawnTimer.Update()
	if p.bulletSpawnTimer.IsReady() {
		p.bulletSpawnTimer.Reset()
		shootCooldown -= 20
	}

	/* Bullets firing logic */
	p.shootCooldown.Update()
	if p.shootCooldown.IsReady() && ebiten.IsKeyPressed(ebiten.KeySpace) {
		p.shootCooldown.Reset()

		bounds := p.sprite.Bounds()
		halfW := float64(bounds.Dx()) / 2
		halfH := float64(bounds.Dy()) / 2

		spawnPos := Vector{
			p.position.X + halfW + math.Sin(p.rotation)*bulletSpawnOffset,
			p.position.Y + halfH + math.Cos(p.rotation)*-bulletSpawnOffset,
		}

		fire := NewFire(spawnPos, p.rotation)
		bullet := NewBullet(spawnPos, p.rotation)
		p.game.AddBullet(bullet)
		p.game.AddFire(fire)
	}
}

func (p *Player) Collider() Circle {
	bounds := p.sprite.Bounds()

	return NewCircle(
		p.position.X,
		p.position.Y,
		math.Sqrt(math.Pow(float64(bounds.Dx()), 2)+math.Pow(float64(bounds.Dy()), 2))/2,
	)
}
