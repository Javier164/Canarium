package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

var color = sdl.Color{R: 255, G: 255, B: 255, A: 255}

func Text(renderer *sdl.Renderer, text, path string, size int, x, y int32, center bool, mode uint8) {
	font, err := ttf.OpenFont(path, size)
	if err != nil {
		log.Fatalf("Failed to load font: %s", err)
	}
	defer font.Close()

	surface, err := font.RenderUTF8Blended(text, color)
	if err != nil {
		log.Fatalf("Failed to create surface: %s", err)
	}
	defer surface.Free()

	texture, err := renderer.CreateTextureFromSurface(surface)
	if err != nil {
		log.Fatalf("Failed to create texture: %s", err)
	}
	defer texture.Destroy()

	rect := &sdl.Rect{X: x, Y: y, W: surface.W, H: surface.H}
	if center {
		switch mode {
		case 0:
			rect.X = (960 - surface.W) / 2
			break
		case 1:
			rect.Y = (720 - surface.H) / 2
		case 2:
			rect.X = (960 - surface.W) / 2
			rect.Y = (720 - surface.H) / 2
		}
	}
	renderer.Copy(texture, nil, rect)
}

func WordWrap(text string, width uint8) []string {
	words := strings.Split(text, " ")
	var lines []string
	var line string

	for _, word := range words {
		if len(line)+len(word) <= int(width) {
			if line == "" {
				line = word
			} else {
				line += " " + word
			}
		} else {
			lines = append(lines, line)
			line = word
		}
	}
	lines = append(lines, line)

	return lines
}

func ParseHexColor(hex string) (uint8, uint8, uint8, error) {
	var r, g, b uint32
	_, err := fmt.Sscanf(hex, "#%02x%02x%02x", &r, &g, &b)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("Failed to parse color string: %s", err)
	}
	return uint8(r), uint8(g), uint8(b), nil
}

func InterruptHandler() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		log.Println("Ctrl+C pressed. Exiting.")
		os.Exit(0)
	}()
}
