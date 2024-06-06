package assets

import (
	"embed"
	"image"
	_ "image/png"
	"io/fs"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

const sampleRate = 44100

//go:embed player.png meteors/*.png laser.png font.ttf background.png fire.png gunshot.mp3 fireball.mp3 gamemusic.mp3
var assets embed.FS
var PlayerSprite = mustLoadImage("player.png")
var MeteorSprites = mustLoadImages("meteors/*.png")
var LaserSprite = mustLoadImage("laser.png")
var ScoreFont = mustLoadFont("font.ttf", 48, 72)
var SmallScoreFont = mustLoadFont("font.ttf", 24, 72)
var BackgoundSprite = mustLoadImage("background.png")
var FireSprite = mustLoadImage("fire.png")

var audioContext = audio.NewContext(44100)
var GunShot = mustLoadAudio("gunshot.mp3")
var Gamemusic = mustLoadAudio("gamemusic.mp3")
var Fireball = mustLoadAudio("fireball.mp3")

func mustLoadAudio(name string) *audio.Player {
	f, err := assets.Open(name)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	d, err := mp3.Decode(audioContext, f)
	if err != nil {
		panic(err)
	}

	p, err := audio.NewPlayer(audioContext, d)
	if err != nil {
		panic(err)
	}
	return p
}

func mustLoadImage(name string) *ebiten.Image {
	f, err := assets.Open(name)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	image, _, err := image.Decode(f)
	if err != nil {
		panic(err)
	}
	return ebiten.NewImageFromImage(image)
}

func mustLoadImages(path string) []*ebiten.Image {
	matches, err := fs.Glob(assets, path)
	if err != nil {
		panic(err)
	}

	images := make([]*ebiten.Image, len(matches))
	for i, match := range matches {
		images[i] = mustLoadImage(match)
	}

	return images
}

func mustLoadFont(name string, size float64, dpi float64) font.Face {
	f, err := assets.ReadFile(name)
	if err != nil {
		panic(err)
	}

	tt, err := opentype.Parse(f)
	if err != nil {
		panic(err)
	}

	face, err := opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    size,
		DPI:     dpi,
		Hinting: font.HintingVertical,
	})
	if err != nil {
		panic(err)
	}

	return face
}
