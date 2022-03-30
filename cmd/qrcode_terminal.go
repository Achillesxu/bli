// Package cmd
// Time    : 2022/3/29 23:00
// Author  : xushiyin
// contact : yuqingxushiyin@gmail.com
package cmd

import (
	myQr "github.com/Achillesxu/bli/qrcode"
	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/qrcode"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"image"
	"os"
)

var (
	picPathFlag string
)

func init() {
	rootCmd.AddCommand(qrCodeTerminalCmd)
	qrCodeTerminalCmd.Flags().StringVarP(&picPathFlag, "pic_path", "p", "", "qrcode picture path")

	_ = playCountCmd.MarkFlagRequired("pic_path")

}

var qrCodeTerminalCmd = &cobra.Command{
	Use:   "qrcode_terminal",
	Short: "display qrcode in terminal",
	Long:  "scan qrcode picture and display it in terminal",
	Run: func(cmd *cobra.Command, args []string) {
		_, err := os.Stat(picPathFlag)
		if err != nil {
			log.Errorf("stat file error: %v", err)
			return
		}

		file, _ := os.Open(picPathFlag)
		img, _, _ := image.Decode(file)

		bmp, _ := gozxing.NewBinaryBitmapFromImage(img)
		// decode image
		qrReader := qrcode.NewQRCodeReader()
		result, _ := qrReader.Decode(bmp, nil)
		log.Infof("qrcodde result: %v", result.GetText())

		_ = myQr.PrintQrCode2TTY("white", "black", "l", result.GetText())
	},
}
