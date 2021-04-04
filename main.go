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
	boardImage      = ebiten.NewImage(sWidth, sHeight)
	symbolImage     *ebiten.Image
	gameImage       = ebiten.NewImage(sWidth, sWidth)
)

type Game struct {
	playing   string
	state     int
	gameBoard [3][3]int
	round     int
}

func (g *Game) Update() error {
	switch g.state {
	case 0:
		g.Load()
		break
	case 1:
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			mx, my := ebiten.CursorPosition()
			if mx/160 < 3 && mx >= 0 && my/160 < 3 && my >= 0 && g.gameBoard[mx/160][my/160] == 0 {
				g.round++
				if g.round%2 == 0 {
					g.DrawSymbol(mx/160, my/160, "O")
					g.gameBoard[mx/160][my/160] = 1
					g.playing = "X"
				} else {
					g.DrawSymbol(mx/160, my/160, "X")
					g.gameBoard[mx/160][my/160] = 2
					g.playing = "O"
				}
			}
		}
		break
	}
	if inpututil.KeyPressDuration(ebiten.KeyEscape) == 60 {
		os.Exit(0)
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {

	screen.DrawImage(boardImage, nil)
	screen.DrawImage(gameImage, nil)
	mx, my := ebiten.CursorPosition()

	msgTFPS := fmt.Sprintf("TPS: %0.2f\nFPS: %0.2f", ebiten.CurrentTPS(), ebiten.CurrentFPS())
	text.Draw(screen, msgTFPS, mplusNormalFont, 0, sHeight-30, color.White)
	if inpututil.KeyPressDuration(ebiten.KeyEscape) > 1 {
		msgClosing := fmt.Sprintf("CLOSING...")
		colorChangeToExit := 255 - (255 / 60 * uint8(inpututil.KeyPressDuration(ebiten.KeyEscape)))
		text.Draw(screen, msgClosing, mplusNormalFont, sWidth/2, sHeight-30, color.RGBA{R: 255, G: colorChangeToExit, B: colorChangeToExit, A: 255})
	}
	msgXY := fmt.Sprintf("(%v, %v)", mx, my)
	text.Draw(screen, msgXY, mplusNormalFont, sWidth/2, sHeight-5, color.White)
	msg := fmt.Sprintf("%v", g.playing)
	text.Draw(screen, msg, mplusNormalFont, mx, my, color.RGBA{G: 255, A: 255})
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return sWidth, sHeight
}

func (g *Game) DrawSymbol(x, y int, sym string) {
	fileName := fmt.Sprintf("images/%v.png", sym)
	var err error
	symbolImage, _, err = ebitenutil.NewImageFromFile(fileName)
	if err != nil {
		log.Fatal(err)
	}
	opSymbol := &ebiten.DrawImageOptions{}
	opSymbol.GeoM.Translate(float64((160*(x+1)-160)+7), float64((160*(y+1)-160)+7))

	gameImage.DrawImage(symbolImage, opSymbol)
}

func (g *Game) Load() {

	var err error
	boardImage, _, err = ebitenutil.NewImageFromFile("images/board.png")
	if err != nil {
		log.Fatal(err)
	}
	g.state = 1
	g.gameBoard[0][0], g.gameBoard[0][1], g.gameBoard[0][2] = 0, 0, 0
	g.gameBoard[1][0], g.gameBoard[1][1], g.gameBoard[1][2] = 0, 0, 0
	g.gameBoard[2][0], g.gameBoard[2][1], g.gameBoard[2][2] = 0, 0, 0
	g.playing = "O"
	g.round = 1

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
