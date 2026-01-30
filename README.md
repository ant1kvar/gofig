<p align="center">
  <h1 align="center">🎨 GoFig</h1>
  <p align="center">
    <strong>Beautiful block text and animations for your terminal</strong>
  </p>
  <p align="center">
    A Go library for generating large block-character text with colors and animations
  </p>
</p>

<p align="center">
  <a href="#installation">Installation</a> •
  <a href="#quick-start">Quick Start</a> •
  <a href="#cli-usage">CLI</a> •
  <a href="#animations">Animations</a> •
  <a href="#api">API</a>
</p>

---
![demo](https://github.com/user-attachments/assets/ebd22497-00be-4d17-a439-d4ab8f53dcbd)
---
## Features

- 🔤 **Block Text** — Convert text to large block characters (█)
- 🎨 **Colors** — Full ANSI color support
- ✨ **Animations** — 7 built-in animation types
- 📐 **Scaling** — Adjustable text size
- 🎯 **Customizable** — Custom characters and spacing
- 🚀 **Zero Dependencies** — Pure Go, no external packages

## Installation

```bash
go get github.com/ant1kvar/gofig
```

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/ant1kvar/gofig"
)

func main() {
    fmt.Println(gofig.Render("HELLO"))
}
```

Output:
```
█   █ █████ █     █      ███
█   █ █     █     █     █   █
█████ ████  █     █     █   █
█   █ █     █     █     █   █
█   █ █████ █████ █████  ███
```

## CLI Usage

### Build

```bash
go build -o gofig ./example/
```

### Basic Usage

```bash
# Simple text
./gofig HELLO

# With color
./gofig -color=green OK
./gofig -color=red ERROR

# Scaled up
./gofig -scale=2 BIG

# Custom characters
./gofig -char='#' -space='.' TEXT
```

### Animations

```bash
# Blinking letters
./gofig -anim=blink ERROR

# Pulsing text
./gofig -anim=pulse -interval=500 ALERT

# Wave effect
./gofig -anim=wave LOADING

# Typewriter effect
./gofig -anim=typing HELLO

# Glitch effect
./gofig -anim=glitch SYSTEM

# Sequential blinking
./gofig -anim=sequence SCAN

# Random animation switching
./gofig -anim=random CHAOS
```

### All Options

```bash
./gofig -help
```

| Flag | Description | Default |
|------|-------------|---------|
| `-scale` | Scale factor (1-5) | 1 |
| `-char` | Block character | █ |
| `-space` | Space character | (space) |
| `-color` | Text color | (none) |
| `-anim` | Animation type | (none) |
| `-interval` | Frame interval (ms) | 100 |
| `-chance` | Effect probability (0.0-1.0) | 0.3 |
| `-min` | Min affected letters | 1 |
| `-max` | Max affected letters | 3 |
| `-wave-width` | Wave width | 3 |
| `-switch` | Frames before animation switch (random) | 30 |
| `-duration` | Duration in seconds (0 = infinite) | 0 |

## Animations

| Type | Description |
|------|-------------|
| `blink` | Random letter blinking |
| `pulse` | Whole text pulses on/off |
| `wave` | Wave moves through text |
| `typing` | Typewriter effect |
| `glitch` | Glitch/corruption effect |
| `sequence` | Sequential letter blinking |
| `random` | Randomly switches between all animations |

## Colors

Available colors: `red`, `green`, `yellow`, `blue`, `magenta`, `cyan`, `white`

```bash
./gofig -color=green SUCCESS
./gofig -color=red -anim=blink ERROR
./gofig -color=cyan -anim=wave LOADING
```

## API

### Basic Rendering

```go
// Simple render
text := gofig.Render("HELLO")

// With scale
text := gofig.RenderWithScale("BIG", 2)

// Full configuration
config := gofig.DefaultConfig()
config.Scale = 2
config.Char = "#"
config.Space = "."
config.Color = gofig.ColorGreen

bf := gofig.NewWithConfig(config)
fmt.Println(bf.Render("OK"))
```

### Animations

```go
// Quick functions (blocking, Ctrl+C to exit)
gofig.Blink("ERROR")
gofig.Pulse("ALERT")
gofig.Wave("LOADING")
gofig.Typing("HELLO")
gofig.Glitch("SYSTEM")
gofig.Sequence("SCAN")
gofig.Random("CHAOS")

// Timed animation
gofig.BlinkFor("LOADING", 5*time.Second)
```

### Custom Animation

```go
fontConfig := gofig.DefaultConfig()
fontConfig.Scale = 2
fontConfig.Color = gofig.ColorCyan

animConfig := gofig.DefaultAnimConfig()
animConfig.Type = gofig.AnimWave
animConfig.Interval = 100 * time.Millisecond
animConfig.WaveWidth = 4

gofig.NewAnimationWithConfig("LOADING", fontConfig, animConfig).Start()
```

### Background Animation

```go
anim := gofig.NewAnimation("STATUS")
anim.StartAsync()  // Non-blocking

// Do other work...
time.Sleep(3 * time.Second)

anim.Stop()
```

### Types

```go
// Font configuration
type Config struct {
    Scale int    // Size multiplier
    Char  string // Block character (default: █)
    Space string // Space character (default: " ")
    Color string // ANSI color code
}

// Animation configuration
type AnimConfig struct {
    Type               AnimationType
    Interval           time.Duration
    Chance             float64
    Min                int
    Max                int
    GlitchChars        string
    WaveWidth          int
    RandomSwitchFrames int
}
```

### Color Constants

```go
gofig.ColorRed
gofig.ColorGreen
gofig.ColorYellow
gofig.ColorBlue
gofig.ColorMagenta
gofig.ColorCyan
gofig.ColorWhite
gofig.ColorBrightRed
gofig.ColorBrightGreen
gofig.ColorBrightYellow
gofig.ColorBrightBlue
gofig.ColorBrightMagenta
gofig.ColorBrightCyan
```

## Supported Characters

- Letters: `A-Z` (auto-converts to uppercase)
- Numbers: `0-9`
- Symbols: `! ? . , : - _ / ( ) < > = + # @ * % $ & ' "`

## Examples

### Startup Banner

```go
func main() {
    config := gofig.DefaultConfig()
    config.Color = gofig.ColorCyan

    bf := gofig.NewWithConfig(config)
    fmt.Println(bf.Render("MYAPP"))
    fmt.Println("v1.0.0")
}
```

### Loading Animation

```go
func showLoading() {
    gofig.Wave("LOADING")
}
```

### Error Display

```go
func showError(msg string) {
    config := gofig.DefaultConfig()
    config.Color = gofig.ColorRed
    config.Scale = 2

    bf := gofig.NewWithConfig(config)
    fmt.Println(bf.Render("ERROR"))
    fmt.Println(msg)
}
```

## License

MIT License - see [LICENSE](LICENSE) for details.

---

<p align="center">
  Made with ❤️ in Go
</p>
