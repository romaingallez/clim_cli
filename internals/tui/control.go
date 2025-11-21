package tui

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/romaingallez/clim_cli/internals/api"
	"github.com/romaingallez/clim_cli/internals/storage"
)

var (
	ctrlTitleStyle  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("39")).MarginBottom(1)
	ctrlLabelStyle  = lipgloss.NewStyle().Bold(true)
	ctrlValueStyle  = lipgloss.NewStyle()
	ctrlDimStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("244")).Italic(true)
	ctrlErrStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
	ctrlSelectStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("46")).Background(lipgloss.Color("235"))
)

type controlFocus int

const (
	focusPower controlFocus = iota
	focusMode
	focusTemp
	focusFanRate
	focusFanDir
	focusDevice
)

type fetchMsg struct{}

type controlModel struct {
	devices      []*storage.DeviceHistory
	cursorDevice int
	focus        controlFocus
	current      map[string]string
	pending      api.Clim
	showHelp     bool
	showConfirm  bool
	showResults  bool
	applyResults []string
	err          error
	quitting     bool
}

func newControlModel(devs []*storage.DeviceHistory) controlModel {
	// ensure stable order by name
	sorted := make([]*storage.DeviceHistory, len(devs))
	copy(sorted, devs)
	sort.Slice(sorted, func(i, j int) bool {
		return strings.ToLower(sorted[i].Device.Name) < strings.ToLower(sorted[j].Device.Name)
	})
	m := controlModel{devices: sorted, current: map[string]string{}}
	if len(sorted) > 0 {
		ip := sorted[0].Device.IP
		m.pending.IP = ip
		m.pending.Power = "1"
		m.pending.Mode = "0"
		m.pending.Temp = "24"
		m.pending.FanRate = "A"
		m.pending.FanDir = ""
	}
	return m
}

func (m controlModel) Init() tea.Cmd {
	return tea.Tick(2*time.Second, func(time.Time) tea.Msg { return fetchMsg{} })
}

func (m controlModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		case "esc":
			if m.showHelp {
				m.showHelp = false
				return m, nil
			}
			if m.showConfirm {
				m.showConfirm = false
				return m, nil
			}
			if m.showResults {
				m.showResults = false
				return m, nil
			}
		case "h":
			m.showHelp = !m.showHelp
		case "tab":
			m.focus = (m.focus + 1) % 6
		case "shift+tab":
			m.focus = (m.focus + 5) % 6
		case "r":
			return m, func() tea.Msg { return fetchMsg{} }
		case "+":
			m.incrementTemp(1)
		case "-":
			m.incrementTemp(-1)
		case "up", "k":
			if m.focus == focusDevice {
				if m.cursorDevice > 0 {
					m.cursorDevice--
					m.pending.IP = m.devices[m.cursorDevice].Device.IP
					return m, func() tea.Msg { return fetchMsg{} }
				}
			} else if m.focus == focusTemp {
				m.incrementTemp(1)
			}
		case "down", "j":
			if m.focus == focusDevice {
				if m.cursorDevice < len(m.devices)-1 {
					m.cursorDevice++
					m.pending.IP = m.devices[m.cursorDevice].Device.IP
					return m, func() tea.Msg { return fetchMsg{} }
				}
			} else if m.focus == focusTemp {
				m.incrementTemp(-1)
			}
		case "left", "right":
			m.cycleField(msg.String() == "right")
		case "p":
			m.togglePower()
		case "m":
			m.cycleMode(true)
		case "f":
			m.cycleFanRate(true)
		case "d":
			m.cycleFanDir(true)
		case "a":
			if len(m.devices) == 0 {
				break
			}
			m.showConfirm = true
		case "y":
			if m.showConfirm {
				m.showConfirm = false
				return m, m.applyAll()
			}
		case "n":
			if m.showConfirm {
				m.showConfirm = false
			}
		}
	case fetchMsg:
		if len(m.devices) == 0 {
			return m, tea.Tick(2*time.Second, func(time.Time) tea.Msg { return fetchMsg{} })
		}
		ip := m.devices[m.cursorDevice].Device.IP
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		info, err := api.FetchControlInfo(ctx, ip)
		if err != nil {
			m.err = err
		} else {
			m.err = nil
			m.current = info
			if v, ok := info["pow"]; ok {
				m.pending.Power = v
			}
			if v, ok := info["mode"]; ok {
				m.pending.Mode = v
			}
			if v, ok := info["stemp"]; ok {
				m.pending.Temp = v
			}
			if v, ok := info["f_rate"]; ok {
				m.pending.FanRate = v
			}
			if v, ok := info["f_dir"]; ok {
				m.pending.FanDir = v
			}
		}
		return m, tea.Tick(2*time.Second, func(time.Time) tea.Msg { return fetchMsg{} })
	}
	return m, nil
}

func (m *controlModel) incrementTemp(delta int) {
	val, err := strconv.Atoi(m.pending.Temp)
	if err != nil {
		val = 24
	}
	val += delta
	if val < 16 {
		val = 16
	}
	if val > 30 {
		val = 30
	}
	m.pending.Temp = strconv.Itoa(val)
}

func (m *controlModel) togglePower() {
	if m.pending.Power == "1" {
		m.pending.Power = "0"
	} else {
		m.pending.Power = "1"
	}
}

func (m *controlModel) cycleMode(forward bool) {
	seq := []string{"0", "1", "2", "3", "4"}
	m.pending.Mode = cycle(seq, m.pending.Mode, forward)
}

func (m *controlModel) cycleFanRate(forward bool) {
	seq := []string{"A", "3", "4", "5", "6", "7"}
	m.pending.FanRate = cycle(seq, m.pending.FanRate, forward)
}

func (m *controlModel) cycleFanDir(forward bool) {
	seq := []string{"0", "1", "2", "3"}
	m.pending.FanDir = cycle(seq, m.pending.FanDir, forward)
}

func (m *controlModel) cycleField(forward bool) {
	switch m.focus {
	case focusMode:
		m.cycleMode(forward)
	case focusFanRate:
		m.cycleFanRate(forward)
	case focusFanDir:
		m.cycleFanDir(forward)
	case focusPower:
		m.togglePower()
	}
}

func cycle(seq []string, cur string, forward bool) string {
	idx := 0
	for i, v := range seq {
		if v == cur {
			idx = i
			break
		}
	}
	if forward {
		idx = (idx + 1) % len(seq)
	} else {
		idx = (idx - 1 + len(seq)) % len(seq)
	}
	return seq[idx]
}

func (m controlModel) View() string {
	if len(m.devices) == 0 {
		return ctrlDimStyle.Render("No devices selected. Press q to quit.")
	}
	var b strings.Builder
	b.WriteString(ctrlTitleStyle.Render("Climate Control"))
	b.WriteString("\n")
	b.WriteString(m.renderTopStatus())
	b.WriteString("\n\n")
	b.WriteString(m.renderBody())
	b.WriteString("\n\n")
	b.WriteString(m.renderHelp())
	if m.showConfirm {
		b.WriteString("\n\n")
		b.WriteString(m.renderConfirm())
	}
	if m.showResults {
		b.WriteString("\n\n")
		b.WriteString(m.renderResults())
	}
	return b.String()
}

func (m controlModel) renderTopStatus() string {
	ip := m.devices[m.cursorDevice].Device.IP
	name := m.devices[m.cursorDevice].Device.Name
	cur := func(k string) string {
		if v, ok := m.current[k]; ok {
			return v
		}
		return "?"
	}
	line := fmt.Sprintf("Focus: %s (%s) | pow=%s mode=%s stemp=%s f_rate=%s f_dir=%s",
		name, ip, cur("pow"), cur("mode"), cur("stemp"), cur("f_rate"), cur("f_dir"))
	if m.err != nil {
		line += "  " + ctrlErrStyle.Render(m.err.Error())
	}
	return line
}

func (m controlModel) renderBody() string {
	left := m.renderControls()
	right := m.renderDeviceList()
	return lipgloss.JoinHorizontal(lipgloss.Top, left, lipgloss.NewStyle().Width(4).Render(""), right)
}

func (m controlModel) renderControls() string {
	var sb strings.Builder
	sb.WriteString(ctrlLabelStyle.Render("Staged Settings"))
	sb.WriteString("\n")
	row := func(label, cur, val string, focused bool) {
		line := fmt.Sprintf("%-9s current:%-4s  ➜  %-4s", label+":", cur, val)
		if focused {
			sb.WriteString(ctrlSelectStyle.Render(line))
		} else {
			sb.WriteString(ctrlValueStyle.Render(line))
		}
		sb.WriteString("\n")
	}
	row("Power", dispPower(m.pending.Power), dispPower(m.pending.Power), m.focus == focusPower)
	row("Mode", m.current["mode"], m.pending.Mode, m.focus == focusMode)
	row("Temp", m.current["stemp"], m.pending.Temp, m.focus == focusTemp)
	row("FanRate", m.current["f_rate"], m.pending.FanRate, m.focus == focusFanRate)
	row("FanDir", m.current["f_dir"], m.pending.FanDir, m.focus == focusFanDir)
	return sb.String()
}

func (m controlModel) renderDeviceList() string {
	var sb strings.Builder
	sb.WriteString(ctrlLabelStyle.Render("Selected Devices"))
	sb.WriteString("\n")
	for i, d := range m.devices {
		cursor := " "
		if i == m.cursorDevice {
			cursor = ">"
		}
		line := fmt.Sprintf("%s %s (%s)", cursor, d.Device.Name, d.Device.IP)
		if m.focus == focusDevice && i == m.cursorDevice {
			sb.WriteString(ctrlSelectStyle.Render(line))
		} else {
			sb.WriteString(line)
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

func (m controlModel) renderHelp() string {
	if !m.showHelp {
		return ctrlDimStyle.Render("Tab/Shift+Tab focus • ↑/↓ navigate • ←/→ cycle • p power • m mode • +/- temp • f fan rate • d fan dir • r refresh • a apply • h help • q quit")
	}
	return ctrlDimStyle.Render("Keybindings:\n  Navigation: Tab/Shift+Tab, ↑/↓, ←/→\n  Power: p\n  Mode: m (or ←/→ on Mode)\n  Temp: +/- or ↑/↓ when Temp focused\n  Fan Rate: f (or ←/→ on FanRate)\n  Fan Dir: d (or ←/→ on FanDir)\n  Refresh: r\n  Apply all: a (then y/n)\n  Help: h\n  Close modal: Esc\n  Quit: q")
}

func (m controlModel) renderConfirm() string {
	count := len(m.devices)
	msg := fmt.Sprintf("Apply to %d device(s)? pow=%s mode=%s stemp=%s f_rate=%s f_dir=%s  [y/n]",
		count, m.pending.Power, m.pending.Mode, m.pending.Temp, m.pending.FanRate, m.pending.FanDir)
	return ctrlSelectStyle.Render(msg)
}

func (m controlModel) renderResults() string {
	var sb strings.Builder
	sb.WriteString(ctrlLabelStyle.Render("Apply Results"))
	sb.WriteString("\n")
	for _, r := range m.applyResults {
		sb.WriteString(r)
		sb.WriteString("\n")
	}
	sb.WriteString(ctrlDimStyle.Render("Press Esc to close"))
	return sb.String()
}

func dispPower(p string) string {
	if p == "1" {
		return "On"
	}
	return "Off"
}

func (m controlModel) applyAll() tea.Cmd {
	return func() tea.Msg {
		var wg sync.WaitGroup
		mu := sync.Mutex{}
		res := make([]string, 0, len(m.devices))
		wg.Add(len(m.devices))
		for _, d := range m.devices {
			devIP := d.Device.IP
			go func(ip string) {
				defer wg.Done()
				ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
				defer cancel()
				cl := m.pending
				cl.IP = ip
				if err := api.SetClim(ctx, cl); err != nil {
					mu.Lock()
					res = append(res, fmt.Sprintf("ERR %s: %v", ip, err))
					mu.Unlock()
					return
				}
				mu.Lock()
				res = append(res, fmt.Sprintf("OK %s", ip))
				mu.Unlock()
			}(devIP)
		}
		wg.Wait()
		m.applyResults = res
		m.showResults = true
		return nil
	}
}

// RunControlScreen runs the interactive control UI
func RunControlScreen(devs []*storage.DeviceHistory) error {
	model := newControlModel(devs)
	p := tea.NewProgram(model)
	_, err := p.Run()
	return err
}
