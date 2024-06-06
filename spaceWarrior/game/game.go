package game

import (
	"fmt"
	"image/color"
	"log"
	"math/rand"
	"os"
	"time"

	"spaceWarrior/assets"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
)

const (
	GameStateStartMenu GameState = iota
	GameStatePlaying
	screenWidth  = 800
	screenHeight = 600

	meteorSpawnTime = 1 * time.Second

	baseMeteorVelocity  = 0.25
	meteorSpeedUpAmount = 0.1
	meteorSpeedUpTime   = 5 * time.Second
)

type GameState int

type Game struct {
	state            GameState
	background       *Background
	player           *Player
	meteorSpawnTimer *Timer
	meteors          []*Meteor
	bullets          []*Bullet
	fire             *Fire

	playerScoresFile *os.File

	playerName      string
	nameInputBuffer string
	score           int

	highestScorer string
	highestScore  int

	baseVelocity  float64
	velocityTimer *Timer

	gameMusic *audio.Player
}

func NewGame() *Game {
	g := &Game{
		state:            GameStateStartMenu,
		meteorSpawnTimer: NewTimer(meteorSpawnTime),
		baseVelocity:     baseMeteorVelocity,
		velocityTimer:    NewTimer(meteorSpeedUpTime),
		playerName:       "",
		nameInputBuffer:  "",
	}
	g.player = NewPlayer(g)
	g.background = NewBackground(g)
	g.gameMusic = assets.Gamemusic

	return g
}

func (g *Game) Update() error {
	switch g.state {
	case GameStateStartMenu:
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			if g.nameInputBuffer != "" {
				g.playerName = g.nameInputBuffer
				g.state = GameStatePlaying
			}
		}
		var name []rune
		for _, r := range ebiten.AppendInputChars(name) {
			g.nameInputBuffer += string(r)
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) {
			if len(g.nameInputBuffer) > 0 {
				g.nameInputBuffer = g.nameInputBuffer[:len(g.nameInputBuffer)-1]
			}
		}

		/* Player Score Data storage and retrieval*/
		_, err := os.Stat("player-scores.txt")
		if os.IsNotExist(err) {
			_, err := os.Create("player-scores.txt")
			if err != nil {
				panic(err)
			}
		}
		g.playerScoresFile, err = os.OpenFile("player-scores.txt", os.O_RDWR, 0644)
		if err != nil {
			panic(err)
		}

		var data []byte
		_, readError := g.playerScoresFile.Read(data)
		if readError != nil {
			log.Fatal(readError)
		}

		_, er := g.playerScoresFile.Write([]byte(fmt.Sprintf("%s %d", g.playerName, g.score)))
		if er != nil {
			log.Fatal(er)
		}
		defer g.playerScoresFile.Close()

	case GameStatePlaying:
		g.updateGame()
	}
	return nil
}

func (g *Game) updateGame() {
	// Play game music
	g.gameMusic.SetVolume(0.5) // Optional: adjust volume
	g.gameMusic.Play()
	go func() {
		c := time.Tick(30 * time.Second)
		for range c {
			g.gameMusic.Rewind()
			g.gameMusic.Play()
		}
	}()

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

	/* Render the fire */
	if g.fire != nil {
		g.fire.Update()
	}

	/* Check for meteor-bullet collision */
	for i := len(g.meteors) - 1; i >= 0; i-- {
		m := g.meteors[i]
		for j := len(g.bullets) - 1; j >= 0; j-- {
			b := g.bullets[j]
			if b.Collider().IntersectsRect(m.Collider()) {

				// Play fireball sound
				assets.Fireball.Rewind()
				assets.Fireball.Play()

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
			if g.score > g.highestScore {
				g.highestScore = g.score
				g.playerName = g.highestScorer

				// Play fireball sound
				assets.Fireball.Rewind()
				assets.Fireball.Play()
			}
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
	text.Draw(screen, "SPACE WARRIOR", assets.ScoreFont, screenWidth/2-350, screenHeight/2-250, color.White)
	text.Draw(screen, "Enter Player Name", assets.ScoreFont, screenWidth/2-350, screenHeight/2-200, color.White)
	text.Draw(screen, g.nameInputBuffer, assets.ScoreFont, screenWidth/2-350, screenHeight/2-150, color.White)
	text.Draw(screen, "Press Enter to Start", assets.ScoreFont, screenWidth/2-350, screenHeight/2-100, color.White)
}

func (g *Game) drawGame(screen *ebiten.Image) {
	/* Draw the background */
	g.background.Draw(screen)

	fps := fmt.Sprintf("FPS: %.f", ebiten.ActualFPS())
	text.Draw(screen, g.playerName, assets.SmallScoreFont, screenWidth/2-380, screenHeight/2-270, color.White)
	text.Draw(screen, fps, assets.SmallScoreFont, screenWidth/2+270, screenHeight/2-270, color.White)

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
	if g.fire != nil && g.fire.isVisible {
		g.fire.Draw(screen)
	}

	text.Draw(screen, fmt.Sprintf("%06d", g.score), assets.ScoreFont, screenWidth/2-100, 50, color.White)
}

func (g *Game) AddBullet(b *Bullet) {
	g.bullets = append(g.bullets, b)
}

func (g *Game) AddFire(f *Fire) {
	g.fire = f
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
