package core

import "github.com/fatih/color"

// ErrorColorUnderline - underlined bold red color useful for strong error messages
var ErrorColorUnderline = color.New(color.FgRed).Add(color.Underline)

// ErrorColorBold - bold red color useful for error messages
var ErrorColorBold = color.New(color.FgRed).Add(color.Bold)

// WarningColorBold - bold yellow color useful for warning messages
var WarningColorBold = color.New(color.FgYellow).Add(color.Bold)

// SuccessColorBold - bold green color useful for success messages
var SuccessColorBold = color.New(color.FgGreen).Add(color.Bold)

// GreenColor - green color for text output
var GreenColor = color.New(color.FgGreen)

// RedColor - red color for text output
var RedColor = color.New(color.FgRed)

// ErrorColorBoldIns - insert variant for variables, bold red color useful for error messages
var ErrorColorBoldIns = color.New(color.FgRed).Add(color.Bold).SprintFunc()

// GreenColorIns - insert variant for variables, green color for text output
var GreenColorIns = color.New(color.FgGreen).SprintFunc()

// RedColorIns - insert variant for variables, red color for text output
var RedColorIns = color.New(color.FgRed).SprintFunc()
