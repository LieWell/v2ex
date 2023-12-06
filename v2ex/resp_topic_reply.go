package v2ex

type Topic struct {
	Node            Node   `json:"node"`             // 所处节点
	Member          Member `json:"member"`           // 发帖人
	Id              int    `json:"id"`               // 帖子编号
	LastReplyBy     string `json:"last_reply_by"`    // 最后回复会员
	LastTouched     int    `json:"last_touched"`     // 最后回复时间,某些帖子可能不准确
	Title           string `json:"title"`            // 帖子标题
	Url             string `json:"url"`              // 帖子访问列表
	Created         int    `json:"created"`          // 创建时间戳
	Deleted         int    `json:"deleted"`          // 是否被删除
	Content         string `json:"content"`          // 原始内容
	ContentRendered string `json:"content_rendered"` // html 渲染之后内容
	LastModified    int    `json:"last_modified"`    // TODO 作用未知
	Replies         int    `json:"replies"`          // 回复数,此回复数是 AntiFlood 之后的回复数(顶/赞之类的无效回复会被屏蔽)
}

type Reply struct {
	Id              int    `json:"id"`               // 编号
	TopicId         int    `json:"topic_id"`         // 回复主题编号
	MemberId        int    `json:"member_id"`        // 回复人编号
	Member          Member `json:"member"`           // 回复人详情
	Created         int    `json:"created"`          // 回复创建时间戳
	Content         string `json:"content"`          // 回复内容
	ContentRendered string `json:"content_rendered"` // html 渲染后内容
	LastModified    int    `json:"last_modified"`    // TODO 作用未知
}
