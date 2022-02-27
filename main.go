package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gdamore/tcell/v2"
)

//global var
var level1 []int
var direction string = "up"
var idir string = "left"
var pdir string = "up"
var bdir string = "right"
var cdir string = "down"
var ionpoint int = 0
var ponpoint int = 0
var bonpoint int = 0
var conpoint int = 0
var death int = 0
var win int = 0
var lifes int = 3
var points int = 259

//a struct used in ghost ai pathfinding algorithm
//the algorithm locates current position of ghost and checks all surrounding cells to see if a ghost can go there
//if it can, it adds it to the list with the struct below containing its location, parent and distance from the ghost cell
//i wont go in detail here, but the algorithm checks next cells, checks if they've been visited, assigns them the right values etc etc
type ainode struct {
	loc, parent, dist int
}

const (
	// Black         = "\u001b[30m"
	// Red           = "\u001b[31m"
	// Green         = "\u001b[32m"
	// Yellow        = "\u001b[33m"
	// Blue          = "\u001b[34m"
	// Magenta       = "\u001b[35m"
	// Cyan          = "\u001b[36m"
	// White         = "\u001b[37m"
	// Reset         = "\u001b[0m"

	//unicode codes for printing with tcell
	Wall   int32 = '\u0023'
	Ghost  int32 = '\u0026'
	Point  int32 = '\u00A4'
	Vup    int32 = '\u0056'
	Vdown  int32 = '\u0041'
	Vright int32 = '\u003C'
	Vleft  int32 = '\u003E'
	Space  int32 = '\u0020'
)

func main() {

	// Initialize screen
	s, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	if err := s.Init(); err != nil {
		log.Fatalf("%+v", err)
	}
	// Set default text style
	defStyle := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)
	wallStyle := tcell.StyleDefault.Background(tcell.ColorGray).Foreground(tcell.ColorDimGray)
	inkyStyle := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorLightCyan)
	pinkyStyle := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorPink)
	blinkyStyle := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorRed)
	clydeStyle := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorOrange)
	// Clear screen
	s.Clear()

	//key event setup
	quit := make(chan struct{})
	go func() {
		for {
			ev := s.PollEvent()
			switch ev := ev.(type) {
			case *tcell.EventKey:
				switch ev.Key() {
				case tcell.KeyEscape, tcell.KeyEnter:
					close(quit)
					return
				case tcell.KeyCtrlL:
					s.Sync()
				case tcell.KeyDown:
					direction = "down"
				case tcell.KeyUp:
					direction = "up"
				case tcell.KeyLeft:
					direction = "left"
				case tcell.KeyRight:
					direction = "right"
				case tcell.KeyTAB:
					lifes = 0
					death = 1
				}

			case *tcell.EventResize:
				s.Sync()
			}
		}
	}()

	//level creation
	level1 = append(level1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 8)
	level1 = append(level1, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 1, 1, 1, 1, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 8)
	level1 = append(level1, 0, 1, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 0, 0, 1, 0, 8)
	level1 = append(level1, 0, 1, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 0, 0, 1, 0, 8)
	level1 = append(level1, 0, 1, 1, 1, 1, 1, 1, 0, 0, 1, 1, 1, 1, 0, 0, 1, 1, 1, 1, 0, 0, 1, 1, 1, 1, 1, 1, 0, 8)
	level1 = append(level1, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 8)
	level1 = append(level1, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 8)
	level1 = append(level1, 0, 1, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 1, 0, 8)
	level1 = append(level1, 0, 1, 0, 0, 1, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 1, 0, 0, 1, 0, 8)
	level1 = append(level1, 0, 1, 0, 0, 1, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 1, 0, 0, 1, 0, 8)
	level1 = append(level1, 2, 1, 1, 1, 1, 1, 1, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 1, 1, 1, 1, 1, 1, 2, 8)
	level1 = append(level1, 0, 1, 0, 0, 0, 0, 1, 0, 0, 2, 0, 0, 2, 0, 0, 2, 0, 0, 2, 0, 0, 1, 0, 0, 0, 0, 1, 0, 8)
	level1 = append(level1, 0, 1, 0, 0, 0, 0, 1, 0, 0, 2, 2, 2, 4, 0, 0, 7, 2, 2, 2, 0, 0, 1, 0, 0, 0, 0, 1, 0, 8)
	level1 = append(level1, 0, 1, 0, 0, 0, 0, 1, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 1, 0, 0, 0, 0, 1, 0, 8)
	level1 = append(level1, 0, 1, 1, 1, 1, 1, 1, 0, 0, 2, 2, 2, 5, 0, 0, 6, 2, 2, 2, 0, 0, 1, 1, 1, 1, 1, 1, 0, 8)
	level1 = append(level1, 0, 2, 0, 0, 0, 0, 1, 0, 0, 2, 0, 0, 2, 0, 0, 2, 0, 0, 2, 0, 0, 1, 0, 0, 0, 0, 2, 0, 8)
	level1 = append(level1, 0, 2, 0, 0, 0, 0, 1, 0, 0, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 0, 0, 1, 0, 0, 0, 0, 2, 0, 8)
	level1 = append(level1, 0, 2, 2, 2, 0, 0, 1, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 1, 0, 0, 2, 2, 2, 0, 8)
	level1 = append(level1, 0, 0, 0, 2, 0, 0, 1, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 1, 0, 0, 2, 0, 0, 0, 8)
	level1 = append(level1, 0, 0, 0, 2, 0, 0, 1, 1, 1, 1, 1, 1, 1, 0, 0, 1, 1, 1, 1, 1, 1, 1, 0, 0, 2, 0, 0, 0, 8)
	level1 = append(level1, 2, 2, 2, 2, 0, 0, 1, 0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 0, 0, 0, 1, 0, 0, 2, 2, 2, 2, 8)
	level1 = append(level1, 0, 2, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 2, 0, 8)
	level1 = append(level1, 0, 2, 0, 0, 0, 0, 1, 0, 0, 1, 1, 1, 1, 3, 1, 1, 1, 1, 1, 0, 0, 1, 0, 0, 0, 0, 2, 0, 8)
	level1 = append(level1, 0, 1, 1, 1, 1, 1, 1, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 1, 1, 1, 1, 1, 1, 0, 8)
	level1 = append(level1, 0, 1, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 1, 0, 8)
	level1 = append(level1, 0, 1, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 1, 0, 8)
	level1 = append(level1, 0, 1, 0, 0, 1, 1, 1, 1, 1, 1, 0, 0, 1, 1, 1, 1, 0, 0, 1, 1, 1, 1, 1, 1, 0, 0, 1, 0, 8)
	level1 = append(level1, 0, 1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 8)
	level1 = append(level1, 0, 1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 8)
	level1 = append(level1, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 8)
	level1 = append(level1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 8)
	game := 1
gameloop:
	for game == 1 {
		time.Sleep(1 * time.Second)
		movement(level1)
		if win == 1 {
			break gameloop
		}
		if death == 1 {
			if lifes <= 0 {
				break gameloop
			}
			lifes--
			direction = "up"
			idir = "left"
			pdir = "down"
			bdir = "right"
			cdir = "up"
			time.Sleep(3 * time.Second)
			for i, v := range []int(level1) {
				switch v {
				case 3:
					level1[i] = 2
				case 4:
					if ionpoint == 1 {
						level1[i] = 1
						ionpoint = 0
					} else {
						level1[i] = 2
					}
				case 5:
					if ponpoint == 1 {
						level1[i] = 1
						ponpoint = 0
					} else {
						level1[i] = 2
					}
				case 6:
					if bonpoint == 1 {
						level1[i] = 1
						bonpoint = 0
					} else {
						level1[i] = 2
					}
				case 7:
					if conpoint == 1 {
						level1[i] = 1
						conpoint = 0
					} else {
						level1[i] = 2
					}
				}
			}
			level1[360] = 4
			level1[363] = 7
			level1[418] = 5
			level1[421] = 6
			level1[651] = 3
			death = 0
		}
		levelprinter(s, 1, 1, defStyle, inkyStyle, pinkyStyle, blinkyStyle, clydeStyle, wallStyle, level1)
		s.Show()
	}
	s.Fini()
	if win == 1 {
		fmt.Println("You win! Congratulations!")
		return
	} else {
		fmt.Println("You lose... Better luck next time!")
		return
	}
}

func levelprinter(s tcell.Screen, x, y int, def, inky, pinky, blinky, clyde, wall tcell.Style, level []int) {
	row := y
	col := x
	for _, r := range []int(level) {
		// s.SetContent(col, row, r, nil, style)
		// col++
		// if col >= x2 {
		// 	row++
		// 	col = x1
		// }
		// if row > y2 {
		// 	break
		// }
		switch r {
		case 0:
			s.SetContent(col, row, '#', nil, wall)
			col++
		case 1:
			s.SetContent(col, row, Point, nil, def)
			col++
		case 2:
			s.SetContent(col, row, ' ', nil, def)
			col++
		case 3:
			switch direction {
			case "up":
				s.SetContent(col, row, 'V', nil, def)
			case "down":
				s.SetContent(col, row, 'A', nil, def)
			case "left":
				s.SetContent(col, row, '>', nil, def)
			case "right":
				s.SetContent(col, row, '<', nil, def)
			}
			col++
		case 4:
			s.SetContent(col, row, '&', nil, inky)
			col++
		case 5:
			s.SetContent(col, row, '&', nil, pinky)
			col++
		case 6:
			s.SetContent(col, row, '&', nil, blinky)
			col++
		case 7:
			s.SetContent(col, row, '&', nil, clyde)
			col++
		case 8:
			row++
			col = 1
		}
	}
}

func movement(level []int) {
	var vmanx int
	var vmany int
	var inkyx int
	var inkyy int
	var blinkyx int
	var blinkyy int
	var clydex int
	var clydey int
	//an array which holds information about whether a ghost has to turn
	//and what which directions are available
	//in order: shouldturn, up, down, left, right
	var ghostturn []int
	ghostturn = append(ghostturn, 1, 0, 0, 0, 0)
	//Player movement
out:
	for i, v := range []int(level) {
		if v == 3 {
			switch direction {
			case "up":
				switch level[i-29] {
				case 0:
					vmanx = findx(i)
					vmany = findy(i)
					break out
				case 1:
					vmanx = findx(i)
					vmany = findy(i)
					level[i-29] = 3
					level[i] = 2
					points--
					break out
				case 2:
					vmanx = findx(i)
					vmany = findy(i)
					level[i-29] = 3
					level[i] = 2
					break out
				case 4:
					vmanx = findx(i)
					vmany = findy(i)
					death = 1
					break out
				case 5:
					vmanx = findx(i)
					vmany = findy(i)
					death = 1
					break out
				case 6:
					vmanx = findx(i)
					vmany = findy(i)
					death = 1
					break out
				case 7:
					vmanx = findx(i)
					vmany = findy(i)
					death = 1
					break out
				case 8:
					vmanx = findx(i)
					vmany = findy(i)
					level[i-27] = 3
					level[i] = 2
					break out
				}
			case "down":
				switch level[i+29] {
				case 0:
					vmanx = findx(i)
					vmany = findy(i)
					break out
				case 1:
					vmanx = findx(i)
					vmany = findy(i)
					level[i+29] = 3
					level[i] = 2
					points--
					break out
				case 2:
					vmanx = findx(i)
					vmany = findy(i)
					level[i+29] = 3
					level[i] = 2
					break out
				case 4:
					vmanx = findx(i)
					vmany = findy(i)
					death = 1
					break out
				case 5:
					vmanx = findx(i)
					vmany = findy(i)
					death = 1
					break out
				case 6:
					vmanx = findx(i)
					vmany = findy(i)
					death = 1
					break out
				case 7:
					vmanx = findx(i)
					vmany = findy(i)
					death = 1
					break out
				case 8:
					vmanx = findx(i)
					vmany = findy(i)
					level[i-27] = 3
					level[i] = 2
					break out
				}
			case "left":
				switch level[i-1] {
				case 0:
					vmanx = findx(i)
					vmany = findy(i)
					break out
				case 1:
					vmanx = findx(i)
					vmany = findy(i)
					level[i-1] = 3
					level[i] = 2
					points--
					break out
				case 2:
					vmanx = findx(i)
					vmany = findy(i)
					level[i-1] = 3
					level[i] = 2
					break out
				case 4:
					vmanx = findx(i)
					vmany = findy(i)
					death = 1
					break out
				case 5:
					vmanx = findx(i)
					vmany = findy(i)
					death = 1
					break out
				case 6:
					vmanx = findx(i)
					vmany = findy(i)
					death = 1
					break out
				case 7:
					vmanx = findx(i)
					vmany = findy(i)
					death = 1
					break out
				case 8:
					vmanx = findx(i)
					vmany = findy(i)
					level[i+27] = 3
					level[i] = 2
					break out
				}
			case "right":
				switch level[i+1] {
				case 0:
					vmanx = findx(i)
					vmany = findy(i)
					break out
				case 1:
					vmanx = findx(i)
					vmany = findy(i)
					level[i+1] = 3
					level[i] = 2
					points--
					break out
				case 2:
					vmanx = findx(i)
					vmany = findy(i)
					level[i+1] = 3
					level[i] = 2
					break out
				case 4:
					vmanx = findx(i)
					vmany = findy(i)
					death = 1
					break out
				case 5:
					vmanx = findx(i)
					vmany = findy(i)
					death = 1
					break out
				case 6:
					vmanx = findx(i)
					vmany = findy(i)
					death = 1
					break out
				case 7:
					vmanx = findx(i)
					vmany = findy(i)
					death = 1
					break out
				case 8:
					vmanx = findx(i)
					vmany = findy(i)
					level[i-27] = 3
					level[i] = 2
					break out
				}
			}
		}
	}
	//inky
out2:
	for i, v := range []int(level) {
		if v == 4 {
			ghostturncheck(ghostturn, level1, i, idir)
			//ghost ai for which turn to choose

			inkyx = findx(i)
			inkyy = findy(i)
			inkyx -= vmanx
			inkyy -= vmany
			if ghostturn[0] == 0 {
				if inkyx == 0 {
					if inkyy > 0 {
						if ghostturn[1] == 0 {
							idir = "up"
						} else if ghostturn[3] == 0 {
							idir = "left"
						} else if ghostturn[4] == 0 {
							idir = "right"
						} else if ghostturn[2] == 0 {
							idir = "down"
						}
					} else {
						if ghostturn[2] == 0 {
							idir = "down"
						} else if ghostturn[3] == 0 {
							idir = "left"
						} else if ghostturn[4] == 0 {
							idir = "right"
						} else if ghostturn[1] == 0 {
							idir = "up"
						}
					}
				} else if inkyx > 0 {
					if inkyy == 0 {
						if ghostturn[3] == 0 {
							idir = "left"
						} else if ghostturn[1] == 0 {
							idir = "up"
						} else if ghostturn[2] == 0 {
							idir = "down"
						} else if ghostturn[4] == 0 {
							idir = "right"
						}
					} else if inkyy > 0 {
						if inkyx >= inkyy {
							if ghostturn[1] == 0 {
								idir = "up"
							} else if ghostturn[3] == 0 {
								idir = "left"
							} else if ghostturn[4] == 0 {
								idir = "right"
							} else if ghostturn[2] == 0 {
								idir = "down"
							}
						} else {
							if ghostturn[3] == 0 {
								idir = "left"
							} else if ghostturn[1] == 0 {
								idir = "up"
							} else if ghostturn[2] == 0 {
								idir = "down"
							} else if ghostturn[4] == 0 {
								idir = "right"
							}
						}
					} else {
						if inkyx >= -inkyy {
							if ghostturn[2] == 0 {
								idir = "down"
							} else if ghostturn[3] == 0 {
								idir = "left"
							} else if ghostturn[4] == 0 {
								idir = "right"
							} else if ghostturn[1] == 0 {
								idir = "up"
							}
						} else {
							if ghostturn[3] == 0 {
								idir = "left"
							} else if ghostturn[2] == 0 {
								idir = "down"
							} else if ghostturn[1] == 0 {
								idir = "up"
							} else if ghostturn[4] == 0 {
								idir = "right"
							}
						}
					}
				} else {
					if inkyy == 0 {
						if ghostturn[4] == 0 {
							idir = "right"
						} else if ghostturn[1] == 0 {
							idir = "up"
						} else if ghostturn[2] == 0 {
							idir = "down"
						} else if ghostturn[3] == 0 {
							idir = "left"
						}
					} else if inkyy > 0 {
						if -inkyx >= inkyy {
							if ghostturn[1] == 0 {
								idir = "up"
							} else if ghostturn[4] == 0 {
								idir = "right"
							} else if ghostturn[3] == 0 {
								idir = "left"
							} else if ghostturn[2] == 0 {
								idir = "down"
							}
						} else {
							if ghostturn[4] == 0 {
								idir = "right"
							} else if ghostturn[1] == 0 {
								idir = "up"
							} else if ghostturn[2] == 0 {
								idir = "down"
							} else if ghostturn[3] == 0 {
								idir = "left"
							}
						}
					} else {
						if -inkyx >= -inkyy {
							if ghostturn[2] == 0 {
								idir = "down"
							} else if ghostturn[4] == 0 {
								idir = "right"
							} else if ghostturn[3] == 0 {
								idir = "left"
							} else if ghostturn[1] == 0 {
								idir = "up"
							}
						}
					}
				}
			}
			//end of ghost ai direction choosing
			switch idir {
			case "up":
				switch level[i-29] {
				case 1:
					if ionpoint == 0 {
						ionpoint = 1
						level[i-29] = 4
						level[i] = 2
					} else {
						level[i-29] = 4
						level[i] = 1
					}
					break out2
				case 2:
					if ionpoint == 1 {
						ionpoint = 0
						level[i-29] = 4
						level[i] = 1
					} else {
						level[i-29] = 4
						level[i] = 2
					}
					break out2
				case 3:
					death = 1
					break out2
				}
			case "down":
				switch level[i+29] {
				case 1:
					if ionpoint == 0 {
						ionpoint = 1
						level[i+29] = 4
						level[i] = 2
					} else {
						level[i+29] = 4
						level[i] = 1
					}
					break out2
				case 2:
					if ionpoint == 1 {
						ionpoint = 0
						level[i+29] = 4
						level[i] = 1
					} else {
						level[i+29] = 4
						level[i] = 2
					}
					break out2
				case 3:
					death = 1
					break out2
				}
			case "left":
				switch level[i-1] {
				case 1:
					if ionpoint == 0 {
						ionpoint = 1
						level[i-1] = 4
						level[i] = 2
					} else {
						level[i-1] = 4
						level[i] = 1
					}
					break out2
				case 2:
					if ionpoint == 1 {
						ionpoint = 0
						level[i-1] = 4
						level[i] = 1
					} else {
						level[i-1] = 4
						level[i] = 2
					}
					break out2
				case 3:
					death = 1
					break out2
				case 8:
					if ionpoint == 1 {
						if level[i+27] == 2 {
							ionpoint = 0
						}
						level[i+27] = 4
						level[i] = 1
					} else {
						if level[i+27] == 1 {
							ionpoint = 1
						}
						level[i+27] = 4
						level[i] = 2
					}
					break out2
				}
			case "right":
				switch level[i+1] {
				case 1:
					if ionpoint == 0 {
						ionpoint = 1
						level[i+1] = 4
						level[i] = 2
					} else {
						level[i+1] = 4
						level[i] = 1
					}
					break out2
				case 2:
					if ionpoint == 1 {
						ionpoint = 0
						level[i+1] = 4
						level[i] = 1
					} else {
						level[i+1] = 4
						level[i] = 2
					}
					break out2
				case 3:
					death = 1
					break out2
				case 8:
					if ionpoint == 1 {
						if level[i-27] == 2 {
							ionpoint = 0
						}
						level[i-27] = 4
						level[i] = 1
					} else {
						if level[i-27] == 1 {
							ionpoint = 1
						}
						level[i-27] = 4
						level[i] = 2
					}
					break out2
				}
			}
		}
	}

	//pinky
out3:
	for i, v := range []int(level) {
		if v == 5 {
			var pqueue []ainode     //queue for cells to check for algorithm
			var pchecked []ainode   //list of already checked cells
			var ptemp []ainode      //temporary list used in appending
			var pfound bool = false //bool which switches true once pacman's location is found
			var pcomptemp int       //int which holds the position of struct in pchecked returned by compareI()

			//the following code is a pathfinding algorithm
			//starting with the current location of our ghost, the algorithm checks the cells above, below and on the sides
			//if they can be "walked into", the algorithm checks if it's a cell that's been already checked
			//if it wasn't it adds it to the queue as a ainode structure.
			//if it was, it checks its distance. If the written distance is larger than what the current path has, the cell is rewritten
			//it's parent becoming the current cell and the distance becoming the shorter one
			//once all 4 directions are checked, the cell is checked if it is the pacman one
			//if it isnt, it's added to the checked list and the algorithm moves onto next element of the queue
			//if it is, the pathfinding algorithm stops (pfound becomes true)
			ptemp = []ainode{{i, -1, 0}}
			pqueue = append(pqueue, ptemp...)
			for !pfound {
				if level[pqueue[0].loc-1] != 0 {
					if level[pqueue[0].loc-1] != 8 {
						pcomptemp = compareI(pchecked, level[pqueue[0].loc-1])
						if pcomptemp >= 0 {
							if pchecked[pcomptemp].dist > pqueue[0].dist+1 {
								pchecked[pcomptemp].parent = pqueue[0].loc
								pchecked[pcomptemp].dist = pqueue[0].dist + 1
							}
						}
						ptemp = []ainode{{pqueue[0].loc - 1, pqueue[0].loc, pqueue[0].dist + 1}}
						pqueue = append(pqueue, ptemp...)
					} else {
						pcomptemp = compareI(pchecked, level[pqueue[0].loc+27])
						if pcomptemp >= 0 {
							if pchecked[pcomptemp].dist > pqueue[0].dist+1 {
								pchecked[pcomptemp].parent = pqueue[0].loc
								pchecked[pcomptemp].dist = pqueue[0].dist + 1
							}
						}
						ptemp = []ainode{{pqueue[0].loc + 27, pqueue[0].loc, pqueue[0].dist + 1}}
						pqueue = append(pqueue, ptemp...)
					}
				}
				if level[pqueue[0].loc+1] != 0 {
					if level[pqueue[0].loc+1] != 8 {
						pcomptemp = compareI(pchecked, level[pqueue[0].loc+1])
						if pcomptemp >= 0 {
							if pchecked[pcomptemp].dist > pqueue[0].dist+1 {
								pchecked[pcomptemp].parent = pqueue[0].loc
								pchecked[pcomptemp].dist = pqueue[0].dist + 1
							}
						}
						ptemp = []ainode{{pqueue[0].loc + 1, pqueue[0].loc, pqueue[0].dist + 1}}
						pqueue = append(pqueue, ptemp...)
					} else {
						pcomptemp = compareI(pchecked, level[pqueue[0].loc-27])
						if pcomptemp >= 0 {
							if pchecked[pcomptemp].dist > pqueue[0].dist+1 {
								pchecked[pcomptemp].parent = pqueue[0].loc
								pchecked[pcomptemp].dist = pqueue[0].dist + 1
							}
						}
						ptemp = []ainode{{pqueue[0].loc - 27, pqueue[0].loc, pqueue[0].dist + 1}}
						pqueue = append(pqueue, ptemp...)
					}

				}
				if level[pqueue[0].loc+29] != 0 {
					pcomptemp = compareI(pchecked, level[pqueue[0].loc+29])
					if pcomptemp >= 0 {
						if pchecked[pcomptemp].dist > pqueue[0].dist+1 {
							pchecked[pcomptemp].parent = pqueue[0].loc
							pchecked[pcomptemp].dist = pqueue[0].dist + 1
						}
					}
					ptemp = []ainode{{pqueue[0].loc + 29, pqueue[0].loc, pqueue[0].dist + 1}}
					pqueue = append(pqueue, ptemp...)
				}
				if level[pqueue[0].loc-29] != 0 {
					pcomptemp = compareI(pchecked, level[pqueue[0].loc-29])
					if pcomptemp >= 0 {
						if pchecked[pcomptemp].dist > pqueue[0].dist+1 {
							pchecked[pcomptemp].parent = pqueue[0].loc
							pchecked[pcomptemp].dist = pqueue[0].dist + 1
						}
					}
					ptemp = []ainode{{pqueue[0].loc - 29, pqueue[0].loc, pqueue[0].dist + 1}}
					pqueue = append(pqueue, ptemp...)
				}
				ptemp = []ainode{{pqueue[0].loc, pqueue[0].parent, pqueue[0].dist}}
				pchecked = append(pchecked, ptemp...)
				if level[i] == 3 {
					pfound = true
				}
				pqueue = pqueue[1:]
			}
			//this is where the pathfinding stops
			//below the code goes backwards from the pacman's cell to the ghost's cell
			//in order to determine which direction to pick
			//pfound variable is reused till the code finds its way back
			pfound = false
			ptemp = []ainode{pchecked[len(pchecked)-1]}
			pchecked = pchecked[:len(pchecked)-1]
			for !pfound {

				if pchecked[len(pchecked)-1].parent == ptemp[0].loc {
					if level[pchecked[len(pchecked)-1].loc] == 5 {
						if ptemp[0].loc == pchecked[len(pchecked)-1].loc-1 {
							pdir = "left"
						} else if ptemp[0].loc == pchecked[len(pchecked)-1].loc+1 {
							pdir = "right"
						} else if ptemp[0].loc == pchecked[len(pchecked)-1].loc-29 {
							pdir = "up"
						} else if ptemp[0].loc == pchecked[len(pchecked)-1].loc+29 {
							pdir = "down"
						}
						pfound = true
					}
					ptemp = []ainode{pchecked[len(pchecked)-1]}
					pchecked = pchecked[:len(pchecked)-1]
				} else {
					pchecked = pchecked[:len(pchecked)-1]
				}

			}

			//end of ghost ai direction choosing based on pathfinding
			//below is the ghosts actual movement based on the current direction
			//if the path it's supposed to move in is blocked, it will automatically switch direction
			//funnily the pmove bool is always set to false and only there to ensure the switch goes on till the ghost moves
			var pmove bool = false
			for !pmove {
				switch pdir {
				case "up":
					switch level[i-29] {
					case 1:
						if ponpoint == 0 {
							ponpoint = 1
							level[i-29] = 5
							level[i] = 2
						} else {
							level[i-29] = 5
							level[i] = 1
						}
						break out3
					case 2:
						if ponpoint == 1 {
							ponpoint = 0
							level[i-29] = 5
							level[i] = 1
						} else {
							level[i-29] = 5
							level[i] = 2
						}
						break out3
					case 3:
						death = 1
						break out3
					default:
						if (level[i-1] < 4 && level[i-1] != 0) || level[i-1] == 8 {
							pdir = "left"
						} else if (level[i+1] < 4 && level[i+1] != 0) || level[i+1] == 8 {
							pdir = "right"
						} else if level[i+29] < 4 && level[i+29] != 0 {
							pdir = "down"
						}
					}
				case "down":
					switch level[i+29] {
					case 1:
						if ponpoint == 0 {
							ponpoint = 1
							level[i+29] = 5
							level[i] = 2
						} else {
							level[i+29] = 5
							level[i] = 1
						}
						break out3
					case 2:
						if ponpoint == 1 {
							ponpoint = 0
							level[i+29] = 5
							level[i] = 1
						} else {
							level[i+29] = 5
							level[i] = 2
						}
						break out3
					case 3:
						death = 1
						break out3
					default:
						if (level[i-1] < 4 && level[i-1] != 0) || level[i-1] == 8 {
							pdir = "left"
						} else if (level[i+1] < 4 && level[i+1] != 0) || level[i+1] == 8 {
							pdir = "right"
						} else if level[i-29] < 4 && level[i-29] != 0 {
							pdir = "up"
						}
					}
				case "left":
					switch level[i-1] {
					case 1:
						if ponpoint == 0 {
							ponpoint = 1
							level[i-1] = 5
							level[i] = 2
						} else {
							level[i-1] = 5
							level[i] = 1
						}
						break out3
					case 2:
						if ponpoint == 1 {
							ponpoint = 0
							level[i-1] = 5
							level[i] = 1
						} else {
							level[i-1] = 5
							level[i] = 2
						}
						break out3
					case 3:
						death = 1
						break out3
					case 8:
						if ponpoint == 1 {
							if level[i+27] == 2 {
								ponpoint = 0
							}
							level[i+27] = 5
							level[i] = 1
						} else {
							if level[i+27] == 1 {
								ponpoint = 1
							}
							level[i+27] = 5
							level[i] = 2
						}
						break out3
					default:
						if level[i+29] < 4 && level[i+29] != 0 {
							pdir = "down"
						} else if level[i-29] < 4 && level[i-29] != 0 {
							pdir = "up"
						} else if (level[i+1] < 4 && level[i+1] != 0) || level[i+1] == 8 {
							pdir = "right"
						}
					}
				case "right":
					switch level[i+1] {
					case 1:
						if ponpoint == 0 {
							ponpoint = 1
							level[i+1] = 5
							level[i] = 2
						} else {
							level[i+1] = 5
							level[i] = 1
						}
						break out3
					case 2:
						if ponpoint == 1 {
							ponpoint = 0
							level[i+1] = 5
							level[i] = 1
						} else {
							level[i+1] = 5
							level[i] = 2
						}
						break out3
					case 3:
						death = 1
						break out3
					case 8:
						if ponpoint == 1 {
							if level[i-27] == 2 {
								ponpoint = 0
							}
							level[i-27] = 5
							level[i] = 1
						} else {
							if level[i-27] == 1 {
								ponpoint = 1
							}
							level[i-27] = 5
							level[i] = 2
						}
						break out3
					default:
						if level[i+29] < 4 && level[i+29] != 0 {
							pdir = "down"
						} else if level[i-29] < 4 && level[i-29] != 0 {
							pdir = "up"
						} else if (level[i-1] < 4 && level[i-1] != 0) || level[i-1] == 8 {
							pdir = "left"
						}
					}
				}
			}
		}
	}
	//blinky
out4:
	for i, v := range []int(level) {
		if v == 6 {
			ghostturncheck(ghostturn, level1, i, bdir)
			//ghost ai for which turn to choose
			if ghostturn[0] == 0 {
				blinkyx = findx(i)
				blinkyy = findy(i)
				blinkyx -= vmanx
				blinkyy -= vmany
				if blinkyx == 0 {
					if ghostturn[1] == 0 {
						bdir = "up"
					} else if ghostturn[3] == 0 {
						bdir = "left"
					} else if ghostturn[4] == 0 {
						bdir = "right"
					} else if ghostturn[2] == 0 {
						bdir = "down"
					}
				}
				if blinkyx > 0 {
					if blinkyy == 0 {
						if ghostturn[3] == 0 {
							bdir = "left"
						} else if ghostturn[1] == 0 {
							bdir = "up"
						} else if ghostturn[2] == 0 {
							bdir = "down"
						} else if ghostturn[4] == 0 {
							bdir = "right"
						}
					}
				}
				if blinkyx < 0 {
					if blinkyy == 0 {
						if ghostturn[3] == 0 {
							bdir = "left"
						} else if ghostturn[1] == 0 {
							bdir = "up"
						} else if ghostturn[2] == 0 {
							bdir = "down"
						} else if ghostturn[4] == 0 {
							bdir = "right"
						}
					}
				}
			}
			//end of ghost ai direction choosing
			switch bdir {
			case "up":
				switch level[i-29] {
				case 1:
					if bonpoint == 0 {
						bonpoint = 1
						level[i-29] = 6
						level[i] = 2
					} else {
						level[i-29] = 6
						level[i] = 1
					}
					break out4
				case 2:
					if bonpoint == 1 {
						bonpoint = 0
						level[i-29] = 6
						level[i] = 1
					} else {
						level[i-29] = 6
						level[i] = 2
					}
					break out4
				case 3:
					death = 1
					break out4
				}
			case "down":
				switch level[i+29] {
				case 1:
					if bonpoint == 0 {
						bonpoint = 1
						level[i+29] = 6
						level[i] = 2
					} else {
						level[i+29] = 6
						level[i] = 1
					}
					break out4
				case 2:
					if bonpoint == 1 {
						bonpoint = 0
						level[i+29] = 6
						level[i] = 1
					} else {
						level[i+29] = 6
						level[i] = 2
					}
					break out4
				case 3:
					death = 1
					break out4
				}
			case "left":
				switch level[i-1] {
				case 1:
					if bonpoint == 0 {
						bonpoint = 1
						level[i-1] = 6
						level[i] = 2
					} else {
						level[i-1] = 6
						level[i] = 1
					}
					break out4
				case 2:
					if bonpoint == 1 {
						bonpoint = 0
						level[i-1] = 6
						level[i] = 1
					} else {
						level[i-1] = 6
						level[i] = 2
					}
					break out4
				case 3:
					death = 1
					break out4
				case 8:
					if bonpoint == 1 {
						if level[i+27] == 2 {
							bonpoint = 0
						}
						level[i+27] = 6
						level[i] = 1
					} else {
						if level[i+27] == 1 {
							bonpoint = 1
						}
						level[i+27] = 6
						level[i] = 2
					}
					break out4
				}
			case "right":
				switch level[i+1] {
				case 1:
					if bonpoint == 0 {
						bonpoint = 1
						level[i+1] = 6
						level[i] = 2
					} else {
						level[i+1] = 6
						level[i] = 1
					}
					break out4
				case 2:
					if bonpoint == 1 {
						bonpoint = 0
						level[i+1] = 6
						level[i] = 1
					} else {
						level[i+1] = 6
						level[i] = 2
					}
					break out4
				case 3:
					death = 1
					break out4
				case 8:
					if bonpoint == 1 {
						if level[i-27] == 2 {
							bonpoint = 0
						}
						level[i-27] = 6
						level[i] = 1
					} else {
						if level[i-27] == 1 {
							bonpoint = 1
						}
						level[i-27] = 6
						level[i] = 2
					}
					break out4
				}
			}
		}
	}
	//clyde
out5:
	for i, v := range []int(level) {
		if v == 7 {
			ghostturncheck(ghostturn, level1, i, cdir)
			//ghost ai for which turn to choose
			if ghostturn[0] == 0 {
				clydex = findx(i)
				clydey = findy(i)
				clydex -= vmanx
				clydey -= vmany
				if clydex == 0 {
					if ghostturn[1] == 0 {
						cdir = "up"
					} else if ghostturn[3] == 0 {
						cdir = "left"
					} else if ghostturn[4] == 0 {
						cdir = "right"
					} else if ghostturn[2] == 0 {
						cdir = "down"
					}
				}
				if clydex > 0 {
					if clydey == 0 {
						if ghostturn[3] == 0 {
							cdir = "left"
						} else if ghostturn[1] == 0 {
							cdir = "up"
						} else if ghostturn[2] == 0 {
							cdir = "down"
						} else if ghostturn[4] == 0 {
							cdir = "right"
						}
					}
				}
				if clydex < 0 {
					if clydey == 0 {
						if ghostturn[3] == 0 {
							cdir = "left"
						} else if ghostturn[1] == 0 {
							cdir = "up"
						} else if ghostturn[2] == 0 {
							cdir = "down"
						} else if ghostturn[4] == 0 {
							cdir = "right"
						}
					}
				}
			}
			//end of ghost ai direction choosing
			switch cdir {
			case "up":
				switch level[i-29] {
				case 1:
					if conpoint == 0 {
						conpoint = 1
						level[i-29] = 7
						level[i] = 2
					} else {
						level[i-29] = 7
						level[i] = 1
					}
					break out5
				case 2:
					if conpoint == 1 {
						conpoint = 0
						level[i-29] = 7
						level[i] = 1
					} else {
						level[i-29] = 7
						level[i] = 2
					}
					break out5
				case 3:
					death = 1
					break out5
				}
			case "down":
				switch level[i+29] {
				case 1:
					if conpoint == 0 {
						conpoint = 1
						level[i+29] = 7
						level[i] = 2
					} else {
						level[i+29] = 7
						level[i] = 1
					}
					break out5
				case 2:
					if conpoint == 1 {
						conpoint = 0
						level[i+29] = 7
						level[i] = 1
					} else {
						level[i+29] = 7
						level[i] = 2
					}
					break out5
				case 3:
					death = 1
					break out5
				}
			case "left":
				switch level[i-1] {
				case 1:
					if conpoint == 0 {
						conpoint = 1
						level[i-1] = 7
						level[i] = 2
					} else {
						level[i-1] = 7
						level[i] = 1
					}
					break out5
				case 2:
					if conpoint == 1 {
						conpoint = 0
						level[i-1] = 7
						level[i] = 1
					} else {
						level[i-1] = 7
						level[i] = 2
					}
					break out5
				case 3:
					death = 1
					break out5
				case 8:
					if conpoint == 1 {
						if level[i+27] == 2 {
							conpoint = 0
						}
						level[i+27] = 7
						level[i] = 1
					} else {
						if level[i+27] == 1 {
							conpoint = 1
						}
						level[i+27] = 7
						level[i] = 2
					}
					break out5
				}
			case "right":
				switch level[i+1] {
				case 1:
					if conpoint == 0 {
						conpoint = 1
						level[i+1] = 7
						level[i] = 2
					} else {
						level[i+1] = 7
						level[i] = 1
					}
					break out5
				case 2:
					if conpoint == 1 {
						conpoint = 0
						level[i+1] = 7
						level[i] = 1
					} else {
						level[i+1] = 7
						level[i] = 2
					}
					break out5
				case 3:
					death = 1
					break out5
				case 8:
					if conpoint == 1 {
						if level[i-27] == 2 {
							conpoint = 0
						}
						level[i-27] = 7
						level[i] = 1
					} else {
						if level[i-27] == 1 {
							conpoint = 1
						}
						level[i-27] = 7
						level[i] = 2
					}
					break out5
				}
			}
		}
	}
}

func findx(i int) int {
	for found := 0; found == 0; {
		if i > 28 {
			i = i - 29
		}
		if i < 29 {
			found = 1
		}
	}
	return i
}

func findy(i int) int {
	var y int = 0
	for found := 0; found == 0; {
		if i > 28 {
			i -= 29
			y += 1
		}
		if i < 29 {
			found = 1
		}
	}
	return y
}

func ghostturncheck(ghostturn, level []int, i int, dir string) {
	switch dir {
	case "up":
		if level[i-29] == 0 || level[i-29] == 5 || level[i-29] == 6 || level[i-29] == 7 {
			ghostturn[0] = 0
			ghostturn[1] = 1

			if level[i+29] == 0 || level[i+29] == 5 || level[i+29] == 6 || level[i+29] == 7 {
				ghostturn[2] = 1
			} else {
				ghostturn[2] = 0
			}
			if level[i-1] == 0 || level[i-1] == 5 || level[i-1] == 6 || level[i-1] == 7 {
				ghostturn[3] = 1
			} else {
				ghostturn[3] = 0
			}
			if level[i+1] == 0 || level[i+1] == 5 || level[i+1] == 6 || level[i+1] == 7 {
				ghostturn[4] = 1
			} else {
				ghostturn[4] = 0
			}
		} else {
			if (level[i-1] == 1 || level[i-1] == 2 || level[i-1] == 3 || level[i-1] == 8) && (level[i+1] == 1 || level[i+1] == 2 || level[i+1] == 3 || level[i+1] == 8) {
				ghostturn[0] = 0
			} else {
				ghostturn[0] = 1
			}
			ghostturn[1] = 0
		}
	case "down":
		if level[i+29] == 0 || level[i+29] == 5 || level[i+29] == 6 || level[i+29] == 7 {
			ghostturn[0] = 0
			ghostturn[2] = 1

			if level[i-29] == 0 || level[i-29] == 5 || level[i-29] == 6 || level[i-29] == 7 {
				ghostturn[1] = 1
			} else {
				ghostturn[1] = 0
			}
			if level[i-1] == 0 || level[i-1] == 5 || level[i-1] == 6 || level[i-1] == 7 {
				ghostturn[3] = 1
			} else {
				ghostturn[3] = 0
			}
			if level[i+1] == 0 || level[i+1] == 5 || level[i+1] == 6 || level[i+1] == 7 {
				ghostturn[4] = 1
			} else {
				ghostturn[4] = 0
			}
		} else {
			if (level[i-1] == 1 || level[i-1] == 2 || level[i-1] == 3 || level[i-1] == 8) && (level[i+1] == 1 || level[i+1] == 2 || level[i+1] == 3 || level[i+1] == 8) {
				ghostturn[0] = 0
			} else {
				ghostturn[0] = 1
			}
			ghostturn[2] = 0
		}
	case "left":
		if level[i-1] == 0 || level[i-1] == 5 || level[i-1] == 6 || level[i-1] == 7 {
			ghostturn[0] = 0
			ghostturn[3] = 1

			if level[i+29] == 0 || level[i+29] == 5 || level[i+29] == 6 || level[i+29] == 7 {
				ghostturn[2] = 1
			} else {
				ghostturn[2] = 0
			}
			if level[i-29] == 0 || level[i-29] == 5 || level[i-29] == 6 || level[i-29] == 7 {
				ghostturn[1] = 1
			} else {
				ghostturn[1] = 0
			}
			if level[i+1] == 0 || level[i+1] == 5 || level[i+1] == 6 || level[i+1] == 7 {
				ghostturn[4] = 1
			} else {
				ghostturn[4] = 0
			}
		} else {
			if (level[i-29] == 1 || level[i-29] == 2 || level[i-29] == 3 || level[i-29] == 8) && (level[i+29] == 1 || level[i+29] == 2 || level[i+29] == 3 || level[i+29] == 8) {
				ghostturn[0] = 0
			} else {
				ghostturn[0] = 1
			}
			ghostturn[3] = 0
		}
	case "right":
		if level[i+1] == 0 || level[i+1] == 5 || level[i+1] == 6 || level[i+1] == 7 {
			ghostturn[0] = 0
			ghostturn[4] = 1

			if level[i+29] == 0 || level[i+29] == 5 || level[i+29] == 6 || level[i+29] == 7 {
				ghostturn[2] = 1
			} else {
				ghostturn[2] = 0
			}
			if level[i-1] == 0 || level[i-1] == 5 || level[i-1] == 6 || level[i-1] == 7 {
				ghostturn[3] = 1
			} else {
				ghostturn[3] = 0
			}
			if level[i-29] == 0 || level[i-29] == 5 || level[i-29] == 6 || level[i-29] == 7 {
				ghostturn[1] = 1
			} else {
				ghostturn[1] = 0
			}
		} else {
			if (level[i-29] == 1 || level[i-29] == 2 || level[i-29] == 3 || level[i-29] == 8) && (level[i+29] == 1 || level[i+29] == 2 || level[i+29] == 3 || level[i+29] == 8) {
				ghostturn[0] = 0
			} else {
				ghostturn[0] = 1
			}
			ghostturn[4] = 0
		}
	}
}

func compareI(table []ainode, tocompare int) int {
	for i, v := range []ainode(table) {
		if v.loc == tocompare {
			return i
		}
	}
	return -1
}
