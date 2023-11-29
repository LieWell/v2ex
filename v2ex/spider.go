package v2ex

import (
	"liewell.fun/v2ex/core"
	"liewell.fun/v2ex/models"
	"net/http"
	"time"
)

var (
	thirtySecond = time.Second * 30
	tenMinute    = time.Minute * 10
	oneDay       = time.Hour * 24
)

func StartSpider() {
	core.Logger.Info("[spider] start to crawl v2ex member, first query total member count.")
	for {

		// 查询当前存在多少注册会员
		// 出错则等待一分钟后重试
		siteState, err := QuerySiteState()
		if err != nil {
			core.Logger.Infof("[spider] query v2ex total member count error, wait 30 seconds and retry.\nerror: %v", err)
			time.Sleep(thirtySecond)
			continue
		}
		core.Logger.Infof("[spider] crawl v2ex total member count success, count: %d", siteState.MemberMax)

		// 查询上次中断的编号
		// 理论上查询数据库不应该出错,象征性的等待下
		latestMember, err := models.FindLastMember()
		if err != nil {
			core.Logger.Infof("[spider] query db latest member error, wait 30 seconds and retry. \nerror: %v", err)
			time.Sleep(thirtySecond)
			continue
		}

		// 检查是否已经爬取到了结尾
		latestMemberNumber := latestMember.Number
		if latestMemberNumber >= siteState.MemberMax {
			core.Logger.Infof("[spider] there is no new member found, wait 1 day and retry.")
			time.Sleep(oneDay)
			continue
		}

		// 从中断编号的下一个开始查询
		nextNumber := latestMemberNumber + 1
		core.Logger.Infof("[spider] from db find last member[%d]", latestMemberNumber)

		// 每次仅处理新增的数据
		for i := nextNumber; i < siteState.MemberMax; i++ {
			core.Logger.Infof("[spider] start to process member[%d]", i)

			// 调用接口查询详情,可能遇到的错误:
			// 1. 真的出错了
			// 2. 数据不存在
			// 3. 触发限速
			rm, re := QueryMemberById(i)
			if re != nil {

				// 处理情况 2
				// 数据不存在时,插入一条假数据代替,保证数据的连续性
				if re.StatusCode() == http.StatusNotFound {
					core.Logger.Infof("[spider] member[%d] not found, insert fake member instead.", i)
					_, e := insertMember(nil, i)
					if e != nil {
						core.Logger.Infof("[spider] insert fake member[%d] error, retry now.\nerror: %v", i, err)
						i--
						continue
					}
					// 继续处理下一条数据
					core.Logger.Infof("[spider] insert fake member[%d] success", i)
					continue
				}

				// 处理情况 3
				// 如果已达到限速,则等待到最近的 ResetAt 时间
				if re.RateLimit() {
					// 理论上需要等到 re.RateLimitReset() 时间点
					// 但是代理切换或者其他原因可能导致提前开始,因此每 10 分钟检查一次
					core.Logger.Infof("[spider] member[%d] rate limit occured, retry after 10 mintus.", i)
					time.Sleep(tenMinute)
					i--
					continue
				}

				// 发生了不可预知错误,等待1分钟后重试
				core.Logger.Infof("[spider] query member[%d] error, wait 30 seconds and retry.\nerror: %v", i, err)
				time.Sleep(thirtySecond)
				i--
				continue
			}

			// 数据入库
			member, e := insertMember(rm, 0)
			if e != nil {
				core.Logger.Infof("[spider] insert real member[%d] error, retry now. error: %v", i, e)
				i--
				continue
			}
			core.Logger.Infof("[spider] insert real member[%d-%s] success, go next.", i, member.Name)
		}
	}
}

func insertMember(rm *ResponseMember, fakeNumber int) (*models.Member, error) {
	var member *models.Member
	if fakeNumber != 0 {
		member = models.NewFakeMember(fakeNumber)
	} else {
		member = rm.toModel()
	}
	_, err := models.SaveMember(member)
	return member, err
}
