package config

import (
	"os"

	"github.com/aditya-K2/utils"
	"github.com/fsnotify/fsnotify"
	"github.com/gdamore/tcell/v2"
	"github.com/spf13/viper"
)

var (
	configDir, configErr   = os.UserConfigDir()
	userCacheDir, cacheErr = os.UserCacheDir()
	ConfigPath             = configDir + "/gspt"
	Config                 = NewConfigS()
	OnConfigChange         func()
	DefaultMappings        = map[string]map[string]map[Key]string{
		"recently_played_view": {
			"normal": {
				{K: tcell.KeyEnter}: "open_entry",
				{R: 'a'}:            "add_to_playlist",
			},
			"visual": {
				{R: 'a'}: "add_to_playlist",
			},
		},
		"nav_menu": {
			"normal": {
				{K: tcell.KeyEnter}: "open_entry",
			},
		},
		"search_view": {
			"normal": {
				{K: tcell.KeyEnter}: "open_entry",
				{K: tcell.KeyCtrlP}: "play_entry",
			},
		},
		"global": {
			"normal": {
				{R: 'd'}:            "choose_device",
				{R: '1'}:            "focus_nav",
				{R: '2'}:            "focus_playlists",
				{R: '3'}:            "focus_main_view",
				{R: '?'}:            "focus_search",
				{R: ' '}:            "toggle_playback",
				{R: 'o'}:            "open_current_track_album",
				{R: 'O'}:            "open_current_track_artist",
				{R: 'n'}:            "next",
				{R: 'p'}:            "previous",
				{K: tcell.KeyCtrlO}: "open_current_context",
			},
		},
		"playlist_nav": {
			"normal": {
				{K: tcell.KeyEnter}: "open_entry",
				{K: tcell.KeyCtrlP}: "play_entry",
			},
		},
		"playlist_view": {
			"normal": {
				{K: tcell.KeyEnter}: "open_entry",
				{R: 'a'}:            "add_to_playlist",
			},
			"visual": {
				{R: 'a'}: "add_to_playlist",
			},
		},
		"top_tracks_view": {
			"normal": {
				{K: tcell.KeyEnter}: "open_entry",
				{K: tcell.KeyCtrlP}: "play_entry",
			},
		},
		"liked_songs_view": {
			"normal": {
				{K: tcell.KeyEnter}: "open_entry",
				{R: 'a'}:            "add_to_playlist",
			},
			"visual": {
				{R: 'a'}: "add_to_playlist",
			},
		},
		"artists_view": {
			"normal": {
				{K: tcell.KeyEnter}: "open_entry",
			},
		},
		"artist_view": {
			"normal": {
				{K: tcell.KeyEnter}: "open_entry",
				{K: tcell.KeyCtrlP}: "play_entry",
			},
		},
		"albums_view": {
			"normal": {
				{K: tcell.KeyEnter}: "open_entry",
				{K: tcell.KeyCtrlP}: "play_entry",
				{R: 'i'}:            "queue_entry",
			},
		},
		"album_view": {
			"normal": {
				{K: tcell.KeyEnter}: "open_entry",
				{R: 'a'}:            "add_to_playlist",
			},
			"visual": {
				{R: 'a'}: "add_to_playlist",
			},
		},
	}
)

type ConfigS struct {
	CacheDir           string  `yaml:"cache_dir" mapstructure:"cache_dir"`
	RedrawInterval     int     `yaml:"redraw_interval" mapstructure:"redraw_interval"`
	Colors             *Colors `mapstructure:"colors"`
	AdditionalPaddingX int     `yaml:"additional_padding_x" mapstructure:"additional_padding_x"`
	AdditionalPaddingY int     `yaml:"additional_padding_y" mapstructure:"additional_padding_y"`
	ImageWidthExtraX   int     `yaml:"image_width_extra_x" mapstructure:"image_width_extra_x"`
	ImageWidthExtraY   int     `yaml:"image_width_extra_y" mapstructure:"image_width_extra_y"`
	HideImage          bool    `yaml:"hide_image" mapstructure:"hide_image"`
	RoundedCorners     bool    `yaml:"rounded_corners" mapstructure:"rounded_corners"`
}

func NewConfigS() *ConfigS {
	return &ConfigS{
		CacheDir:       utils.CheckDirectoryFmt(userCacheDir),
		RedrawInterval: 500,
		Colors:         NewColors(),
		HideImage:      false,
		RoundedCorners: false,
	}
}

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

func GenerateMappings() map[string]map[string]map[Key]string {
	all := viper.GetStringMap("mappings")
	keys := DefaultMappings
	for view, mode := range all {
		if keys[view] == nil {
			keys[view] = make(map[string]map[Key]string)
		}
		for modeName, mappings := range mode.(map[string]map[string]string) {
			for key, function := range mappings {
				keys[view][modeName][NewKey(key)] = function
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
