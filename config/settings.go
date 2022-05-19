package config

import (
	"reflect"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

type bootnodeConfig struct {
	Name             string
	ListeningAddress string
	CacheSize        int
}

type chainConfig struct {
	ID                    int
	StateDbPath           string
	StateDbInMemory       bool
	TransactionDbPath     string
	TransactionDbInMemory bool
	MonitorDbPath         string
	MonitorDbInMemory     bool
}

type ethereumConfig struct {
	DefaultAccount            string
	DeployAccount             string
	Endpoint                  string
	EndpointMinimumPeers      int
	FinalityDelay             int
	Keystore                  string
	MerkleProofContract       string
	Passcodes                 string
	RegistryAddress           string
	RetryCount                int
	RetryDelay                time.Duration
	StartingBlock             uint64
	TestEther                 string
	Timeout                   time.Duration
	TxFeePercentageToIncrease int
	TxMaxFeeThresholdInGwei   uint64
	TxCheckFrequency          time.Duration
	TxTimeoutForReplacement   time.Duration
}

type monitorConfig struct {
	BatchSize int
	Interval  time.Duration
	Timeout   time.Duration
}

type transportConfig struct {
	Size                       int
	Timeout                    time.Duration
	OriginLimit                int
	PeerLimitMin               int
	PeerLimitMax               int
	FirewallMode               bool
	FirewallHost               string
	Whitelist                  string
	PrivateKey                 string
	BootNodeAddresses          string
	P2PListeningAddress        string
	DiscoveryListeningAddress  string
	LocalStateListeningAddress string
	UPnP                       bool
}

type deployConfig struct {
	Migrations     bool
	TestMigrations bool
}

type utilsConfig struct {
	Status bool
}

type validatorConfig struct {
	Repl            bool
	RewardAccount   string
	RewardCurveSpec int
	SymmetricKey    string
}

type loglevelConfig struct {
	Madnet     string
	Consensus  string
	Transport  string
	App        string
	Db         string
	Gossipbus  string
	Badger     string
	PeerMan    string
	LocalRPC   string
	Dman       string
	Peer       string
	Yamux      string
	Ethereum   string
	Main       string
	Deploy     string
	Utils      string
	Monitor    string
	Dkg        string
	Services   string
	Settings   string
	Validator  string
	MuxHandler string
	Bootnode   string
	P2pmux     string
	Status     string
	Test       string
}

type LogfileConfig struct {
	FileName   string
	MinLevel   string
	MaxAgeDays float64
}

type firewalldConfig struct {
	Enabled    bool
	SocketFile string
}

type configuration struct {
	ConfigurationFileName string
	Loglevel              loglevelConfig
	Logfile               LogfileConfig
	Deploy                deployConfig
	Ethereum              ethereumConfig
	Monitor               monitorConfig
	Transport             transportConfig
	Utils                 utilsConfig
	Validator             validatorConfig
	Firewalld             firewalldConfig
	Chain                 chainConfig
	BootNode              bootnodeConfig
}

// Configuration contains all active settings
var Configuration configuration

func (t transportConfig) BootNodes() []string {
	bootNodeAddresses := strings.Split(t.BootNodeAddresses, ",")
	for idx := range bootNodeAddresses {
		bootNodeAddresses[idx] = strings.TrimSpace(bootNodeAddresses[idx])
	}
	return bootNodeAddresses
}

func LogLevelMap() map[string]logrus.Level {
	llr := reflect.ValueOf(Configuration.Loglevel)
	t := llr.Type()
	len := llr.NumField()
	levels := make(map[string]logrus.Level, len)

	for i := 0; i < len; i++ {
		logName := strings.ToLower(t.Field(i).Name)
		logLevel := strings.ToLower(llr.Field(i).String())
		lvl, err := logrus.ParseLevel(logLevel)
		if err == nil {
			levels[logName] = lvl
		}
	}

	return levels
}
