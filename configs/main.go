package configs

import (
	"aif/utils/log"
	"flag"
	"os"
	"os/user"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Configuration .
type Configuration interface {
	Get(string) (string, error)
	Init(set *flag.FlagSet)
}

// ViperConfiguration .
type ViperConfiguration struct {
}

func (vc *ViperConfiguration) setDefaults() {
	usr, err := user.Current()
	if err != nil {
		log.Fatalf("setDefaults: %+v\n", err)
	}
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatalf("setDefaults: %+v\n", err)
	}
	viper.SetDefault("user", usr.Name+"@"+hostname)
}

// Init .
func (vc *ViperConfiguration) Init() {

	vc.setDefaults()

	// config paths
	viper.AddConfigPath(".")
	viper.SetConfigType("json")
	viper.SetConfigName("configuration")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Fatal("Config file not found; ignore error if desired")
		} else {
			log.Fatalf("Config file error %s", err.Error())
		}
	}

	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	err := viper.BindPFlags(pflag.CommandLine)
	if err != nil {
		log.Fatalf("an error occured while running viper.BindPFlags(): %+v\n", err)
	}
}

// CheckExists .
func (vc *ViperConfiguration) CheckExists(param string) bool {
	return viper.IsSet(param)
}

// Get .
func (vc *ViperConfiguration) Get(param string) string {
	return viper.GetString(param)
}

// Get .
func (vc *ViperConfiguration) GetInterface(param string) interface{} {
	return viper.Get(param)
}

// GetInt .
func (vc *ViperConfiguration) GetInt(param string) int {
	return viper.GetInt(param)
}

// GetBool .
func (vc *ViperConfiguration) GetBool(param string) bool {
	return viper.GetBool(param)
}

// NewConfiguration .
func NewConfiguration() (cfg *ViperConfiguration) {
	cfg = &ViperConfiguration{}
	return
}

// Set .
func (vc *ViperConfiguration) Set(key string, value interface{}) {
	viper.Set(key, value)
	err := viper.WriteConfig()
	if err != nil {
		log.Fatal(err)
	}
}
