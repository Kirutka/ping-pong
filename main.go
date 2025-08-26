package main

import (
	"fmt"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
)

const (
	screenWidth  = 800
	screenHeight = 600
	paddleWidth  = 15
	paddleHeight = 100
	ballSize     = 15
	paddleSpeed  = 8
	ballSpeed    = 5
)

type GameState int

const (
	MenuState GameState = iota
	PlayingState
	PauseState
)

type Button struct {
	x, y, width, height float32
	text                string
	color               color.RGBA
	hoverColor          color.RGBA
	isHovered           bool
}

type Game struct {
	state         GameState
	leftScore     int
	rightScore    int
	leftPaddleY   float64
	rightPaddleY  float64
	ballX         float64
	ballY         float64
	ballVelX      float64
	ballVelY      float64
	font          font.Face
	playButton    Button
	exitButton    Button
	resumeButton  Button
	menuButton    Button
	escapePressed bool
	mousePressed  bool
}

func NewGame() *Game {
	game := &Game{
		state:        MenuState,
		leftPaddleY:  screenHeight/2 - paddleHeight/2,
		rightPaddleY: screenHeight/2 - paddleHeight/2,
		ballX:        screenWidth/2 - ballSize/2,
		ballY:        screenHeight/2 - ballSize/2,
		ballVelX:     ballSpeed,
		ballVelY:     ballSpeed * 0.7,
		font:         basicfont.Face7x13,
	}

	game.playButton = Button{
		x:          screenWidth/2 - 80,
		y:          screenHeight/2 + 50,
		width:      160,
		height:     40,
		text:       "PLAY",
		color:      color.RGBA{100, 200, 100, 255},
		hoverColor: color.RGBA{120, 220, 120, 255},
	}

	game.exitButton = Button{
		x:          screenWidth/2 - 80,
		y:          screenHeight/2 + 100,
		width:      160,
		height:     40,
		text:       "EXIT",
		color:      color.RGBA{200, 100, 100, 255},
		hoverColor: color.RGBA{220, 120, 120, 255},
	}

	game.resumeButton = Button{
		x:          screenWidth/2 - 80,
		y:          screenHeight/2 - 50,
		width:      160,
		height:     40,
		text:       "RESUME",
		color:      color.RGBA{100, 200, 100, 255},
		hoverColor: color.RGBA{120, 220, 120, 255},
	}

	game.menuButton = Button{
		x:          screenWidth/2 - 80,
		y:          screenHeight/2 + 20,
		width:      160,
		height:     40,
		text:       "MAIN MENU",
		color:      color.RGBA{200, 100, 100, 255},
		hoverColor: color.RGBA{220, 120, 120, 255},
	}

	return game
}

func (g *Game) Update() error {
	mouseX, mouseY := ebiten.CursorPosition()
	mousePressed := ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)
	ebiten.IsKeyPressed(ebiten.KeySpace)
	escapePressed := ebiten.IsKeyPressed(ebiten.KeyEscape)

	g.updateButtonHover(float32(mouseX), float32(mouseY))

	switch g.state {
	case MenuState:
		return g.updateMenu(float32(mouseX), float32(mouseY), mousePressed)
	case PlayingState:
		return g.updatePlaying(escapePressed)
	case PauseState:
		return g.updatePause(float32(mouseX), float32(mouseY), mousePressed, escapePressed)
	}
	return nil
}

func (g *Game) updateMenu(mouseX, mouseY float32, mousePressed bool) error {
	if mousePressed && !g.mousePressed {
		if g.isInButton(mouseX, mouseY, g.playButton) {
			g.resetGame()
			g.state = PlayingState
		}
		if g.isInButton(mouseX, mouseY, g.exitButton) {
		}
	}
	g.mousePressed = mousePressed
	return nil
}

func (g *Game) updatePlaying(escapePressed bool) error {
	if escapePressed && !g.escapePressed {
		g.state = PauseState
	}
	g.escapePressed = escapePressed

	// Paddle control
	if ebiten.IsKeyPressed(ebiten.KeyW) && g.leftPaddleY > 0 {
		g.leftPaddleY -= paddleSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) && g.leftPaddleY < screenHeight-paddleHeight {
		g.leftPaddleY += paddleSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) && g.rightPaddleY > 0 {
		g.rightPaddleY -= paddleSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) && g.rightPaddleY < screenHeight-paddleHeight {
		g.rightPaddleY += paddleSpeed
	}

	g.ballX += g.ballVelX
	g.ballY += g.ballVelY

	if g.ballY < 0 || g.ballY > screenHeight-ballSize {
		g.ballVelY = -g.ballVelY
	}

	if g.checkPaddleCollision() {
		g.ballVelX = -g.ballVelX * 1.05
	}

	if g.ballX < 0 {
		g.rightScore++
		g.resetBall()
	}
	if g.ballX > screenWidth-ballSize {
		g.leftScore++
		g.resetBall()
	}

	return nil
}

func (g *Game) updatePause(mouseX, mouseY float32, mousePressed, escapePressed bool) error {
	if escapePressed && !g.escapePressed {
		g.state = PlayingState
	}
	g.escapePressed = escapePressed

	if mousePressed && !g.mousePressed {
		if g.isInButton(mouseX, mouseY, g.resumeButton) {
			g.state = PlayingState
		}
		if g.isInButton(mouseX, mouseY, g.menuButton) {
			g.state = MenuState
		}
	}
	g.mousePressed = mousePressed
	return nil
}

func (g *Game) updateButtonHover(mouseX, mouseY float32) {
	g.playButton.isHovered = g.isInButton(mouseX, mouseY, g.playButton)
	g.exitButton.isHovered = g.isInButton(mouseX, mouseY, g.exitButton)
	g.resumeButton.isHovered = g.isInButton(mouseX, mouseY, g.resumeButton)
	g.menuButton.isHovered = g.isInButton(mouseX, mouseY, g.menuButton)
}

func (g *Game) Draw(screen *ebiten.Image) {
	switch g.state {
	case MenuState:
		g.drawMenu(screen)
	case PlayingState:
		g.drawPlaying(screen)
	case PauseState:
		g.drawPlaying(screen)
		g.drawPause(screen)
	}
}

func (g *Game) Layout(_, _ int) (int, int) {
	return screenWidth, screenHeight
}

func (g *Game) drawMenu(screen *ebiten.Image) {
	vector.DrawFilledRect(screen, 0, 0, screenWidth, screenHeight, color.RGBA{30, 30, 60, 255}, false)

	text.Draw(screen, "PING-PONG", g.font, screenWidth/2-55, 100, color.White)

	text.Draw(screen, "Controls", g.font, screenWidth/2-55, 180, color.White)
	text.Draw(screen, "Left paddle: W/S", g.font, screenWidth/2-80, 210, color.RGBA{255, 100, 100, 255})
	text.Draw(screen, "Right paddle: Arrow Keys", g.font, screenWidth/2-80, 240, color.RGBA{100, 100, 255, 255})
	text.Draw(screen, "Pause: ESC", g.font, screenWidth/2-80, 270, color.White)

	g.drawButton(screen, g.playButton)
	g.drawButton(screen, g.exitButton)
}

func (g *Game) drawPlaying(screen *ebiten.Image) {
	// Background
	vector.DrawFilledRect(screen, 0, 0, screenWidth, screenHeight, color.RGBA{20, 20, 40, 255}, false)

	// Net
	for y := 0; y < screenHeight; y += 20 {
		vector.DrawFilledRect(screen, screenWidth/2-2, float32(y), 4, 10, color.RGBA{100, 100, 100, 255}, false)
	}

	// Paddles
	vector.DrawFilledRect(screen, 20, float32(g.leftPaddleY), paddleWidth, paddleHeight, color.RGBA{255, 80, 80, 255}, false)
	vector.DrawFilledRect(screen, screenWidth-20-paddleWidth, float32(g.rightPaddleY), paddleWidth, paddleHeight, color.RGBA{80, 80, 255, 255}, false)

	// Ball
	vector.DrawFilledRect(screen, float32(g.ballX), float32(g.ballY), ballSize, ballSize, color.RGBA{255, 255, 100, 255}, false)

	// Score
	text.Draw(screen, fmt.Sprintf("%d", g.leftScore), g.font, 50, 50, color.RGBA{255, 100, 100, 255})
	text.Draw(screen, fmt.Sprintf("%d", g.rightScore), g.font, screenWidth-70, 50, color.RGBA{100, 100, 255, 255})
}

func (g *Game) drawPause(screen *ebiten.Image) {
	vector.DrawFilledRect(screen, 0, 0, screenWidth, screenHeight, color.RGBA{0, 0, 0, 180}, false)

	// Pause text
	text.Draw(screen, "PAUSE", g.font, screenWidth/2-40, screenHeight/2-120, color.White)

	// Buttons
	g.drawButton(screen, g.resumeButton)
	g.drawButton(screen, g.menuButton)
}

func (g *Game) drawButton(screen *ebiten.Image, btn Button) {
	btnColor := btn.color
	if btn.isHovered {
		btnColor = btn.hoverColor
	}

	// Draw button background
	vector.DrawFilledRect(screen, btn.x, btn.y, btn.width, btn.height, btnColor, false)

	// Draw button border
	vector.StrokeRect(screen, btn.x, btn.y, btn.width, btn.height, 2, color.White, false)

	// Draw button text
	textWidth := len(btn.text) * 7
	textX := btn.x + (btn.width-float32(textWidth))/2
	textY := btn.y + btn.height/2 + 5
	text.Draw(screen, btn.text, g.font, int(textX), int(textY), color.White)
}

func (g *Game) resetGame() {
	g.leftScore = 0
	g.rightScore = 0
	g.resetBall()
	g.leftPaddleY = screenHeight/2 - paddleHeight/2
	g.rightPaddleY = screenHeight/2 - paddleHeight/2
}

func (g *Game) resetBall() {
	g.ballX = screenWidth/2 - ballSize/2
	g.ballY = screenHeight/2 - ballSize/2
	g.ballVelX = ballSpeed
	g.ballVelY = ballSpeed * 0.7
}

func (g *Game) checkPaddleCollision() bool {
	// Left paddle
	if g.ballX <= 20+paddleWidth && g.ballY+ballSize >= g.leftPaddleY && g.ballY <= g.leftPaddleY+paddleHeight {
		return true
	}
	// Right paddle
	if g.ballX >= screenWidth-20-paddleWidth-ballSize && g.ballY+ballSize >= g.rightPaddleY && g.ballY <= g.rightPaddleY+paddleHeight {
		return true
	}
	return false
}

func (g *Game) isInButton(x, y float32, btn Button) bool {
	return x >= btn.x && x <= btn.x+btn.width && y >= btn.y && y <= btn.y+btn.height
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Ping-Pong Game")
	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}
