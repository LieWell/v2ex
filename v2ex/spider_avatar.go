package v2ex

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"liewell.fun/v2ex/core"
	"liewell.fun/v2ex/models"
)

/**
 * 头像存在 CDN 上,没有速率限制,因此可以持续运行
 */
func StartAvatarSpider() {
	startTime := time.Now()
	core.Logger.Infof("[StartAvatarSpider] start at [%s]", startTime.Format(core.DefaultTimeFormat))

	// 保存头像根目录
	directory := core.GlobalConfig.V2ex.AvatarDir

	// 使用三个线程抓取
	step := 10000 // 每个线程抓取数量
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

				// 不处理默认头像
				pass, newImageUrl := replaceAndFilterAvatar(imageURL)
				if !pass {
					core.Logger.Infof("[StartAvatarSpider] member[%s] is default avatar!", name)
					continue
				}

				// 尝试保存新头像
				filename := fmt.Sprintf("%s.png", name)
				err := GetImageAndSave(newImageUrl, filename, directory)
				if err != nil {
					core.Logger.Errorf("[StartAvatarSpider] query member[%s] avatar error: %s", name, err)
					// 如果出错了,就记录到数据库后续重试
					// TODO 意义不大
					// saveAvatar(kv.StringOne, imageURL)
					continue
				}
				core.Logger.Infof("[StartAvatarSpider] query member[%s] avatar success!", name)
			}
		}(i)
	}
	wg.Wait()
	core.Logger.Infof("[StartAvatarSpider] end cost: %v", time.Since(startTime))
}

func replaceAndFilterAvatar(rawUrl string) (bool, string) {
	// 默认头像不尽兴处理
	if strings.Contains(rawUrl, "/gravatar/") {
		return false, ""
	}
	// 查询 large 尺寸的图像
	return true, strings.Replace(rawUrl, "_normal.png", "_large.png", 1)
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
	core.Logger.Infof("[CheckMissingAvatarSpider] end cost: %v", time.Since(startTime))
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
