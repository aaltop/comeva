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

// TODO: Could utilise the std image/color package here?

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

// Apply returns the string formatted with the colours.
func (printColors *ColorScheme) Apply(s string) string {
	return printColors.String() + s + SGRReset
}

// Applyf works like [fmt.Sprintf] and formats the result with the colours.
func (printColors *ColorScheme) Applyf(format string, a ...any) string {
	return printColors.Apply(fmt.Sprintf(format, a...))
}

// Fore returns the foreground code.
func (printColors *ColorScheme) Fore() string {
	return fmt.Sprintf("%v38;2;%vm", CSI, printColors.Foreground.ToRGB())
}

// ApplyFore returns the string formatted with the foreground colour.
func (printColors *ColorScheme) ApplyFore(s string) string {
	return printColors.Fore() + s + SGRReset
}

// ApplyForef works like [fmt.Sprintf] and formats the result with the foreground colour.
func (printColors *ColorScheme) ApplyForef(format string, a ...any) string {
	return printColors.ApplyFore(fmt.Sprintf(format, a...))
}

// Back returns the background code.
func (printColors *ColorScheme) Back() string {
	return fmt.Sprintf("%v48;2;%vm", CSI, printColors.Background.ToRGB())
}

// ApplyBack returns the string formatted with the background colour.
func (printColors *ColorScheme) ApplyBack(s string) string {
	return printColors.Back() + s + SGRReset
}

// ApplyBackf works like [fmt.Sprintf] and formats the result with the background colour.
func (printColors *ColorScheme) ApplyBackf(format string, a ...any) string {
	return printColors.ApplyBack(fmt.Sprintf(format, a...))
}

// BasicColors contains basic [Color]s for ease of use.
var BasicColors = struct {
	Red,
	Green,
	Blue, SkyBlue,
	Black, White, Grey, LightGrey,
	Orange,
	Magenta Color
}{
	Red:       NewRGBColor(255, 0, 0),
	Green:     NewRGBColor(0, 255, 0),
	Blue:      NewRGBColor(0, 0, 255),
	SkyBlue:   NewRGBColor(135, 206, 235),
	Black:     NewRGBColor(0, 0, 0),
	White:     NewRGBColor(255, 255, 255),
	Grey:      NewRGBColor(128, 128, 128),
	LightGrey: NewRGBColor(211, 211, 211),
	Orange:    NewRGBColor(255, 157, 0),
	Magenta:   NewRGBColor(255, 0, 255),
}

// BasicColorSchemes contains basic [ColorScheme]s for ease of use.
var BasicColorSchemes = struct {
	// Debug is formatting for a debug message.
	Debug *ColorScheme
	// Info is formatting for an info message.
	Info *ColorScheme
	// Warning is formatting for a warning message.
	Warning *ColorScheme
	// Error is formatting for an error message.
	Error *ColorScheme
	// Critical is formatting for a critical message.
	Critical *ColorScheme
}{
	Debug:    NewColorScheme(BasicColors.SkyBlue, BasicColors.Black),
	Info:     NewColorScheme(BasicColors.White, BasicColors.Black),
	Warning:  NewColorScheme(BasicColors.Orange, BasicColors.Black),
	Error:    NewColorScheme(BasicColors.Red, BasicColors.Black),
	Critical: NewColorScheme(BasicColors.Magenta, BasicColors.Black),
}
