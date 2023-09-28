package models

import (
	"time"
	"sync"
	"math"
	
	"github.com/faiface/pixel/pixelgl"
)


func Player(win *pixelgl.Window, wg *sync.WaitGroup, hits *int, mu *sync.Mutex, gameEnded *bool, misses *int, darts *int, speedLevel *int, lastSpeedChange time.Time, targetX *float64) {
	defer wg.Done()

	for !win.Closed() {
		mu.Lock()
		if *gameEnded && win.JustPressed(pixelgl.KeyR) {
			*gameEnded = false
			*hits = 0
			*misses = 0
			*darts = 5
		}

		if *darts > 0 && !*gameEnded {
			if win.JustPressed(pixelgl.MouseButtonLeft) {
				mousePos := win.MousePosition()
				if mousePos.X > 10 && mousePos.X < 90 && mousePos.Y > 10 && mousePos.Y < 90 {
					if time.Since(lastSpeedChange) > time.Millisecond*500 {
						if *speedLevel > 0 {
							*speedLevel--
						}
						lastSpeedChange = time.Now()
					}
				} else if mousePos.X > 110 && mousePos.X < 190 && mousePos.Y > 10 && mousePos.Y < 90 {
					if time.Since(lastSpeedChange) > time.Millisecond*500 {
						if *speedLevel < 6 {
							*speedLevel++
						}
						lastSpeedChange = time.Now()
					}
				} else {
					if distance(mousePos.X, mousePos.Y, *targetX, 300) < 2500 {
						*hits++
					} else {
						*misses++
					}
					*darts--
				}
			}
		} else {
			*gameEnded = true
		}
		mu.Unlock()

		time.Sleep(time.Millisecond * 16)
	}
}

func distance(x1, y1, x2, y2 float64) float64 {
	return math.Pow(x2-x1, 2) + math.Pow(y2-y1, 2)
}

func SpeedControl(wg *sync.WaitGroup, mu *sync.Mutex, targetX *float64, speeds []float64, speedLevel *int) {
	defer wg.Done()
	for {
		mu.Lock()
		newTargetX := *targetX + speeds[*speedLevel]
		if newTargetX > 750 || newTargetX < 50 {
			speeds[*speedLevel] *= -1
			newTargetX = *targetX
		}
		*targetX = newTargetX
		mu.Unlock()

		time.Sleep(time.Millisecond * 16)
	}
}