package model

type ViewReq struct {
	Url       string `json:"Url"`
	Type      string `json:"Type"`      //判断是图片展示，还是pdf展示
	FileWay   string `json:"FileWay"`   // 判断是否是本地文件
	WaterMark string `json:"watermark"` // 水印
}
