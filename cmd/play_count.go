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
	userIdFlag   string
	videoCntFlag int
)

func init() {
	rootCmd.AddCommand(playCountCmd)
	playCountCmd.Flags().StringVarP(&userIdFlag, "user_id", "u", "94816944", "valid user id of bilibili (required)")
	playCountCmd.Flags().IntVarP(&videoCntFlag, "video_count", "c", 10, "count of video to play ")

	_ = playCountCmd.MarkFlagRequired("user_id")
}

var playCountCmd = &cobra.Command{
	Use:   "play_count",
	Short: "add video play count",
	Long:  "go https://space.bilibili.com/${uid}, find all newest videos, then add play count",
	Run: func(cmd *cobra.Command, args []string) {
		pLog := log.WithFields(log.Fields{
			"user_id": userIdFlag,
			"command": "play_count",
		})
		pLog.Infoln("start chrome to get newest videos list")
		url := url.URL{
			Scheme: "https",
			Host:   "space.bilibili.com",
			Path:   fmt.Sprintf("/%s", userIdFlag),
		}
		path, _ := launcher.LookPath()
		u := launcher.New().Bin(path).Logger(pLog.Writer()).Headless(false).MustLaunch()
		b := rod.New().ControlURL(u).MustConnect()
		page := b.MustPage(url.String()).Context(cmd.Context())
		page.MustWaitLoad()
		dataAids := page.MustElementsX("//div/a[@href and @target='_blank' and @class='cover']")

		Urls := make([]string, 1)

		for n, aid := range dataAids {
			if n > videoCntFlag {
				break
			}
			uStr := aid.MustProperty("href").String()
			Urls = append(Urls, uStr)
			b.MustPage(uStr).MustWaitLoad()
			log.Infoln("play video: ", uStr)
		}
		pages, _ := b.Pages()
		for _, u := range Urls {
			pages.MustFindByURL(u).MustActivate()
			time.Sleep(time.Second * 5)
		}

		time.Sleep(time.Second * 15)
	},
}
