// A simple console color package.

package color

import (
	"fmt"
	"strings"
)

type Col struct {
	R, G, B uint8
}

func NewCol(r, g, b int) Col {
	return Col{
		R: uint8(r),
		G: uint8(g),
		B: uint8(b),
	}
}

var Background = Col{40, 42, 54}
var Current = Col{68, 71, 90}
var Foreground = Col{248, 248, 242}

var Black = Col{0, 0, 0}
var Comment = Col{98, 114, 164}
var Cyan = Col{139, 233, 253}
var Green = Col{80, 250, 123}
var Orange = Col{255, 184, 108}
var Pink = Col{255, 121, 198}
var Purple = Col{189, 147, 249}
var Red = Col{255, 85, 85}
var Yellow = Col{241, 250, 140}

func Fg(c Col) string {
	return fmt.Sprintf("\x1b[38;2;%d;%d;%dm", c.R, c.G, c.B)
}
func Bg(c Col) string {
	return fmt.Sprintf("\x1b[48;2;%d;%d;%dm", c.R, c.G, c.B)
}
func Res() string {
	return "\x1b[0m"
}
func Bld() string {
	return "\x1b[1m"
}
func Ul() string {
	return "\x1b[4m"
}

type ColorReplace struct {
	Name  string
	Color string
}

type ColorManager struct {
	Colors []ColorReplace
}

var CM ColorManager

func (cm *ColorManager) Register(name string, color string) {
	cm.Colors = append(cm.Colors, ColorReplace{Name: name, Color: color})
}

func (cm *ColorManager) Printf(str string, args ...any) *ColorManager {
	tmp := cm.Sprintf(str, args...)
	fmt.Print(tmp)
	return cm
}
func (cm *ColorManager) Sprintf(str string, args ...any) string {
	tmp := fmt.Sprintf(str, args...)
	for _, cr := range cm.Colors {
		tmp = strings.ReplaceAll(tmp, cr.Name, cr.Color)
	}
	return tmp
}

func init() {

	CM.Register("[bg]", Fg(Background))
	CM.Register("[comment]", Fg(Comment))
	CM.Register("[cyan]", Fg(Cyan))
	CM.Register("[green]", Fg(Green))
	CM.Register("[orange]", Fg(Orange))
	CM.Register("[pink]", Fg(Pink))
	CM.Register("[purple]", Fg(Purple))
	CM.Register("[red]", Fg(Red))
	CM.Register("[yellow]", Fg(Yellow))

	CM.Register("[_bg]", Bg(Background))
	CM.Register("[_comment]", Bg(Comment))
	CM.Register("[_cyan]", Bg(Cyan))
	CM.Register("[_green]", Bg(Green))
	CM.Register("[_orange]", Bg(Orange))
	CM.Register("[_pink]", Bg(Pink))
	CM.Register("[_purple]", Bg(Purple))
	CM.Register("[_red]", Bg(Red))
	CM.Register("[_yellow]", Bg(Yellow))

	CM.Register("[ul]", Ul())
	CM.Register("[bld]", Bld())

	CM.Register("[res]", Res())

}

func Show() {
	for _, cr1 := range CM.Colors {
		if !strings.Contains(cr1.Name, "_") {
			for _, cr2 := range CM.Colors {
				if strings.Contains(cr2.Name, "_") {
					CM.Printf("%s%s(%s/[bld]%s)[res] ", cr1.Name, cr2.Name, cr1.Name[1:4], cr2.Name[1:5])
				}
			}
			fmt.Printf("\n")
		}
	}
}
