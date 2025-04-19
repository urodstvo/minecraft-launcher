package minecraft

import (
	"time"
)

type ClientJsonRule struct{
	Action 	string `json:"action"`
	Os 		struct {
		Name *string `json:"name"`
		Arch *string `json:"arch"`
		Version *string `json:"version"`	
	} `json:"os"`
	Features *struct {
		HasCustomResolution 	*bool `json:"has_custom_resolution"`
		IsDemoUser 				*bool `json:"is_demo_user"`
		HasQuickPlaysSupport 	*bool `json:"has_quick_plays_support"`
		IsQuickPlaySingleplayer *bool `json:"is_quick_play_singleplayer"`
		IsQuickPlayMultiplayer 	*bool `json:"is_quick_play_multiplayer"`
		IsQuickPlayRealms 		*bool `json:"is_quick_play_realms"`
	} `json:"features"`
}

type ClientJsonArgumentRule struct{
	CompatibilityRules 	[]ClientJsonRule 	`json:"-"`
	Rules 				[]ClientJsonRule 	`json:"rules"`
	Value 				any 				`json:"value"`
}

type clientJsonAssetIndex struct{
	Id 			string 	`json:"id"`
	Sha1 		string 	`json:"sha1"`
	Size 		int 	`json:"size"`
	TotalSize 	int 	`json:"totalSize"`
	Url 		string 	`json:"url"`
}

type clientJsonDownloads struct{
	Sha1 	string 	`json:"sha1"`
	Size 	int 	`json:"size"`
	Url 	string 	`json:"url"`
}

type clientJsonJavaVersion struct{
	Component		string 	`json:"component"`
	MajorVersion 	int 	`json:"majorVersion"`
}

type clientJsonLibraryDownloadsArtifact struct{
	Path 	string 	`json:"path"`
	Url 	string 	`json:"url"`
	Sha1 	string	`json:"sha1"`
	Size 	int 	`json:"size"`
}

type clientJsonLibraryDownloads struct{
	Artifact *clientJsonLibraryDownloadsArtifact `json:"artifact"`
	Classifiers map[string]clientJsonLibraryDownloadsArtifact `json:"classifiers"`
}

type ClientJsonLibrary struct{
	Name 		string `json:"name"`
	Downloads 	clientJsonLibraryDownloads `json:"downloads"`
	Extract *struct {
		Exclude []string `json:"extract"`
	} `json:"exclude"`
	Rules 	[]ClientJsonRule `json:"rules"`
	Natives *struct{
		Linux 	*string `json:"linux"`
		Osx 	*string `json:"osx"`
		Windows *string `json:"windows"`
	} `json:"natives"`
	Url *string `json:"url"`
}

type clientJsonLoggingFile struct{
	Id 		string 	`json:"id"`
	Sha1 	string 	`json:"sha1"`
	Size 	int 	`json:"size"`
	Url 	string 	`json:"url"`
}

type clientJsonLogging struct{
	Argument 	string 					`json:"argument"`
	File 		clientJsonLoggingFile 	`json:"file"`
	Type 		string 					`json:"type"`
}

type ClientJson struct{
	Id 		string `json:"id"`
	Jar 	string `json:"jar"`
	Arguments *struct {
		Game []any `json:"game"`
		Jvm []any `json:"jvm"`
	} `json:"arguments"`
	MinecraftArguments string `json:"minecraftArguments"`
	AssetIndex 		*clientJsonAssetIndex `json:"assetIndex"`
	Assets 			string `json:"assets"`
	Downloads struct {
		Client 			clientJsonDownloads `json:"client"`
		ClientMappings 	clientJsonDownloads `json:"client_mappings"`
		Server 			clientJsonDownloads `json:"server"`
		ServerMappings 	clientJsonDownloads `json:"server_mappings"`
	} `json:"downloads"`
	JavaVersion 	clientJsonJavaVersion `json:"javaVersion"`
	Libraries 		[]ClientJsonLibrary `json:"libraries"`
	Logging *struct {
		Client clientJsonLogging `json:"client"`
	} `json:"logging"`
	MainClass 				string 	`json:"mainClass"`
	MinimumLauncherVersion 	int 	`json:"minimumLauncherVersion"`
	ReleaseTime 			string 	`json:"releaseTime"`
	Time 					string 	`json:"time"`
	Type 					string 	`json:"type"`
	ComplianceLevel 		int 	`json:"complianceLevel"`
	InheritsFrom 			string 	`json:"-"`
}

type versionListManifestJsonVersion struct {
	Id 				string `json:"id"`
	Type 			string `json:"type"`
	Url 			string `json:"url"`
	Time 			string `json:"time"`
	ReleaseTime 	string `json:"releaseTime"`
	Sha1 			string `json:"sha1"`
	ComplianceLevel int `json:"complianceLevel"`
}

type VersionListManifestJson struct{
	Latest struct {
		Release 	string `json:"release"`
		Snapshot 	string `json:"snapshot"`
	} `json:"latest"`
	Versions []versionListManifestJsonVersion `json:"versions"`
}

type assetsJsonObject struct{
	Hash string `json:"hash"`
	Size int `json:"size"`
}

type AssetsJson struct {
	Objects map[string]assetsJsonObject `json:"objects"`
}

type JavaInformation struct {
	Path string		
    Name string
    Version string
    JavaPath string
    JavawPath *string
    Is64bit bool
    Openjdk bool
}

type JVMRuntimeInformation struct {
	Name     string 	`json:"name"`
	Released time.Time	`json:"released"`
}

type VersionRuntimeInformation struct {
	Name             string 
	JavaMajorVersion int 
}

type runtimeListJsonEntryManifest struct {
	SHA1 string `json:"sha1"`
    Size int `json:"size"`
    Url string `json:"url"`
}

type runtimeListJsonEntry struct {
	Availability struct {
		Group 		int	`json:"group"`
		Progress 	int	 `json:"progress"`
	} `json:"availability"`
    Manifest runtimeListJsonEntryManifest  `json:"manifest"`
    Version struct {
		Name 		string  `json:"name"`
		Released 	string  `json:"released"`
	}  `json:"version"`
}

type RuntimeListJson map[string]map[string][]runtimeListJsonEntry

type platformManifestJsonFileDownloads struct {
	SHA1 string `json:"sha1"`
    Size int `json:"size"`
    Url string `json:"url"`
}

type platformManifestJsonFile struct{
	Downloads map[string] platformManifestJsonFileDownloads `json:"downloads"`
    Type string `json:"type"`
    Executable bool `json:"executable"`
    Target string `json:"target"`
}

type PlatformManifestJson struct{
	Files map[string]platformManifestJsonFile `json:"files"`
}

type LatestMinecraftVersions struct {
	Release string
	Snapshot string
}

type MinecraftVersionInfo struct {
	Id string `json:"id"`
    Type string `json:"type"`
    ReleaseTime string `json:"releaseTime"`
    ComplianceLevel int `json:"complianceLevel"`
}

type requestsResponseCache struct {
	Response []byte
	Datetime time.Time
}

type ProgressCallback func(message string)

type Callback struct {
	Status ProgressCallback
	Progress ProgressCallback
	Max ProgressCallback
}

type Resolution struct {
	Height int `json:"height"`
	Width  int `json:"width"`
}

type VanillaLauncherProfilesJsonProfile struct {
	Created       string     `json:"created,omitempty"`
	GameDir       string     `json:"gameDir,omitempty"`
	Icon          string     `json:"icon,omitempty"`
	JavaArgs      string     `json:"javaArgs,omitempty"`
	JavaDir       string     `json:"javaDir,omitempty"`
	LastUsed      string     `json:"lastUsed,omitempty"`
	LastVersionID string     `json:"lastVersionId,omitempty"`
	Name          string     `json:"name,omitempty"`
	Resolution    *Resolution `json:"resolution,omitempty"`
	Type          string     `json:"type,omitempty"`
}

type VanillaLauncherProfilesJson struct {
	Profiles map[string]VanillaLauncherProfilesJsonProfile `json:"profiles"`
	Version  int                                            `json:"version"`
}

type VanillaLauncherProfile struct {
	Name            string     `json:"name"`
	VersionType     string     `json:"versionType"` // latest-release | latest-snapshot | custom
	Version         *string    `json:"version,omitempty"`
	GameDirectory   *string    `json:"gameDirectory,omitempty"`
	JavaExecutable  *string    `json:"javaExecutable,omitempty"`
	JavaArguments   []string   `json:"javaArguments,omitempty"`
	CustomResolution *Resolution `json:"customResolution,omitempty"`
}


type AuthorizationTokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	Scope        string `json:"scope"`
	RefreshToken string `json:"refresh_token"`
}

type Xui struct {
	Uhs string `json:"uhs"`
}

type DisplayClaims struct {
	Xui []Xui `json:"xui"`
}

type XBLResponse struct {
	IssueInstant   string        `json:"IssueInstant"`
	NotAfter       string        `json:"NotAfter"`
	Token          string        `json:"Token"`
	DisplayClaims  DisplayClaims `json:"DisplayClaims"`
}

type XSTSResponse struct {
	IssueInstant    string        `json:"IssueInstant"`
	NotAfter        string        `json:"NotAfter"`
	Token           string        `json:"Token"`
	DisplayClaims   DisplayClaims `json:"DisplayClaims"`
}

type MinecraftStoreItem struct {
	Name      string `json:"name"`
	Signature string `json:"signature"`
}

type MinecraftStoreResponse struct {
	Items    []MinecraftStoreItem `json:"items"`
	Signature string              `json:"signature"`
	KeyID     string              `json:"keyId"`
}

type MinecraftAuthenticateResponse struct {
	Username    string        `json:"username"`
	Roles       []interface{} `json:"roles"`
	AccessToken string        `json:"access_token"`
	TokenType   string        `json:"token_type"`
	ExpiresIn   int           `json:"expires_in"`
}

type MinecraftProfileInfo struct {
	ID    string `json:"id"`
	State string `json:"state"` // ACTIVE or INACTIVE
	URL   string `json:"url"`
}

type MinecraftProfileSkin struct {
	MinecraftProfileInfo
	Variant string `json:"variant"`
}

type MinecraftProfileCape struct {
	MinecraftProfileInfo
	Alias string `json:"alias"`
}

type MinecraftProfileResponse struct {
	ID            string                `json:"id"`
	Name          string                `json:"name"`
	Skins         []MinecraftProfileSkin `json:"skins"`
	Capes         []MinecraftProfileCape `json:"capes"`
	Error         string                `json:"error"`
	ErrorMessage  string                `json:"errorMessage"`
}

type CompleteLoginResponse struct {
	MinecraftProfileResponse
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type NewsEntryPlayPageImage struct {
	Title string `json:"title"`
	URL   string `json:"url"`
}

type NewsEntryNewsPageImageDimensions struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

type NewsEntryNewsPageImage struct {
	Title      string                          `json:"title"`
	URL        string                          `json:"url"`
	Dimensions NewsEntryNewsPageImageDimensions `json:"dimensions"`
}

type NewsEntry struct {
	Title         string                    `json:"title"`
	Category      string                    `json:"category"`
	Date          string                    `json:"date"`
	Text          string                    `json:"text"`
	PlayPageImage NewsEntryPlayPageImage    `json:"playPageImage"`
	NewsPageImage NewsEntryNewsPageImage    `json:"newsPageImage"`
	ReadMoreLink  string                    `json:"readMoreLink"`
	NewsType      []string                  `json:"newsType"`
	ID            string                    `json:"id"`
}

type MinecraftNews struct {
	Version int         `json:"version"`
	Entries []NewsEntry `json:"entries"`
}

type JavaPatchNoteEntryImage struct {
	URL   string `json:"url"`
	Title string `json:"title"`
}

type JavaPatchNoteEntry struct {
	Title       string                  `json:"title"`
	Type        string                  `json:"type"` // "release" | "snapshot"
	Version     string                  `json:"version"`
	Image       JavaPatchNoteEntryImage `json:"image"`
	Body        string                  `json:"body"`
	ContentPath string                  `json:"contentPath"`
}

type JavaPatchNotes struct {
	Version int                 `json:"version"`
	Entries []JavaPatchNoteEntry `json:"entries"`
}