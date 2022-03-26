// Package chrome
// Time    : 2022/3/25 23:24
// Author  : xushiyin
// contact : yuqingxushiyin@gmail.com
package chrome

import (
	"fmt"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/input"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

type VideoInfo struct {
	Id        string
	Title     string
	ImageUrl  string
	Duration  time.Duration
	PlayCount int64
	PlayUrl   string
}

type PlayVideo struct {
	UserId    string
	Count     int64
	VInfoList []*VideoInfo
	Log       *log.Entry
}

func NewPlayVideo(userId string, count int64) *PlayVideo {
	return &PlayVideo{
		UserId: userId,
		Count:  count,
		Log: log.WithFields(log.Fields{
			"user_id": userId,
			"command": "play_count",
		}),
	}
}

// Str2TimeDuration 将字符串转换为时间长度
func Str2TimeDuration(str string) (time.Duration, error) {
	items := strings.Split(str, ":")

	switch len(items) {
	case 2:
		return time.ParseDuration(fmt.Sprintf("%sm%ss", items[0], items[1]))
	case 3:
		return time.ParseDuration(fmt.Sprintf("%sh%sm%ss", items[0], items[1], items[2]))
	default:
		return 0, fmt.Errorf("时间格式错误")
	}
}

// GetVideoInfoList 获取视频信息列表 from https://space.bilibili.com/{uid}/video, 仅获取当前页的视频信息，单页总共30条
func (pv *PlayVideo) GetVideoInfoList(page *rod.Page) error {

	ulElement, err := page.Element("#submit-video-list > ul.clearfix.cube-list")
	if err != nil {
		pv.Log.Errorf(`get element #submit-video-list > ul.clearfix.cube-list, error: %v`, err)
		return err
	}

	lElements, err := ulElement.Elements("li")
	if err != nil {
		pv.Log.Errorf("get video li element, error: %v", err)
		return err
	}

	for n, el := range lElements {
		vi := &VideoInfo{}
		nLog := pv.Log.WithFields(log.Fields{
			"video_index": n,
		})
		// video id
		vid, err := el.Attribute("data-aid")
		if err != nil {
			nLog.Errorf("get video id, error: %v", err)
			continue
		} else {
			vi.Id = *vid
		}

		// video title
		te, err := el.Element("a.title")
		if err != nil {
			nLog.Errorf("get %s, error: %v", "a.title", err)
		} else {
			vi.Title = te.MustText()
		}

		// video image url
		imgElement, err := el.Element("a.cover > img")
		if err != nil {
			nLog.Errorf("get %s , error: %v", "a.cover > img", err)
		} else {
			vi.ImageUrl = fmt.Sprintf(
				"https:%s", strings.Split(*imgElement.MustAttribute("src"), "@")[0])
		}

		// video duration
		le, err := el.Element("a.cover > span.length")
		if err != nil {
			nLog.Errorf("get %s , error: %v", "a.cover > span.length", err)
		} else {
			vi.Duration, err = Str2TimeDuration(le.MustText())
			if err != nil {
				nLog.Errorf("parse video lenth %s, %v, err: %v",
					le.MustText(), vi.Duration, err)
			}
		}

		// video play count
		pe, err := el.ElementX("div/span[1]")
		if err != nil {
			nLog.Errorf("get %s , error: %v", "div/span[1]", err)
		} else {
			vi.PlayCount, err = strconv.ParseInt(pe.MustText(), 10, 64)
			if err != nil {
				nLog.Warnf("parse video play count %s, %v, err: %v",
					pe.MustText(), vi.PlayCount, err)
			}
		}

		vi.PlayUrl = fmt.Sprintf("https://www.bilibili.com/video/%s", vi.Id)

		pv.VInfoList = append(pv.VInfoList, vi)
	}

	return nil
}

// VideoFastForward 视频快进, 方向键控制快进
func (pv *PlayVideo) VideoFastForward(p *rod.Page, n int) {
	totalSeconds := pv.VInfoList[n].Duration.Seconds()
	min := 1
	max := 4
	firstSleep := rand.Intn(max-min) + min
	time.Sleep(time.Duration(firstSleep) * time.Second)
	fiveCnt := int(totalSeconds) / (5 * 2)

	if fiveCnt > 10 {
		fiveCnt = 10
	}

	for i := 0; i < fiveCnt; i++ {
		p.Keyboard.MustPress(input.ArrowRight)
		time.Sleep(time.Duration(rand.Intn(max-2*min)+min) * time.Second)
	}
}

func (pv *PlayVideo) PlayVideos(cmd *cobra.Command, b *rod.Browser) {
	p := b.MustPage()

	select {
	case <-cmd.Context().Done():
		return
	default:
		for n, vi := range pv.VInfoList {
			if int64(n) >= pv.Count {
				break
			} else {
				p.MustNavigate(vi.PlayUrl)
				p.MustWaitLoad()
				pv.Log.Infoln(pv.VInfoList[n].Title)
				pv.Log.Infoln(pv.VInfoList[n].PlayUrl)
				pv.VideoFastForward(p, n)
			}
		}
	}
}
