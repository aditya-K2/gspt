package config

import (
	"fmt"
	"os"
	"path/filepath"

	"errors"

	"github.com/aditya-K2/utils"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var (
	configDir, configErr   = os.UserConfigDir()
	userCacheDir, cacheErr = os.UserCacheDir()
	userConfigPath         = filepath.Join(configDir, "gspt")
	Config                 = newConfigS()
	OnConfigChange         func()
	version                = "unknown"
	buildDate              = "unknown"
)

const (
	ImageHidden    string = "hidden"
	ImageShow      string = "show"
	CornersRounded string = "rounded"
	CornersDefault string = "default"
)

type configS struct {
	CacheDir           string  `yaml:"cache_dir" mapstructure:"cache_dir"`
	RedrawInterval     int     `yaml:"redraw_interval" mapstructure:"redraw_interval"`
	Colors             *Colors `mapstructure:"colors"`
	Icons              *Icons  `mapstructure:"icons"`
	AdditionalPaddingX int     `yaml:"additional_padding_x" mapstructure:"additional_padding_x"`
	AdditionalPaddingY int     `yaml:"additional_padding_y" mapstructure:"additional_padding_y"`
	ImageWidthExtraX   int     `yaml:"image_width_extra_x" mapstructure:"image_width_extra_x"`
	ImageWidthExtraY   int     `yaml:"image_width_extra_y" mapstructure:"image_width_extra_y"`
	Image              string  `yaml:"image" mapstructure:"image"`
	Corners            string  `yaml:"corners" mapstructure:"corners"`
	UseIcons           bool    `yaml:"use_icons" mapstructure:"use_icons"`
}

func newConfigS() *configS {
	return &configS{
		CacheDir:       utils.CheckDirectoryFmt(userCacheDir),
		RedrawInterval: 500,
		Colors:         NewColors(),
		Icons:          NewIcons(),
		Corners:        CornersDefault,
		Image:          ImageShow,
	}
}

func ReadConfig() error {
	parseFlags()

	if Flags.Version {
		fmt.Printf("gspt: %s \nBuild Date: %s\n", version, buildDate)
		os.Exit(0)
	}

	// If config path is provided through command-line use that
	if Flags.ConfigPath != "" {
		userConfigPath = Flags.ConfigPath
	}

	if configErr != nil {
		utils.Print("RED", "Couldn't get $XDG_CONFIG!")
		return configErr
	}

	if cacheErr != nil {
		utils.Print("RED", "Couldn't get CACHE DIR!")
		return cacheErr
	}

	viper.SetConfigName("config")
	viper.AddConfigPath(utils.ExpandHomeDir(userConfigPath))

	if err := viper.ReadInConfig(); err != nil {
		utils.Print("RED", "Could Not Read Config file.\n")
	} else {
		viper.Unmarshal(Config)
	}

	// Expanding ~ to the User's Home Directory
	expandHome := func() {
		Config.CacheDir = utils.ExpandHomeDir(Config.CacheDir)
	}

	if Flags.Image != "" {
		if Flags.Image == ImageHidden || Flags.Image == ImageShow {
			Config.Image = Flags.Image
		} else {
			return errors.New(fmt.Sprintf("Undefined value provided to --image flag: '%s' ( accepted: %s | %s )", Flags.Image, ImageHidden, ImageShow))
		}
	}

	if Flags.Corners != "" {
		if Flags.Corners == CornersRounded || Flags.Corners == CornersDefault {
			Config.Corners = Flags.Corners
		} else {
			return errors.New(fmt.Sprintf("Undefined value provided to --corners flag: '%s' ( accepted: %s | %s )", Flags.Corners, CornersRounded, CornersDefault))
		}
		Config.Corners = Flags.Corners
	}

	if Flags.UseIcons != false {
		Config.UseIcons = Flags.UseIcons
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
	return nil
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
