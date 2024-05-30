package game

import (
	"spaceWarrior/assets"

	"github.com/hajimehoshi/ebiten/v2"
)

type Background struct {
	sprite *ebiten.Image
}

func NewBackground(g *Game) *Background {
	sprite := assets.BackgoundSprite

	return &Background{
		sprite: sprite,
	}
}

func (b *Background) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}

	screenWidth, screenHeight := screen.Size()

	bounds := b.sprite.Bounds()
	imgWidth := bounds.Dx()
	imgHeight := bounds.Dy()

	scaleX := float64(screenWidth) / float64(imgWidth)
	scaleY := float64(screenHeight) / float64(imgHeight)

	op.GeoM.Scale(scaleX, scaleY)

	screen.DrawImage(b.sprite, op)
}
