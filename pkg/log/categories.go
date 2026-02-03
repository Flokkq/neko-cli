package log

/*
@Author     Benjamin Senekowitsch
@Contact    senekowitsch@nekoman.at
@Since      20.12.2025
*/

type Category string

const (
	Init      Category = "init"
	Config    Category = "config"
	Preflight Category = "pre-flight"
	Guard     Category = "guard"
	Exec      Category = "exec"
)

var categoryColors = map[Category]string{
	Init:      ColorBrightYellow,
	Config:    ColorBrightCyan,
	Preflight: ColorBrightYellow,
	Guard:     ColorBrightBlue,
	Exec:      ColorBrightGreen,
}
