package player

type Article struct {
	// 文章id
	Id string `json:"id,required"`
	// 文章标签
	Tag string `json:"tag,required"`
	// 文章标题
	Title string `json:"title,required"`
	// 文章内容
	Content string `json:"content,required"`
	// 文章作者
	Author string `json:"author,required"`
	// 阅读量
	ViewCount int64 `json:"viewCount,required"`
	// 点赞数
	StarCount int64 `json:"starCount,required"`
	// 发布时间
	PublishTime string `json:"publishTime,required"`
}
