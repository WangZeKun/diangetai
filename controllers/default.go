package controllers

import (
	"github.com/astaxie/beego"
	"diangetai/models"
	"encoding/json"
	"strconv"
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
	s := models.Schedule{
		SendSong:   c.GetString("sendSong"),
		Song:       c.GetString("Song"),
		Num:        num,
		ArriveSong: c.GetString("arriveSong"),
		Something:  c.GetString("Something"),
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
	_, h, _ := c.GetFile("picture")                  //获取上传的文件
	path := "static/img/" + h.Filename						//文件目录
	c.SaveToFile("picture", path)                    //存文件

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
	c.SaveToFile("file", "static/music/"+strconv.Itoa(num)+".zip") //存文件
	beego.AppConfig.Set("num", strconv.Itoa(num+1))
	c.ServeJSON()
}
