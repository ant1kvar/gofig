package gofig

import (
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

// AnimationType тип анимации
type AnimationType string

const (
	AnimBlink    AnimationType = "blink"    // Случайное мигание букв
	AnimPulse    AnimationType = "pulse"    // Весь текст мигает
	AnimWave     AnimationType = "wave"     // Волна по буквам
	AnimTyping   AnimationType = "typing"   // Печатание по буквам
	AnimGlitch   AnimationType = "glitch"   // Глитч-эффект
	AnimSequence AnimationType = "sequence" // Последовательное мигание
	AnimRandom   AnimationType = "random"   // Случайная смена анимаций
)

// AnimConfig общие настройки анимации
type AnimConfig struct {
	// Type тип анимации
	Type AnimationType
	// Interval между кадрами
	Interval time.Duration
	// Chance шанс эффекта (0.0 - 1.0)
	Chance float64
	// Min минимум затронутых букв
	Min int
	// Max максимум затронутых букв
	Max int
	// GlitchChars символы для глитча
	GlitchChars string
	// WaveWidth ширина волны (сколько букв видно)
	WaveWidth int
	// RandomSwitchFrames сколько кадров до смены анимации в random режиме
	RandomSwitchFrames int
}

// DefaultAnimConfig возвращает настройки по умолчанию
func DefaultAnimConfig() AnimConfig {
	return AnimConfig{
		Type:               AnimBlink,
		Interval:           100 * time.Millisecond,
		Chance:             0.3,
		Min:                1,
		Max:                3,
		GlitchChars:        "░▒▓█▄▀■□●○",
		WaveWidth:          3,
		RandomSwitchFrames: 30,
	}
}

// BlinkConfig для обратной совместимости
type BlinkConfig = AnimConfig

// DefaultBlinkConfig для обратной совместимости
func DefaultBlinkConfig() BlinkConfig {
	return DefaultAnimConfig()
}

// Animation управляет анимацией текста
type Animation struct {
	text              string
	blockFont         *BlockFont
	config            AnimConfig
	running           bool
	stopChan          chan struct{}
	rng               *rand.Rand
	frameCount        int
	currentRandomType AnimationType
}

// NewAnimation создаёт новую анимацию
func NewAnimation(text string) *Animation {
	return NewAnimationWithConfig(text, DefaultConfig(), DefaultAnimConfig())
}

// NewAnimationWithConfig создаёт анимацию с настройками
func NewAnimationWithConfig(text string, fontConfig Config, animConfig AnimConfig) *Animation {
	return &Animation{
		text:      strings.ToUpper(text),
		blockFont: NewWithConfig(fontConfig),
		config:    animConfig,
		stopChan:  make(chan struct{}),
		rng:       rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// SetType устанавливает тип анимации
func (a *Animation) SetType(t AnimationType) {
	a.config.Type = t
}

// SetChance устанавливает шанс эффекта (0.0 - 1.0)
func (a *Animation) SetChance(chance float64) {
	a.config.Chance = chance
}

// SetBlinkChance алиас для совместимости
func (a *Animation) SetBlinkChance(chance float64) {
	a.config.Chance = chance
}

// SetInterval устанавливает интервал между кадрами
func (a *Animation) SetInterval(interval time.Duration) {
	a.config.Interval = interval
}

// SetRange устанавливает диапазон затронутых букв
func (a *Animation) SetRange(min, max int) {
	a.config.Min = min
	a.config.Max = max
}

// SetBlinkRange алиас для совместимости
func (a *Animation) SetBlinkRange(min, max int) {
	a.SetRange(min, max)
}

// Frame генерирует один кадр анимации
func (a *Animation) Frame() string {
	a.frameCount++

	switch a.config.Type {
	case AnimPulse:
		return a.framePulse()
	case AnimWave:
		return a.frameWave()
	case AnimTyping:
		return a.frameTyping()
	case AnimGlitch:
		return a.frameGlitch()
	case AnimSequence:
		return a.frameSequence()
	case AnimRandom:
		return a.frameRandom()
	default:
		return a.frameBlink()
	}
}

// frameRandom случайно переключает типы анимаций
func (a *Animation) frameRandom() string {
	// Список доступных анимаций (кроме random)
	types := []AnimationType{AnimBlink, AnimPulse, AnimWave, AnimTyping, AnimGlitch, AnimSequence}

	// Переключаем тип каждые N кадров
	if a.frameCount == 1 || a.frameCount%a.config.RandomSwitchFrames == 0 {
		a.currentRandomType = types[a.rng.Intn(len(types))]
	}

	// Вызываем соответствующий метод
	switch a.currentRandomType {
	case AnimPulse:
		return a.framePulse()
	case AnimWave:
		return a.frameWave()
	case AnimTyping:
		return a.frameTyping()
	case AnimGlitch:
		return a.frameGlitch()
	case AnimSequence:
		return a.frameSequence()
	default:
		return a.frameBlink()
	}
}

// frameBlink случайное мигание букв
func (a *Animation) frameBlink() string {
	blinkCount := a.config.Min
	if a.config.Max > a.config.Min {
		blinkCount += a.rng.Intn(a.config.Max - a.config.Min + 1)
	}

	blinkPositions := make(map[int]bool)
	textLen := len([]rune(a.text))

	for i := 0; i < blinkCount && len(blinkPositions) < textLen; i++ {
		if a.rng.Float64() < a.config.Chance {
			pos := a.rng.Intn(textLen)
			blinkPositions[pos] = true
		}
	}

	return a.renderText(blinkPositions, nil)
}

// framePulse весь текст мигает
func (a *Animation) framePulse() string {
	// Мигаем каждый N-й кадр
	pulseRate := int(1.0 / a.config.Chance)
	if pulseRate < 1 {
		pulseRate = 1
	}

	if a.frameCount%pulseRate == 0 {
		// Выключен
		return a.renderEmpty()
	}
	return a.renderFull()
}

// frameWave волна по буквам
func (a *Animation) frameWave() string {
	textLen := len([]rune(a.text))
	wavePos := a.frameCount % (textLen + a.config.WaveWidth)

	hidePositions := make(map[int]bool)
	for i := 0; i < textLen; i++ {
		// Буква видна если она в "окне" волны
		if i < wavePos-a.config.WaveWidth || i >= wavePos {
			hidePositions[i] = true
		}
	}

	return a.renderText(hidePositions, nil)
}

// frameTyping эффект печатания
func (a *Animation) frameTyping() string {
	textLen := len([]rune(a.text))
	visibleCount := a.frameCount % (textLen + 5) // +5 для паузы в конце

	hidePositions := make(map[int]bool)
	for i := visibleCount; i < textLen; i++ {
		hidePositions[i] = true
	}

	return a.renderText(hidePositions, nil)
}

// frameGlitch глитч-эффект
func (a *Animation) frameGlitch() string {
	glitchCount := a.config.Min
	if a.config.Max > a.config.Min {
		glitchCount += a.rng.Intn(a.config.Max - a.config.Min + 1)
	}

	glitchPositions := make(map[int]rune)
	textLen := len([]rune(a.text))
	glitchRunes := []rune(a.config.GlitchChars)

	for i := 0; i < glitchCount && len(glitchPositions) < textLen; i++ {
		if a.rng.Float64() < a.config.Chance {
			pos := a.rng.Intn(textLen)
			glitchPositions[pos] = glitchRunes[a.rng.Intn(len(glitchRunes))]
		}
	}

	return a.renderText(nil, glitchPositions)
}

// frameSequence последовательное мигание
func (a *Animation) frameSequence() string {
	textLen := len([]rune(a.text))
	blinkPos := a.frameCount % textLen

	hidePositions := make(map[int]bool)
	hidePositions[blinkPos] = true

	return a.renderText(hidePositions, nil)
}

// renderText рендерит текст с эффектами
func (a *Animation) renderText(hidePositions map[int]bool, glitchPositions map[int]rune) string {
	runes := []rune(a.text)
	height := a.blockFont.height * a.blockFont.config.Scale
	lines := make([]string, height)

	for i, ch := range runes {
		var pattern []string

		if hidePositions != nil && hidePositions[i] && ch != ' ' {
			pattern = a.getBlankPattern()
		} else if glitchPositions != nil && glitchPositions[i] != 0 {
			pattern = a.getGlitchPattern(glitchPositions[i])
		} else {
			pattern = a.getCharPattern(ch)
		}

		for lineIdx := 0; lineIdx < height; lineIdx++ {
			patternLine := ""
			if lineIdx < len(pattern) {
				patternLine = pattern[lineIdx]
			}
			lines[lineIdx] += patternLine + strings.Repeat(a.blockFont.config.Space, a.blockFont.config.Scale)
		}
	}

	result := strings.Join(lines, "\n")

	if a.blockFont.config.Color != "" {
		result = a.blockFont.config.Color + result + ColorReset
	}

	return result
}

// renderEmpty рендерит пустой текст (для пульса)
func (a *Animation) renderEmpty() string {
	runes := []rune(a.text)
	height := a.blockFont.height * a.blockFont.config.Scale
	lines := make([]string, height)

	for range runes {
		pattern := a.getBlankPattern()
		for lineIdx := 0; lineIdx < height; lineIdx++ {
			lines[lineIdx] += pattern[lineIdx] + strings.Repeat(a.blockFont.config.Space, a.blockFont.config.Scale)
		}
	}

	return strings.Join(lines, "\n")
}

// renderFull рендерит полный текст
func (a *Animation) renderFull() string {
	result := a.renderText(nil, nil)
	return result
}

// getCharPattern возвращает паттерн символа
func (a *Animation) getCharPattern(ch rune) []string {
	pattern, ok := a.blockFont.chars[ch]
	if !ok {
		pattern = a.blockFont.chars[' ']
	}

	scaled := make([]string, a.blockFont.height*a.blockFont.config.Scale)
	for i := 0; i < a.blockFont.height; i++ {
		scaledLine := a.blockFont.scaleLine(pattern[i])
		for s := 0; s < a.blockFont.config.Scale; s++ {
			scaled[i*a.blockFont.config.Scale+s] = scaledLine
		}
	}
	return scaled
}

// getBlankPattern возвращает пустой паттерн
func (a *Animation) getBlankPattern() []string {
	width := 5 * a.blockFont.config.Scale
	blankLine := strings.Repeat(a.blockFont.config.Space, width)

	pattern := make([]string, a.blockFont.height*a.blockFont.config.Scale)
	for i := range pattern {
		pattern[i] = blankLine
	}
	return pattern
}

// getGlitchPattern возвращает паттерн глитча
func (a *Animation) getGlitchPattern(ch rune) []string {
	width := 5 * a.blockFont.config.Scale
	glitchChar := string(ch)
	glitchLine := strings.Repeat(glitchChar, width)

	pattern := make([]string, a.blockFont.height*a.blockFont.config.Scale)
	for i := range pattern {
		// Случайное заполнение
		line := ""
		for j := 0; j < width; j++ {
			if a.rng.Float64() < 0.7 {
				line += glitchChar
			} else {
				line += a.blockFont.config.Space
			}
		}
		pattern[i] = line
	}
	_ = glitchLine // подавить предупреждение
	return pattern
}

// Start запускает анимацию (блокирующий вызов, выход по Ctrl+C)
func (a *Animation) Start() {
	a.running = true
	a.stopChan = make(chan struct{})
	a.frameCount = 0

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Очистить экран и скрыть курсор
	fmt.Print("\033[2J\033[H\033[?25l")
	defer fmt.Print("\033[?25h\033[0m\n")

	firstFrame := a.Frame()
	lineCount := strings.Count(firstFrame, "\n") + 1
	fmt.Print(firstFrame)

	ticker := time.NewTicker(a.config.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-sigChan:
			a.running = false
			return
		case <-a.stopChan:
			a.running = false
			return
		case <-ticker.C:
			fmt.Printf("\033[%dA\033[G\033[J", lineCount)
			fmt.Print(a.Frame())
		}
	}
}

// StartAsync запускает анимацию в фоне
func (a *Animation) StartAsync() {
	go a.Start()
}

// Stop останавливает анимацию
func (a *Animation) Stop() {
	if a.running {
		close(a.stopChan)
	}
}

// IsRunning проверяет запущена ли анимация
func (a *Animation) IsRunning() bool {
	return a.running
}

// === Удобные функции ===

// Blink запускает мигающую анимацию
func Blink(text string) {
	NewAnimation(text).Start()
}

// BlinkWithConfig запускает мигающую анимацию с настройками
func BlinkWithConfig(text string, fontConfig Config, animConfig AnimConfig) {
	NewAnimationWithConfig(text, fontConfig, animConfig).Start()
}

// BlinkFor запускает анимацию на указанное время
func BlinkFor(text string, duration time.Duration) {
	anim := NewAnimation(text)
	go func() {
		time.Sleep(duration)
		anim.Stop()
	}()
	anim.Start()
}

// Pulse запускает пульсирующую анимацию
func Pulse(text string) {
	config := DefaultAnimConfig()
	config.Type = AnimPulse
	config.Interval = 500 * time.Millisecond
	config.Chance = 0.3
	NewAnimationWithConfig(text, DefaultConfig(), config).Start()
}

// Wave запускает волновую анимацию
func Wave(text string) {
	config := DefaultAnimConfig()
	config.Type = AnimWave
	config.Interval = 150 * time.Millisecond
	config.WaveWidth = 3
	NewAnimationWithConfig(text, DefaultConfig(), config).Start()
}

// Typing запускает анимацию печатания
func Typing(text string) {
	config := DefaultAnimConfig()
	config.Type = AnimTyping
	config.Interval = 200 * time.Millisecond
	NewAnimationWithConfig(text, DefaultConfig(), config).Start()
}

// Glitch запускает глитч-анимацию
func Glitch(text string) {
	config := DefaultAnimConfig()
	config.Type = AnimGlitch
	config.Interval = 80 * time.Millisecond
	config.Chance = 0.4
	config.Max = 4
	NewAnimationWithConfig(text, DefaultConfig(), config).Start()
}

// Sequence запускает последовательную анимацию
func Sequence(text string) {
	config := DefaultAnimConfig()
	config.Type = AnimSequence
	config.Interval = 200 * time.Millisecond
	NewAnimationWithConfig(text, DefaultConfig(), config).Start()
}

// Random запускает анимацию со случайной сменой типов
func Random(text string) {
	config := DefaultAnimConfig()
	config.Type = AnimRandom
	config.Interval = 100 * time.Millisecond
	config.RandomSwitchFrames = 30
	NewAnimationWithConfig(text, DefaultConfig(), config).Start()
}
