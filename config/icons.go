package config

type Icons struct {
	Computer    string `mapstructure:"computer"`
	Tablet      string `mapstructure:"tablet"`
	Smartphone  string `mapstructure:"smartphone"`
	Speaker     string `mapstructure:"speaker"`
	Tv          string `mapstructure:"tv"`
	Avr         string `mapstructure:"avr"`
	Stb         string `mapstructure:"stb"`
	AudioDongle string `yaml:"audio_dongle" mapstructure:"audio_dongle"`
	GameConsole string `yaml:"game_console" mapstructure:"game_console"`
	CastVideo   string `yaml:"cast_video" mapstructure:"cast_video"`
	CastAudio   string `yaml:"cast_audio" mapstructure:"cast_audio"`
	Automobile  string `mapstructure:"automobile"`
	Playing     string `mapstructure:"playing"`
	Paused      string `mapstructure:"paused"`
	ShuffleOn   string `yaml:"shuffle_on" mapstructure:"shuffle_on"`
	ShuffleOff  string `yaml:"shuffle_off" mapstructure:"shuffle_off"`
	RepeatOne   string `yaml:"repeat_one" mapstructure:"repeat_one"`
	RepeatAll   string `yaml:"repeat_all" mapstructure:"repeat_all"`
	RepeatOff   string `yaml:"repeat_off" mapstructure:"repeat_off"`
}

func NewIcons() *Icons {
	return &Icons{
		"󰍹",
		"",
		"󰄜",
		"󰓃",
		"",
		"󰤽",
		"󰤽",
		"󱡬",
		"󰺵",
		"󰄙",
		"󰄙",
		"",
		"",
		"",
		"󰒟",
		"󰒞",
		"󰑘",
		"󰑗",
		"󰑖",
	}
}
