package api

import (
	"GoViewFile/app/model"
	"GoViewFile/app/service"
	"GoViewFile/library/logger"
	"GoViewFile/library/response"
	"GoViewFile/library/utils"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/gogf/gf/net/ghttp"
)

var View = new(ViewApi)

//本地文件路径
var filePath string

type ViewApi struct{}

// @summary 预览文件入口
// @tags    预览
// @produce json
// @param   entity "
// @router  /view/view [POST]
// @success 200 {object} response.JsonResponse "执行结果"
func (a *ViewApi) View(r *ghttp.Request) {
	var (
		reqData *model.ViewReq
	)
	//解析参数
	if err := r.Parse(&reqData); err != nil {
		logger.Errorf("View ->   execution failed. err", err.Error())
		response.JsonExit(r, 1, "参数解析错误")

	}

	if reqData.FileWay == "local" { //本地文件预览
		filePath = reqData.Url
	} else {
		//下载文件
		file, err := service.DownloadFile(reqData.Url, "cache/download/"+path.Base(reqData.Url))
		if err != nil {
			logger.Print(err.Error())
			response.JsonExit(r, -1, "文件下载失败")
		}
		filePath = file
	}
	fileType := strings.ToLower(path.Ext(filePath))

	//MD文件预览
	if fileType == ".md" {
		dataByte := service.MdPage(filePath)
		r.Response.Writer.Header().Set("content-length", strconv.Itoa(len(dataByte)))
		r.Response.Writer.Header().Set("content-type:", "text/html;charset=UTF-8")
		r.Response.Writer.Write([]byte(dataByte))
		return
	}

	//MD文件预览
	if fileType == ".msg" || fileType == ".eml" {
		pdfPath := utils.MsgToPdf(filePath)
		if pdfPath == "" {
			response.JsonExit(r, -1, "转pdf失败")
		}
		dataByte := service.PdfPage("cache/pdf/" + path.Base(pdfPath))
		r.Response.Writer.Header().Set("content-length", strconv.Itoa(len(dataByte)))
		r.Response.Writer.Header().Set("content-type:", "text/html;charset=UTF-8")
		r.Response.Writer.Write([]byte(dataByte))
		return
	}

	//后缀是pdf直接读取文件类容返回
	if fileType == ".pdf" {
		dataByte := service.PdfPageDownload(filePath)
		r.Response.Writer.Header().Set("content-length", strconv.Itoa(len(dataByte)))
		r.Response.Writer.Header().Set("content-type:", "text/html;charset=UTF-8")
		r.Response.Writer.Write([]byte(dataByte))
		return
	}
	//后缀png , jpg ,gif
	if utils.IsInArr(fileType, service.AllImageEtx) {
		dataByte := service.ImagePage(filePath)
		r.Response.Writer.Header().Set("content-length", strconv.Itoa(len(dataByte)))
		r.Response.Writer.Header().Set("content-type:", "text/html;charset=UTF-8")
		r.Response.Writer.Write([]byte(dataByte))
		return
	}

	// 后缀xlsx
	if (fileType == ".xlsx" || fileType == ".xls") && reqData.Type != "pdf" {
		dataByte := service.ExcelPage(filePath)
		r.Response.Writer.Header().Set("content-length", strconv.Itoa(len(dataByte)))
		r.Response.Writer.Header().Set("content-type:", "text/html;charset=UTF-8")
		r.Response.Writer.Write([]byte(dataByte))
		return
	}

	// 除了PDF外的其他word文件  (如果没有安装ImageMagick，可以将这个分支去掉)
	if utils.IsInArr(fileType, service.AllOfficeEtx) && reqData.Type != "pdf" {
		pdfPath := utils.ConvertToPDF(filePath)
		if pdfPath == "" {
			response.JsonExit(r, -1, "转pdf失败")
		}
		imgPath := utils.ConvertToImg(pdfPath)
		if imgPath == "" {
			response.JsonExit(r, -1, "转图片失败")
		}
		dataByte := service.OfficePage("cache/convert/" + path.Base(imgPath))
		r.Response.Writer.Header().Set("content-length", strconv.Itoa(len(dataByte)))
		r.Response.Writer.Header().Set("content-type:", "text/html;charset=UTF-8")
		r.Response.Writer.Write([]byte(dataByte))
		return
	}

	// 除了PDF外的其他word文件
	if utils.IsInArr(fileType, service.AllOfficeEtx) {
		pdfPath := utils.ConvertToPDF(filePath)
		if pdfPath == "" {
			response.JsonExit(r, -1, "转pdf失败")
		}

		dataByte := service.PdfPage("cache/pdf/" + path.Base(pdfPath))
		r.Response.Writer.Header().Set("content-length", strconv.Itoa(len(dataByte)))
		r.Response.Writer.Header().Set("content-type:", "text/html;charset=UTF-8")
		r.Response.Writer.Write([]byte(dataByte))
		return
	}

	response.JsonExit(r, 0, "ok", "暂不支持该类型文件预览！")

}

// @summary 返回文件类容-img
// @tags    预览
// @produce json
// @param   entity "
// @router  /view/view [POST]
// @success 200 {object} response.JsonResponse "执行结果"
func (a *ViewApi) Img(r *ghttp.Request) {
	var (
		reqData *model.ViewReq
	)
	//解析参数
	if err := r.Parse(&reqData); err != nil {
		logger.Errorf("View ->   execution failed. err", err.Error())
		response.JsonExit(r, 1, "参数解析错误")

	}
	imgPath := reqData.Url
	DataByte, err := ioutil.ReadFile("cache/download/" + imgPath)
	if err != nil { //如果是本地预览，则文件在local木下下
		DataByte, err = ioutil.ReadFile("cache/local/" + imgPath)
		if err != nil {
			r.Response.Writer.Header().Set("content-length", strconv.Itoa(len("404")))
			r.Response.Writer.Header().Set("content-type:", "text/html;charset=UTF-8")
			r.Response.Writer.Write([]byte("出现了一些问题,导致File View无法获取您的数据!"))
			return
		}
	}
	r.Response.Writer.Header().Set("content-length", strconv.Itoa(len(DataByte)))
	r.Response.Writer.Header().Set("content-type:", "text/html;charset=UTF-8")
	r.Response.Writer.Write(DataByte)
}

// @summary 返回文件类容-（转换后的pdf）
// @tags    预览
// @produce json
// @param   entity "
// @router  /view/view [POST]
// @success 200 {object} response.JsonResponse "执行结果"
func (a *ViewApi) Pdf(r *ghttp.Request) {
	var (
		reqData *model.ViewReq
	)
	//解析参数
	if err := r.Parse(&reqData); err != nil {
		logger.Errorf("View ->   execution failed. err", err.Error())
		response.JsonExit(r, 1, "参数解析错误")

	}
	imgPath := reqData.Url
	DataByte, err := ioutil.ReadFile("cache/pdf/" + imgPath)
	if err != nil {
		r.Response.Writer.Header().Set("content-length", strconv.Itoa(len("404")))
		r.Response.Writer.Header().Set("content-type:", "text/html;charset=UTF-8")
		r.Response.Writer.Write([]byte("出现了一些问题,导致File View无法获取您的数据!"))
		return
	}
	r.Response.Writer.Header().Set("content-length", strconv.Itoa(len(DataByte)))
	r.Response.Writer.Write(DataByte)
}

// @summary 返回文件类容-（转换后的图片）
// @tags    预览
// @produce json
// @param   entity "
// @router  /view/view [POST]
// @success 200 {object} response.JsonResponse "执行结果"
func (a *ViewApi) Office(r *ghttp.Request) {
	var (
		reqData *model.ViewReq
	)
	//解析参数
	if err := r.Parse(&reqData); err != nil {
		logger.Errorf("View ->   execution failed. err", err.Error())
		response.JsonExit(r, 1, "参数解析错误")

	}
	imgPath := reqData.Url
	DataByte, err := ioutil.ReadFile("cache/convert/" + imgPath)
	if err != nil {
		r.Response.Writer.Header().Set("content-length", strconv.Itoa(len("404")))
		r.Response.Writer.Header().Set("content-type:", "text/html;charset=UTF-8")
		r.Response.Writer.Write([]byte("出现了一些问题,导致File View无法获取您的数据!"))
		return
	}
	r.Response.Writer.Header().Set("content-length", strconv.Itoa(len(DataByte)))
	r.Response.Writer.Header().Set("content-type:", "text/html;charset=UTF-8")
	r.Response.Writer.Write(DataByte)
}

// --------------------------------------------首页预览----------------------------------------

// @summary 上传文件（用于测试预览）
// @tags    预览
// @produce json
// @param   entity "
// @router  /view/Upload [POST]
// @success 200 {object} response.JsonResponse "执行结果"
func (a *ViewApi) Upload(r *ghttp.Request) {
	files := r.GetUploadFile("upload-file")
	_, _ = files.Save("cache/local/")

	allFile, _ := service.GetAllFile("cache/local/")
	logger.Print(allFile)
	view := r.GetView()
	view.Assign("AllFile", allFile)
	r.Response.WriteTpl("/index.html")
}

// @summary 删除本地上传的文件
// @tags    预览
// @produce json
// @param   entity "
// @router  /view/delete [POST]
// @success 200 {object} response.JsonResponse "执行结果"
func (a *ViewApi) Delete(r *ghttp.Request) {
	var (
		reqData *model.ViewReq
	)
	//解析参数
	if err := r.Parse(&reqData); err != nil {
		logger.Errorf("View ->   execution failed. err", err.Error())
		response.JsonExit(r, 1, "参数解析错误")

	}
	file := reqData.Url
	err := os.Remove(file)
	if err != nil {
		logger.Error(err.Error())
	}
	allFile, _ := service.GetAllFile("cache/local/")
	view := r.GetView()
	view.Assign("AllFile", allFile)
	r.Response.WriteTpl("/index.html")
}
