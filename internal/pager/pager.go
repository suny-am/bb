package pager

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

var (
	baseStyle = lipgloss.NewStyle().
			PaddingLeft(1).
			PaddingTop(1)
	headerStyle = func(width int) lipgloss.Style {
		return lipgloss.NewStyle().
			BorderBottom(true).
			BorderStyle(lipgloss.NormalBorder()).Width(width)
	}
	topOffset int = 4
)

type binderModel struct {
	pages       []pageModel
	currentPage int

	currentPageSize int

	verticalScroll int

	windowHeight int
	windowWidth  int
}

type pageModel struct {
	content string
	id      int
}

func (bm binderModel) Init() tea.Cmd {
	return nil
}

func (bm binderModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		bm.windowWidth = msg.Width
		bm.windowHeight = msg.Height

	case tea.KeyMsg:
		switch msg.String() {

		case "q", "ctrl+c", "esc":
			return bm, tea.Quit

		case "right":
			if bm.currentPage < len(bm.pages)-1 {
				bm.currentPage++
				bm.verticalScroll = 0
				bm.currentPageSize = len(bm.splitPage())
			}

		case "left":
			if bm.currentPage > 0 {
				bm.currentPage--
				bm.verticalScroll = 0
				bm.currentPageSize = len(bm.splitPage())
			}

		case "down":
			if bm.currentPageSize > bm.windowHeight {
				bm.verticalScroll++
			}
		case "up":
			if bm.verticalScroll > 0 && bm.currentPageSize > bm.windowHeight {
				bm.verticalScroll--
			}
		}

	case tea.MouseMsg:
		if tea.MouseEvent(msg).IsWheel() && bm.currentPageSize > bm.windowHeight {
			switch tea.MouseButton(msg.Button) {
			case tea.MouseButtonWheelUp:
				if bm.verticalScroll > 0 {
					bm.verticalScroll--
				}
			case tea.MouseButtonWheelDown:
				bm.verticalScroll++
			}
		}
	}

	return bm, nil
}

func (bm *binderModel) splitPage() []string {
	return strings.Split(bm.pages[bm.currentPage].content, "\n")
}

func (bm binderModel) View() string {
	var sb strings.Builder

	header := headerStyle(bm.windowWidth).Render(fmt.Sprintf("Page: %d",
		bm.currentPage,
	))

	split := bm.splitPage()

	for i, v := range split {
		if i < bm.windowHeight-topOffset && i > bm.verticalScroll {
			sb.WriteString(fmt.Sprintf("%s\n", v))
		}
	}

	body := baseStyle.Render(sb.String())
	return fmt.Sprintf("%s\n%s", header, body)
}

func Draw(pagesContent []string) {
	width, height, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		panic(err)
	}

	pager := binderModel{}

	for i, v := range pagesContent {
		pager.pages = append(pager.pages, pageModel{
			v,
			i,
		})
	}

	pager.currentPage = 0
	pager.windowWidth = width
	pager.windowHeight = height
	pager.currentPageSize = height

	if _, err := tea.NewProgram(pager,
		tea.WithAltScreen(),
		tea.WithMouseAllMotion()).
		Run(); err != nil {
		panic(err)
	}
}
