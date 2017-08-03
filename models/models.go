package models

import (
	"chain"
	"strconv"

	"github.com/astaxie/beego"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

type Schedule struct {
	Num        int
	Asc        bool
	SongType   string
	SendSong   string
	ArriveSong string
	Something  string
	Song       chain.Music `xorm:"json"	`
}

type Periods struct {
	Num     int    `xorm:"pk autoincr"`
	Picture string `xorm:"LONGTEXT"`
	IsEnd   bool
}

type SP struct {
	S []Schedule
	P Periods
}

var engine *xorm.Engine

func init() {
	engine, _ = xorm.NewEngine("mysql", "Vod:AKZzP66CmZ@tcp(47.94.91.118:3306)/Vod?charset=utf8")
	engine.Sync2(new(Schedule), new(Periods))
	var p Periods
	bool2, _ := engine.Where("is_end=false").Get(&p)
	if !bool2 {
		Periods{
			Num:1,
			IsEnd: false,
		}.Insert()
		beego.AppConfig.Set("num", "1")
	} else {
		beego.AppConfig.Set("num", strconv.Itoa(p.Num))
	}
}

func (s Schedule) Insert() (err error) {
	_, err = engine.Insert(s)
	return
}

func EndRead() (s map[int]SP, err error) {
	var i []int
	s = make(map[int]SP)
	err = engine.Table("periods").Where("is_end=true").Cols("num").Find(&i)
	for _, a := range i {
		var j []Schedule
		var h Periods
		engine.Where("`asc`<>0").And("num=?", a).Find(&j)
		engine.Where("num=?", a).Get(&h)
		s[a] = SP{
			S: j,
			P: h,
		}
	}
	return
}

func NowRead(num int) (s []Schedule, err error) {
	err = engine.Where("num=?", num).Find(&s)
	return
}

func Delete(num int, s []Schedule) (err error) {
	for _, i := range s {
		i.Asc = true
		_,err = engine.Where("arrive_song=? and send_song=? and num=?", i.ArriveSong, i.SendSong, num).Cols("asc").Update(i)
	}
	return
}

func (p Periods) Insert() (err error) {
	_, err = engine.Insert(&p)
	return
}

func (p Periods) End() (err error) {
	_, err = engine.ID(p.Num).Cols("is_end", "picture").Update(&p)
	if err != nil {
		return
	}
	p1 := Periods{
		Num:   p.Num + 1,
		IsEnd: false,
	}
	err = p1.Insert()
	return
}
