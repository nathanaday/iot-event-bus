package utils

var Reset = "\033[0m"
var Red = "\033[31m"
var Green = "\033[32m"
var Yellow = "\033[33m"
var Blue = "\033[34m"
var Magenta = "\033[35m"
var Cyan = "\033[36m"
var Gray = "\033[37m"
var White = "\033[97m"

func StrToRed(s string) string {
	return Red + s + Reset
}

func StrToGreen(s string) string {
	return Green + s + Reset
}

func StrToYellow(s string) string {
	return Yellow + s + Reset
}

func StrToBlue(s string) string {
	return Blue + s + Reset
}
