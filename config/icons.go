package config

type Icons struct {
	Computer    string `mapstructure:"computer"`
	Tablet      string `mapstructure:"tablet"`
	Smartphone  string `mapstructure:"smartphone"`
	Speaker     string `mapstructure:"speaker"`
	Tv          string `mapstructure:"tv"`
	Avr         string `mapstructure:"avr"`
	Stb         string `mapstructure:"stb"`
	AudioDongle string `mapstructure:"audio_dongle"`
	GameConsole string `mapstructure:"game_console"`
	CastVideo   string `mapstructure:"cast_video"`
	CastAudio   string `mapstructure:"cast_audio"`
	Automobile  string `mapstructure:"automobile"`
	Playing     string `mapstructure:"playing"`
	Paused      string `mapstructure:"paused"`
	ShuffleOn   string `mapstructure:"shuffle_on"`
	ShuffleOff  string `mapstructure:"shuffle_off"`
	RepeatOne   string `mapstructure:"repeat_one"`
	RepeatAll   string `mapstructure:"repeat_all"`
	RepeatOff   string `mapstructure:"repeat_off"`
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
