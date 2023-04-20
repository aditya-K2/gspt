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
	DefaultMappings        = map[string]map[Key]string{
		"recently_played_view": {
			{K: tcell.KeyEnter}: "open_entry",
		},
		"nav_menu": {
			{K: tcell.KeyEnter}: "open_entry",
		},
		"search_view": {
			{K: tcell.KeyEnter}: "open_entry",
		},
		"global": {
			{R: 'd'}:            "choose_device",
			{R: '1'}:            "focus_nav",
			{R: '2'}:            "focus_playlists",
			{R: '3'}:            "focus_main_view",
			{R: '?'}:            "focus_search",
			{R: ' '}:            "toggle_playback",
			{R: 'o'}:            "open_current_track_album",
			{R: 'O'}:            "open_current_track_artist",
			{K: tcell.KeyCtrlO}: "open_current_context",
		},
		"playlist_nav": {
			{K: tcell.KeyEnter}: "open_entry",
			{K: tcell.KeyCtrlP}: "play_entry",
		},
		"playlist_view": {
			{K: tcell.KeyEnter}: "open_entry",
		},
		"top_tracks_view": {
			{K: tcell.KeyEnter}: "open_entry",
			{K: tcell.KeyCtrlP}: "play_entry",
		},
		"liked_songs_view": {
			{K: tcell.KeyEnter}: "open_entry",
		},
		"artists_view": {
			{K: tcell.KeyEnter}: "open_entry",
		},
		"artist_view": {
			{K: tcell.KeyEnter}: "open_entry",
			{K: tcell.KeyCtrlP}: "play_entry",
		},
		"albums_view": {
			{K: tcell.KeyEnter}: "open_entry",
			{K: tcell.KeyCtrlP}: "play_entry",
		},
		"album_view": {
			{K: tcell.KeyEnter}: "open_entry",
		},
	}
)

type ConfigS struct {
	CacheDir           string  `mapstructure:"cache_dir"`
	RedrawInterval     int     `mapstructure:"redraw_interval"`
	Colors             *Colors `mapstructure:"colors"`
	AdditionalPaddingX int     `mapstructure:"additional_padding_x"`
	AdditionalPaddingY int     `mapstructure:"additional_padding_y"`
	ImageWidthExtraX   int     `mapstructure:"image_width_extra_x"`
	ImageWidthExtraY   int     `mapstructure:"image_width_extra_y"`
	HideImage          bool    `mapstructure:"hide_image"`
	RoundedCorners     bool    `mapstructure:"rounded_corners"`
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

func GenerateMappings() map[string]map[Key]string {
	all := viper.GetStringMap("mappings")
	keys := DefaultMappings
	for view, mappings := range all {
		if keys[view] == nil {
			keys[view] = make(map[Key]string)
		}
		if mappings != nil {
			for function, key := range mappings.(map[string]interface{}) {
				keys[view][NewKey(key.(string))] = function
			}
		}
	}
	for k := range keys {
		if k != "global" {
			keys[k] = utils.MergeMaps(keys["global"], keys[k])
		}
	}
	return keys
}
