package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"image"
	"image/color"
	_ "image/png"
	"log"
	"os"
)

const (
	sWidth   = 480
	sHeight  = 600
	fontSize = 15
	dpi      = 72
)

var (
	mplusNormalFont font.Face
	round           int
	reset           int
	canvasImage     = ebiten.NewImage(sWidth, sHeight)
)

type Game struct {
	playing string
	state int
	gameBoard [3][3]int
	round int
}

func (g *Game) Update() error {
	switch g.state {
		case 0:
			g.Load()
			break
		case 1:
			mx, my := ebiten.CursorPosition()
			if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
				if mx/160 < 3 && mx >= 0 && my/160 < 3 && my >= 0 {
					g.round++
					if g.round%2 == 0 {
						g.playing = "O"
					} else {
						g.playing = "X"
					}
				}
			}
			fmt.Printf("X: %v, Y: %v \n", mx/160, my/160)
			break
	}
	if inpututil.KeyPressDuration(ebiten.KeyEscape) == 60 {
		os.Exit(0)
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {

	screen.DrawImage(canvasImage, nil)
	mx, my := ebiten.CursorPosition()

	msgTFPS := fmt.Sprintf("TPS: %0.2f\nFPS: %0.2f", ebiten.CurrentTPS(), ebiten.CurrentFPS())
	text.Draw(screen, msgTFPS, mplusNormalFont, 0, sHeight-30, color.White)
	if inpututil.KeyPressDuration(ebiten.KeyEscape) > 1 {
		msgClosing := fmt.Sprintf("CLOSING...")
		colorChangeToExit := 255-(255/60*uint8(inpututil.KeyPressDuration(ebiten.KeyEscape)))
		text.Draw(screen, msgClosing, mplusNormalFont, sWidth/2, sHeight-30, color.RGBA{R: 255, G: colorChangeToExit, B: colorChangeToExit, A: 255})
	}
	msgXY := fmt.Sprintf("(%v, %v)", mx, my)
	text.Draw(screen, msgXY, mplusNormalFont, sWidth/2, sHeight-5, color.White)
	msg := fmt.Sprintf("%v", g.playing)
	text.Draw(screen, msg, mplusNormalFont, mx, my, color.RGBA{G: 255, A: 255})
}

func (g *Game) DrawBoard() {
	picFile, _ := ebitenutil.OpenFile("images/board.png")
	img, _, err := image.Decode(picFile)
	if err != nil {
		log.Fatal(err)
	}
	canvasImage = ebiten.NewImageFromImage(img)
}

func (g *Game) DrawSymbol(symbol string) {
	fil := fmt.Sprintf("images/%v.png", symbol)
	picFile, _ := ebitenutil.OpenFile(fil)
	img, _, err := image.Decode(picFile)
	if err != nil {
		log.Fatal(err)
	}
	canvasImage = ebiten.NewImageFromImage(img)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return sWidth, sHeight
}

func (g *Game) Load() {
	g.DrawBoard()
	g.state = 1
	g.gameBoard[0][0], g.gameBoard[0][1], g.gameBoard[0][2] = 0,0,0
	g.gameBoard[1][0], g.gameBoard[1][1], g.gameBoard[1][2] = 0,0,0
	g.gameBoard[2][0], g.gameBoard[2][1], g.gameBoard[2][2] = 0,0,0
	g.playing = "O"
	g.round = 0

}

func init() {
	tt, err := opentype.Parse(fonts.MPlus1pRegular_ttf)
	if err != nil {
		log.Fatal(err)
	}
	mplusNormalFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    fontSize,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	game := &Game{}
	ebiten.SetWindowSize(sWidth, sHeight)
	ebiten.SetWindowTitle("TicTacToe")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
