package utils

import "github.com/mileusna/useragent"

// ClientAgent 表示从 User-Agent 中提取出的核心客户端信息。
type ClientAgent struct {
	Device  string
	OS      string
	Browser string
}

// ParseClientAgent 解析 User-Agent，返回设备、系统与浏览器信息。
func ParseClientAgent(uaString string) ClientAgent {
	ua := useragent.Parse(uaString)

	os := ua.OS
	if os == "" {
		os = "Unknown"
	}

	browser := ua.Name
	if browser == "" {
		browser = "Unknown"
	}

	device := "PC"
	switch {
	case ua.Bot:
		device = "Bot"
	case ua.Tablet:
		device = "Tablet"
	case ua.Mobile:
		device = "Mobile"
	case ua.Device != "":
		device = ua.Device
	}

	return ClientAgent{
		Device:  device,
		OS:      os,
		Browser: browser,
	}
}
