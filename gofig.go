package gofig

import (
	"strings"
)

// ANSI color codes
const (
	ColorReset   = "\033[0m"
	ColorRed     = "\033[31m"
	ColorGreen   = "\033[32m"
	ColorYellow  = "\033[33m"
	ColorBlue    = "\033[34m"
	ColorMagenta = "\033[35m"
	ColorCyan    = "\033[36m"
	ColorWhite   = "\033[37m"
	ColorBrightRed     = "\033[91m"
	ColorBrightGreen   = "\033[92m"
	ColorBrightYellow  = "\033[93m"
	ColorBrightBlue    = "\033[94m"
	ColorBrightMagenta = "\033[95m"
	ColorBrightCyan    = "\033[96m"
)

// Config holds settings for block text rendering
type Config struct {
	// Scale multiplies the size (1 = normal, 2 = double, etc.)
	Scale int
	// Char is the block character to use (default: █)
	Char string
	// Space is the space character (default: space)
	Space string
	// Color is the ANSI color code (e.g., ColorGreen)
	Color string
}

// DefaultConfig returns default configuration
func DefaultConfig() Config {
	return Config{
		Scale: 1,
		Char:  "█",
		Space: " ",
	}
}

// BlockFont generates large text using block characters
type BlockFont struct {
	config Config
	chars  map[rune][]string
	height int
}

// New creates a new block font with default config
func New() *BlockFont {
	return NewWithConfig(DefaultConfig())
}

// NewWithConfig creates a new block font with custom config
func NewWithConfig(config Config) *BlockFont {
	if config.Scale < 1 {
		config.Scale = 1
	}
	if config.Char == "" {
		config.Char = "█"
	}
	if config.Space == "" {
		config.Space = " "
	}

	bf := &BlockFont{
		config: config,
		chars:  make(map[rune][]string),
		height: 5,
	}
	bf.initChars()
	return bf
}

// Render converts text to block characters
func (bf *BlockFont) Render(text string) string {
	text = strings.ToUpper(text)
	lines := make([]string, bf.height*bf.config.Scale)

	for _, ch := range text {
		pattern, ok := bf.chars[ch]
		if !ok {
			pattern = bf.chars[' ']
		}
		for i := 0; i < bf.height; i++ {
			// Scale the line horizontally
			scaledLine := bf.scaleLine(pattern[i])
			// Repeat for vertical scaling
			for s := 0; s < bf.config.Scale; s++ {
				lines[i*bf.config.Scale+s] += scaledLine + strings.Repeat(bf.config.Space, bf.config.Scale)
			}
		}
	}

	result := strings.Join(lines, "\n")

	// Apply color if set
	if bf.config.Color != "" {
		result = bf.config.Color + result + ColorReset
	}

	return result
}

// SetColor sets the text color
func (bf *BlockFont) SetColor(color string) {
	bf.config.Color = color
}

// scaleLine scales a single line horizontally and applies custom char
func (bf *BlockFont) scaleLine(line string) string {
	var result strings.Builder
	for _, ch := range line {
		char := bf.config.Space
		if ch == '█' {
			char = bf.config.Char
		}
		result.WriteString(strings.Repeat(char, bf.config.Scale))
	}
	return result.String()
}

// SetScale changes the scale factor
func (bf *BlockFont) SetScale(scale int) {
	if scale < 1 {
		scale = 1
	}
	bf.config.Scale = scale
}

// SetChar changes the block character
func (bf *BlockFont) SetChar(char string) {
	bf.config.Char = char
	bf.initChars()
}

func (bf *BlockFont) initChars() {
	b := "█"
	s := " "

	bf.chars['A'] = []string{
		s + b + b + b + s,
		b + s + s + s + b,
		b + b + b + b + b,
		b + s + s + s + b,
		b + s + s + s + b,
	}
	bf.chars['B'] = []string{
		b + b + b + b + s,
		b + s + s + s + b,
		b + b + b + b + s,
		b + s + s + s + b,
		b + b + b + b + s,
	}
	bf.chars['C'] = []string{
		s + b + b + b + b,
		b + s + s + s + s,
		b + s + s + s + s,
		b + s + s + s + s,
		s + b + b + b + b,
	}
	bf.chars['D'] = []string{
		b + b + b + b + s,
		b + s + s + s + b,
		b + s + s + s + b,
		b + s + s + s + b,
		b + b + b + b + s,
	}
	bf.chars['E'] = []string{
		b + b + b + b + b,
		b + s + s + s + s,
		b + b + b + b + s,
		b + s + s + s + s,
		b + b + b + b + b,
	}
	bf.chars['F'] = []string{
		b + b + b + b + b,
		b + s + s + s + s,
		b + b + b + b + s,
		b + s + s + s + s,
		b + s + s + s + s,
	}
	bf.chars['G'] = []string{
		s + b + b + b + b,
		b + s + s + s + s,
		b + s + b + b + b,
		b + s + s + s + b,
		s + b + b + b + s,
	}
	bf.chars['H'] = []string{
		b + s + s + s + b,
		b + s + s + s + b,
		b + b + b + b + b,
		b + s + s + s + b,
		b + s + s + s + b,
	}
	bf.chars['I'] = []string{
		b + b + b + b + b,
		s + s + b + s + s,
		s + s + b + s + s,
		s + s + b + s + s,
		b + b + b + b + b,
	}
	bf.chars['J'] = []string{
		s + s + b + b + b,
		s + s + s + b + s,
		s + s + s + b + s,
		b + s + s + b + s,
		s + b + b + s + s,
	}
	bf.chars['K'] = []string{
		b + s + s + s + b,
		b + s + s + b + s,
		b + b + b + s + s,
		b + s + s + b + s,
		b + s + s + s + b,
	}
	bf.chars['L'] = []string{
		b + s + s + s + s,
		b + s + s + s + s,
		b + s + s + s + s,
		b + s + s + s + s,
		b + b + b + b + b,
	}
	bf.chars['M'] = []string{
		b + s + s + s + b,
		b + b + s + b + b,
		b + s + b + s + b,
		b + s + s + s + b,
		b + s + s + s + b,
	}
	bf.chars['N'] = []string{
		b + s + s + s + b,
		b + b + s + s + b,
		b + s + b + s + b,
		b + s + s + b + b,
		b + s + s + s + b,
	}
	bf.chars['O'] = []string{
		s + b + b + b + s,
		b + s + s + s + b,
		b + s + s + s + b,
		b + s + s + s + b,
		s + b + b + b + s,
	}
	bf.chars['P'] = []string{
		b + b + b + b + s,
		b + s + s + s + b,
		b + b + b + b + s,
		b + s + s + s + s,
		b + s + s + s + s,
	}
	bf.chars['Q'] = []string{
		s + b + b + b + s,
		b + s + s + s + b,
		b + s + s + s + b,
		b + s + s + b + s,
		s + b + b + s + b,
	}
	bf.chars['R'] = []string{
		b + b + b + b + s,
		b + s + s + s + b,
		b + b + b + b + s,
		b + s + s + b + s,
		b + s + s + s + b,
	}
	bf.chars['S'] = []string{
		s + b + b + b + b,
		b + s + s + s + s,
		s + b + b + b + s,
		s + s + s + s + b,
		b + b + b + b + s,
	}
	bf.chars['T'] = []string{
		b + b + b + b + b,
		s + s + b + s + s,
		s + s + b + s + s,
		s + s + b + s + s,
		s + s + b + s + s,
	}
	bf.chars['U'] = []string{
		b + s + s + s + b,
		b + s + s + s + b,
		b + s + s + s + b,
		b + s + s + s + b,
		s + b + b + b + s,
	}
	bf.chars['V'] = []string{
		b + s + s + s + b,
		b + s + s + s + b,
		b + s + s + s + b,
		s + b + s + b + s,
		s + s + b + s + s,
	}
	bf.chars['W'] = []string{
		b + s + s + s + b,
		b + s + s + s + b,
		b + s + b + s + b,
		b + b + s + b + b,
		b + s + s + s + b,
	}
	bf.chars['X'] = []string{
		b + s + s + s + b,
		s + b + s + b + s,
		s + s + b + s + s,
		s + b + s + b + s,
		b + s + s + s + b,
	}
	bf.chars['Y'] = []string{
		b + s + s + s + b,
		s + b + s + b + s,
		s + s + b + s + s,
		s + s + b + s + s,
		s + s + b + s + s,
	}
	bf.chars['Z'] = []string{
		b + b + b + b + b,
		s + s + s + b + s,
		s + s + b + s + s,
		s + b + s + s + s,
		b + b + b + b + b,
	}

	// Numbers
	bf.chars['0'] = []string{
		s + b + b + b + s,
		b + s + s + b + b,
		b + s + b + s + b,
		b + b + s + s + b,
		s + b + b + b + s,
	}
	bf.chars['1'] = []string{
		s + s + b + s + s,
		s + b + b + s + s,
		s + s + b + s + s,
		s + s + b + s + s,
		b + b + b + b + b,
	}
	bf.chars['2'] = []string{
		s + b + b + b + s,
		b + s + s + s + b,
		s + s + b + b + s,
		s + b + s + s + s,
		b + b + b + b + b,
	}
	bf.chars['3'] = []string{
		b + b + b + b + s,
		s + s + s + s + b,
		s + b + b + b + s,
		s + s + s + s + b,
		b + b + b + b + s,
	}
	bf.chars['4'] = []string{
		b + s + s + s + b,
		b + s + s + s + b,
		b + b + b + b + b,
		s + s + s + s + b,
		s + s + s + s + b,
	}
	bf.chars['5'] = []string{
		b + b + b + b + b,
		b + s + s + s + s,
		b + b + b + b + s,
		s + s + s + s + b,
		b + b + b + b + s,
	}
	bf.chars['6'] = []string{
		s + b + b + b + s,
		b + s + s + s + s,
		b + b + b + b + s,
		b + s + s + s + b,
		s + b + b + b + s,
	}
	bf.chars['7'] = []string{
		b + b + b + b + b,
		s + s + s + s + b,
		s + s + s + b + s,
		s + s + b + s + s,
		s + s + b + s + s,
	}
	bf.chars['8'] = []string{
		s + b + b + b + s,
		b + s + s + s + b,
		s + b + b + b + s,
		b + s + s + s + b,
		s + b + b + b + s,
	}
	bf.chars['9'] = []string{
		s + b + b + b + s,
		b + s + s + s + b,
		s + b + b + b + b,
		s + s + s + s + b,
		s + b + b + b + s,
	}

	// Symbols
	bf.chars[' '] = []string{
		s + s + s + s + s,
		s + s + s + s + s,
		s + s + s + s + s,
		s + s + s + s + s,
		s + s + s + s + s,
	}
	bf.chars['!'] = []string{
		s + s + b + s + s,
		s + s + b + s + s,
		s + s + b + s + s,
		s + s + s + s + s,
		s + s + b + s + s,
	}
	bf.chars['?'] = []string{
		s + b + b + b + s,
		b + s + s + s + b,
		s + s + s + b + s,
		s + s + s + s + s,
		s + s + b + s + s,
	}
	bf.chars['.'] = []string{
		s + s + s + s + s,
		s + s + s + s + s,
		s + s + s + s + s,
		s + s + s + s + s,
		s + s + b + s + s,
	}
	bf.chars[','] = []string{
		s + s + s + s + s,
		s + s + s + s + s,
		s + s + s + s + s,
		s + s + b + s + s,
		s + b + s + s + s,
	}
	bf.chars[':'] = []string{
		s + s + s + s + s,
		s + s + b + s + s,
		s + s + s + s + s,
		s + s + b + s + s,
		s + s + s + s + s,
	}
	bf.chars['-'] = []string{
		s + s + s + s + s,
		s + s + s + s + s,
		b + b + b + b + b,
		s + s + s + s + s,
		s + s + s + s + s,
	}
	bf.chars['_'] = []string{
		s + s + s + s + s,
		s + s + s + s + s,
		s + s + s + s + s,
		s + s + s + s + s,
		b + b + b + b + b,
	}
	bf.chars['/'] = []string{
		s + s + s + s + b,
		s + s + s + b + s,
		s + s + b + s + s,
		s + b + s + s + s,
		b + s + s + s + s,
	}
	bf.chars['('] = []string{
		s + s + b + s + s,
		s + b + s + s + s,
		s + b + s + s + s,
		s + b + s + s + s,
		s + s + b + s + s,
	}
	bf.chars[')'] = []string{
		s + s + b + s + s,
		s + s + s + b + s,
		s + s + s + b + s,
		s + s + s + b + s,
		s + s + b + s + s,
	}
	bf.chars['<'] = []string{
		s + s + s + b + s,
		s + s + b + s + s,
		s + b + s + s + s,
		s + s + b + s + s,
		s + s + s + b + s,
	}
	bf.chars['>'] = []string{
		s + b + s + s + s,
		s + s + b + s + s,
		s + s + s + b + s,
		s + s + b + s + s,
		s + b + s + s + s,
	}
	bf.chars['='] = []string{
		s + s + s + s + s,
		b + b + b + b + b,
		s + s + s + s + s,
		b + b + b + b + b,
		s + s + s + s + s,
	}
	bf.chars['+'] = []string{
		s + s + s + s + s,
		s + s + b + s + s,
		s + b + b + b + s,
		s + s + b + s + s,
		s + s + s + s + s,
	}
	bf.chars['#'] = []string{
		s + b + s + b + s,
		b + b + b + b + b,
		s + b + s + b + s,
		b + b + b + b + b,
		s + b + s + b + s,
	}
	bf.chars['@'] = []string{
		s + b + b + b + s,
		b + s + b + s + b,
		b + s + b + b + b,
		b + s + s + s + s,
		s + b + b + b + b,
	}
	bf.chars['*'] = []string{
		s + s + s + s + s,
		b + s + b + s + b,
		s + s + b + s + s,
		b + s + b + s + b,
		s + s + s + s + s,
	}
	bf.chars['%'] = []string{
		b + b + s + s + b,
		b + b + s + b + s,
		s + s + b + s + s,
		s + b + s + b + b,
		b + s + s + b + b,
	}
	bf.chars['$'] = []string{
		s + b + b + b + b,
		b + s + b + s + s,
		s + b + b + b + s,
		s + s + b + s + b,
		b + b + b + b + s,
	}
	bf.chars['&'] = []string{
		s + b + b + s + s,
		b + s + s + b + s,
		s + b + b + s + b,
		b + s + s + b + s,
		s + b + b + s + b,
	}
	bf.chars['\''] = []string{
		s + s + b + s + s,
		s + s + b + s + s,
		s + s + s + s + s,
		s + s + s + s + s,
		s + s + s + s + s,
	}
	bf.chars['"'] = []string{
		s + b + s + b + s,
		s + b + s + b + s,
		s + s + s + s + s,
		s + s + s + s + s,
		s + s + s + s + s,
	}
}

// Render is a convenience function to render text with default config
func Render(text string) string {
	return New().Render(text)
}

// RenderWithScale renders text with specified scale
func RenderWithScale(text string, scale int) string {
	config := DefaultConfig()
	config.Scale = scale
	return NewWithConfig(config).Render(text)
}

