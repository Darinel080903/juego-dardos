package models

import (
	"fmt"
	"sync"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
)

func Render(wg sync.WaitGroup, win *pixelgl.Window, hits *int, misses *int, targetX *float64, speedLevel *int, gameEnded *bool, leftSprite *pixel.Sprite, rightSprite *pixel.Sprite, mu *sync.Mutex, backgroundSprite *pixel.Sprite, targetSprite *pixel.Sprite) {
	defer wg.Done()
	atlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
	for !win.Closed() {
		win.Clear(colornames.White)

		// Obtener la posición centrada del fondo
		bgPos := pixel.V(400, 300)

		// Escalar el fondo para que encaje en la ventana
		backgroundSprite.Draw(win, pixel.IM.Scaled(pixel.ZV, 0.66666667).Moved(bgPos))

		targetSprite.Draw(win, pixel.IM.Scaled(pixel.ZV, 0.0833334).Moved(pixel.V(*targetX, 300)))

		// Dibujar el texto con contorno negro
		txt := text.New(pixel.V(10, 580), atlas)
		txt.Color = colornames.White
		mu.Lock()
		fmt.Fprintf(txt, "Aciertos: %d, Fallos: %d", *hits, *misses)
		mu.Unlock()
		txt.Draw(win, pixel.IM.Scaled(txt.Orig, 2).Moved(pixel.V(1, -1))) // Contorno negro

		// Dibujar el texto blanco encima del texto con contorno negro
		txt.Color = colornames.Black
		txt.Draw(win, pixel.IM.Scaled(txt.Orig, 2))

		leftSprite.Draw(win, pixel.IM.Scaled(pixel.ZV, 0.25).Moved(pixel.V(50, 50)))
		rightSprite.Draw(win, pixel.IM.Scaled(pixel.ZV, 0.09765625).Moved(pixel.V(150, 50)))

		txt = text.New(pixel.V(10, 120), atlas)
		txt.Color = colornames.White
		fmt.Fprintf(txt, "Control de velocidad: %d", *speedLevel)
		txt.Draw(win, pixel.IM)

		if *gameEnded {
			txt := text.New(pixel.V(100, 300), atlas)
			txt.Color = colornames.White
			fmt.Fprintf(txt, "Juego terminado. Presione R para reiniciar")
			txt.Draw(win, pixel.IM.Scaled(txt.Orig, 2))
		}

		win.Update()

		time.Sleep(time.Millisecond * 16)
	}
}
