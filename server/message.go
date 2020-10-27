package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
)

// OnActivate register the plugin command
func (p *Plugin) MessageHasBeenPosted(c *plugin.Context, post *model.Post) {
	configuration := p.getConfiguration()
	if p.stringInSlice(post.ChannelId, strings.Split(configuration.ChannelIDList, ",")) && post.RootId == "" && post.Message[strings.LastIndex(post.Message, ",")+1:] != "admin." && post.Message[strings.LastIndex(post.Message, ",")+1:] != "channel." {
		user, err := p.API.GetUser(post.UserId)
		if err != nil {
			p.API.LogInfo("Unable to get User" + err.Error())
			return
		}
		channel, err3 := p.API.GetChannel(post.ChannelId)
		if err3 != nil {
			p.API.LogInfo("Unable to get Channel" + err3.Error())
			return
		}
		team, err4 := p.API.GetTeam(channel.TeamId)
		if err4 != nil {
			p.API.LogInfo("Unable to get Team" + err4.Error())
			return
		}
		serverConfig := p.API.GetConfig()
		permalink, err2 := url.Parse(*serverConfig.ServiceSettings.SiteURL)
		if err2 != nil {
			p.API.LogInfo("Unable to get url" + err2.Error())
			return
		}
		permalink.Path = path.Join(permalink.Path, team.Name, "pl", post.Id)
		loc, _ := time.LoadLocation("Asia/Kolkata")
		values := map[string]string{"text": "Username: " + user.Username + "\nMessage: " + post.Message + "\nLink: " + permalink.String() + "\nSent At: " + time.Unix(int64(post.CreateAt/1000), 0).In(loc).Format(time.ANSIC)}
		jsonValue, _ := json.Marshal(values)
		client := &http.Client{}
		req, _ := http.NewRequest("POST", configuration.HookURL, bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-type", "application/json")
		_, err1 := client.Do(req)
		if err1 != nil {
			p.API.LogInfo("Unable to send message" + err1.Error())
			return
		}
	}
}

func (p *Plugin) stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if strings.Trim(b, " ") == a {
			return true
		}
	}
	return false
}
