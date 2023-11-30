package v2ex

import (
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"
	"github.com/robfig/cron"
	"liewell.fun/v2ex/core"
	"liewell.fun/v2ex/models"
	"os"
)

// StartDrawCharts 绘制各个图表
func StartDrawCharts() {
	c := cron.New()

	// 每天 00:30 执行
	_ = c.AddFunc("*/10 * * * *", func() {
		DrawMemberCountBar()
	})

	c.Start()
	select {}
}

func DrawMemberCountBar() {
	kvList, err := models.CountMemberByYear()
	if err != nil {
		core.Logger.Error(err)
	}
	var xAxis []string
	var series []opts.BarData
	for _, kv := range kvList {
		xAxis = append(xAxis, kv.Date)
		series = append(series, opts.BarData{
			Value: kv.Count,
			Label: &opts.Label{
				Show: true,
			},
		})
	}
	bar := charts.NewBar()
	bar.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{
			Theme: types.ThemeMacarons,
		}),
		charts.WithTitleOpts(opts.Title{
			Title: "注册会员数统计(年)",
		}),
	)
	bar.SetXAxis(xAxis).AddSeries("会员数量", series)
	f, _ := os.Create("./static/template/members_count.html")
	_ = bar.Render(f)
}
