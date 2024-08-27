package entities

import (
	"net"
	"net/url"
	"time"
)

type App struct {
	Version   string
	GoVersion string
	BuildTime time.Time
	Git       Git
	Server    Server
}

type Git struct {
	Branch     string
	Commit     string
	Repository *url.URL
}

type Server struct {
	OS        string
	Arch      string
	Hostname  string
	URL       *url.URL
	IP        net.IP
	Port      int
	Interface string
	Uptime    time.Duration
}
