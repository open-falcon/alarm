package api

import (
	"fmt"
	"github.com/open-falcon/alarm/g"
	"github.com/toolkits/container/set"
	"github.com/toolkits/net/httplib"
	"log"
	"strings"
	"sync"
	"time"
)

type User struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

type UsersWrap struct {
	Msg   string  `json:"msg"`
	Users []*User `json:"users"`
}

type UserMsg struct {
    User *User  `json:user`
}

type UsersCache struct {
	sync.RWMutex
	M map[string][]*User
}

var Users = &UsersCache{M: make(map[string][]*User)}

func (this *UsersCache) Get(team string) []*User {
	this.RLock()
	defer this.RUnlock()
	val, exists := this.M[team]
	if !exists {
		return nil
	}

	return val
}

func (this *UsersCache) Set(team string, users []*User) {
	this.Lock()
	defer this.Unlock()
	this.M[team] = users
}

func UsersOf(team string) []*User {
	users := CurlUic(team)

	if users != nil {
		Users.Set(team, users)
	} else {
		users = Users.Get(team)
	}

	return users
}

func GetUsers(teams string) map[string]*User {
	userMap := make(map[string]*User)
	arr := strings.Split(teams, ",")
	for _, team := range arr {
		if team == "" {
			continue
		}

		users := UsersOf(team)
		if users == nil {
			continue
		}

		for _, user := range users {
			userMap[user.Name] = user
		}
	}
	return userMap
}

// return phones, emails
func ParseTeams(teams string) ([]string, []string) {
	if teams == "" {
		return []string{}, []string{}
	}

	userMap := GetUsers(teams)
	phoneSet := set.NewStringSet()
	mailSet := set.NewStringSet()
	for _, user := range userMap {
		phoneSet.Add(user.Phone)
		mailSet.Add(user.Email)
	}
	return phoneSet.ToSlice(), mailSet.ToSlice()
}

func CurlUic(team string) []*User {
	if team == "" {
		return []*User{}
	}

	uri := fmt.Sprintf("%s/team/users", g.Config().Api.Uic)
	req := httplib.Get(uri).SetTimeout(2*time.Second, 10*time.Second)
	req.Param("name", team)
	req.Param("token", g.Config().UicToken)

	var usersWrap UsersWrap
	err := req.ToJson(&usersWrap)
	if err != nil {
		log.Printf("curl %s fail: %v", uri, err)
		return nil
	}

	if usersWrap.Msg != "" {
		log.Printf("curl %s return msg: %v", uri, usersWrap.Msg)
		return nil
	}

	return usersWrap.Users
}

func GenSig() (string, error) {
    uri := fmt.Sprintf("%s/sso/sig",g.Config().Api.Uic)   
    req := httplib.Get(uri).SetTimeout(2*time.Second, 10*time.Second)
    sig,err := req.String()
    return sig,err   
}

func LoginUrl(sig string, callback string) string {
    return fmt.Sprintf("%s/auth/login?sig=%s&callback=%s",g.Config().Api.Uic,sig,callback)
} 

func UsernameFromSso( sig string ) string {
    uri := fmt.Sprintf("%s/sso/user/%s?token=%s",g.Config().Api.Uic,sig,g.Config().UicToken)
    req := httplib.Get(uri).SetTimeout(2*time.Second, 10*time.Second)
    var userMsg UserMsg
    err := req.ToJson( &userMsg )
    if err != nil {
        log.Printf("curl %s fail: %v",uri,err)
        return ""
    }   

    if strings.TrimSpace(userMsg.User.Name) == "" {
        log.Printf("curl %s return none user!", uri)
        return ""
    }

    return userMsg.User.Name
}

func CheckUserInTeam( username string, team string ) bool {
    uri := fmt.Sprintf("%s/user/in",g.Config().Api.Uic)
    req := httplib.Get(uri).SetTimeout(2*time.Second, 10*time.Second)
    req.Param("name",username)
    req.Param("teams",team)
    checkRes,err := req.String()

    if err != nil {
		log.Printf("curl %s fail: %v", uri, err)
		return false
    }

    if strings.TrimSpace(checkRes) == "1" {
        return true
    }

    return false
}
