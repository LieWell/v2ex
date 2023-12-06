package v2ex

type Node struct {
	Root             bool          `json:"root"`              // 是否是跟节点
	ParentNodeName   string        `json:"parent_node_name"`  // 父节点名称
	Id               int           `json:"id"`                // 节点编号
	Name             string        `json:"name"`              // 节点名称,是一个英文,用来过滤、查询、跳转等
	Title            string        `json:"title"`             // 节点展示名称,可以是中文
	TitleAlternative string        `json:"title_alternative"` // TODO 含义未知
	Url              string        `json:"url"`               // 节点访问地址
	Topics           int           `json:"topics"`            // 节点主体总数
	Header           string        `json:"header"`            // 节点介绍
	Footer           string        `json:"footer"`            // TODO 含义未知
	AvatarMini       string        `json:"avatar_mini"`       // 小尺寸头像
	AvatarNormal     string        `json:"avatar_normal"`     // 正常尺寸头像
	AvatarLarge      string        `json:"avatar_large"`      // 大尺寸头像
	Stars            int           `json:"stars"`             // 收藏人数
	Aliases          []interface{} `json:"aliases"`           // TODO 含义未知
}
