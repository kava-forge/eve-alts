package app

import (
	"bytes"
	_ "embed"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/hashicorp/go-multierror"
	"github.com/kava-forge/eve-alts/lib/deferutil"
	"github.com/kava-forge/eve-alts/lib/errors"
	"github.com/kirsle/configdir"
	"github.com/spf13/viper"

	"github.com/kava-forge/eve-alts/pkg/keys"
)

const (
	ConfDirName             = "evealts"
	PromNamespace           = "evealts"
	DefaultCallbackHostport = "localhost:8619"
	DefaultCallbackPath     = "/callback"
	DefaultCallbackScheme   = "http"
	DefaultPProfHostport    = "localhost:8089"
	DefaultStatsHostport    = "localhost:8091"
)

var ensureConfDir sync.Once

func GetConfigDir() string {
	configDir := configdir.LocalConfig(ConfDirName)
	ensureConfDir.Do(func() {
		err := configdir.MakePath(configDir) // Ensure it exists.
		if err != nil {
			panic(errors.Wrap(err, "could not ensure that the configDir exists"))
		}
	})

	return configDir
}

//go:embed default-config.toml
var defaultConfig []byte

var ErrMissingSecret = errors.New("missing secret")

type TelemeterConf struct {
	JaegerHostPort      string  `mapstructure:"jaeger_hostport"`
	PrometheusNamespace string  `mapstructure:"prometheus_namespace"`
	TraceProbability    float64 `mapstructure:"trace_probability"`
}

func (c *TelemeterConf) FillDefaults() error {
	if c.PrometheusNamespace == "" {
		c.PrometheusNamespace = PromNamespace
	}

	return nil
}

type LoggingConf struct {
	Format    string `mapstructure:"format"`
	Level     string `mapstructure:"level"`
	Directory string `mapstructure:"directory"`
}

func (c *LoggingConf) FillDefaults() error {
	if c.Level == "" {
		c.Level = "error"
	}

	if c.Directory == "" {
		c.Directory = filepath.Join(GetConfigDir(), "logs")
		err := configdir.MakePath(c.Directory) // Ensure it exists.
		if err != nil {
			return errors.Wrap(err, "could not ensure that the configDir exists")
		}
	}

	return nil
}

type DatabaseConf struct {
	Location       string `mapstructure:"location"`
	StaticLocation string `mapstructure:"static_location"`
	Database       string `mapstructure:"database"`
}

func (c *DatabaseConf) FillDefaults() error {
	if c.Location == "" {
		c.Location = filepath.Join(GetConfigDir(), "database.db")
	}

	if c.StaticLocation == "" {
		c.StaticLocation = filepath.Join(GetConfigDir(), "staticdata.db")
	}

	return nil
}

type PProfConf struct {
	Enabled bool `mapstructure:"enabled"`
}

type ServingConf struct {
	HostPort           string `mapstructure:"hostport"`
	CallbackPath       string `mapstructure:"callback_path"`
	CallbackScheme     string `mapstructure:"callback_scheme"`
	PProfHostPort      string `mapstructure:"pprof_hostport"`
	PrometheusHostPort string `mapstructure:"prometheus_hostport"`
}

func (c *ServingConf) FillDefaults() error {
	if c.HostPort == "" {
		c.HostPort = DefaultCallbackHostport
	}

	if c.CallbackPath == "" {
		c.CallbackPath = DefaultCallbackPath
	}

	if c.CallbackScheme == "" {
		c.CallbackScheme = DefaultCallbackScheme
	}

	if c.PProfHostPort == "" {
		c.PProfHostPort = DefaultPProfHostport
	}

	if c.PrometheusHostPort == "" {
		c.PrometheusHostPort = DefaultStatsHostport
	}

	return nil
}

type Config struct {
	Database  DatabaseConf  `mapstructure:"database"`
	Logging   LoggingConf   `mapstructure:"logging"`
	PProf     PProfConf     `mapstructure:"pprof"`
	Serving   ServingConf   `mapstructure:"serving"`
	Telemetry TelemeterConf `mapstructure:"telemetry"`

	ConfigFile string `mapstructure:"-"`
}

func (c *Config) FillDefaults() error {
	var errs error

	if err := c.Database.FillDefaults(); err != nil {
		errs = multierror.Append(errs, err)
	}

	if err := c.Logging.FillDefaults(); err != nil {
		errs = multierror.Append(errs, err)
	}

	if err := c.Serving.FillDefaults(); err != nil {
		errs = multierror.Append(errs, err)
	}

	if err := c.Telemetry.FillDefaults(); err != nil {
		errs = multierror.Append(errs, err)
	}

	return errs
}

func SetupConfig(configFile string) (conf Config, err error) {
	v := viper.New()

	v.SetConfigType("toml")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if configFile == "" {
		configFile = filepath.Join(GetConfigDir(), "config.toml")
	}
	v.SetConfigFile(configFile)
	conf.ConfigFile = configFile

	if _, err := os.Stat(configFile); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			if err := writeDefaultConfig(configFile); err != nil {
				return conf, errors.Wrap(err, "could not write default config file", keys.Path, configFile)
			}
		} else {
			return conf, errors.Wrap(err, "could not access config file", keys.Path, configFile)
		}
	}

	v.SetEnvPrefix("EVE_ALTS")
	v.AutomaticEnv()

	if err := v.ReadConfig(bytes.NewReader(defaultConfig)); err != nil {
		return conf, errors.Wrap(err, "could not read in default config")
	}

	if err := v.MergeInConfig(); err != nil {
		return conf, errors.Wrap(err, "could not read in config file")
	}

	if err := v.Unmarshal(&conf); err != nil {
		return conf, errors.Wrap(err, "could not unmarshal config into struct")
	}

	if err := conf.FillDefaults(); err != nil {
		return conf, errors.Wrap(err, "could not FillDefaults")
	}

	return conf, nil
}

func writeDefaultConfig(path string) error {
	fp, err := os.Create(path)
	if err != nil {
		return errors.Wrap(err, "could not create file")
	}
	defer deferutil.CheckDefer(fp.Close)

	_, err = fp.Write(defaultConfig)
	return errors.Wrap(err, "could not write config defaults")
}
