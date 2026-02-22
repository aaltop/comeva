// Ansi colour code handling.

package ansi

import (
	"fmt"
)

// CSI is the ANSI Control Sequence Introducer.
const CSI = "\x1b["

// SGRReset is the ANSI control sequence for resetting the Select Graphic Rendition
// (SGR) sequence.
const SGRReset = "\x1b[m"

// TODO: replace with [color.RGBA]

// RGBColor represents an rgb colour.
type RGBColor struct {
	R, G, B uint8
}

func NewDefaultRGBColor() (rgbColor *RGBColor) {
	return &RGBColor{}
}

func NewRGBColor(r, g, b uint8) (rgbColor *RGBColor) {
	rgbColor = NewDefaultRGBColor()
	rgbColor.R = r
	rgbColor.G = g
	rgbColor.B = b
	return
}

func (rgbcolor *RGBColor) String() string {
	return fmt.Sprintf("%d;%d;%d", rgbcolor.R, rgbcolor.G, rgbcolor.B)
}

func (rgbColor *RGBColor) ToRGB() *RGBColor {
	return rgbColor
}

// Color represents a generic colour with transforms to other colours.
type Color interface {
	ToRGB() *RGBColor
}

// ColorScheme represents the colours used for printing.
type ColorScheme struct {
	Foreground Color
	Background Color
}

func NewDefaultColorScheme() (scheme *ColorScheme) {
	scheme = &ColorScheme{}
	scheme.Foreground = NewRGBColor(255, 255, 255)
	scheme.Background = NewRGBColor(0, 0, 0)
	return
}

func NewColorScheme(foreground, background Color) (scheme *ColorScheme) {
	scheme = NewDefaultColorScheme()
	scheme.Foreground = foreground
	scheme.Background = background
	return
}

func NewColorSchemeSimple(fr, fg, fb, br, bg, bb uint8) (scheme *ColorScheme) {
	return NewColorScheme(
		NewRGBColor(fr, fg, fb),
		NewRGBColor(br, bg, bb),
	)
}

func (printColors *ColorScheme) String() string {
	return fmt.Sprintf("%v38;2;%v;48;2;%vm", CSI, printColors.Foreground.ToRGB(), printColors.Background.ToRGB())
}

// Apply returns a string formatted with the colours.
func (printColors *ColorScheme) Apply(format string, a ...any) string {
	return printColors.String() + fmt.Sprintf(format, a...) + SGRReset
}

// Fore returns the foreground code.
func (printColors *ColorScheme) Fore() string {
	return fmt.Sprintf("%v38;2;%vm", CSI, printColors.Foreground.ToRGB())
}

// ApplyFore returns a string formatted with the foreground colour.
func (printColors *ColorScheme) ApplyFore(format string, a ...any) string {
	return printColors.Fore() + fmt.Sprintf(format, a...) + SGRReset
}

// Back returns the background code.
func (printColors *ColorScheme) Back() string {
	return fmt.Sprintf("%v48;2;%vm", CSI, printColors.Background.ToRGB())
}

// ApplyBack returns a string formatted with the background colour.
func (printColors *ColorScheme) ApplyBack(format string, a ...any) string {
	return printColors.Back() + fmt.Sprintf(format, a...) + SGRReset
}

// BasicColors contains basic [Color]s for ease of use.
var BasicColors = struct {
	Red, Green, Blue, Black, White, Orange Color
}{
	Red:    NewRGBColor(255, 0, 0),
	Green:  NewRGBColor(0, 255, 0),
	Blue:   NewRGBColor(0, 0, 255),
	Black:  NewRGBColor(0, 0, 0),
	White:  NewRGBColor(255, 255, 255),
	Orange: NewRGBColor(255, 157, 0),
}

// BasicColorSchemes contains basic [ColorScheme]s for ease of use.
var BasicColorSchemes = struct {
	// Error is formatting for an error message.
	Error *ColorScheme
	// Warning is formatting for a warning message.
	Warning *ColorScheme
}{
	Error:   NewColorScheme(BasicColors.Red, BasicColors.Black),
	Warning: NewColorScheme(BasicColors.Orange, BasicColors.Black),
}
