// Package cmd
// Time    : 2022/3/17 22:01
// Author  : xushiyin
// contact : yuqingxushiyin@gmail.com
package cmd

import (
	"fmt"
	"github.com/Achillesxu/bli/chrome"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"net/url"
	"time"
)

var (
	userIdFlag   string
	videoCntFlag int64
)

func init() {
	rootCmd.AddCommand(playCountCmd)
	playCountCmd.Flags().StringVarP(&userIdFlag, "user_id", "u", "94816944", "valid user id of bilibili (required)")
	playCountCmd.Flags().Int64VarP(&videoCntFlag, "video_count", "c", 10, "count of the newest video to play")

	_ = playCountCmd.MarkFlagRequired("user_id")
}

var playCountCmd = &cobra.Command{
	Use:   "play_count",
	Short: "add video play count",
	Long:  "go https://space.bilibili.com/${uid}/video, find all newest videos, then add play count",
	Run: func(cmd *cobra.Command, args []string) {
		pLog := log.WithFields(log.Fields{
			"user_id": userIdFlag,
			"command": "play_count",
		})

		url := url.URL{
			Scheme: "https",
			Host:   "space.bilibili.com",
			Path:   fmt.Sprintf("/%s/video", userIdFlag),
		}

		// l := launcher.NewManaged("")

		var browser *rod.Browser

		if path, err := launcher.LookPath(); err != false {
			l := launcher.New().Bin(path).Logger(pLog.Writer()).Set("autoplay-policy", "no-user-gesture-required")
			if isHeadless {
				l.Headless(isHeadless)
				l.Set("disable-gpu", "true")
			}
			u := l.MustLaunch()
			browser = rod.New().ControlURL(u).MustConnect()
		} else {
			pLog.Errorf("look path error: %v", err)
			return
		}
		pLog.Infof("start chrome to get newest videos from %s", url.String())

		page := browser.MustPage(url.String()).Context(cmd.Context())
		page.MustWaitLoad()
		time.Sleep(time.Second * 2)

		pv := chrome.NewPlayVideo(userIdFlag, videoCntFlag)

		if err := pv.GetVideoInfoList(page); err != nil {
			pLog.Errorf("get video info list error: %v", err)
			return
		} else {
			pLog.Printf("get %d videos\n", len(pv.VInfoList))
		}

		pv.PlayVideos(cmd, browser)

	},
}
