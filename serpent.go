package main

import (
	"math/rand"
	"time"

	tl "github.com/JoelOtter/termloop"
)

type Snake struct {
	segments  []tl.Entity
	direction string
	tickCount int
	growth    int
}

type Food struct {
	*tl.Entity
	placed bool
}

func NewFood() *Food {
	return &Food{
		Entity: tl.NewEntityFromCanvas(2, 2, tl.CanvasFromString("¤")),
		placed: false,
	}
}

func (f *Food) Tick(event tl.Event) {
	// Check if food has been placed, if not, place the food
	if !f.placed {
		width, height := game.Screen().Size()
		if width > 0 && height > 0 {
			f.PlaceFood(width, height)
			f.placed = true
		}
	}
}

func (f *Food) PlaceFood(levelWidth, levelHeight int) {
	rand.Seed(time.Now().UnixNano())
	foodX := rand.Intn(levelWidth-4) + 2 // Pad to avoid placing food on the border
	foodY := rand.Intn(levelHeight-4) + 2
	f.SetPosition(foodX, foodY)
}

func (f *Food) Draw(screen *tl.Screen) {
	// Draw food after it has been placed
	if f.placed {
		f.Entity.Draw(screen)
	}
}

func NewSnake(x, y int) *Snake {
	snake := &Snake{
		segments:  make([]tl.Entity, 0),
		direction: "right", // initial direction
		tickCount: 0,
		growth:    0,
	}
	// Initialize snake with 3 segments
	for i := 0; i < 3; i++ {
		snake.segments = append(snake.segments, *tl.NewEntity(1, 1, x-i*2, y)) // We are moving 2 cells per tick
	}
	return snake
}

func (snake *Snake) Draw(screen *tl.Screen) {
	for _, segment := range snake.segments {
		x, y := segment.Position()
		screen.RenderCell(x, y, &tl.Cell{Fg: tl.ColorGreen, Ch: '■'})
	}
}

func (snake *Snake) Tick(event tl.Event) {
	// Handle direction change
	if event.Type == tl.EventKey {
		switch event.Key {
		case tl.KeyArrowRight:
			if snake.direction != "left" {
				snake.direction = "right"
			}
		case tl.KeyArrowLeft:
			if snake.direction != "right" {
				snake.direction = "left"
			}
		case tl.KeyArrowUp:
			if snake.direction != "down" {
				snake.direction = "up"
			}
		case tl.KeyArrowDown:
			if snake.direction != "up" {
				snake.direction = "down"
			}
		}
	}

	// Update snake every two ticks
	snake.tickCount++
	if snake.tickCount >= 2 {
		snake.tickCount = 0

		head := snake.segments[0]
		x, y := head.Position()

		// Determine the new head position based on the direction
		newX, newY := x, y
		switch snake.direction {
		case "right":
			newX = x + 2
		case "left":
			newX = x - 2
		case "up":
			newY = y - 1
		case "down":
			newY = y + 1
		}

		// Check for collision with food
		foodX, foodY := food.Position()
		if newX == foodX && newY == foodY {
			snake.growth += 1   // Grow by one segment for each food eaten
			food.placed = false // Signal to re-place the food

			// Re-place the food immediately
			width, height := game.Screen().Size()
			food.PlaceFood(width, height)
		}

		// Move and grow the snake as necessary
		if snake.growth > 0 {
			// Append a new segment at the front (new head position)
			newSegment := tl.NewEntity(1, 1, newX, newY) // Create a new segment
			// Create a temporary slice with the new segment
			tempSegments := make([]tl.Entity, len(snake.segments)+1)
			tempSegments[0] = *newSegment // Add the new segment at the beginning
			// Copy the old segments into the new slice
			for i, segment := range snake.segments {
				tempSegments[i+1] = segment
			}
			snake.segments = tempSegments // Replace the old slice with the new slice
			snake.growth--
		} else {
			// Move the snake forward: set new positions for all segments
			for i := len(snake.segments) - 1; i > 0; i-- {
				prevX, prevY := snake.segments[i-1].Position()
				snake.segments[i].SetPosition(prevX, prevY)
			}
			snake.segments[0].SetPosition(newX, newY) // Move the head to the new position
		}
	}
}

var food *Food
var game *tl.Game

func main() {
	game = tl.NewGame()
	game.Screen().SetFps(10) // Set FPS

	level := tl.NewBaseLevel(tl.Cell{
		Bg: tl.ColorBlack,
		Fg: tl.ColorWhite,
		Ch: ' ',
	})

	// Initialize entities
	snake := NewSnake(20, 20)
	food = NewFood() // Food will place itself during the game loop

	// Add entities to level
	level.AddEntity(snake)
	level.AddEntity(food)

	// Configure and start the game
	game.Screen().SetLevel(level)
	game.Start()
}
