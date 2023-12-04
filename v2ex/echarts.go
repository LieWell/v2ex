package v2ex

import (
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
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
	_ = c.AddFunc("30 0 * * *", func() {
		DrawMemberCountBar()
	})

	c.Start()
	select {}
}

func DrawMemberCountBar() {

	// 页面布局
	page := components.NewPage()
	page.SetLayout(components.PageFlexLayout)
	page.PageTitle = core.ServerGroup

	// 新增图标
	bar := _drawMemberCountBar()
	page.AddCharts(bar)

	// 页面渲染
	f, _ := os.Create("./static/template/members_count.html")
	_ = page.Render(f)
}

func _drawMemberCountBar() *charts.Bar {
	kvList, err := models.CountMember()
	if err != nil {
		core.Logger.Error(err)
	}
	var xAxis []string
	var series []opts.BarData
	for _, kv := range kvList {
		// 排除非法数据
		if kv.Date == "1970-01" {
			continue
		}
		xAxis = append(xAxis, kv.Date)
		series = append(series, opts.BarData{
			Value: kv.Count,
			Tooltip: &opts.Tooltip{
				Show:      true,
				TriggerOn: "mousemove",
			},
		})
	}
	bar := charts.NewBar()
	bar.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{
			Theme: types.ThemeWalden,
			Width: "1800px",
		}),
		charts.WithTitleOpts(opts.Title{
			Title: "会员数量统计(月)",
		}),
	)
	bar.SetXAxis(xAxis).AddSeries("数量", series)
	return bar
}
