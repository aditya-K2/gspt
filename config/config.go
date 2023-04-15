package config

import (
	"os"

	"github.com/aditya-K2/utils"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type ConfigS struct {
	CacheDir           string  `mapstructure:"cache_dir"`
	RedrawInterval     int     `mapstructure:"redraw_interval"`
	Colors             *Colors `mapstructure:"colors"`
	AdditionalPaddingX int     `mapstructure:"additional_padding_x"`
	AdditionalPaddingY int     `mapstructure:"additional_padding_y"`
	ExtraImageWidthX   float64 `mapstructure:"image_width_extra_x"`
	ExtraImageWidthY   float64 `mapstructure:"image_width_extra_y"`
}

func NewConfigS() *ConfigS {
	return &ConfigS{
		AdditionalPaddingX: 12,
		AdditionalPaddingY: 16,
		ExtraImageWidthX:   -1.5,
		ExtraImageWidthY:   -3.75,
		CacheDir:           utils.CheckDirectoryFmt(userCacheDir),
		RedrawInterval:     500,
		Colors:             NewColors(),
	}
}

var (
	configDir, configErr   = os.UserConfigDir()
	userCacheDir, cacheErr = os.UserCacheDir()
	ConfigPath             = configDir + "/gspt"
	Config                 = NewConfigS()
	OnConfigChange         func()
)

func ReadConfig() {
	parseFlags()
	if configErr != nil {
		utils.Print("RED", "Couldn't get $XDG_CONFIG!")
		panic(configErr)
	}

	if cacheErr != nil {
		utils.Print("RED", "Couldn't get CACHE DIR!")
		panic(cacheErr)
	}

	viper.SetConfigName("config")
	viper.AddConfigPath(ConfigPath)

	if err := viper.ReadInConfig(); err != nil {
		utils.Print("RED", "Could Not Read Config file.\n")
	} else {
		viper.Unmarshal(Config)
	}

	// Expanding ~ to the User's Home Directory
	expandHome := func() {
		Config.CacheDir = utils.ExpandHomeDir(Config.CacheDir)
	}

	viper.OnConfigChange(func(e fsnotify.Event) {
		viper.Unmarshal(Config)
		expandHome()
		if OnConfigChange != nil {
			OnConfigChange()
		}
	})
	viper.WatchConfig()

	expandHome()
}

func GenerateMappings() map[string]map[Key]string {
	all := viper.GetStringMap("mappings")
	keys := make(map[string]map[Key]string)
	for view, mappings := range all {
		if keys[view] == nil {
			keys[view] = make(map[Key]string)
		}
		for function, key := range mappings.(map[string]interface{}) {
			keys[view][NewKey(key.(string))] = function
		}
	}
	return keys
}
