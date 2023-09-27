package main

import (
	"fmt"
	"image"
	_ "image/png"
	"math"
	"os"
	"sync"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
)

var hits int = 0
var misses int = 0
var darts int = 5
var targetX float64 = 400.0
var gameEnded bool = false
var speedLevel int = 3
var speeds = []float64{0, 5, 10, 20, 40, 60, 80}
var mu sync.Mutex
var lastSpeedChange time.Time

func main() {
	pixelgl.Run(run)
}

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Juego de dardos",
		Bounds: pixel.R(0, 0, 800, 600),
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	// Cargar la imagen de fondo desde la carpeta "assets"
	backgroundPic, err := loadPicture("assets/background.png")
	if err != nil {
		panic(err)
	}
	backgroundSprite := pixel.NewSprite(backgroundPic, backgroundPic.Bounds())

	var wg sync.WaitGroup

	// Cargar la imagen de la diana desde la carpeta "assets"
	targetPic, err := loadPicture("assets/diana.png")
	if err != nil {
		panic(err)
	}
	targetSprite := pixel.NewSprite(targetPic, targetPic.Bounds())

	// Cargar la imagen "left.png" desde la carpeta "assets"
	leftPic, err := loadPicture("assets/left.png")
	if err != nil {
		panic(err)
	}
	leftSprite := pixel.NewSprite(leftPic, leftPic.Bounds())

	// Cargar la imagen "right.png" desde la carpeta "assets"
	rightPic, err := loadPicture("assets/right.png")
	if err != nil {
		panic(err)
	}
	rightSprite := pixel.NewSprite(rightPic, rightPic.Bounds())

	//Go routine 1
	wg.Add(1)
	go func() {
		defer wg.Done()
		atlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
		for !win.Closed() {
			win.Clear(colornames.White)

			// Obtener la posición centrada del fondo
			bgPos := pixel.V(400, 300)

			// Escalar el fondo para que encaje en la ventana
			backgroundSprite.Draw(win, pixel.IM.Scaled(pixel.ZV, 0.66666667).Moved(bgPos))

			targetSprite.Draw(win, pixel.IM.Scaled(pixel.ZV, 0.0833334).Moved(pixel.V(targetX, 300)))

			// Dibujar el texto con contorno negro
			txt := text.New(pixel.V(10, 580), atlas)
			txt.Color = colornames.White
			mu.Lock()
			fmt.Fprintf(txt, "Aciertos: %d, Fallos: %d", hits, misses)
			mu.Unlock()
			txt.Draw(win, pixel.IM.Scaled(txt.Orig, 2).Moved(pixel.V(1, -1))) // Contorno negro

			// Dibujar el texto blanco encima del texto con contorno negro
			txt.Color = colornames.Black
			txt.Draw(win, pixel.IM.Scaled(txt.Orig, 2))

			leftSprite.Draw(win, pixel.IM.Scaled(pixel.ZV, 0.25).Moved(pixel.V(50, 50)))
			rightSprite.Draw(win, pixel.IM.Scaled(pixel.ZV, 0.09765625).Moved(pixel.V(150, 50)))

			txt = text.New(pixel.V(10, 120), atlas)
			txt.Color = colornames.White
			fmt.Fprintf(txt, "Control de velocidad: %d", speedLevel)
			txt.Draw(win, pixel.IM)

			if gameEnded {
				txt := text.New(pixel.V(100, 300), atlas)
				txt.Color = colornames.White
				fmt.Fprintf(txt, "Juego terminado. Presione R para reiniciar")
				txt.Draw(win, pixel.IM.Scaled(txt.Orig, 2))
			}

			win.Update()

			time.Sleep(time.Millisecond * 16)
		}
	}()

	//Go routine 2
	wg.Add(1)
	go player(win, &wg)

	//Go routine 3
	wg.Add(1)
	go speedControl(&wg)

	wg.Wait()
}

func player(win *pixelgl.Window, wg *sync.WaitGroup) {
	defer wg.Done()

	for !win.Closed() {
		mu.Lock()
		if gameEnded && win.JustPressed(pixelgl.KeyR) {
			gameEnded = false
			hits = 0
			misses = 0
			darts = 5
		}

		if darts > 0 && !gameEnded {
			if win.JustPressed(pixelgl.MouseButtonLeft) {
				mousePos := win.MousePosition()
				if mousePos.X > 10 && mousePos.X < 90 && mousePos.Y > 10 && mousePos.Y < 90 {
					if time.Since(lastSpeedChange) > time.Millisecond*500 {
						if speedLevel > 0 {
							speedLevel--
						}
						lastSpeedChange = time.Now()
					}
				} else if mousePos.X > 110 && mousePos.X < 190 && mousePos.Y > 10 && mousePos.Y < 90 {
					if time.Since(lastSpeedChange) > time.Millisecond*500 {
						if speedLevel < 6 {
							speedLevel++
						}
						lastSpeedChange = time.Now()
					}
				} else {
					if distance(mousePos.X, mousePos.Y, targetX, 300) < 2500 {
						hits++
					} else {
						misses++
					}
					darts--
				}
			}
		} else {
			gameEnded = true
		}
		mu.Unlock()

		time.Sleep(time.Millisecond * 16)
	}
}

func speedControl(wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		mu.Lock()
		newTargetX := targetX + speeds[speedLevel]
		if newTargetX > 750 || newTargetX < 50 {
			speeds[speedLevel] *= -1
			newTargetX = targetX
		}
		targetX = newTargetX
		mu.Unlock()

		time.Sleep(time.Millisecond * 16)
	}
}

func distance(x1, y1, x2, y2 float64) float64 {
	return math.Pow(x2-x1, 2) + math.Pow(y2-y1, 2)
}

func loadPicture(path string) (pixel.Picture, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return pixel.PictureDataFromImage(img), nil
}
