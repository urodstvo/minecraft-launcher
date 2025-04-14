package minecraft

const maxWorkers int = 10

type MinecraftOptions struct {
	Username              string   `json:"username,omitempty"`
	Uuid                  string   `json:"uuid,omitempty"`
	Token                 string   `json:"token,omitempty"`
	ExecutablePath        string   `json:"executablePath,omitempty"`
	DefaultExecutablePath string   `json:"defaultExecutablePath,omitempty"`
	JvmArguments          []string `json:"jvmArguments,omitempty"`
	LauncherName          string   `json:"launcherName,omitempty"`
	LauncherVersion       string   `json:"launcherVersion,omitempty"`
	GameDirectory         string   `json:"gameDirectory,omitempty"`
	Demo                  bool     `json:"demo,omitempty"`
	CustomResolution      bool     `json:"customResolution,omitempty"`
	ResolutionWidth       string   `json:"resolutionWidth,omitempty"`
	ResolutionHeight      string   `json:"resolutionHeight,omitempty"`
	Server                string   `json:"server,omitempty"`
	Port                  string   `json:"port,omitempty"`
	NativesDirectory      string   `json:"nativesDirectory,omitempty"`
	EnableLoggingConfig   bool     `json:"enableLoggingConfig,omitempty"`
	DisableMultiplayer    bool     `json:"disableMultiplayer,omitempty"`
	DisableChat           bool     `json:"disableChat,omitempty"`
	QuickPlayPath         *string  `json:"quickPlayPath,omitempty"`
	QuickPlaySingleplayer *string  `json:"quickPlaySingleplayer,omitempty"`
	QuickPlayMultiplayer  *string  `json:"quickPlayMultiplayer,omitempty"`
	QuickPlayRealms       *string  `json:"quickPlayRealms,omitempty"`
}

type Minecraft struct {
	Config struct {
		Directory string
	}
}

type IMinecraft interface {
}