// Package qrcode
// Time    : 2022/3/29 23:28
// Author  : xushiyin
// contact : yuqingxushiyin@gmail.com
package qrcode

import (
	"bytes"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	qrcode "github.com/skip2/go-qrcode"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const (
	NormalBlack   = "\033[38;5;0m  \033[0m"
	NormalRed     = "\033[38;5;1m  \033[0m"
	NormalGreen   = "\033[38;5;2m  \033[0m"
	NormalYellow  = "\033[38;5;3m  \033[0m"
	NormalBlue    = "\033[38;5;4m  \033[0m"
	NormalMagenta = "\033[38;5;5m  \033[0m"
	NormalCyan    = "\033[38;5;6m  \033[0m"
	NormalWhite   = "\033[38;5;7m  \033[0m"

	BrightBlack   = "\033[48;5;0m  \033[0m"
	BrightRed     = "\033[48;5;1m  \033[0m"
	BrightGreen   = "\033[48;5;2m  \033[0m"
	BrightYellow  = "\033[48;5;3m  \033[0m"
	BrightBlue    = "\033[48;5;4m  \033[0m"
	BrightMagenta = "\033[48;5;5m  \033[0m"
	BrightCyan    = "\033[48;5;6m  \033[0m"
	BrightWhite   = "\033[48;5;7m  \033[0m"
)

// PrintQrCode2TTY support as following:
// Supported background colors: [black, red, green, yellow, blue, magenta, cyan, white]
// Supported front colors: [black, red, green, yellow, blue, magenta, cyan, white]
// Supported error correction levels: [L, M, Q, H]
func PrintQrCode2TTY(frontColor, backgroundColor, levelStr, content string) error {

	level, err := level2Level(levelStr)
	if err != nil {
		log.Errorf("levelStr(%s) should [L, M, Q, H] in error: %v", levelStr, err)
		return err
	}

	qr, err := qrcode.New(content, level)
	if err != nil {
		log.Errorf("PrintQrCode2TTY error: %v", err)
		return err
	}

	_, _ = getTTYSize()
	front, err := color2ansi(frontColor)
	if err != nil {
		log.Errorf("frontColor: %s %s, err: %s", frontColor, "[black, red, green, yellow, blue, magenta, cyan, white]", err)
		return err
	}

	back, err := color2ansi(backgroundColor)
	if err != nil {
		log.Errorf("backgroundColor: %s %s, err: %s", frontColor, "[black, red, green, yellow, blue, magenta, cyan, white]", err)
		return err

	}

	bitmap := qr.Bitmap()
	output := bytes.NewBuffer([]byte{})
	for ir, row := range bitmap {
		lr := len(row)

		if ir == 0 || ir == 1 || ir == 2 ||
			ir == lr-1 || ir == lr-2 || ir == lr-3 {
			continue
		}

		for ic, col := range row {
			lc := len(bitmap)
			if ic == 0 || ic == 1 || ic == 2 ||
				ic == lc-1 || ic == lc-2 || ic == lc-3 {
				continue
			}
			if col {
				output.WriteString(front)
			} else {
				output.WriteString(back)
			}
		}
		output.WriteByte('\n')
	}
	_, _ = output.WriteTo(os.Stdout)
	return nil
}

func color2ansi(c string) (color string, err error) {
	s := strings.ToUpper(c)
	switch s {
	case "BLACK":
		color = BrightBlack
	case "RED":
		color = BrightRed
	case "GREEN":
		color = BrightGreen
	case "YELLOW":
		color = BrightYellow
	case "BLUE":
		color = BrightBlue
	case "MAGENTA":
		color = BrightMagenta
	case "CYAN":
		color = BrightCyan
	case "WHITE":
		color = BrightWhite
	default:
		err = errors.New(fmt.Sprintf("'%s' is not support.", c))
	}
	return
}

func level2Level(str string) (level qrcode.RecoveryLevel, err error) {
	s := strings.ToUpper(str)
	switch s {
	case "L":
		level = qrcode.Low
	case "M":
		level = qrcode.Medium
	case "Q":
		level = qrcode.High
	case "H":
		level = qrcode.Highest
	default:
		err = errors.New(fmt.Sprintf("'%s' is not support.", str))
	}

	return
}

func getTTYSize() (int, int) {
	cmd := exec.Command("stty", "size")
	cmd.Stdin = os.Stdin
	out, err := cmd.Output()
	if err != nil {
		return 0, 0
	}
	outStr := strings.Replace(string(out), "\n", "", -1)
	cols, err := strconv.Atoi(strings.Split(outStr, " ")[1])
	if err != nil {
		return 0, 0
	}
	rows, err := strconv.Atoi(strings.Split(outStr, " ")[0])
	if err != nil {
		return 0, 0
	}
	return cols, rows
}
