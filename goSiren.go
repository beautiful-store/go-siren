package goSiren

import (
	"encoding/json"
	"fmt"
	"github.com/labstack/echo"
	"github.com/mssola/user_agent"
	"github.com/sirupsen/logrus"
	"golang.org/x/exp/slices"
	"net/http"
)

const EnvironmentStaging = "staging"
const EnvironmentProduction = "production"
const EnvironmentLocal = "local"
const EnvironmentDev = "dev"

type GoSiren struct {
	TargetCodes        []int
	TargetEnvironments []string
	TargetHost         []string
}
type Content interface {
	make() (string, string, error)
}

func (g *GoSiren) IsSending(code int, environment string, host string) bool {
	return slices.Contains(g.TargetEnvironments, environment) &&
		slices.Contains(g.TargetCodes, code) &&
		slices.Contains(g.TargetHost, host)
}

func (g *GoSiren) Setting(content Content) (string, string, error) {
	return content.make()
}

type Alarm struct {
	Request    *http.Request
	Code       int
	StackTrace string
	UserID     int
	UserName   string
}

func (g *Alarm) make() (string, string, error) {
	token := g.Request.Header.Get(echo.HeaderAuthorization)

	title := g.StackTrace
	if len(g.StackTrace) > 40 {
		title = g.StackTrace[:40]
	}

	content := fmt.Sprintf("# 로그인 사용자\n"+
		"* memberId:<span style=\"color:#9933ff\"> %v</span>\n"+
		"* Name:<span style=\"color:#9933ff\"> %v</span>\n\n"+
		"# HTTP Request\n"+
		"* HOST :<span style=\"color:#9933ff\"> %v</span>\n"+
		"* URL :<span style=\"color:#9933ff\"> %v</span>\n"+
		"* Method :<span style=\"color:#9933ff\"> %v</span>\n* Header\n    "+
		"* Content-Type:<span style=\"color:#9933ff\"> %v</span>\n    "+
		"* Authorization\n```\n%v\n```\n"+
		"* BrowserUserAgent\n```\n%v\n```\n"+
		"* Body\n```\n%v\n```\n\n"+
		"# HTTP Response\n"+
		"* status code : <span style=\"color:#9933ff\">%v</span>\n"+
		"# Error 원인\n```\n%v\n```",
		g.UserID,
		g.UserName,
		g.Request.Host,
		g.Request.RequestURI,
		g.Request.Method,
		g.Request.Header.Get(echo.HeaderContentType),
		token,
		getHumanizeBrowserUserAgent(g.Request.UserAgent()),
		g.Request.Body,
		g.Code,
		fmt.Sprintf("[ERROR] %v \n", g.StackTrace))

	return title, content, nil
}

func getHumanizeBrowserUserAgent(browserUserAgent string) string {
	if len(browserUserAgent) == 0 {
		return ""
	}
	ua := user_agent.New(browserUserAgent)

	humanizeUserAgent := map[string]interface{}{}
	humanizeUserAgent["mobile"] = ua.Mobile()
	humanizeUserAgent["platform"] = ua.Platform()
	humanizeUserAgent["os"] = ua.OS()

	name, version := ua.Engine()
	engine := map[string]interface{}{}
	engine["name"] = name
	engine["version"] = version
	humanizeUserAgent["engine"] = engine

	name, version = ua.Browser()
	browser := map[string]interface{}{}
	browser["name"] = name
	browser["version"] = version
	humanizeUserAgent["browser"] = browser

	jsonString, err := json.Marshal(humanizeUserAgent)
	if err != nil {
		logrus.Info(err)
		return ""
	}

	return string(jsonString)
}
