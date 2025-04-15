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

type minecraftConfig struct {
	Config struct {
		Directory string
	}
}

type API interface {
	InstallMinecraftVersion(versionId string) error

	GetMinecraftCommand(version string, options MinecraftOptions) ([]string, error)

	GetMinecraftDirectory() string
	GetLatestVersion() (LatestMinecraftVersions, error)
	GetVersionList() ([]MinecraftVersionInfo, error)
	GetInstalledVersions(minecraftDirectory string) ([]MinecraftVersionInfo, error)
	GetAvailableVersions(minecraftDirectory string) ([]MinecraftVersionInfo, error)
	GenerateTestOptions() MinecraftOptions
	IsPlatformSupported() bool
	IsMinecraftInstalled(minecraftDirectory string) bool

	FindSystemJavaVersions(additionalDirectories []string) ([]string, error)
	GetJavaInformation(path string) (JavaInformation, error)
	GetSystemJavaVersionInformation(additionalDirectories []string) ([]JavaInformation, error)

	GetLibraryVersion() string
}

type Opts struct {
	MinecraftDirectory string
}

func NewAPI(opts Opts) API {
	return &minecraftConfig{
		Config: struct{ Directory string }{
			Directory: opts.MinecraftDirectory,
		},
	}
}