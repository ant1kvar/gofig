package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"gofig"
)

var colors = map[string]string{
	"red":     gofig.ColorRed,
	"green":   gofig.ColorGreen,
	"yellow":  gofig.ColorYellow,
	"blue":    gofig.ColorBlue,
	"magenta": gofig.ColorMagenta,
	"cyan":    gofig.ColorCyan,
	"white":   gofig.ColorWhite,
}

var animTypes = map[string]gofig.AnimationType{
	"blink":    gofig.AnimBlink,
	"pulse":    gofig.AnimPulse,
	"wave":     gofig.AnimWave,
	"typing":   gofig.AnimTyping,
	"glitch":   gofig.AnimGlitch,
	"sequence": gofig.AnimSequence,
	"random":   gofig.AnimRandom,
}

func main() {
	// Основные настройки
	scale := flag.Int("scale", 1, "Scale factor (1-5)")
	char := flag.String("char", "█", "Block character to use")
	space := flag.String("space", " ", "Space character (e.g., '.', '_')")
	color := flag.String("color", "", "Color: red, green, yellow, blue, magenta, cyan, white")

	// Настройки анимации
	anim := flag.String("anim", "", "Animation: blink, pulse, wave, typing, glitch, sequence, random")
	interval := flag.Int("interval", 100, "Animation interval in ms")
	chance := flag.Float64("chance", 0.3, "Effect chance (0.0-1.0)")
	min := flag.Int("min", 1, "Minimum affected letters")
	max := flag.Int("max", 3, "Maximum affected letters")
	waveWidth := flag.Int("wave-width", 3, "Wave width (for wave animation)")
	switchFrames := flag.Int("switch", 30, "Frames before switching animation (for random mode)")
	duration := flag.Int("duration", 0, "Animation duration in seconds (0 = infinite)")

	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		fmt.Println("Usage: textblock [options] <text>")
		fmt.Println("\nOptions:")
		flag.PrintDefaults()
		fmt.Println("\nColors: red, green, yellow, blue, magenta, cyan, white")
		fmt.Println("\nAnimations:")
		fmt.Println("  blink    - Random letter blinking")
		fmt.Println("  pulse    - Whole text pulses")
		fmt.Println("  wave     - Wave moves through text")
		fmt.Println("  typing   - Typewriter effect")
		fmt.Println("  glitch   - Glitch/corruption effect")
		fmt.Println("  sequence - Sequential letter blinking")
		fmt.Println("  random   - Randomly switches between animations")
		fmt.Println("\nExamples:")
		fmt.Println("  textblock Hello")
		fmt.Println("  textblock -scale=2 -color=green OK")
		fmt.Println("  textblock -char='#' -space='.' DOTS")
		fmt.Println("  textblock -anim=blink ERROR")
		fmt.Println("  textblock -anim=wave -color=cyan LOADING")
		fmt.Println("  textblock -anim=typing -interval=150 HELLO")
		fmt.Println("  textblock -anim=glitch -chance=0.5 -max=5 SYSTEM")
		fmt.Println("  textblock -anim=pulse -interval=500 ALERT")
		os.Exit(1)
	}

	text := strings.Join(args, " ")

	// Настройки шрифта
	fontConfig := gofig.DefaultConfig()
	fontConfig.Scale = *scale
	fontConfig.Char = *char
	fontConfig.Space = *space
	if c, ok := colors[*color]; ok {
		fontConfig.Color = c
	}

	// Если анимация не задана - просто вывести текст
	if *anim == "" {
		bf := gofig.NewWithConfig(fontConfig)
		fmt.Println(bf.Render(text))
		return
	}

	// Настройки анимации
	animConfig := gofig.DefaultAnimConfig()
	if animType, ok := animTypes[*anim]; ok {
		animConfig.Type = animType
	} else {
		fmt.Printf("Unknown animation: %s\n", *anim)
		fmt.Println("Available: blink, pulse, wave, typing, glitch, sequence, random")
		os.Exit(1)
	}

	animConfig.Interval = time.Duration(*interval) * time.Millisecond
	animConfig.Chance = *chance
	animConfig.Min = *min
	animConfig.Max = *max
	animConfig.WaveWidth = *waveWidth
	animConfig.RandomSwitchFrames = *switchFrames

	// Запуск анимации
	animation := gofig.NewAnimationWithConfig(text, fontConfig, animConfig)

	if *duration > 0 {
		go func() {
			time.Sleep(time.Duration(*duration) * time.Second)
			animation.Stop()
		}()
	}

	animation.Start()
}
