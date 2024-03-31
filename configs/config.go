package configs

import (
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	projectcalicov3 "github.com/projectcalico/api/pkg/client/clientset_generated/clientset/typed/projectcalico/v3"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/projectcalico/api/pkg/client/clientset_generated/clientset"
	"github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

const (
	PolicyEngineType = "MICROSEG_POLICY_ENGINE_TYPE"
)

const (
	CalicoEngine = "calico"
	K8sEngine    = "kubernetes"
)

var (
	ApiServerURL   string
	Cfg            *Config
	MicrosegK8sCli *kubernetes.Clientset
	CalicoCli      projectcalicov3.ProjectcalicoV3Interface
	PolicyEngine   string
)

type Config struct {
	RunMode string `mapstructure:"runMode"` // dev or prod
}

func InitConfig() (*rest.Config, error) {
	var err error
	ApiServerURL = fmt.Sprintf("https://%s:%s", os.Getenv("KUBERNETES_SERVICE_HOST"), os.Getenv("KUBERNETES_SERVICE_PORT"))
	logrus.SetFormatter(&LogFormatter{})

	config := &rest.Config{
		Host: ApiServerURL,
		//QPS: 20,
	}
	config.Insecure = true

	MicrosegK8sCli, err = kubernetes.NewForConfig(config)
	if err != nil {
		return nil, errors.Wrap(err, "failed to init k8s clientset")
	}

	PolicyEngine = os.Getenv(PolicyEngineType)
	switch PolicyEngine {
	case CalicoEngine:
		CalicoCliSet, err := clientset.NewForConfig(config)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create calico client by using k8s as datastore")
		}
		CalicoCli = CalicoCliSet.ProjectcalicoV3()
	case K8sEngine:
	default:
		return nil, errors.Errorf("unsupported policy engine type: %s", PolicyEngine)
	}

	logrus.Infof("init application with config: %+v", Cfg)
	return config, nil
}

type LogFormatter struct{}

func (s *LogFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	timestamp := time.Now().Local().Format("2006/01/02 - 15:04:05.000")
	var file string
	var len int
	if entry.Caller != nil {
		file = filepath.Base(entry.Caller.File)
		len = entry.Caller.Line
	}
	//fmt.Println(entry.Data)
	msg := fmt.Sprintf("%s [%s:%d][GOID:%d][%s] %s\n", timestamp, file, len, getGID(), strings.ToUpper(entry.Level.String()), entry.Message)
	return []byte(msg), nil
}

func getGID() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}
