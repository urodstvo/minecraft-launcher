package minecraft

import (
	"time"
)


type ClientJsonRule struct{
	Action 	string `json:"action"`
	Os 		struct {
		Name string `json:"name"`
		Arch string `json:"arch"`
		Version string `json:"version"`	
	} `json:"os"`
	Features *struct {
		HasCustomResolution 	bool `json:"has_Custom_Resolution"`
		IsDemoUser 				bool `json:"is_demo_user"`
		HasQuickPlaysSupport 	bool `json:"has_quick_plays_support"`
		IsQuickPlaySingleplayer bool `json:"is_quick_play_singleplayer"`
		IsQuickPlayMultiplayer 	bool `json:"is_quick_play_multiplayer"`
		IsQuickPlayRealms 		bool `json:"is_quick_play_realms"`
	} `json:"features"`
}

type ClientJsonArgumentRule struct{
	CompatibilityRules 	[]ClientJsonRule 	`json:"-"`
	Rules 				[]ClientJsonRule 	`json:"rules"`
	Value 				any 				`json:"value"`
}

type _ClientJsonAssetIndex struct{
	Id 			string 	`json:"id"`
	Sha1 		string 	`json:"sha1"`
	Size 		int 	`json:"size"`
	TotalSize 	int 	`json:"totalSize"`
	Url 		string 	`json:"url"`
}

type _ClientJsonDownloads struct{
	Sha1 	string 	`json:"sha1"`
	Size 	int 	`json:"size"`
	Url 	string 	`json:"url"`
}

type _ClientJsonJavaVersion struct{
	Component		string 	`json:"component"`
	MajorVersion 	int 	`json:"majorVersion"`
}

type _ClientJsonLibraryDownloadsArtifact struct{
	Path 	string 	`json:"path"`
	Url 	string 	`json:"url"`
	Sha1 	string	`json:"sha1"`
	Size 	int 	`json:"size"`
}

type _ClientJsonLibraryDownloads struct{
	Artifact _ClientJsonLibraryDownloadsArtifact `json:"artifact"`
	Classifiers map[string]_ClientJsonLibraryDownloadsArtifact
}

type ClientJsonLibrary struct{
	Name 		string `json:"name"`
	Downloads 	_ClientJsonLibraryDownloads `json:"downloads"`
	Extract struct {
		Exclude []string `json:"-"`
	} `json:"-"`
	Rules 	[]ClientJsonRule `json:"rules"`
	Natives struct{
		Linux 	string `json:"linux"`
		Osx 	string `json:"osx"`
		Windows string `json:"windows"`
	} `json:"natives"`
	Url string `json:"url"`
}

type _ClientJsonLoggingFile struct{
	Id 		string 	`json:"id"`
	Sha1 	string 	`json:"sha1"`
	Size 	int 	`json:"size"`
	Url 	string 	`json:"url"`
}

type _ClientJsonLogging struct{
	Argument 	string 					`json:"argument"`
	File 		_ClientJsonLoggingFile 	`json:"file"`
	Type 		string 					`json:"type"`
}

type ClientJson struct{
	Id 		string `json:"id"`
	Jar 	string `json:"jar"`
	Arguments *struct {
		Game []any `json:"game"`
		Jvm []any `json:"jvm"`
	} `json:"arguments"`
	MinecraftArguments string `json:"-"`
	AssetIndex 		*_ClientJsonAssetIndex `json:"assetIndex"`
	Assets 			string `json:"assets"`
	Downloads struct {
		Client 			_ClientJsonDownloads `json:"client"`
		ClientMappings 	_ClientJsonDownloads `json:"client_mappings"`
		Server 			_ClientJsonDownloads `json:"server"`
		ServerMappings 	_ClientJsonDownloads `json:"server_mappings"`
	} `json:"downloads"`
	JavaVersion 	_ClientJsonJavaVersion `json:"javaVersion"`
	Libraries 		[]ClientJsonLibrary `json:"libraries"`
	Logging struct {
		Client _ClientJsonLogging `json:"client"`
	} `json:"logging"`
	MainClass 				string 	`json:"mainClass"`
	MinimumLauncherVersion 	int 	`json:"minimumLauncherVersion"`
	ReleaseTime 			string 	`json:"releaseTime"`
	Time 					string 	`json:"time"`
	Type 					string 	`json:"type"`
	ComplianceLevel 		int 	`json:"complianceLevel"`
	InheritsFrom 			string 	`json:"-"`
}

type _VersionListManifestJsonVersion struct {
	Id 				string `json:"id"`
	Type 			string `json:"type"`
	Url 			string `json:"url"`
	Time 			time.Time `json:"time"`
	ReleaseTime 	time.Time `json:"releaseTime"`
	Sha1 			string `json:"sha1"`
	ComplianceLevel uint `json:"complianceLevel"`
}

type VersionListManifestJson struct{
	Latest struct {
		Release 	string `json:"release"`
		Snapshot 	string `json:"snapshot"`
	} `json:"latest"`
	Versions []_VersionListManifestJsonVersion `json:"versions"`
}

type _AssetsJsonObject struct{
	Hash string `json:"hash"`
	Size int `json:"size"`
}

type AssetsJson struct {
	Objects map[string]_AssetsJsonObject `json:"objects"`
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

type _RuntimeListJsonEntryManifest struct {
	SHA1 string `json:"sha1"`
    Size int `json:"size"`
    Url string `json:"url"`
}


type _RuntimeListJsonEntry struct {
	Availability struct {
		Group 		int	`json:"group"`
		Progress 	int	 `json:"progress"`
	} `json:"availability"`
    Manifest _RuntimeListJsonEntryManifest  `json:"manifest"`
    Version struct {
		Name 		string  `json:"name"`
		Released 	string  `json:"released"`
	}  `json:"version"`
}


type RuntimeListJson map[string]map[string][]_RuntimeListJsonEntry


type _PlatformManifestJsonFileDownloads struct {
	SHA1 string `json:"sha1"`
    Size int `json:"size"`
    Url string `json:"url"`
}


type _PlatformManifestJsonFile struct{
	Downloads map[string] _PlatformManifestJsonFileDownloads `json:"downloads"`
    Type string `json:"type"`
    Executable bool `json:"executable"`
    Target string `json:"target"`
}


type PlatformManifestJson struct{
	Files map[string]_PlatformManifestJsonFile
}