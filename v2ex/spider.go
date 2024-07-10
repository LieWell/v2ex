package v2ex

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"liewell.fun/v2ex/core"
	"liewell.fun/v2ex/models"
)

var (
	thirtySecond = time.Second * 30
	tenMinute    = time.Minute * 10
	oneDay       = time.Hour * 24
)

func StartMemberSpider() {
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
		latestMember, dbError := models.FindLastMember()
		if dbError != nil {
			core.Logger.Infof("[spider] query db latest member error, wait 30 seconds and retry. \nerror: %v", dbError)
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
		for i := nextNumber; i <= siteState.MemberMax; i++ {
			core.Logger.Infof("[spider] start to process member[%d]", i)

			// 调用接口查询详情,可能遇到的错误:
			// 1. 真的出错了
			// 2. 数据不存在
			// 3. 触发限速
			rm, re := QueryMember(i)
			if re != nil {

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

				// 发生了不可预知错误,等待1分钟后重试
				core.Logger.Infof("[spider] query member[%d] error, wait 30 seconds and retry.\nerror: %v", i, re)
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
			core.Logger.Infof("[spider] insert real member[%d-%s] success.", i, member.Name)
		}
	}
}

func insertMember(rm *Member, fakeNumber int) (*models.Member, error) {
	var member *models.Member
	if fakeNumber != 0 {
		member = models.NewFakeMember(fakeNumber)
	} else {
		member = rm.toModel()
	}
	_, err := models.SaveMember(member)
	return member, err
}

func StartTopicSpider() {

}

/**
 * 头像存在 CDN 上,没有速率限制,因此可以持续运行
 */
func StartAvatarSpider() {
	startTime := time.Now()
	core.Logger.Infof("[StartAvatarSpider] start at [%s]", startTime.Format(core.DefaultTimeFormat))

	// TODO 从数据库查询已经扫描到哪个id

	// 保存头像根目录
	directory := core.GlobalConfig.V2ex.AvatarDir

	// 使用三个线程抓取
	step := 200 // 每个线程抓取数量
	var wg sync.WaitGroup
	wg.Add(3)
	for i := 0; i < 3; i++ {
		go func(p int) {
			defer wg.Done()

			// 确定各自负责的 number 范围
			s := p * step
			start := s + 1
			end := s + step
			memberList, err := models.FindCustomerAvatarMember(start, end)
			if err != nil {
				core.Logger.Errorf("[StartAvatarSpider] query member avatar address scope[%d:%d] error: %s", start, end, err)
				return
			}

			// 遍历结果抓取头像并保存
			for _, kv := range memberList {
				name := kv.StringOne
				imageURL := kv.StringTwo
				filename := fmt.Sprintf("%s.png", name)
				err := GetImageAndSave(imageURL, filename, directory)
				if err != nil {
					core.Logger.Errorf("[StartAvatarSpider] query member[%s] avatar error: %s", name, err)
					// 如果出错了,就记录到数据库充后续重试
					saveAvatar(kv.StringOne, imageURL)
					continue
				}
				core.Logger.Infof("[StartAvatarSpider] query member[%s] avatar success!", name)
			}
		}(i)
	}
	wg.Wait()
	core.Logger.Infof("[StartAvatarSpider] end cost: %vs", time.Now().Sub(startTime).Seconds())
}

func CheckMissingAvatarSpider() {

	startTime := time.Now()
	core.Logger.Infof("[CheckMissingAvatarSpider] start at [%s]", startTime.Format(core.DefaultTimeFormat))

	// 保存头像根目录
	directory := core.GlobalConfig.V2ex.AvatarDir

	avatars, err := models.FindAllAvatar()
	if err != nil {
		core.Logger.Errorf("[CheckMissingAvatarSpider] query from db error: %s", err)
		return
	}

	for _, avatar := range avatars {
		id := avatar.Id
		name := avatar.Name
		imageURL := avatar.Avatar
		filename := fmt.Sprintf("%s.png", name)
		err := GetImageAndSave(imageURL, filename, directory)
		if err != nil {
			// 如果出错了,不从表中删除
			core.Logger.Errorf("[CheckMissingAvatarSpider] query member[%s] avatar error: %s", name, err)
			continue
		}
		// 成功保存了头像,从表中删除,后续不会再次查询
		_ = models.DeleteAvatar(id)
		core.Logger.Infof("[CheckMissingAvatarSpider] query member[%s] avatar success!", name)
	}
	core.Logger.Infof("[CheckMissingAvatarSpider] end cost: %vs", time.Now().Sub(startTime).Seconds())
}

func saveAvatar(name, imageURL string) {
	avatar := &models.Avatar{
		Name:   name,
		Avatar: imageURL,
	}
	_, err := models.SaveAvatar(avatar)
	if err != nil {
		core.Logger.Errorf("[StartAvatarSpider] save member[%s] to table avatar error: %s", name, err)
	}
}
