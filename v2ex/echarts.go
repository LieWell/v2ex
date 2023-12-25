package v2ex

import (
	"fmt"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"
	"github.com/robfig/cron"
	"liewell.fun/v2ex/core"
	"liewell.fun/v2ex/models"
	"os"
	"time"
)

var (
	rootRenderDir  = "./static/template/"
	initialization = opts.Initialization{
		Theme: types.ThemeWalden,
		Width: "1800px",
	}
)

// StartDrawCharts 绘制各个图表
func StartDrawCharts() {
	c := cron.New()

	// 每天 00:30 执行
	_ = c.AddFunc("30 0 * * *", func() {

		// 绘制图表
		DrawMemberCountBar("members_count.html")
		DrawMemberTrendLine("members_trend.html")

		// 更新绘制时间
		_ = models.UpdateSystemConfig(models.SystemConfigKeyLastDrawTime, time.Now().Format(core.DefaultTimeFormat))
	})

	c.Start()
	select {}
}

// ====================== 会员数量统计 start ======================

func DrawMemberCountBar(fileName string) {

	page := components.NewPage()
	page.SetLayout(components.PageFlexLayout)
	page.PageTitle = core.ServerGroup

	// 会员统计柱状图
	bar := _drawMemberCountBar()
	page.AddCharts(bar)

	// 页面渲染
	f, _ := os.Create(fmt.Sprintf("%s%s", rootRenderDir, fileName))
	_ = page.Render(f)
}

func _drawMemberCountBar() *charts.Bar {
	kvList, err := models.StatisticsMember()
	if err != nil {
		core.Logger.Errorf("DrawMemberCountBar error: %v", err)
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
		charts.WithInitializationOpts(initialization),
		charts.WithTitleOpts(opts.Title{
			Title: "会员数量统计(月)",
		}),
	)
	bar.SetXAxis(xAxis).AddSeries("数量", series)
	return bar
}

// ====================== 会员数量统计 end ======================

// ====================== 会员增长趋势统计 start ======================

func DrawMemberTrendLine(fileName string) {

	page := components.NewPage()
	page.SetLayout(components.PageFlexLayout)
	page.PageTitle = core.ServerGroup

	// 增长趋势折线图
	line := _drawMemberTrendLine()
	page.AddCharts(line)

	// 页面渲染
	f, _ := os.Create(fmt.Sprintf("%s%s", rootRenderDir, fileName))
	_ = page.Render(f)
}

func _drawMemberTrendLine() *charts.Line {

	kvList, err := models.StatisticsMemberTrend()
	if err != nil {
		core.Logger.Errorf("DrawMemberTrendLine error: %v", err)
	}
	var xAxis []string
	var series []opts.LineData
	for _, kv := range kvList {
		// 排除非法数据
		if kv.Date == "1970-01" {
			continue
		}
		xAxis = append(xAxis, kv.Date)
		series = append(series, opts.LineData{
			Value: kv.Count,
		})
	}

	line := charts.NewLine()
	line.SetGlobalOptions(
		charts.WithInitializationOpts(initialization),
		charts.WithTitleOpts(opts.Title{
			Title: "会员增长趋势统计(月)",
		}),
		charts.WithTooltipOpts(opts.Tooltip{
			Show:    true,
			Trigger: "axis",
		}),
	)

	line.SetXAxis(xAxis).AddSeries("数量", series).
		SetSeriesOptions(charts.WithLineChartOpts(
			opts.LineChart{
				Smooth:     true,
				ShowSymbol: true,
				SymbolSize: 10,
			}),
		)
	return line
}

// ====================== 会员增长趋势统计 end ========================
