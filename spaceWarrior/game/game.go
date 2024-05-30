package game

import (
	"fmt"
	"image/color"
	"math/rand"
	"time"

	"github.com/ThreeDotsLabs/meteors/assets"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
)

const (
	screenWidth  = 800
	screenHeight = 600

	meteorSpawnTime = 1 * time.Second

	baseMeteorVelocity  = 0.25
	meteorSpeedUpAmount = 0.1
	meteorSpeedUpTime   = 5 * time.Second
)

type GameState int

const (
	GameStateStartMenu GameState = iota
	GameStatePlaying
)

type Game struct {
	state            GameState
	background       *Background
	player           *Player
	meteorSpawnTimer *Timer
	meteors          []*Meteor
	bullets          []*Bullet
	fire             []*Fire

	score int

	baseVelocity  float64
	velocityTimer *Timer
}

func NewGame() *Game {
	g := &Game{
		state:            GameStateStartMenu,
		meteorSpawnTimer: NewTimer(meteorSpawnTime),
		baseVelocity:     baseMeteorVelocity,
		velocityTimer:    NewTimer(meteorSpeedUpTime),
	}
	g.player = NewPlayer(g)
	g.background = NewBackground(g)

	return g
}

func (g *Game) Update() error {
	switch g.state {
	case GameStateStartMenu:
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			g.state = GameStatePlaying
		}
	case GameStatePlaying:
		g.updateGame()
	}
	return nil
}

func (g *Game) updateGame() {
	/* Speedup the meteors as game progresses */
	g.velocityTimer.Update()
	if g.velocityTimer.IsReady() {
		g.velocityTimer.Reset()
		g.baseVelocity += meteorSpeedUpAmount
	}

	/* Render the player */
	g.player.Update()

	/* Render the meteors */
	g.meteorSpawnTimer.Update()
	if g.meteorSpawnTimer.IsReady() {
		g.meteorSpawnTimer.Reset()

		m := NewMeteor(g.baseVelocity, g.player)
		g.meteors = append(g.meteors, m)
	}
	for _, m := range g.meteors {
		m.Update()
	}

	/* Render the bullets */
	for _, b := range g.bullets {
		b.Update()
	}

	/* Render the fire*/
	for _, f := range g.fire {
		f.Update()
	}

	/* Check for meteor-bullet collision */
	for i := len(g.meteors) - 1; i >= 0; i-- {
		m := g.meteors[i]
		for j := len(g.bullets) - 1; j >= 0; j-- {
			b := g.bullets[j]
			if b.Collider().IntersectsRect(m.Collider()) {
				g.meteors = append(g.meteors[:i], g.meteors[i+1:]...)
				g.bullets = append(g.bullets[:j], g.bullets[j+1:]...)
				g.score++
				break
			}
		}
	}

	/* Check for meteor-player collision */
	for _, m := range g.meteors {
		if m.Collider().IntersectsCircle(g.player.Collider()) {
			g.Reset()
			break
		}
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	switch g.state {
	case GameStateStartMenu:
		g.drawStartMenu(screen)
	case GameStatePlaying:
		g.drawGame(screen)
	}
}

func (g *Game) drawStartMenu(screen *ebiten.Image) {
	text.Draw(screen, "Enter Player Name", assets.ScoreFont, screenWidth/2-350, screenWidth/2-350, color.White)
	text.Draw(screen, "SPACE WARRIOR", assets.ScoreFont, screenWidth/2-350, screenHeight/2-50, color.White)
	text.Draw(screen, "Press SPACE to Start", assets.ScoreFont, screenWidth/2-350, screenHeight/2, color.White)
}

func (g *Game) drawGame(screen *ebiten.Image) {
	/* Draw the background */
	g.background.Draw(screen)

	/* Draw the Player */
	g.player.Draw(screen)

	/* Draw the meteors */
	rand.Shuffle(len(g.meteors), func(i, j int) {
		g.meteors[i], g.meteors[j] = g.meteors[j], g.meteors[i]
	})
	for _, m := range g.meteors {
		m.Draw(screen)
	}

	/* Draw the bullets */
	for _, bullet := range g.bullets {
		bullet.Draw(screen)
	}

	/* Draw the fire */
	for _, fire := range g.fire {
		fire.Draw(screen)
	}

	text.Draw(screen, fmt.Sprintf("%06d", g.score), assets.ScoreFont, screenWidth/2-100, 50, color.White)
}

func (g *Game) AddBullet(b *Bullet) {
	g.bullets = append(g.bullets, b)
}

func (g *Game) AddFire(f *Fire) {
	g.fire = append(g.fire, f)
}

func (g *Game) Reset() {
	g.state = GameStateStartMenu
	g.player = NewPlayer(g)
	g.meteors = nil
	g.bullets = nil
	g.score = 0
	g.meteorSpawnTimer.Reset()
	g.baseVelocity = baseMeteorVelocity
	g.velocityTimer.Reset()
}

func (g *Game) Layout(outsideWidth int, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}
