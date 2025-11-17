// ============================================================================
// File: main.go
// Author: Ding Lionel, Darx Christian
// Email: lionel.ding@hes-so.ch, christia.darx@hes-so.ch
// Date: 07.11.2025
// Description: TicTacToe game implementation using Ebiten game engine.
//              Two-player turn-based game with score tracking.
// Version: 1.0
//
// Original Project:
//   Author: LempekPL
//   Source: https://github.com/LempekPL/GoTicTacToe
//
// License: MIT
// Copyright 2025, School of Engineering and Architecture of Fribourg
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
// ============================================================================

package main

import (
	"bytes"
	"embed"
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

const (
	StateInit = iota
	StatePlaying
	StateGameOver

	sWidth      = 480
	sHeight     = 600
	fontSize    = 15
	bigFontSize = 100
	dpi         = 72

	// Board constants
	boardSize      = 3
	cellSize       = 160
	symbolOffset   = 7
	maxRounds      = 8
	keyPressFrames = 60

	// Player symbols
	playerO = "O"
	playerX = "X"
)

//go:embed images/*
var imageFS embed.FS

var (
	normalText font.Face
	bigText    font.Face
	boardImage *ebiten.Image
	gameImage  = ebiten.NewImage(sWidth, sWidth)
	randomGen  = rand.New(rand.NewSource(time.Now().UnixNano()))
)

// Game represents the TicTacToe game state
type Game struct {
	currentPlayer string                       // Currently active player (O or X)
	state         int                          // Current game state (Init, Playing, GameOver)
	gameBoard     [boardSize][boardSize]string // 3x3 board grid
	round         int                          // Current round number (0-8, max 9 moves)
	pointsO       int                          // Score for player O
	pointsX       int                          // Score for player X
	winner        string                       // Winner of current game (empty if no winner yet)
	firstPlayer   int                          // First player indicator: 0 for O, 1 for X
}

// Update handles game logic updates each frame
func (g *Game) Update() error {
	switch g.state {
	case StateInit:
		// Initialize game on first run
		g.Init()

	case StatePlaying:
		// Handle player moves during gameplay
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			mouseX, mouseY := ebiten.CursorPosition()
			// Convert pixel coordinates to board cell coordinates
			// Each cell is 160x160 pixels, so divide to get cell index
			cellX := mouseX / cellSize
			cellY := mouseY / cellSize

			// Check if click is valid (within bounds and cell is empty)
			if g.isValidCell(cellX, cellY) && g.gameBoard[cellX][cellY] == "" {
				// Determine which symbol to place based on current round
				symbol := g.getCurrentSymbol()
				g.DrawSymbol(cellX, cellY, symbol)
				g.gameBoard[cellX][cellY] = symbol
				g.switchPlayer()          // Update display indicator
				g.handleWin(g.CheckWin()) // Check for win condition
				g.round++                 // Increment round counter
			}
		}

	case StateGameOver:
		// Allow restarting game with mouse click
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			g.Load()
		}
	}

	// Global keyboard shortcuts (work in any state)
	// R key: Reset game and clear scores (must hold for 60 frames)
	if inpututil.KeyPressDuration(ebiten.KeyR) == keyPressFrames {
		g.Load()
		g.ResetPoints()
	}
	// Escape key: Exit game (must hold for 60 frames)
	if inpututil.KeyPressDuration(ebiten.KeyEscape) == keyPressFrames {
		os.Exit(0)
	}
	return nil
}

// drawKeyPressFeedback displays visual feedback when keys are held down
// The color fades from bright to dark as the key is held longer
func drawKeyPressFeedback(key ebiten.Key, screen *ebiten.Image) {
	duration := inpututil.KeyPressDuration(key)
	if duration > 1 {
		var msgText string
		var colorText color.RGBA
		// Calculate color fade: starts at 255, decreases as duration approaches keyPressFrames
		// This creates a visual countdown effect
		colorChange := 255 - (255 / keyPressFrames * uint8(duration))
		if key == ebiten.KeyEscape {
			msgText = "CLOSING..."
			// Red fades to black (G and B decrease)
			colorText = color.RGBA{R: 255, G: colorChange, B: colorChange, A: 255}
		} else if key == ebiten.KeyR {
			msgText = "RESETTING..."
			// Cyan fades to black (R decreases)
			colorText = color.RGBA{R: colorChange, G: 255, B: 255, A: 255}
		}
		text.Draw(screen, msgText, normalText, sWidth/2, sHeight-30, colorText)
	}
}

// Draw renders the game screen each frame
func (g *Game) Draw(screen *ebiten.Image) {
	// Draw background board and game pieces
	screen.DrawImage(boardImage, nil)
	screen.DrawImage(gameImage, nil) // Contains all X and O symbols
	mouseX, mouseY := ebiten.CursorPosition()

	// Display performance metrics (top-left)
	msgFPS := fmt.Sprintf("TPS: %0.2f\nFPS: %0.2f", ebiten.ActualTPS(), ebiten.ActualFPS())
	text.Draw(screen, msgFPS, normalText, 0, sHeight-30, color.White)

	// Show key press feedback (if keys are being held)
	drawKeyPressFeedback(ebiten.KeyEscape, screen)
	drawKeyPressFeedback(ebiten.KeyR, screen)

	// Display current scores (bottom-center)
	msgOX := fmt.Sprintf("O: %v | X: %v", g.pointsO, g.pointsX)
	text.Draw(screen, msgOX, normalText, sWidth/2, sHeight-5, color.White)

	// Display winner message if game is over (large text in center)
	if g.winner != "" {
		msgWin := fmt.Sprintf("%v wins!", g.winner)
		text.Draw(screen, msgWin, bigText, 70, 200, color.RGBA{G: 50, B: 200, A: 255})
	}

	// Display current player indicator at mouse cursor position
	msg := fmt.Sprintf("%v", g.currentPlayer)
	text.Draw(screen, msg, normalText, mouseX, mouseY, color.RGBA{G: 255, A: 255})
}

// loadImageFromFS loads an image from the embedded filesystem
func loadImageFromFS(path string) (*ebiten.Image, error) {
	imageBytes, err := imageFS.ReadFile(path)
	if err != nil {
		return nil, err
	}
	decoded, _, err := image.Decode(bytes.NewReader(imageBytes))
	if err != nil {
		return nil, err
	}
	return ebiten.NewImageFromImage(decoded), nil
}

// DrawSymbol draws a player symbol (X or O) at the specified board cell
// x, y are cell coordinates (0-2), not pixel coordinates
func (g *Game) DrawSymbol(x, y int, symbol string) {
	symbolImage, err := loadImageFromFS(fmt.Sprintf("images/%s.png", symbol))
	if err != nil {
		log.Fatal(err)
	}
	opSymbol := &ebiten.DrawImageOptions{}
	// Convert cell coordinates to pixel coordinates
	// Multiply by cellSize to get cell position, then add offset for centering
	xPos := float64(x*cellSize + symbolOffset)
	yPos := float64(y*cellSize + symbolOffset)
	opSymbol.GeoM.Translate(xPos, yPos)

	// Draw to gameImage (persistent layer) so symbol stays after frame ends
	gameImage.DrawImage(symbolImage, opSymbol)
}

// Init initializes the game on first startup
// Loads the board image and randomly selects which player goes first
func (g *Game) Init() {
	var err error
	boardImage, err = loadImageFromFS("images/board.png")
	if err != nil {
		log.Fatal(err)
	}

	// Randomly choose first player for the initial game
	// Subsequent games will alternate (handled in Load())
	randomChoice := randomGen.Intn(2)
	if randomChoice == 0 {
		g.currentPlayer = playerO
		g.firstPlayer = 0
	} else {
		g.currentPlayer = playerX
		g.firstPlayer = 1
	}
	g.Load()        // Set up the first game board
	g.ResetPoints() // Initialize scores to zero
}

// Load resets the game board for a new round while keeping scores
// Alternates the starting player so each player gets to go first
func (g *Game) Load() {
	gameImage.Clear()                            // Clear all drawn symbols
	g.gameBoard = [boardSize][boardSize]string{} // Reset board to empty
	g.round = 0
	// Alternate starting player: if O started last game, X starts this game
	// This ensures fair play between rounds
	if g.firstPlayer == 0 {
		g.currentPlayer = playerX
		g.firstPlayer = 1
	} else {
		g.currentPlayer = playerO
		g.firstPlayer = 0
	}
	g.winner = ""
	g.state = StatePlaying
}

func (g *Game) handleWin(winner string) {
	if winner == playerO {
		g.winner = playerO
		g.pointsO++
		g.state = StateGameOver
	} else if winner == playerX {
		g.winner = playerX
		g.pointsX++
		g.state = StateGameOver
	} else if winner == "tie" {
		g.winner = "No one"
		g.state = StateGameOver
	}
}

// CheckWin checks if there's a winner or tie condition
// Returns: "O" or "X" for winner, "tie" for draw, "" for no winner yet
func (g *Game) CheckWin() string {
	// Check rows: look for three identical symbols in a horizontal line
	for row := 0; row < boardSize; row++ {
		if g.gameBoard[row][0] != "" &&
			g.gameBoard[row][0] == g.gameBoard[row][1] &&
			g.gameBoard[row][1] == g.gameBoard[row][2] {
			return g.gameBoard[row][0]
		}
	}

	// Check columns: look for three identical symbols in a vertical line
	for col := 0; col < boardSize; col++ {
		if g.gameBoard[0][col] != "" &&
			g.gameBoard[0][col] == g.gameBoard[1][col] &&
			g.gameBoard[1][col] == g.gameBoard[2][col] {
			return g.gameBoard[0][col]
		}
	}

	// Check diagonals: only need to check if center is filled
	center := g.gameBoard[1][1]
	if center != "" {
		// Main diagonal (top-left to bottom-right): [0,0] [1,1] [2,2]
		if g.gameBoard[0][0] == center && center == g.gameBoard[2][2] {
			return center
		}
		// Anti-diagonal (top-right to bottom-left): [0,2] [1,1] [2,0]
		if g.gameBoard[0][2] == center && center == g.gameBoard[2][0] {
			return center
		}
	}

	// Check for tie: after 8 rounds (0-8), all 9 cells are filled with no winner
	// Round 0 = first move, round 8 = ninth move (board full)
	if g.round == maxRounds {
		return "tie"
	}

	return ""
}

func (g *Game) ResetPoints() {
	g.pointsO = 0
	g.pointsX = 0
}

func init() {
	tt, err := opentype.Parse(fonts.MPlus1pRegular_ttf)
	if err != nil {
		log.Fatal(err)
	}
	normalText, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    fontSize,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
	bigText, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    bigFontSize,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
}

// Helper methods for Game

func (g *Game) isValidCell(x, y int) bool {
	return x >= 0 && x < boardSize && y >= 0 && y < boardSize
}

// getCurrentSymbol determines which symbol should be placed based on the current round
// Uses modulo arithmetic: (round + firstPlayer) % 2 determines whose turn it is
// If firstPlayer is 0 (O starts), rounds 0,2,4,6,8 are O's turns
// If firstPlayer is 1 (X starts), rounds 1,3,5,7 are O's turns (X goes first)
func (g *Game) getCurrentSymbol() string {
	if (g.round+g.firstPlayer)%2 == 0 {
		return playerO
	}
	return playerX
}

func (g *Game) switchPlayer() {
	if g.currentPlayer == playerO {
		g.currentPlayer = playerX
	} else {
		g.currentPlayer = playerO
	}
}

func (g *Game) Layout(int, int) (int, int) {
	return sWidth, sHeight
}

func main() {
	game := &Game{}
	ebiten.SetWindowSize(sWidth, sHeight)
	ebiten.SetWindowTitle("TicTacToe")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
