package tui

import (
	"fmt"
	"os"
	"time"

	"cli-tutor/src/printer"
	"cli-tutor/src/tui/lessonui"
	"cli-tutor/src/tui/menuui"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/termenv"
)

type sessionState int

const (
	menuView sessionState = iota
	lessonView
)

type MainModel struct {
	state      sessionState
	menu       tea.Model
	lesson     tea.Model
	quitting   bool
	windowsize tea.WindowSizeMsg
}

func (m MainModel) Init() tea.Cmd {
	termenv.ClearScreen()
	printer.Print("Welcome to Chistole", "welcome")
	time.Sleep(1 * time.Second)
	return nil
}

func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.windowsize = msg // pass this along to the entry view so it uses the full window size when it's initialized

	case lessonui.BackMsg:
		m.state = menuView
	}

	switch m.state {
	case menuView:
		menu, cmd := m.menu.Update(msg)
		cmds = append(cmds, cmd)
		menuModel, ok := menu.(menuui.MenuModel)
		if !ok {
			panic("something went wrong with the menu ui ")
		}
		m.menu = menuModel

	case lessonView:
		m.lesson.Update(msg)
		// lesson, _ := m.lesson.Update(msg)
		// lessonModel, ok := lesson.(lessonui.LessonModel)
		// if !ok {
		// 	panic("something went wrong with the lesson ui ")
		// }
		// m.menu = lessonModel
	}
	return m, tea.Batch(cmds...)
}

func (m MainModel) View() string {
	switch m.state {

	case lessonView:
		return m.lesson.View()

	default:
		return m.menu.View()

	}
}

func New() MainModel {
	return MainModel{
		state:      menuView,
		menu:       menuui.New(),
		lesson:     nil,
		quitting:   false,
		windowsize: tea.WindowSizeMsg{},
	}
}

func StartUI() {
	m := New()
	if err := tea.NewProgram(m).Start(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
