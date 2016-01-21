package http

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/open-falcon/alarm/g"
    "github.com/open-falcon/alarm/api"
	"github.com/toolkits/file"
	"sort"
	"strings"
    "log"
	"time"
)

type MainController struct {
	beego.Controller
}

func (this *MainController) Version() {
	this.Ctx.WriteString(g.VERSION)
}

func (this *MainController) Health() {
	this.Ctx.WriteString("ok")
}

func (this *MainController) Workdir() {
	this.Ctx.WriteString(fmt.Sprintf("%s", file.SelfDir()))
}

func (this *MainController) ConfigReload() {
	remoteAddr := this.Ctx.Input.Request.RemoteAddr
	if strings.HasPrefix(remoteAddr, "127.0.0.1") {
		g.ParseConfig(g.ConfigFile)
		this.Data["json"] = g.Config()
		this.ServeJson()
	} else {
		this.Ctx.WriteString("no privilege")
	}
}

func (this *MainController) Index() {
	events := g.Events.Clone()

	defer func() {
		this.Data["Now"] = time.Now().Unix()
		this.TplNames = "index.html"
	}()

    username := getLoginUser( this )
    this.Data["username"] = username

	if len(events) == 0 {
		this.Data["Events"] = []*g.EventDto{}
		return
	}

	count := len(events)
	if count == 0 {
		this.Data["Events"] = []*g.EventDto{}
		return
	}

	// 按照持续时间排序
	beforeOrder := make([]*g.EventDto, count)
	i := 0

    //筛选event，只有属于同一用户team的才可被展示
	for _, event := range events {
        if checkEventBelongUser( event, username ) {
            beforeOrder[i] = event
		    i++
        }
	}

	sort.Sort(g.OrderedEvents(beforeOrder))
	this.Data["Events"] = beforeOrder
}

func (this *MainController) Solve() {
	ids := this.GetString("ids")
	if ids == "" {
		this.Ctx.WriteString("")
		return
	}

	idArr := strings.Split(ids, ",,")
	for i := 0; i < len(idArr); i++ {
		g.Events.Delete(idArr[i])
	}

	this.Ctx.WriteString("")
}

func getLoginUser( this *MainController ) string {
    sig := this.Ctx.GetCookie("sig")
    if strings.TrimSpace( sig ) == "" {
        redirectToSso( this )
    }
    
    username := api.UsernameFromSso( sig )
    if username == "" {
        redirectToSso( this )
    }

    return username
}

func redirectToSso( this *MainController ) {
    sig,err := api.GenSig()
    if err != nil {
        log.Println("get sig from uic fail", err)
        return
    }
    this.Ctx.SetCookie("sig",sig)
    loginurl := api.LoginUrl(sig,this.Ctx.Input.Scheme()+"://"+this.Ctx.Input.Request.Host+this.Ctx.Input.Request.RequestURI)
    this.Ctx.Redirect(302,loginurl)
}

func checkEventBelongUser( e *g.EventDto, user string) bool {
    //get event action id
    actionId := e.ActionId

    //获取event对应的uic team
    uicTeam := api.GetAction(actionId).Uic

    //当前登录user是否为此team成员
    
    return true
}
