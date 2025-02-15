package info

import (
	"context"
	"net"
	"net/url"
	"os"
	"runtime"
	"time"

	"github.com/green-ecolution/green-ecolution-backend/internal/config"
	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/logger"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage"
)

var version = "development"
var gitCommit = "unknown"
var gitBranch = "develop"
var gitRepository = "https://github.com/green-ecolution/green-ecolution-management"
var buildTime = ""
var runTime = time.Now()

type InfoRepository struct {
	cfg            *config.Config
	localIP        net.IP
	localInterface string
	buildTime      time.Time
	gitRepository  *url.URL
}

func init() {
	if buildTime == "" || buildTime == "unknown" {
		buildTime = time.Now().Format("2006-01-02T15:04:05Z")
	}
}

func NewInfoRepository(cfg *config.Config) (*InfoRepository, error) {
	gitRepository, err := getGitRepository()
	if err != nil {
		return nil, err
	}

	buildTime, err := getBuildTime()
	if err != nil {
		return nil, err
	}

	localIP, err := getIP()
	if err != nil {
		return nil, err
	}

	localInterface, err := getInterface(localIP)
	if err != nil {
		return nil, err
	}

	return &InfoRepository{
		cfg:            cfg,
		localIP:        localIP,
		localInterface: localInterface,
		buildTime:      buildTime,
		gitRepository:  gitRepository,
	}, nil
}

func (r *InfoRepository) GetAppInfo(ctx context.Context) (*entities.App, error) {
	log := logger.GetLogger(ctx)
	hostname, err := r.getHostname()
	if err != nil {
		log.Error("failed to get hostname from host", "error", err)
		return nil, storage.ErrHostnameNotFound
	}

	appURL, err := r.getAppURL()
	if err != nil {
		log.Error("failed to parse configured app url", "error", err, "app_url", r.cfg.Server.AppURL)
		return nil, storage.ErrCannotGetAppURL
	}

	return &entities.App{
		Version:   version,
		GoVersion: r.getGoVersion(),
		BuildTime: r.buildTime,
		Git: entities.Git{
			Branch:     gitBranch,
			Commit:     gitCommit,
			Repository: r.gitRepository,
		},
		Server: entities.Server{
			OS:        r.getOS(),
			Arch:      r.getArch(),
			Hostname:  hostname,
			URL:       appURL,
			IP:        r.localIP,
			Port:      r.getPort(),
			Interface: r.localInterface,
			Uptime:    r.getUptime(),
		},
	}, nil
}

func (r *InfoRepository) getOS() string {
	return runtime.GOOS
}

func (r *InfoRepository) getHostname() (string, error) {
	return os.Hostname()
}

func (r *InfoRepository) getPort() int {
	return r.cfg.Server.Port
}

func (r *InfoRepository) getAppURL() (*url.URL, error) {
	return url.Parse(r.cfg.Server.AppURL)
}

func (r *InfoRepository) getUptime() time.Duration {
	return time.Since(runTime)
}

func (r *InfoRepository) getGoVersion() string {
	return runtime.Version()
}

func (r *InfoRepository) getArch() string {
	return runtime.GOARCH
}

func getBuildTime() (time.Time, error) {
	return time.Parse("2006-01-02T15:04:05Z", buildTime)
}

func getGitRepository() (*url.URL, error) {
	return url.Parse(gitRepository)
}

func getIP() (net.IP, error) {
	conn, err := net.Dial("udp", "green-ecolution.de:80")
	if err != nil {
		return nil, storage.ErrIPNotFound
	}

	defer conn.Close()

	return conn.LocalAddr().(*net.UDPAddr).IP, nil
}

func getInterface(localIP net.IP) (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", storage.ErrIFacesNotFound
	}

	for _, iface := range ifaces {
		address, err := iface.Addrs()
		if err != nil {
			return "", storage.ErrIFacesAddressNotFound
		}

		for _, addr := range address {
			if addr.(*net.IPNet).IP.String() == localIP.String() {
				return iface.Name, nil
			}
		}
	}

	return "", storage.ErrIFacesNotFound
}
