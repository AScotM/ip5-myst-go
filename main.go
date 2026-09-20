package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"os/signal"
	"sort"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
	"unicode/utf8"
	"unsafe"
)

type GlyphType string

const (
	GlyphFlow    GlyphType = "flow"
	GlyphPulse   GlyphType = "pulse"
	GlyphWhisper GlyphType = "whisper"
	GlyphShadow  GlyphType = "shadow"
	GlyphEcho    GlyphType = "echo"
)

const (
	attrNormal = 0
	attrBold   = 1 << iota
	attrDim
)

type Glyph struct {
	Type      GlyphType
	Intensity float64
	X         int
	Y         int
	Age       int
	Char      string
	Color     int
}

type NetworkEntity struct {
	Name          string
	Interface     string
	RXBytes       uint64
	TXBytes       uint64
	RXRate        float64
	TXRate        float64
	Essence       float64
	Aura          string
	Glyphs        []Glyph
	LastSeen      time.Time
	LastSample    time.Time
	Whispers      []string
	IPs           []string
	LastIPRefresh time.Time
}

func newNetworkEntity(interfaceName string, rx, tx uint64, sampleTime time.Time) *NetworkEntity {
	entity := &NetworkEntity{
		Interface:  interfaceName,
		RXBytes:    rx,
		TXBytes:    tx,
		Essence:    randomFloat(0.1, 1.0),
		LastSeen:   sampleTime,
		LastSample: sampleTime,
	}

	entity.Name = generateSpiritName(interfaceName)
	entity.Aura = generateAura()

	return entity
}

func generateSpiritName(interfaceName string) string {
	names := []struct {
		prefix string
		names  []string
	}{
		{"veth", []string{"Bridge Guardian", "Gatekeeper", "Threshold Walker"}},
		{"docker", []string{"Container Spirit", "Boxed Essence", "Isolated Soul"}},
		{"podman", []string{"Container Spirit", "Boxed Essence", "Isolated Soul"}},
		{"virbr", []string{"Bridge Guardian", "Gatekeeper", "Threshold Walker"}},
		{"wlan", []string{"Air Spirit", "Cloud Dancer", "Sky Whisper"}},
		{"wl", []string{"Air Spirit", "Cloud Dancer", "Sky Whisper"}},
		{"tun", []string{"Tunnel Dreamer", "Veil Piercer", "Hidden Path"}},
		{"tap", []string{"Mirror Pool", "Reflection Well", "Surface Tension"}},
		{"eth", []string{"Silent River", "Whispering Wind", "Ethereal Flow"}},
		{"en", []string{"Silent River", "Whispering Wind", "Ethereal Flow"}},
		{"lo", []string{"Inner Echo", "Soul Mirror", "Self Reflection"}},
	}

	for _, entry := range names {
		if strings.HasPrefix(interfaceName, entry.prefix) {
			return randomChoice(entry.names) + " (" + interfaceName + ")"
		}
	}

	prefixes := []string{
		"Whispering ",
		"Echoing ",
		"Veiled ",
		"Ancient ",
		"Forgotten ",
		"Secret ",
	}

	suffixes := []string{
		" the Observer",
		" the Watcher",
		" the Listener",
		" the Silent",
		" the Flowing",
		" the Hidden",
	}

	if rand.Float64() > 0.5 {
		return randomChoice(prefixes) + interfaceName
	}

	return interfaceName + randomChoice(suffixes)
}

func generateAura() string {
	auras := []string{
		"Pale Blue Mist",
		"Flickering Shadow",
		"Soft Golden Glow",
		"Deep Ocean Haze",
		"Forest Whisper Green",
		"Twilight Purple",
		"Dawn Orange",
		"Dusk Red",
		"Moon Silver",
		"Sun Gold",
	}

	return randomChoice(auras)
}

func (e *NetworkEntity) addWhisper(message string) {
	timestamp := time.Now().Format("15:04:05")
	e.Whispers = append(e.Whispers, fmt.Sprintf("[%s] %s", timestamp, message))

	if len(e.Whispers) > 5 {
		e.Whispers = e.Whispers[len(e.Whispers)-5:]
	}
}

func (e *NetworkEntity) updateCounters(rx, tx uint64, sampleTime time.Time) {
	elapsed := sampleTime.Sub(e.LastSample).Seconds()
	if elapsed <= 0 {
		return
	}

	if rx < e.RXBytes || tx < e.TXBytes {
		e.RXRate = 0
		e.TXRate = 0
	} else {
		e.RXRate = float64(rx-e.RXBytes) / elapsed
		e.TXRate = float64(tx-e.TXBytes) / elapsed
	}

	e.RXBytes = rx
	e.TXBytes = tx
	e.LastSample = sampleTime
	e.LastSeen = sampleTime
}

func (e *NetworkEntity) updateGlyphs(density float64) {
	if len(e.Glyphs) > 20 {
		e.Glyphs = e.Glyphs[len(e.Glyphs)-20:]
	}

	active := make([]Glyph, 0, len(e.Glyphs))

	for _, glyph := range e.Glyphs {
		glyph.Age++
		glyph.Intensity *= 0.95

		switch glyph.Type {
		case GlyphFlow:
			glyph.Char = randomChoice([]string{"↗", "↘", "↖", "↙", "↕", "↔"})
		case GlyphPulse:
			glyph.Char = randomChoice([]string{"●", "○", "◎", "◉", "⊙"})
		case GlyphWhisper:
			glyph.Char = randomChoice([]string{"…", "~", "⋮", "⋯"})
		case GlyphShadow:
			glyph.Char = randomChoice([]string{"░", "▒", "▓", "▚", "▞"})
		case GlyphEcho:
			glyph.Char = randomChoice([]string{"⦿", "⟳", "⟲", "↻", "↺"})
		}

		if glyph.Intensity > 0.1 {
			active = append(active, glyph)
		}
	}

	e.Glyphs = active

	trafficRate := math.Max(e.RXRate+e.TXRate, 0)

	flowIntensity := math.Min(
		math.Log1p(trafficRate)/math.Log1p(10_000_000),
		1.0,
	)

	density = math.Max(0, math.Min(density, 1))

	if rand.Float64() < flowIntensity*density {
		e.Glyphs = append(e.Glyphs, Glyph{
			Type:      GlyphFlow,
			Intensity: math.Max(flowIntensity, 0.1),
			X:         rand.Intn(21),
			Y:         rand.Intn(6),
			Color:     1 + rand.Intn(6),
		})
	}

	if rand.Float64() < e.Essence*density*0.7 {
		e.Glyphs = append(e.Glyphs, Glyph{
			Type:      GlyphWhisper,
			Intensity: e.Essence,
			X:         rand.Intn(21),
			Y:         rand.Intn(6),
			Color:     1 + rand.Intn(6),
		})
	}

	if rand.Float64() < density*0.08 {
		types := []GlyphType{
			GlyphPulse,
			GlyphShadow,
			GlyphEcho,
		}

		e.Glyphs = append(e.Glyphs, Glyph{
			Type:      types[rand.Intn(len(types))],
			Intensity: randomFloat(0.3, 1.0),
			X:         rand.Intn(21),
			Y:         rand.Intn(6),
			Color:     1 + rand.Intn(6),
		})
	}

	if len(e.Glyphs) > 20 {
		e.Glyphs = e.Glyphs[len(e.Glyphs)-20:]
	}
}

type MysticConfig struct {
	UpdateInterval float64 `json:"update_interval"`
	GlyphDensity   float64 `json:"glyph_density"`
	ShowWhispers   bool    `json:"show_whispers"`
	ShowAura       bool    `json:"show_aura"`
	AncientScript  bool    `json:"ancient_script"`
	LogFile        string  `json:"log_file"`
	RitualTimeout  float64 `json:"ritual_timeout"`
	MaxSpirits     int     `json:"max_spirits"`
	ShowLoopback   bool    `json:"show_loopback"`
	IPCacheTTL     float64 `json:"ip_cache_ttl"`
}

func defaultConfig() MysticConfig {
	return MysticConfig{
		UpdateInterval: 1.0,
		GlyphDensity:   0.3,
		ShowWhispers:   true,
		ShowAura:       true,
		AncientScript:  true,
		LogFile:        "network_mysteries.jsonl",
		RitualTimeout:  30.0,
		MaxSpirits:     20,
		ShowLoopback:   false,
		IPCacheTTL:     10.0,
	}
}

func (c *MysticConfig) loadFromFile(path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		return
	}

	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		return
	}

	if value, ok := numberValue(raw["update_interval"]); ok && value > 0 {
		c.UpdateInterval = value
	}

	if value, ok := numberValue(raw["glyph_density"]); ok && value >= 0 && value <= 1 {
		c.GlyphDensity = value
	}

	if value, ok := raw["show_whispers"].(bool); ok {
		c.ShowWhispers = value
	}

	if value, ok := raw["show_aura"].(bool); ok {
		c.ShowAura = value
	}

	if value, ok := raw["ancient_script"].(bool); ok {
		c.AncientScript = value
	}

	if value, ok := raw["log_file"].(string); ok {
		c.LogFile = value
	}

	if value, ok := numberValue(raw["ritual_timeout"]); ok && value > 0 {
		c.RitualTimeout = value
	}

	if value, ok := numberValue(raw["max_spirits"]); ok && value > 0 {
		c.MaxSpirits = int(value)
	}

	if value, ok := raw["show_loopback"].(bool); ok {
		c.ShowLoopback = value
	}

	if value, ok := numberValue(raw["ip_cache_ttl"]); ok && value > 0 {
		c.IPCacheTTL = value
	}
}

func numberValue(v any) (float64, bool) {
	switch n := v.(type) {
	case float64:
		return n, true
	case float32:
		return float64(n), true
	case int:
		return float64(n), true
	case int64:
		return float64(n), true
	case json.Number:
		value, err := n.Float64()
		return value, err == nil
	default:
		return 0, false
	}
}

type AncientScript struct{}

func (AncientScript) encodeNumber(num float64) string {
	symbols := []string{
		"Ⅰ", "Ⅱ", "Ⅲ", "Ⅳ", "Ⅴ",
		"Ⅵ", "Ⅶ", "Ⅷ", "Ⅸ", "Ⅹ",
		"Ⅺ", "Ⅻ", "ↀ", "ↁ", "ↂ",
		"Ↄ", "ↅ", "ↆ", "ↇ", "ↈ",
	}

	if num < 1 {
		return "↊"
	}

	value := int(num)
	result := make([]string, 0)

	for value > 0 {
		index := (value - 1) % len(symbols)
		result = append(result, symbols[index])
		value /= len(symbols)
	}

	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}

	return strings.Join(result, "")
}

func (a AncientScript) encodeRate(rate float64) string {
	rate = math.Max(rate, 0)

	units := []string{
		"Ⓑ/ⓢ",
		"ⓀⒷ/ⓢ",
		"ⓂⒷ/ⓢ",
		"ⒼⒷ/ⓢ",
		"ⓉⒷ/ⓢ",
	}

	divisors := []float64{
		1,
		1024,
		1024 * 1024,
		1024 * 1024 * 1024,
		1024 * 1024 * 1024 * 1024,
	}

	for i := len(units) - 1; i >= 0; i-- {
		if rate >= divisors[i] {
			return fmt.Sprintf(
				"%s %s",
				a.encodeNumber(rate/divisors[i]),
				units[i],
			)
		}
	}

	return "↊ " + units[0]
}

func (AncientScript) mysticTime() string {
	now := time.Now()

	hours := []string{
		"子", "丑", "寅", "卯",
		"辰", "巳", "午", "未",
		"申", "酉", "戌", "亥",
	}

	minutes := []string{
		"初", "壹", "贰", "叁",
		"肆", "伍", "陆", "柒",
		"捌", "玖", "拾",
	}

	hourIndex := now.Hour() % 12
	minuteIndex := now.Minute() / 6

	if minuteIndex >= len(minutes) {
		minuteIndex = len(minutes) - 1
	}

	return hours[hourIndex] + minutes[minuteIndex] + "刻"
}

type WhisperCollector struct {
	mu       sync.Mutex
	whispers []string
	max      int
}

func newWhisperCollector(maximum int) *WhisperCollector {
	return &WhisperCollector{
		max: maximum,
	}
}

func (w *WhisperCollector) add(source, message, level string) string {
	w.mu.Lock()
	defer w.mu.Unlock()

	timestamp := time.Now().Format("2006-01-02 15:04:05")

	whisper := fmt.Sprintf(
		"[%s] [%s] %s: %s",
		timestamp,
		level,
		source,
		message,
	)

	w.whispers = append(w.whispers, whisper)

	if len(w.whispers) > w.max {
		w.whispers = w.whispers[len(w.whispers)-w.max:]
	}

	return whisper
}

func (w *WhisperCollector) recent(count int) []string {
	w.mu.Lock()
	defer w.mu.Unlock()

	if count <= 0 {
		return nil
	}

	start := len(w.whispers) - count
	if start < 0 {
		start = 0
	}

	result := make([]string, len(w.whispers[start:]))
	copy(result, w.whispers[start:])

	return result
}

type Terminal struct {
	fd       int
	original syscall.Termios
	raw      bool
	mu       sync.Mutex
}

func newTerminal() *Terminal {
	return &Terminal{
		fd: int(os.Stdin.Fd()),
	}
}

func ioctl(fd int, request uintptr, arg unsafe.Pointer) error {
	_, _, errno := syscall.Syscall(
		syscall.SYS_IOCTL,
		uintptr(fd),
		request,
		uintptr(arg),
	)

	if errno != 0 {
		return errno
	}

	return nil
}

func (t *Terminal) enableRaw() error {
	t.mu.Lock()
	defer t.mu.Unlock()

	var current syscall.Termios

	if err := ioctl(
		t.fd,
		syscall.TCGETS,
		unsafe.Pointer(&current),
	); err != nil {
		return err
	}

	t.original = current

	raw := current

	raw.Iflag &^= syscall.BRKINT |
		syscall.ICRNL |
		syscall.INPCK |
		syscall.ISTRIP |
		syscall.IXON

	raw.Oflag &^= syscall.OPOST
	raw.Cflag |= syscall.CS8

	raw.Lflag &^= syscall.ECHO |
		syscall.ICANON |
		syscall.IEXTEN

	raw.Lflag |= syscall.ISIG

	raw.Cc[syscall.VMIN] = 0
	raw.Cc[syscall.VTIME] = 1

	if err := ioctl(
		t.fd,
		syscall.TCSETS,
		unsafe.Pointer(&raw),
	); err != nil {
		return err
	}

	t.raw = true

	fmt.Print("\033[?1049h")
	fmt.Print("\033[?25l")
	fmt.Print("\033[2J")
	fmt.Print("\033[H")

	return nil
}

func (t *Terminal) restore() {
	t.mu.Lock()
	defer t.mu.Unlock()

	if !t.raw {
		return
	}

	_ = ioctl(
		t.fd,
		syscall.TCSETS,
		unsafe.Pointer(&t.original),
	)

	fmt.Print("\033[0m")
	fmt.Print("\033[?25h")
	fmt.Print("\033[?1049l")

	t.raw = false
}

func terminalSize() (int, int) {
	type winsize struct {
		Row    uint16
		Col    uint16
		Xpixel uint16
		Ypixel uint16
	}

	var ws winsize

	_, _, errno := syscall.Syscall(
		syscall.SYS_IOCTL,
		os.Stdout.Fd(),
		syscall.TIOCGWINSZ,
		uintptr(unsafe.Pointer(&ws)),
	)

	if errno != 0 || ws.Row == 0 || ws.Col == 0 {
		return 24, 80
	}

	return int(ws.Row), int(ws.Col)
}

func ansiColor(color int) string {
	switch color {
	case 1:
		return "\033[34m"
	case 2:
		return "\033[36m"
	case 3:
		return "\033[32m"
	case 4:
		return "\033[35m"
	case 5:
		return "\033[33m"
	case 6:
		return "\033[31m"
	case 7:
		return "\033[37m"
	default:
		return "\033[0m"
	}
}

func ansiAttr(attr, color int) string {
	var builder strings.Builder

	builder.WriteString("\033[0m")

	if attr&attrBold != 0 {
		builder.WriteString("\033[1m")
	}

	if attr&attrDim != 0 {
		builder.WriteString("\033[2m")
	}

	if color > 0 {
		builder.WriteString(ansiColor(color))
	}

	return builder.String()
}

type Cell struct {
	Char  string
	Color int
	Attr  int
}

type Screen struct {
	Height int
	Width  int
	Cells  [][]Cell
}

func newScreen(height, width int) *Screen {
	screen := &Screen{}
	screen.resize(height, width)
	return screen
}

func (s *Screen) resize(height, width int) {
	if height < 1 {
		height = 1
	}

	if width < 1 {
		width = 1
	}

	s.Height = height
	s.Width = width
	s.Cells = make([][]Cell, height)

	for y := 0; y < height; y++ {
		s.Cells[y] = make([]Cell, width)

		for x := 0; x < width; x++ {
			s.Cells[y][x].Char = " "
		}
	}
}

func (s *Screen) clear() {
	for y := 0; y < s.Height; y++ {
		for x := 0; x < s.Width; x++ {
			s.Cells[y][x] = Cell{Char: " "}
		}
	}
}

func (s *Screen) put(y, x int, text string, color, attr int) {
	if y < 0 || y >= s.Height || x >= s.Width {
		return
	}

	position := x

	for _, r := range text {
		if position < 0 {
			position++
			continue
		}

		if position >= s.Width {
			break
		}

		s.Cells[y][position] = Cell{
			Char:  string(r),
			Color: color,
			Attr:  attr,
		}

		position++
	}
}

func (s *Screen) flush() {
	var output strings.Builder

	output.WriteString("\033[H")

	lastColor := -1
	lastAttr := -1

	for y := 0; y < s.Height; y++ {
		for x := 0; x < s.Width; x++ {
			cell := s.Cells[y][x]

			if cell.Color != lastColor || cell.Attr != lastAttr {
				output.WriteString(ansiAttr(cell.Attr, cell.Color))
				lastColor = cell.Color
				lastAttr = cell.Attr
			}

			if cell.Char == "" {
				output.WriteByte(' ')
			} else {
				output.WriteString(cell.Char)
			}
		}

		if y != s.Height-1 {
			output.WriteString("\r\n")
		}
	}

	output.WriteString("\033[0m")

	fmt.Print(output.String())
}

type NetworkMystic struct {
	Config      MysticConfig
	Entities    map[string]*NetworkEntity
	Whispers    *WhisperCollector
	Terminal    *Terminal
	Screen      *Screen
	Ancient     AncientScript
	running     bool
	runningMu   sync.RWMutex
	entityMu    sync.RWMutex
	cleanupMu   sync.Mutex
	cleanedUp   bool
	lastScan    time.Time
	lastRitual  time.Time
	input       chan byte
	inputErrors chan error
}

func newNetworkMystic(config MysticConfig) *NetworkMystic {
	height, width := terminalSize()

	return &NetworkMystic{
		Config:      config,
		Entities:    make(map[string]*NetworkEntity),
		Whispers:    newWhisperCollector(100),
		Terminal:    newTerminal(),
		Screen:      newScreen(height, width),
		input:       make(chan byte, 32),
		inputErrors: make(chan error, 1),
	}
}

func (m *NetworkMystic) isRunning() bool {
	m.runningMu.RLock()
	defer m.runningMu.RUnlock()

	return m.running
}

func (m *NetworkMystic) setRunning(value bool) {
	m.runningMu.Lock()
	m.running = value
	m.runningMu.Unlock()
}

func (m *NetworkMystic) stop() {
	m.setRunning(false)
}

func (m *NetworkMystic) ritualMoonCycle() string {
	const lunarCycle = 2_419_200.0

	cycle := math.Mod(
		float64(time.Now().Unix()),
		lunarCycle,
	) / lunarCycle

	switch {
	case cycle < 0.25:
		return "New Moon: Spirits are quiet"
	case cycle < 0.5:
		return "Waxing Moon: Connections strengthen"
	case cycle < 0.75:
		return "Full Moon: Maximum flow"
	default:
		return "Waning Moon: Paths fade"
	}
}

func (m *NetworkMystic) ritualTideChange() string {
	tide := math.Sin(
		float64(time.Now().Unix()) / 3600,
	)

	if tide > 0.7 {
		return "High Tide: Data flows freely"
	}

	if tide < -0.7 {
		return "Low Tide: Channels narrow"
	}

	return "Changing Tide: Flux in patterns"
}

func (m *NetworkMystic) ritualStarAlignment() string {
	return randomChoice([]string{
		"Stars favor communication",
		"Constellations whisper of packets",
		"Celestial paths clear",
		"Starlight reveals hidden flows",
		"Galactic currents shift",
	})
}

func (m *NetworkMystic) ritualWindShift() string {
	return randomChoice([]string{
		"Northern Wind: Cold data streams",
		"Southern Wind: Warm connections",
		"Eastern Wind: Dawn of new packets",
		"Western Wind: Dusk of old routes",
		"Still Air: Silent monitoring",
	})
}

func (m *NetworkMystic) performRituals(force bool) {
	now := time.Now()

	if !force &&
		!m.lastRitual.IsZero() &&
		now.Sub(m.lastRitual).Seconds() < m.Config.RitualTimeout {
		return
	}

	rituals := []func() string{
		m.ritualMoonCycle,
		m.ritualTideChange,
		m.ritualStarAlignment,
		m.ritualWindShift,
	}

	ritual := rituals[rand.Intn(len(rituals))]

	m.Whispers.add(
		"Ritual",
		ritual(),
		"MYSTIC",
	)

	m.lastRitual = now
}

func (m *NetworkMystic) scanSpirits() {
	file, err := os.Open("/proc/net/dev")
	if err != nil {
		m.Whispers.add(
			"Scanner",
			"Failed to scan: "+err.Error(),
			"ERROR",
		)
		return
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNumber := 0
	currentTime := time.Now()
	found := make(map[string]bool)

	m.entityMu.Lock()
	defer m.entityMu.Unlock()

	for scanner.Scan() {
		lineNumber++

		if lineNumber <= 2 {
			continue
		}

		line := scanner.Text()

		index := strings.IndexByte(line, ':')
		if index < 0 {
			continue
		}

		interfaceName := strings.TrimSpace(line[:index])
		if interfaceName == "" {
			continue
		}

		if interfaceName == "lo" && !m.Config.ShowLoopback {
			continue
		}

		fields := strings.Fields(line[index+1:])
		if len(fields) < 16 {
			continue
		}

		rx, err := strconv.ParseUint(fields[0], 10, 64)
		if err != nil {
			continue
		}

		tx, err := strconv.ParseUint(fields[8], 10, 64)
		if err != nil {
			continue
		}

		found[interfaceName] = true

		entity, exists := m.Entities[interfaceName]

		if !exists {
			entity = newNetworkEntity(
				interfaceName,
				rx,
				tx,
				currentTime,
			)

			m.Entities[interfaceName] = entity

			m.Whispers.add(
				"Veil",
				fmt.Sprintf(
					"Spirit '%s' has appeared",
					entity.Name,
				),
				"SPIRIT",
			)
		} else {
			entity.updateCounters(
				rx,
				tx,
				currentTime,
			)
		}

		m.refreshEntityIPs(
			entity,
			currentTime,
		)

		m.updateEntityWhispers(entity)
		entity.updateGlyphs(m.Config.GlyphDensity)
	}

	if err := scanner.Err(); err != nil {
		m.Whispers.add(
			"Scanner",
			"Failed to scan: "+err.Error(),
			"ERROR",
		)
	}

	for interfaceName, entity := range m.Entities {
		if found[interfaceName] {
			continue
		}

		delete(m.Entities, interfaceName)

		m.Whispers.add(
			"Veil",
			fmt.Sprintf(
				"Spirit '%s' has faded",
				entity.Name,
			),
			"SPIRIT",
		)
	}

	m.lastScan = currentTime
}

func (m *NetworkMystic) updateEntityWhispers(entity *NetworkEntity) {
	switch {
	case entity.RXRate > 1_000_000:
		entity.addWhisper("Great flow from beyond")

	case entity.RXRate < 1_000 && entity.TXRate < 1_000:
		entity.addWhisper("Resting in silence")

	case entity.RXRate > entity.TXRate*2:
		entity.addWhisper("Listening more than speaking")

	case entity.TXRate > entity.RXRate*2:
		entity.addWhisper("Whispering to the void")
	}
}

func (m *NetworkMystic) refreshEntityIPs(entity *NetworkEntity, now time.Time) {
	if !entity.LastIPRefresh.IsZero() &&
		now.Sub(entity.LastIPRefresh).Seconds() < m.Config.IPCacheTTL {
		return
	}

	entity.IPs = getInterfaceIPs(entity.Interface)
	entity.LastIPRefresh = now
}

func getInterfaceIPs(interfaceName string) []string {
	command := exec.Command(
		"ip",
		"-o",
		"-4",
		"addr",
		"show",
		"dev",
		interfaceName,
	)

	type result struct {
		output []byte
		err    error
	}

	done := make(chan result, 1)

	go func() {
		output, err := command.Output()

		done <- result{
			output: output,
			err:    err,
		}
	}()

	var output []byte

	select {
	case result := <-done:
		if result.err != nil {
			return nil
		}

		output = result.output

	case <-time.After(2 * time.Second):
		if command.Process != nil {
			_ = command.Process.Kill()
		}

		return nil
	}

	var ips []string

	for _, line := range strings.Split(string(output), "\n") {
		fields := strings.Fields(line)

		if len(fields) < 4 {
			continue
		}

		ip := strings.SplitN(fields[3], "/", 2)[0]

		if ip != "" && ip != "127.0.0.1" {
			ips = append(ips, ip)
		}
	}

	return ips
}

func (m *NetworkMystic) drawBorder() {
	height := m.Screen.Height
	width := m.Screen.Width

	if width < 2 || height < 2 {
		return
	}

	m.Screen.put(0, 0, "╔", 0, attrNormal)

	for x := 1; x < width-1; x++ {
		m.Screen.put(0, x, "═", 0, attrNormal)
	}

	m.Screen.put(0, width-1, "╗", 0, attrNormal)

	for y := 1; y < height-1; y++ {
		m.Screen.put(y, 0, "║", 0, attrNormal)
		m.Screen.put(y, width-1, "║", 0, attrNormal)
	}

	m.Screen.put(height-1, 0, "╚", 0, attrNormal)

	for x := 1; x < width-1; x++ {
		m.Screen.put(height-1, x, "═", 0, attrNormal)
	}

	m.Screen.put(height-1, width-1, "╝", 0, attrNormal)
}

func (m *NetworkMystic) renderGlyph(
	glyph Glyph,
	offsetX,
	offsetY int,
) {
	attr := attrNormal

	if glyph.Intensity > 0.6 {
		attr |= attrBold
	} else if glyph.Intensity < 0.3 {
		attr |= attrDim
	}

	m.Screen.put(
		offsetY+glyph.Y,
		offsetX+glyph.X,
		glyph.Char,
		glyph.Color,
		attr,
	)
}

func rateToHeight(rate float64, height int) int {
	if rate <= 0 {
		return 0
	}

	normalized := math.Log1p(rate) /
		math.Log1p(10_000_000)

	result := int(
		math.Ceil(
			normalized * float64(height),
		),
	)

	if result < 1 {
		result = 1
	}

	if result > height {
		result = height
	}

	return result
}

func (m *NetworkMystic) renderFlowVisualization(
	rx,
	tx float64,
	x,
	y,
	width int,
) {
	if width < 3 {
		return
	}

	rxHeight := rateToHeight(rx, 5)
	txHeight := rateToHeight(tx, 5)

	bars := []rune("▁▂▃▄▅▆▇█")

	for index := 0; index < rxHeight; index++ {
		charIndex := index * len(bars) / 5

		if charIndex >= len(bars) {
			charIndex = len(bars) - 1
		}

		m.Screen.put(
			y+4-index,
			x,
			string(bars[charIndex]),
			2,
			attrNormal,
		)
	}

	for index := 0; index < txHeight; index++ {
		charIndex := index * len(bars) / 5

		if charIndex >= len(bars) {
			charIndex = len(bars) - 1
		}

		m.Screen.put(
			y+4-index,
			x+2,
			string(bars[charIndex]),
			3,
			attrNormal,
		)
	}
}

func formatRate(rate float64) string {
	rate = math.Max(rate, 0)

	units := []string{
		"B/s",
		"KB/s",
		"MB/s",
		"GB/s",
		"TB/s",
	}

	for _, unit := range units {
		if rate < 1024 || unit == units[len(units)-1] {
			return fmt.Sprintf(
				"%.1f %s",
				rate,
				unit,
			)
		}

		rate /= 1024
	}

	return "0.0 B/s"
}

func (m *NetworkMystic) drawSpiritInfo(
	entity *NetworkEntity,
	yStart int,
) {
	height := m.Screen.Height
	width := m.Screen.Width

	if yStart >= height-2 || width < 20 {
		return
	}

	ipDisplay := "[No physical form]"

	if len(entity.IPs) > 0 {
		ipDisplay = strings.Join(entity.IPs, ", ")
	}

	m.Screen.put(
		yStart,
		2,
		"╭─ "+entity.Name,
		0,
		attrBold,
	)

	line := 1

	if m.Config.ShowAura {
		m.Screen.put(
			yStart+line,
			2,
			"│  Aura: "+entity.Aura,
			4,
			attrNormal,
		)

		line++
	}

	var rxDisplay string
	var txDisplay string

	if m.Config.AncientScript {
		rxDisplay = m.Ancient.encodeRate(entity.RXRate)
		txDisplay = m.Ancient.encodeRate(entity.TXRate)
	} else {
		rxDisplay = formatRate(entity.RXRate)
		txDisplay = formatRate(entity.TXRate)
	}

	m.Screen.put(
		yStart+line,
		2,
		"│  From Beyond: "+rxDisplay,
		2,
		attrNormal,
	)

	line++

	m.Screen.put(
		yStart+line,
		2,
		"│  To Void:     "+txDisplay,
		3,
		attrNormal,
	)

	line++

	m.Screen.put(
		yStart+line,
		2,
		fmt.Sprintf(
			"│  Essence: %.2f",
			entity.Essence,
		),
		0,
		attrDim,
	)

	line++

	m.Screen.put(
		yStart+line,
		2,
		"│  Form: "+ipDisplay,
		0,
		attrDim,
	)

	line++

	if width >= 48 {
		m.renderFlowVisualization(
			entity.RXRate,
			entity.TXRate,
			40,
			yStart+1,
			width-40,
		)
	}

	start := len(entity.Whispers) - 2
	if start < 0 {
		start = 0
	}

	for _, whisper := range entity.Whispers[start:] {
		if yStart+line >= height-2 {
			break
		}

		m.Screen.put(
			yStart+line,
			2,
			"│  "+whisper,
			0,
			attrDim,
		)

		line++
	}

	if yStart+line < height-1 {
		m.Screen.put(
			yStart+line,
			2,
			"╰"+strings.Repeat(
				"─",
				maxInt(width-5, 0),
			),
			0,
			attrNormal,
		)
	}

	for _, glyph := range entity.Glyphs {
		m.renderGlyph(
			glyph,
			25,
			yStart-4,
		)
	}
}

func (m *NetworkMystic) drawWhispersPanel() {
	if !m.Config.ShowWhispers {
		return
	}

	height := m.Screen.Height
	width := m.Screen.Width

	if height < 8 || width < 60 {
		return
	}

	panelWidth := minInt(
		50,
		maxInt(width/3, 24),
	)

	panelX := width - panelWidth - 2

	if panelX < 2 {
		return
	}

	m.Screen.put(
		1,
		panelX,
		"╭─[Whispers]─",
		0,
		attrDim,
	)

	count := minInt(8, height-5)
	whispers := m.Whispers.recent(count)

	for index, whisper := range whispers {
		line := truncateRunes(
			whisper,
			maxInt(panelWidth-2, 0),
		)

		m.Screen.put(
			2+index,
			panelX,
			"│"+line,
			0,
			attrDim,
		)
	}

	sideX := panelX + panelWidth

	for y := 1; y < minInt(len(whispers)+2, height-2); y++ {
		m.Screen.put(
			y,
			sideX,
			"│",
			0,
			attrDim,
		)
	}

	endY := minInt(
		len(whispers)+2,
		height-2,
	)

	m.Screen.put(
		endY,
		panelX,
		"╰"+strings.Repeat("─", panelWidth),
		0,
		attrDim,
	)
}

func (m *NetworkMystic) drawStatusBar() {
	height := m.Screen.Height
	width := m.Screen.Width

	if height < 4 || width < 8 {
		return
	}

	var totalRX float64
	var totalTX float64

	for _, entity := range m.Entities {
		totalRX += entity.RXRate
		totalTX += entity.TXRate
	}

	var rxDisplay string
	var txDisplay string

	if m.Config.AncientScript {
		rxDisplay = m.Ancient.encodeRate(totalRX)
		txDisplay = m.Ancient.encodeRate(totalTX)
	} else {
		rxDisplay = formatRate(totalRX)
		txDisplay = formatRate(totalTX)
	}

	status := fmt.Sprintf(
		"Total Flow: %s ╫ %s │ Spirits: %d │ %s",
		rxDisplay,
		txDisplay,
		len(m.Entities),
		m.Ancient.mysticTime(),
	)

	status = truncateRunes(
		status,
		maxInt(width-4, 0),
	)

	statusX := maxInt(
		1,
		(width-runeLength(status))/2,
	)

	m.Screen.put(
		height-3,
		statusX,
		status,
		0,
		attrDim,
	)
}

func (m *NetworkMystic) drawVeil() {
	height, width := terminalSize()

	if height != m.Screen.Height ||
		width != m.Screen.Width {
		m.Screen.resize(height, width)
	}

	m.Screen.clear()

	if height < 8 || width < 30 {
		message := "Terminal too small"

		m.Screen.put(
			maxInt(height/2, 0),
			maxInt(
				(width-runeLength(message))/2,
				0,
			),
			message,
			0,
			attrBold,
		)

		m.Screen.flush()
		return
	}

	m.drawBorder()

	title := "╡ Network Mystic ╞"

	titleX := maxInt(
		1,
		(width-runeLength(title))/2,
	)

	m.Screen.put(
		0,
		titleX,
		title,
		1,
		attrBold,
	)

	mysticTime := m.Ancient.mysticTime()
	timeX := width - runeLength(mysticTime) - 2

	if timeX > titleX+runeLength(title) {
		m.Screen.put(
			0,
			timeX,
			mysticTime,
			0,
			attrDim,
		)
	}

	m.entityMu.RLock()
	defer m.entityMu.RUnlock()

	spirits := make([]*NetworkEntity, 0, len(m.Entities))

	for _, entity := range m.Entities {
		spirits = append(spirits, entity)
	}

	sort.Slice(
		spirits,
		func(i, j int) bool {
			left := spirits[i].RXRate + spirits[i].TXRate
			right := spirits[j].RXRate + spirits[j].TXRate

			return left > right
		},
	)

	y := 2
	reservedBottom := 4

	limit := len(spirits)

	if limit > m.Config.MaxSpirits {
		limit = m.Config.MaxSpirits
	}

	for _, spirit := range spirits[:limit] {
		blockHeight := 8

		if m.Config.ShowAura {
			blockHeight = 9
		}

		if y+blockHeight >= height-reservedBottom {
			break
		}

		m.drawSpiritInfo(spirit, y)
		y += blockHeight
	}

	m.drawWhispersPanel()
	m.drawStatusBar()
	m.Screen.flush()
}

func (m *NetworkMystic) forceRefresh() {
	m.scanSpirits()
	m.performRituals(false)
	m.drawVeil()

	m.Whispers.add(
		"Mystic",
		"Ritual refreshed",
		"INFO",
	)
}

func (m *NetworkMystic) inputLoop() {
	buffer := make([]byte, 1)

	for m.isRunning() {
		n, err := os.Stdin.Read(buffer)

		if err != nil {
			if errors.Is(err, syscall.EINTR) {
				continue
			}

			if !m.isRunning() {
				return
			}

			select {
			case m.inputErrors <- err:
			default:
			}

			return
		}

		if n == 0 {
			continue
		}

		select {
		case m.input <- buffer[0]:
		default:
		}
	}
}

func (m *NetworkMystic) handleKey(key byte) {
	switch key {
	case 'q', 'Q':
		m.stop()

	case ' ':
		m.performRituals(true)
		m.drawVeil()

	case 'r', 'R':
		m.forceRefresh()

	case 'l', 'L':
		m.Config.ShowLoopback = !m.Config.ShowLoopback

		status := "hidden"

		if m.Config.ShowLoopback {
			status = "shown"
		}

		m.Whispers.add(
			"Mystic",
			"Loopback spirits "+status,
			"INFO",
		)

		m.scanSpirits()
		m.drawVeil()

	case 'a', 'A':
		m.Config.AncientScript = !m.Config.AncientScript

		status := "disabled"

		if m.Config.AncientScript {
			status = "enabled"
		}

		m.Whispers.add(
			"Mystic",
			"Ancient script "+status,
			"INFO",
		)

		m.drawVeil()

	case 'w', 'W':
		m.Config.ShowWhispers = !m.Config.ShowWhispers

		status := "disabled"

		if m.Config.ShowWhispers {
			status = "enabled"
		}

		m.Whispers.add(
			"Mystic",
			"Whispers "+status,
			"INFO",
		)

		m.drawVeil()
	}
}

func (m *NetworkMystic) runRitual() error {
	if err := m.Terminal.enableRaw(); err != nil {
		return fmt.Errorf(
			"failed to enter raw terminal mode: %w",
			err,
		)
	}

	defer m.Terminal.restore()

	m.setRunning(true)

	m.Whispers.add(
		"Mystic",
		"Beginning the ritual...",
		"INFO",
	)

	m.scanSpirits()
	m.drawVeil()

	go m.inputLoop()

	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()

	for m.isRunning() {
		select {
		case key := <-m.input:
			m.handleKey(key)

		case err := <-m.inputErrors:
			if err != nil && m.isRunning() {
				m.Whispers.add(
					"Mystic",
					"Input error: "+err.Error(),
					"ERROR",
				)
			}

		case <-ticker.C:
			now := time.Now()

			if m.lastScan.IsZero() ||
				now.Sub(m.lastScan).Seconds() >= m.Config.UpdateInterval {
				m.scanSpirits()
				m.performRituals(false)
				m.drawVeil()
			}
		}
	}

	return nil
}

type MysteryEntity struct {
	Spirit    string   `json:"spirit"`
	Interface string   `json:"interface"`
	Essence   float64  `json:"essence"`
	Aura      string   `json:"aura"`
	TotalRX   uint64   `json:"total_rx"`
	TotalTX   uint64   `json:"total_tx"`
	RXRate    float64  `json:"rx_rate"`
	TXRate    float64  `json:"tx_rate"`
	IPs       []string `json:"ips"`
	Whispers  []string `json:"whispers"`
	Timestamp string   `json:"timestamp"`
}

type MysteryLog struct {
	Timestamp      string          `json:"timestamp"`
	MysticTime     string          `json:"mystic_time"`
	TotalSpirits   int             `json:"total_spirits"`
	Mysteries      []MysteryEntity `json:"mysteries"`
	SystemWhispers []string        `json:"system_whispers"`
}

func (m *NetworkMystic) saveMysteries() {
	if m.Config.LogFile == "" {
		return
	}

	m.entityMu.RLock()

	mysteries := make([]MysteryEntity, 0, len(m.Entities))

	for _, entity := range m.Entities {
		start := len(entity.Whispers) - 5

		if start < 0 {
			start = 0
		}

		whispers := append(
			[]string(nil),
			entity.Whispers[start:]...,
		)

		ips := append(
			[]string(nil),
			entity.IPs...,
		)

		mysteries = append(
			mysteries,
			MysteryEntity{
				Spirit:    entity.Name,
				Interface: entity.Interface,
				Essence:   entity.Essence,
				Aura:      entity.Aura,
				TotalRX:   entity.RXBytes,
				TotalTX:   entity.TXBytes,
				RXRate:    entity.RXRate,
				TXRate:    entity.TXRate,
				IPs:       ips,
				Whispers:  whispers,
				Timestamp: time.Now().Format(time.RFC3339Nano),
			},
		)
	}

	totalSpirits := len(m.Entities)

	m.entityMu.RUnlock()

	entry := MysteryLog{
		Timestamp:      time.Now().Format(time.RFC3339Nano),
		MysticTime:     m.Ancient.mysticTime(),
		TotalSpirits:   totalSpirits,
		Mysteries:      mysteries,
		SystemWhispers: m.Whispers.recent(10),
	}

	data, err := json.Marshal(entry)

	if err != nil {
		m.Whispers.add(
			"Archivist",
			"Failed to encode: "+err.Error(),
			"ERROR",
		)

		return
	}

	file, err := os.OpenFile(
		m.Config.LogFile,
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0o644,
	)

	if err != nil {
		m.Whispers.add(
			"Archivist",
			"Failed to save: "+err.Error(),
			"ERROR",
		)

		return
	}

	defer file.Close()

	if _, err := file.Write(append(data, '\n')); err != nil {
		m.Whispers.add(
			"Archivist",
			"Failed to save: "+err.Error(),
			"ERROR",
		)
	}
}

func (m *NetworkMystic) cleanup() {
	m.cleanupMu.Lock()

	if m.cleanedUp {
		m.cleanupMu.Unlock()
		return
	}

	m.cleanedUp = true
	m.cleanupMu.Unlock()

	m.Terminal.restore()

	m.Whispers.add(
		"Mystic",
		"Ritual complete. Veil closing...",
		"INFO",
	)

	m.saveMysteries()

	if m.Config.LogFile != "" {
		fmt.Fprintf(
			os.Stderr,
			"\nMysteries saved to %s\n",
			m.Config.LogFile,
		)
	}

	fmt.Fprintln(
		os.Stderr,
		"\nMay the flows guide you...",
	)

	fmt.Fprintln(
		os.Stderr,
		"\nControls:",
	)

	fmt.Fprintln(
		os.Stderr,
		"  q - Quit",
	)

	fmt.Fprintln(
		os.Stderr,
		"  Ctrl-C - Quit",
	)

	fmt.Fprintln(
		os.Stderr,
		"  Space - Perform ritual",
	)

	fmt.Fprintln(
		os.Stderr,
		"  r - Refresh",
	)

	fmt.Fprintln(
		os.Stderr,
		"  l - Toggle loopback spirits",
	)

	fmt.Fprintln(
		os.Stderr,
		"  a - Toggle ancient script",
	)

	fmt.Fprintln(
		os.Stderr,
		"  w - Toggle whispers",
	)
}

func randomChoice[T any](values []T) T {
	return values[rand.Intn(len(values))]
}

func randomFloat(minimum, maximum float64) float64 {
	return minimum +
		rand.Float64()*(maximum-minimum)
}

func minInt(a, b int) int {
	if a < b {
		return a
	}

	return b
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}

	return b
}

func runeLength(value string) int {
	return utf8.RuneCountInString(value)
}

func truncateRunes(value string, maximum int) string {
	if maximum <= 0 {
		return ""
	}

	runes := []rune(value)

	if len(runes) <= maximum {
		return value
	}

	return string(runes[:maximum])
}

func main() {
	config := defaultConfig()

	if len(os.Args) > 1 {
		config.loadFromFile(os.Args[1])
	}

	mystic := newNetworkMystic(config)

	signals := make(chan os.Signal, 2)

	signal.Notify(
		signals,
		os.Interrupt,
		syscall.SIGTERM,
	)

	defer signal.Stop(signals)

	go func() {
		<-signals
		mystic.stop()
	}()

	err := mystic.runRitual()

	mystic.cleanup()

	if err != nil {
		fmt.Fprintln(
			os.Stderr,
			"Network Mystic:",
			err,
		)

		os.Exit(1)
	}
}
