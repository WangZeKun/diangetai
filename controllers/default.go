package controllers

import (
	"github.com/astaxie/beego"
	"diangetai/models"
	"encoding/json"
	"strconv"
	"net/http"
	"os"
	"io"
	"chain"
	"archive/zip"
	"bytes"
	"io/ioutil"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {
	data, err := models.EndRead()
	if err != nil {
		beego.Error(err)
		c.Abort("500")
	}
	c.Data["data"] = data
	c.Data["Email"] = "1015190212@qq.com"
	c.TplName = "index.html"
}

func (c *MainController) Post() {
	num, err := beego.AppConfig.Int("num")
	if err != nil {
		beego.Error(err)
		c.Abort("500")
	}
	var path string
	if c.GetString("type") == "自行上传" {
		_, h, err := c.GetFile("music") //获取上传的文件
		if err != nil {
			beego.Error(err)
			c.Abort("401")
		}
		path = "static/ready/" + strconv.Itoa(num) + h.Filename //文件目录
		c.SaveToFile("music", path)                             //存文件
	} else {
		path, err = chain.GetURL(c.GetString("url"))
		if err != nil {
			beego.Error(err)
			c.Abort("401")
		}
	}
	s := models.Schedule{
		SendSong:   c.GetString("sendSong"),
		Song:       c.GetString("Song"),
		SongType:   c.GetString("type"),
		Num:        num,
		ArriveSong: c.GetString("arriveSong"),
		Something:  c.GetString("Something"),
		Url:        path,
	}
	err = s.Insert()
	if err != nil {
		beego.Error(err)
		c.Abort("500")
	}
	c.ServeJSON()
}

type AfterController struct {
	beego.Controller
}

func (c *AfterController) Get() {
	beego.Informational(beego.AppConfig.String("num"))
	num, err := beego.AppConfig.Int("num")
	if err != nil {
		beego.Error(err)
		c.Abort("500")
	}
	data, err := models.NowRead(num)
	if err != nil {
		beego.Error(err)
		c.Abort("500")
	}
	j, _ := json.Marshal(data)
	c.Data["data"] = string(j)
	c.TplName = "after.html"
}

func (c *AfterController) End() {
	num, err := beego.AppConfig.Int("num")
	if err != nil {
		beego.Error(err)
		c.Abort("500")
	}
	beego.Informational(c.GetString("select"))
	var s []models.Schedule
	err = json.Unmarshal([]byte(c.GetString("select")), &s)
	if err != nil {
		beego.Error(err)
		c.Abort("500")
	}
	err = compress(s, "static/music/"+strconv.Itoa(num)+".zip") //存文件
	if err != nil {
		beego.Error(err)
		c.Abort("500")
	}
	_, h, err := c.GetFile("picture") //获取上传的文件
	var path string
	if err != nil {
		path = "static/img/" + strconv.Itoa(num) + ".jpg"
		res, err := http.Get("https://api.dujin.org/bing/1920.php")
		defer res.Body.Close()
		if err != nil {
			beego.Error(err)
			c.Abort("500")
		}
		f, err := os.Create("static/img/" + strconv.Itoa(num) + ".jpg")
		defer f.Close()
		if err != nil {
			beego.Error(err)
			c.Abort("500")
		}
		io.Copy(f, res.Body)
	} else {
		path = "static/img/" + strconv.Itoa(num) + h.Filename //文件目录
		c.SaveToFile("picture", path)                         //存文件
	}
	p := models.Periods{
		Num:     num,
		Picture: path,
		IsEnd:   true,
	}
	err = p.End()
	if err != nil {
		beego.Error(err)
		c.Abort("500")
	}
	err = models.Delete(num, s)
	beego.AppConfig.Set("num", strconv.Itoa(num+1))
	c.ServeJSON()
}

func compress(in []models.Schedule, dest string) (err error) {
	buf := new(bytes.Buffer)
	w := zip.NewWriter(buf)
	for _, s := range in {
		if s.SongType == "自行上传" {
			file, err := os.Open(s.Url)
			if err != nil {
				return err
			}
			buf1, err := ioutil.ReadAll(file)
			if err != nil {
				return err
			}
			f, err := w.Create(s.Song + ".mp3")
			if err != nil {
				return err
			}
			_, err = f.Write(buf1)
			if err != nil {
				return err
			}
		} else {
			res, err := http.Get(s.Url)
			if err != nil {
				return err
			}
			f, err := w.Create(s.Song + ".mp3")
			if err != nil {
				return err
			}
			buf1, err := ioutil.ReadAll(res.Body)
			if err != nil {
				return err
			}
			_, err = f.Write(buf1)
			if err != nil {
				return err
			}
		}
	}
	err = w.Close()
	if err != nil {
		return err
	}
	os.Create(dest)
	f, err := os.OpenFile(dest, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return
	}
	buf.WriteTo(f)
	return
}
