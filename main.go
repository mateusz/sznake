package main

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/text"
	"log"
	"image/color"
	"math"
	"math/rand"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/gomono"
	"fmt"
)
type Point struct {
	X int
	Y int
}

type Snake struct {
	Body   []Point
	DX     int
	DY     int
	Length int
	LifeRemaining int
	LifeMax int
	HowDead int
	Handicap int
}

type Apple struct {
	Position Point
}

const(
	DotSize = 10
	AreaWidth = 100
	AreaHeight = 80
)

var Snake1 Snake

var Apple1 Apple

var SnakeImage1 *ebiten.Image
var SnakeImage2 *ebiten.Image
var AppleImage *ebiten.Image

var Font font.Face

func update(screen *ebiten.Image) error {
	if ebiten.IsKeyPressed(ebiten.KeyR) {
		reboot()
	}

	dirHoriz := Snake1.DX != 0
	dirVert := Snake1.DY != 0
	if dirHoriz && ebiten.IsKeyPressed(ebiten.KeyDown) {
		Snake1.DX = 0
		Snake1.DY = 1
	}
	if dirHoriz && ebiten.IsKeyPressed(ebiten.KeyUp) {
		Snake1.DX = 0
		Snake1.DY = -1
	}
	if dirVert && ebiten.IsKeyPressed(ebiten.KeyLeft) {
		Snake1.DX = -1
		Snake1.DY = 0
	}
	if dirVert && ebiten.IsKeyPressed(ebiten.KeyRight) {
		Snake1.DX = 1
		Snake1.DY = 0
	}

	if Snake1.HowDead>=10000 {
		text.Draw(screen, "Died from starvation. R to restart.", Font, 10, 30, color.White)
	} else if Snake1.HowDead>=Snake1.Handicap {
		text.Draw(screen, "Died from running into things. R to restart.", Font, 10, 30, color.White)
	} else if Snake1.HowDead>0 {
		text.Draw(screen, "Nearly dead - rescue it, quickly!", Font, 10, 30, color.White)
	}

	text.Draw(screen, fmt.Sprintf("%d", Snake1.Length-10), Font, 10, AreaHeight*DotSize - 10, color.White)

	if Snake1.HowDead < Snake1.Handicap {
		head := Snake1.Body[len(Snake1.Body) - 1]
		newHead := Point{
			X:head.X + Snake1.DX,
			Y:head.Y + Snake1.DY,
		}
		newDead := Snake1.HowDead

		if newHead.OutOfBounds(0, 0, AreaWidth - 1, AreaHeight - 1) {
			newDead += 1
		}

		for i, seg := range Snake1.Body {
			if i != len(Snake1.Body) - 1 && len(Snake1.Body) > 1 && newHead.Same(seg) {
				newDead += 1
			}
		}

		if Snake1.LifeRemaining > 0 {
			if Snake1.HowDead == 0 {
				Snake1.LifeRemaining--
			}
		} else {
			newDead += 10000
		}

		if newDead == Snake1.HowDead {
			// Rescued
			Snake1.HowDead = 0
		} else {
			Snake1.HowDead = newDead
		}

		if newHead.Near(Apple1.Position, 1) {
			Snake1.Length = int(float64(Snake1.Length) * 1.3)
			Snake1.LifeRemaining = Snake1.LifeMax
			Apple1.Generate()
		}

		if Snake1.HowDead < 1 {
			Snake1.Body = append(Snake1.Body, newHead)
			l := len(Snake1.Body)
			begin := 0
			if l >= Snake1.Length {
				begin = l - Snake1.Length
			}

			Snake1.Body = Snake1.Body[begin:l]
		}
	}

	alpha := (float32(Snake1.LifeRemaining) / float32(Snake1.LifeMax)) * 255
	bodyColour1 := color.RGBA{R:0, G:255, B:0, A:uint8(alpha)}
	bodyColour2 := color.RGBA{R:0, G:150, B:0, A:uint8(alpha/2)}
	if Snake1.HowDead >= Snake1.Handicap {
		bodyColour1 = color.RGBA{R:100, G:100, B:100, A:uint8(alpha)}
		bodyColour2 = color.RGBA{R:50, G:50, B:50, A:uint8(alpha/2)}
	} else if Snake1.HowDead > 0 {
		deadness := uint8((float64(Snake1.HowDead) / float64(Snake1.Handicap)) * 100)
		bodyColour1 = color.RGBA{R:255 - deadness, G:255 - deadness, B:0, A:uint8(alpha)}
		bodyColour2 = color.RGBA{R:155 - deadness, G:155 - deadness, B:0, A:uint8(alpha/2)}
	}

	SnakeImage1.Fill(bodyColour1)
	SnakeImage2.Fill(bodyColour2)
	for _, seg := range Snake1.Body {
		var snakeImage *ebiten.Image
		if (seg.X+seg.Y)%2==0 {
			snakeImage = SnakeImage1
		} else {
			snakeImage = SnakeImage2
		}
		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Translate(float64(DotSize*seg.X), float64(DotSize*seg.Y))
		screen.DrawImage(snakeImage, opts)
	}

	AppleImage.Fill(color.RGBA{R:255,G:0,B:0,A:uint8(100+rand.Int()%155)})
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(DotSize*(float64(Apple1.Position.X)-0.5), DotSize*(float64(Apple1.Position.Y)-0.5))
	screen.DrawImage(AppleImage, opts)

	return nil
}

func init() {

	SnakeImage1, _ = ebiten.NewImage(DotSize, DotSize, ebiten.FilterNearest)
	SnakeImage2, _ = ebiten.NewImage(DotSize, DotSize, ebiten.FilterNearest)
	AppleImage, _ = ebiten.NewImage(DotSize*2, DotSize*2, ebiten.FilterNearest)

	tt, err := truetype.Parse(gomono.TTF)
	if err != nil {
		log.Fatal(err)
	}

	const dpi = 72
	Font = truetype.NewFace(tt, &truetype.Options{
		Size:    24,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
}

func main() {
	reboot();
	if err := ebiten.Run(update, AreaWidth*DotSize, AreaHeight*DotSize, 1.0, "Sznake"); err != nil {
		log.Fatal(err)
	}
}

func reboot() {
	Snake1 = Snake{
		Body: []Point{
			{X:int(AreaWidth/2),Y:int(AreaHeight/2)},
		},
		Length: 10,
		LifeRemaining: 500,
		LifeMax: 500,
		DX: 1,
		Handicap: 100,
	}
	Apple1.Generate()
}

func (p1 *Point) Same(p2 Point) bool {
	return p1.X==p2.X && p1.Y==p2.Y
}

func (p1 *Point) Near(p2 Point, dist int) bool {
	distX := math.Abs(float64(p1.X-p2.X))
	distY := math.Abs(float64(p1.Y-p2.Y))

	return distX<=float64(dist) && distY<=float64(dist)
}

func (p *Point) OutOfBounds(left, top, right, bottom int) bool {
	if p.X>right {
		return true
	}
	if p.X<left {
		return true
	}
	if p.Y>bottom {
		return true
	}
	if p.Y<top {
		return true
	}

	return false
}

func (a *Apple) Generate() {
	Apple1 = Apple{
		Position: Point{
			X: rand.Int()%AreaWidth,
			Y: rand.Int()%AreaHeight,
		},
	}
}
