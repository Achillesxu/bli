// Package cmd
// Time    : 2022/3/17 22:01
// Author  : xushiyin
// contact : yuqingxushiyin@gmail.com
package cmd

import (
	"fmt"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"net/url"
	"time"
)

var (
	userIdFlag string
)

func init() {
	rootCmd.AddCommand(playCountCmd)
	playCountCmd.Flags().StringVarP(&userIdFlag, "user_id", "u", "94816944", "valid user id of bilibili (required)")

	_ = playCountCmd.MarkFlagRequired("user_id")
}

var playCountCmd = &cobra.Command{
	Use:   "play_count",
	Short: "add video play count",
	Long:  "go https://space.bilibili.com/${uid}, find all newest videos, then add play count",
	Run: func(cmd *cobra.Command, args []string) {
		log.WithFields(log.Fields{
			"command": "play_count",
			"user_id": userIdFlag,
		}).Infoln("start chrome to get newest videos list")
		url := url.URL{
			Scheme: "https",
			Host:   "space.bilibili.com",
			Path:   fmt.Sprintf("/%s", userIdFlag),
		}
		path, _ := launcher.LookPath()
		u := launcher.New().Bin(path).Headless(false).MustLaunch()
		b := rod.New().ControlURL(u).MustConnect()
		page := b.MustPage(url.String())
		page.MustWaitLoad()

		time.Sleep(time.Second * 5)

		page.MustWaitLoad().MustScreenshot("a.png")
	},
}
