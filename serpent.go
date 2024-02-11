package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	tl "github.com/JoelOtter/termloop"
)

type Coordinates struct {
	X, Y int
}

type Snake struct {
	body      []Coordinates
	direction string
	tickCount int
	growth    int
}

type Food struct {
	*tl.Entity
	placed bool
}

var foodPodMappings []FoodPodMapping

type FoodPodMapping struct {
	foodEntity *Food
	podInfo    PodInfo
}

const (
	LevelWidth  = 80
	LevelHeight = 24
)

func NewFood() *Food {
	return &Food{
		Entity: tl.NewEntityFromCanvas(2, 2, tl.CanvasFromString("O")),
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
	foodX := rand.Intn(LevelWidth-4) + 2
	foodY := rand.Intn(LevelHeight-4) + 2

	f.SetPosition(foodX, foodY)

	// Get a random pod name and namespace to associate with this food
	podInfo := getRandomPodInfo()
	if podInfo.Name != "" && podInfo.Namespace != "" {
		foodPodMappings = append(foodPodMappings, FoodPodMapping{
			foodEntity: f,
			podInfo:    podInfo,
		})
	}
}

func (f *Food) Draw(screen *tl.Screen) {
	// Draw food after it has been placed
	if f.placed {
		f.Entity.Draw(screen)
	}
}

func (f *Food) AtPosition(x, y int) bool {
	foodX, foodY := f.Position()
	// Check for collision in a wider range for X to accommodate faster horizontal movement
	return (x == foodX || x == foodX-1 || x == foodX+1) && y == foodY
}

func drawWalls(screen *tl.Screen) {
	// Top and bottom walls
	for x := 0; x < LevelWidth; x++ {
		screen.RenderCell(x, 0, &tl.Cell{Fg: tl.ColorWhite, Ch: '-'})             // Top wall
		screen.RenderCell(x, LevelHeight-1, &tl.Cell{Fg: tl.ColorWhite, Ch: '-'}) // Bottom wall
	}
	// Left and right walls
	for y := 0; y < LevelHeight; y++ {
		screen.RenderCell(0, y, &tl.Cell{Fg: tl.ColorWhite, Ch: '|'})            // Left wall
		screen.RenderCell(LevelWidth-1, y, &tl.Cell{Fg: tl.ColorWhite, Ch: '|'}) // Right wall
	}
}

func (snake *Snake) CollidesWithWalls() bool {
	head := snake.body[0]
	return head.X < 1 || head.Y < 1 || head.X >= LevelWidth-1 || head.Y >= LevelHeight-1
}

func (snake *Snake) CollidesWithSelf() bool {
	head := snake.body[0]
	for _, segment := range snake.body[1:] {
		if head.X == segment.X && head.Y == segment.Y {
			return true
		}
	}
	return false
}

func GameOver() {
	log.Println("Game Over!")
	os.Exit(0) // Exits the game
}

func NewSnake(x, y int) *Snake {
	snake := &Snake{
		direction: "right",
		tickCount: 0,
		growth:    0,
	}
	// Initialize snake with 3 segments
	for i := 0; i < 3; i++ {
		snake.body = append(snake.body, Coordinates{X: x - i*2, Y: y})
	}
	return snake
}

func (snake *Snake) Draw(screen *tl.Screen) {
	drawWalls(screen)
	for _, segment := range snake.body {
		screen.RenderCell(segment.X, segment.Y, &tl.Cell{Fg: tl.ColorGreen, Ch: '■'})
	}
}

var score int

func (snake *Snake) Tick(event tl.Event) {
	// Handle direction change input
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
		newHead := snake.body[0]
		switch snake.direction {
		case "right":
			newHead.X += 2
		case "left":
			newHead.X -= 2
		case "up":
			newHead.Y -= 1
		case "down":
			newHead.Y += 1
		}

		if food.AtPosition(newHead.X, newHead.Y) {
			snake.growth += 1
			food.placed = false
			score++
			scoreText.SetText(fmt.Sprintf("Score: %d", score))

			// Find which pod to delete
			for index, mapping := range foodPodMappings {
				if mapping.foodEntity == food {
					// Delete the pod and print its name
					deletePod(mapping.podInfo)
					deletionMessage := fmt.Sprintf("Oh no! Seems like you ate pod: %s in namespace %s", mapping.podInfo.Name, mapping.podInfo.Namespace)
					deletedPodText.SetText(deletionMessage)
					log.Println(deletionMessage)
					// Remove the mapping as the food has been eaten
					foodPodMappings = append(foodPodMappings[:index], foodPodMappings[index+1:]...)
					break
				}
			}
		}

		// Grow the snake if needed
		if snake.growth > 0 {
			snake.body = append([]Coordinates{newHead}, snake.body...)
			snake.growth--
		} else {
			snake.body = append([]Coordinates{newHead}, snake.body[:len(snake.body)-1]...)
		}

		if snake.CollidesWithWalls() || snake.CollidesWithSelf() {
			GameOver()
		}
	}
}

var food *Food
var game *tl.Game
var scoreText *tl.Text
var deletedPodText *tl.Text

func main() {

	// init k8s client
	initKubeClient()

	logFile, err := os.OpenFile("debug.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()

	log.SetOutput(logFile)

	game = tl.NewGame()
	game.Screen().SetFps(30)

	level := tl.NewBaseLevel(tl.Cell{
		Bg: tl.ColorBlack,
		Fg: tl.ColorWhite,
		Ch: ' ',
	})

	snake := NewSnake(20, 20)
	food = NewFood()

	level.AddEntity(snake)
	level.AddEntity(food)

	scoreText = tl.NewText(1, 0, "Score: 0", tl.ColorWhite, tl.ColorBlack)
	deletedPodText = tl.NewText(1, LevelHeight, "", tl.ColorWhite, tl.ColorBlack)
	level.AddEntity(scoreText)
	level.AddEntity(deletedPodText)

	game.Screen().SetLevel(level)
	game.Start()
}
