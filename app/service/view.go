package service

import (
	"GoViewFile/library/utils"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
	"time"

	"github.com/gogf/gf/util/gconv"
)

type NowFile struct {
	Md5            string
	Ext            string
	LastActiveTime int64
}

var (
	Pattern      string
	Address      string
	AllFile      map[string]*NowFile
	ExpireTime   int64
	AllOfficeEtx = []string{".doc", ".docx", ".xls", ".xlsx", ".ppt", ".pptx", ".txt"}
	AllImageEtx  = []string{".jpg", ".png", ".gif"}
)

func OfficePage(imgPath string) []byte {
	rd, _ := ioutil.ReadDir(imgPath)
	dataByte, _ := ioutil.ReadFile("public/html/office.html")
	dataStr := string(dataByte)
	htmlCode := ""
	for _, fi := range rd {
		if !fi.IsDir() {
			htmlCode = htmlCode + `<img class="my-photo" alt="loading" title="查看大图" style="cursor: pointer;"
									data-src="/view/office?url=` + path.Base(imgPath) + "/" + fi.Name() + `" src="images/loading.gif"
									">`
		}
	}
	dataStr = strings.Replace(dataStr, "{{AllImages}}", htmlCode, -1)
	dataByte = []byte(dataStr)
	return dataByte
}

func ImagePage(filePath string) []byte {
	dataByte, _ := ioutil.ReadFile("public/html/image.html")
	dataStr := string(dataByte)
	imageUrl := "/view/img?url=" + path.Base(filePath)
	htmlCode := `<li>
					<img id="` + imageUrl + `" url="` + imageUrl + `"
						src="` + imageUrl + `" width="1px" height="1px">
				 </li>`
	dataStr = strings.Replace(dataStr, "{{AllImages}}", htmlCode, -1)
	dataStr = strings.Replace(dataStr, "{{FirstPath}}", imageUrl, -1)
	dataByte = []byte(dataStr)
	return dataByte
}

func PdfPage(filePath string) []byte {
	dataByte, _ := ioutil.ReadFile("public/html/pdf.html")
	dataStr := string(dataByte)

	pdfUrl := "/view/pdf?url=" + path.Base(filePath)
	dataStr = strings.Replace(dataStr, "{{url}}", pdfUrl, -1)

	dataByte = []byte(dataStr)
	return dataByte
}

func PdfPageDownload(filePath string) []byte {
	dataByte, _ := ioutil.ReadFile("public/html/pdf.html")
	dataStr := string(dataByte)
	pdfUrl := "/view/img?url=" + path.Base(filePath)
	dataStr = strings.Replace(dataStr, "{{url}}", pdfUrl, -1)
	dataByte = []byte(dataStr)
	return dataByte
}

func IsHave(fileName string) bool {
	fileName = strings.Split(fileName, ".")[0]
	if _, ok := AllFile[fileName]; ok {
		AllFile[fileName].LastActiveTime = time.Now().Unix()
		return true
	} else {
		return false
	}
}

func SetFileMap(fileName string) {
	ext := path.Ext(fileName)
	fileName = strings.Split(fileName, ".")[0]
	if _, ok := AllFile[fileName]; ok {
		AllFile[fileName].LastActiveTime = time.Now().Unix()
		return
	} else {
		temp := &NowFile{
			Md5:            fileName,
			Ext:            ext,
			LastActiveTime: time.Now().Unix(),
		}
		AllFile[fileName] = temp
	}
}

func Monitor() {
	log.Println("Info: Starting Monitor Thread")
	for {
		for _, v := range AllFile {
			if time.Now().Unix()-v.LastActiveTime > ExpireTime {
				if v.Md5 != "" {
					os.RemoveAll("cache/convert/" + v.Md5)
					os.Remove("cache/download/" + v.Md5 + v.Ext)
					os.Remove("cache/pdf/" + v.Md5 + ".pdf")
					log.Println("Cache file ", v.Md5, " delete")
					delete(AllFile, v.Md5)
				} else {
					delete(AllFile, v.Md5)
					log.Println("Cache file ", v.Md5, " delete with error")
				}
			}
		}
		time.Sleep(time.Second * 60)
	}
}

func GetAllFile(pathname string) ([]map[string]string, error) {
	s := []map[string]string{}
	rd, err := ioutil.ReadDir(pathname)
	if err != nil {
		return s, err
	}

	for _, fi := range rd {
		tmp := map[string]string{}
		if !fi.IsDir() {
			fullName := pathname + "/" + fi.Name()
			tmp["path"] = fullName
			tmp["name"] = fi.Name()
			tmp["type"] = path.Ext(fullName)
		}
		s = append(s, tmp)
	}
	return s, nil
}

//将Excel转html
func ExcelPage(filePath string) []byte {
	ret := utils.ExcelParse(filePath)
	html := `
			<!DOCTYPE HTML PUBLIC "-//W3C//DTD HTML 4.0 Transitional//EN"><html><head>
			<meta http-equiv="content-type" content="text/html; charset=utf-8"/>
			<title></title>	<style type="text/css">body,div,table,thead,tbody,tfoot,tr,th,td,p { font-family:"等线"; font-size:x-small }
			a.comment-indicator:hover + comment { background:#ffd; position:absolute; display:block; border:1px solid black; padding:0.5em;  }
			a.comment-indicator { background:red; display:inline-block; border:1px solid black; width:0.5em; height:0.5em;  }
			comment { display:none;  } 	</style>
	`

	//html := ""
	html += "<p><center>		<h1>Overview</h1>"
	for i := 0; i < len(ret); i++ {
		html += "<A HREF=\"#table" + gconv.String(i) + "\" style = \"font-size: 30px;\"  >Sheet" + gconv.String(i+1) + " </A><br>"
	}
	html += "</center></p><hr>"

	for k, v := range ret {
		html += "<A NAME=\"table" + gconv.String(k) + "\"   style = \"color: #337ab7;\">"
		html += "<h1>Sheet" + gconv.String(k+1) + "</h1></A>"
		html += `
			<table  class = "table table-striped" cellspacing ="0" border ="0"  style= "width: 100%;max-width: 100%;"> 
		`
		for _, vs := range gconv.SliceAny(v["title"]) {
			num := len(gconv.String(vs)) * 10
			html += "<colgroup width=\"" + gconv.String(num) + "\"></colgroup>  "
		}

		html += "<tr>"
		for _, vs := range gconv.SliceAny(v["title"]) {
			html += "<td height=\"19\" align=\"left\" valign=bottom><font color=\"#000000\">" + gconv.String(vs) + "</font></td>	 "
		}
		html += "</tr>"

		for _, vs := range gconv.SliceAny(v["resourceArr"]) {
			html += "<tr>"
			for _, vss := range gconv.SliceAny(vs) {
				html += "<td height=\"19\" align=\"left\" valign=bottom><font color=\"#000000\">" + gconv.String(vss) + "</font></td>	 "
			}
			html += "</tr>"
		}
		html += "</table>"

	}
	html += `
		</html>
		<script src="/html/js/jquery-3.0.0.min.js" type="text/javascript">
		</script><script src="/html/js/excel.header.js" type="text/javascript">
		</script><link rel="stylesheet" href="/html/css/bootstrap.min.css">
		`
	// tmpFilePath := "public/excel/" + path.Base(filePath) + ".html"
	// file, err := os.Create(tmpFilePath)
	// if err != nil {
	// 	logger.Error(err.Error())
	// 	return nil
	// }
	// defer file.Close()
	// buf := []byte(html)
	// num := len(buf)
	// _, er := file.Write(buf[0:num])
	// if er != nil {
	// 	logger.Error(er)
	// }

	dataByte, _ := ioutil.ReadFile("public/html/excel.html")
	dataStr := string(dataByte)

	dataStr = strings.Replace(dataStr, "{{url}}", html, -1)
	dataByte = []byte(dataStr)
	return dataByte
}
