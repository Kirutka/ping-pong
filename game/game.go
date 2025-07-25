package game

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"ping-pong/entities"
)

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
	playButton    entities.Button
	exitButton    entities.Button
	resumeButton  entities.Button
	menuButton    entities.Button
	escapePressed bool
	mousePressed  bool
}

func NewGame() *Game {
	game := &Game{
		state:        MenuState,
		leftPaddleY:  ScreenHeight/2 - PaddleHeight/2,
		rightPaddleY: ScreenHeight/2 - PaddleHeight/2,
		ballX:        ScreenWidth/2 - BallSize/2,
		ballY:        ScreenHeight/2 - BallSize/2,
		ballVelX:     BallSpeed,
		ballVelY:     BallSpeed * 0.7,
		font:         basicfont.Face7x13,
	}

	game.playButton = entities.NewButton(
		ScreenWidth/2-80,
		ScreenHeight/2+50,
		160,
		40,
		"PLAY",
		color.RGBA{100, 200, 100, 255},
		color.RGBA{120, 220, 120, 255},
	)

	game.exitButton = entities.NewButton(
		ScreenWidth/2-80,
		ScreenHeight/2+100,
		160,
		40,
		"EXIT",
		color.RGBA{200, 100, 100, 255},
		color.RGBA{220, 120, 120, 255},
	)

	game.resumeButton = entities.NewButton(
		ScreenWidth/2-80,
		ScreenHeight/2-50,
		160,
		40,
		"RESUME",
		color.RGBA{100, 200, 100, 255},
		color.RGBA{120, 220, 120, 255},
	)

	game.menuButton = entities.NewButton(
		ScreenWidth/2-80,
		ScreenHeight/2+20,
		160,
		40,
		"MAIN MENU",
		color.RGBA{200, 100, 100, 255},
		color.RGBA{220, 120, 120, 255},
	)

	return game
}

func (g *Game) Update() error {
	mouseX, mouseY := ebiten.CursorPosition()
	mousePressed := ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)
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
		if entities.IsInButton(mouseX, mouseY, g.playButton) {
			g.resetGame()
			g.state = PlayingState
		}
		if entities.IsInButton(mouseX, mouseY, g.exitButton) {
			return ebiten.Termination
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
		g.leftPaddleY -= PaddleSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) && g.leftPaddleY < ScreenHeight-PaddleHeight {
		g.leftPaddleY += PaddleSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) && g.rightPaddleY > 0 {
		g.rightPaddleY -= PaddleSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) && g.rightPaddleY < ScreenHeight-PaddleHeight {
		g.rightPaddleY += PaddleSpeed
	}

	g.ballX += g.ballVelX
	g.ballY += g.ballVelY

	if g.ballY < 0 || g.ballY > ScreenHeight-BallSize {
		g.ballVelY = -g.ballVelY
	}

	if g.checkPaddleCollision() {
		g.ballVelX = -g.ballVelX * 1.05
	}

	if g.ballX < 0 {
		g.rightScore++
		g.resetBall()
	}
	if g.ballX > ScreenWidth-BallSize {
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
		if entities.IsInButton(mouseX, mouseY, g.resumeButton) {
			g.state = PlayingState
		}
		if entities.IsInButton(mouseX, mouseY, g.menuButton) {
			g.state = MenuState
		}
	}
	g.mousePressed = mousePressed
	return nil
}

func (g *Game) updateButtonHover(mouseX, mouseY float32) {
	g.playButton.SetHovered(entities.IsInButton(mouseX, mouseY, g.playButton))
	g.exitButton.SetHovered(entities.IsInButton(mouseX, mouseY, g.exitButton))
	g.resumeButton.SetHovered(entities.IsInButton(mouseX, mouseY, g.resumeButton))
	g.menuButton.SetHovered(entities.IsInButton(mouseX, mouseY, g.menuButton))
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
	return ScreenWidth, ScreenHeight
}

func (g *Game) drawMenu(screen *ebiten.Image) {
	vector.DrawFilledRect(screen, 0, 0, ScreenWidth, ScreenHeight, color.RGBA{30, 30, 60, 255}, false)

	text.Draw(screen, "PING-PONG", g.font, ScreenWidth/2-55, 100, color.White)

	text.Draw(screen, "Controls", g.font, ScreenWidth/2-55, 180, color.White)
	text.Draw(screen, "Left paddle: W/S", g.font, ScreenWidth/2-80, 210, color.RGBA{255, 100, 100, 255})
	text.Draw(screen, "Right paddle: Arrow Keys", g.font, ScreenWidth/2-80, 240, color.RGBA{100, 100, 255, 255})
	text.Draw(screen, "Pause: ESC", g.font, ScreenWidth/2-80, 270, color.White)

	entities.DrawButton(screen, g.playButton, g.font)
	entities.DrawButton(screen, g.exitButton, g.font)
}

func (g *Game) drawPlaying(screen *ebiten.Image) {
	// Background
	vector.DrawFilledRect(screen, 0, 0, ScreenWidth, ScreenHeight, color.RGBA{20, 20, 40, 255}, false)

	// Net
	for y := 0; y < ScreenHeight; y += 20 {
		vector.DrawFilledRect(screen, ScreenWidth/2-2, float32(y), 4, 10, color.RGBA{100, 100, 100, 255}, false)
	}

	// Paddles
	vector.DrawFilledRect(screen, 20, float32(g.leftPaddleY), PaddleWidth, PaddleHeight, color.RGBA{255, 80, 80, 255}, false)
	vector.DrawFilledRect(screen, ScreenWidth-20-PaddleWidth, float32(g.rightPaddleY), PaddleWidth, PaddleHeight, color.RGBA{80, 80, 255, 255}, false)

	// Ball
	vector.DrawFilledRect(screen, float32(g.ballX), float32(g.ballY), BallSize, BallSize, color.RGBA{255, 255, 100, 255}, false)

	// Score
	text.Draw(screen, fmt.Sprintf("%d", g.leftScore), g.font, 50, 50, color.RGBA{255, 100, 100, 255})
	text.Draw(screen, fmt.Sprintf("%d", g.rightScore), g.font, ScreenWidth-70, 50, color.RGBA{100, 100, 255, 255})
}

func (g *Game) drawPause(screen *ebiten.Image) {
	vector.DrawFilledRect(screen, 0, 0, ScreenWidth, ScreenHeight, color.RGBA{0, 0, 0, 180}, false)

	// Pause text
	text.Draw(screen, "PAUSE", g.font, ScreenWidth/2-40, ScreenHeight/2-120, color.White)

	// Buttons
	entities.DrawButton(screen, g.resumeButton, g.font)
	entities.DrawButton(screen, g.menuButton, g.font)
}

func (g *Game) resetGame() {
	g.leftScore = 0
	g.rightScore = 0
	g.resetBall()
	g.leftPaddleY = ScreenHeight/2 - PaddleHeight/2
	g.rightPaddleY = ScreenHeight/2 - PaddleHeight/2
}

func (g *Game) resetBall() {
	g.ballX = ScreenWidth/2 - BallSize/2
	g.ballY = ScreenHeight/2 - BallSize/2
	g.ballVelX = BallSpeed
	g.ballVelY = BallSpeed * 0.7
}

func (g *Game) checkPaddleCollision() bool {
	// Left paddle
	if g.ballX <= 20+PaddleWidth && g.ballY+BallSize >= g.leftPaddleY && g.ballY <= g.leftPaddleY+PaddleHeight {
		return true
	}
	// Right paddle
	if g.ballX >= ScreenWidth-20-PaddleWidth-BallSize && g.ballY+BallSize >= g.rightPaddleY && g.ballY <= g.rightPaddleY+PaddleHeight {
		return true
	}
	return false
}