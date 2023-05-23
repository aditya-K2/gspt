package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/aditya-K2/utils"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var (
	configDir, configErr   = os.UserConfigDir()
	userCacheDir, cacheErr = os.UserCacheDir()
	UserConfigPath         = filepath.Join(configDir, "gspt")
	Config                 = NewConfigS()
	OnConfigChange         func()
	Version                = "unknown"
	BuildDate              = "unknown"
)

type ConfigS struct {
	CacheDir           string  `yaml:"cache_dir" mapstructure:"cache_dir"`
	RedrawInterval     int     `yaml:"redraw_interval" mapstructure:"redraw_interval"`
	Colors             *Colors `mapstructure:"colors"`
	Icons              *Icons  `mapstructure:"icons"`
	AdditionalPaddingX int     `yaml:"additional_padding_x" mapstructure:"additional_padding_x"`
	AdditionalPaddingY int     `yaml:"additional_padding_y" mapstructure:"additional_padding_y"`
	ImageWidthExtraX   int     `yaml:"image_width_extra_x" mapstructure:"image_width_extra_x"`
	ImageWidthExtraY   int     `yaml:"image_width_extra_y" mapstructure:"image_width_extra_y"`
	HideImage          bool    `yaml:"hide_image" mapstructure:"hide_image"`
	RoundedCorners     bool    `yaml:"rounded_corners" mapstructure:"rounded_corners"`
	UseIcons           bool    `yaml:"use_icons" mapstructure:"use_icons"`
}

func NewConfigS() *ConfigS {
	return &ConfigS{
		CacheDir:       utils.CheckDirectoryFmt(userCacheDir),
		RedrawInterval: 500,
		Colors:         NewColors(),
		Icons:          NewIcons(),
	}
}

func ReadConfig() {
	parseFlags()

	if Flags.Version {
		fmt.Printf("gspt %s \nBuild Date: %s\n", Version, BuildDate)
	}

	// If config path is provided through command-line use that
	if Flags.ConfigPath != "" {
		UserConfigPath = Flags.ConfigPath
	}

	if configErr != nil {
		utils.Print("RED", "Couldn't get $XDG_CONFIG!")
		panic(configErr)
	}

	if cacheErr != nil {
		utils.Print("RED", "Couldn't get CACHE DIR!")
		panic(cacheErr)
	}

	viper.SetConfigName("config")
	viper.AddConfigPath(utils.ExpandHomeDir(UserConfigPath))

	if err := viper.ReadInConfig(); err != nil {
		utils.Print("RED", "Could Not Read Config file.\n")
	} else {
		viper.Unmarshal(Config)
	}

	// Expanding ~ to the User's Home Directory
	expandHome := func() {
		Config.CacheDir = utils.ExpandHomeDir(Config.CacheDir)
	}

	// Flags have precedence over config values
	useFlags := func() {
		if Flags.HideImage != false {
			Config.HideImage = Flags.HideImage
		}
		if Flags.RoundedCorners != false {
			Config.RoundedCorners = Flags.RoundedCorners
		}
		if Flags.UseIcons != false {
			Config.UseIcons = Flags.UseIcons
		}
	}
	useFlags()

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

func GenerateMappings() map[string]map[string]map[Key]string {
	userMaps := viper.GetStringMap("mappings")
	keys := DefaultMappings
	for view, modes := range userMaps {
		if keys[view] == nil {
			keys[view] = make(map[string]map[Key]string)
		}
		for mode, mappings := range modes.(map[string]interface{}) {
			for function, key := range mappings.(map[string]interface{}) {
				keys[view][mode][NewKey(key.(string))] = function
			}
		}
	}
	for k := range keys {
		if k != "global" {
			keys[k]["normal"] = utils.MergeMaps(keys["global"]["normal"], keys[k]["normal"])
		}
	}
	return keys
}

func WriteConfig() error {
	viper.Set("additional_padding_x", Config.AdditionalPaddingX)
	viper.Set("additional_padding_y", Config.AdditionalPaddingY)
	viper.Set("image_width_extra_x", Config.ImageWidthExtraX)
	viper.Set("image_width_extra_y", Config.ImageWidthExtraY)
	return viper.WriteConfig()
}
