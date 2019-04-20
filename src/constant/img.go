package constant

const (

	// 图片做缩略处理：w: 160, h: 160
	ImgOps       = "imageView2/2/w/160/h/160"
	ImgURIPrefix = "image host"
	ImgMicroSize = 160

	ImgSuffix = ".jpg"

	ImgFeedbackImgType        = 1
	ImgPrefixFeedbackImg      = "app/feedback/"
	ImgPrefixMicroFeedbackImg = "app/feedback/micro/"

	TokenQiniuExpire = 7200
)

var (
	ImgPrefix = map[int]string{
		ImgFeedbackImgType: ImgPrefixFeedbackImg,
	}
	ImgPrefixMicro = map[int]string{
		ImgFeedbackImgType: ImgPrefixMicroFeedbackImg,
	}
)
