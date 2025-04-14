package minecraft

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

func getLibraries(data ClientJson, path string) string {
	classpathSeparator := getClasspathSeparator()
	var libstr string

	for _, i := range data.Libraries {
		if len(i.Rules) > 0 && !parseRuleList(i.Rules, nil) {
			continue
		}

		libstr += getLibraryPath(i.Name, path) + classpathSeparator
		native := getNatives(i)
		if native != "" {
			if download, exists := i.Downloads.Classifiers[native]; exists {
				libstr += fmt.Sprintf("%s/libraries/%s", path, download.Path) + classpathSeparator
			} else {
				libstr += getLibraryPath(i.Name+"-"+native, path) + classpathSeparator
			}
		}
	}

	if data.Jar != "" {
		libstr += fmt.Sprintf("%s/versions/%s/%s.jar", path, data.Jar, data.Jar)
	} else {
		libstr += fmt.Sprintf("%s/versions/%s/%s.jar", path, data.Id, data.Id)
	}

	return libstr
}

func replaceArguments(argstr string, versionData ClientJson, path string, options MinecraftOptions, classpath string) string {
	argstr = strings.ReplaceAll(argstr, "${natives_directory}", options.NativesDirectory)
	argstr = strings.ReplaceAll(argstr, "${launcher_name}", options.LauncherName)
	if options.LauncherVersion == "" {
		options.LauncherVersion = getLibraryVersion()
	}
	argstr = strings.ReplaceAll(argstr, "${launcher_version}", options.LauncherVersion)
	argstr = strings.ReplaceAll(argstr, "${classpath}", classpath)
	argstr = strings.ReplaceAll(argstr, "${auth_player_name}", options.Username)
	argstr = strings.ReplaceAll(argstr, "${version_name}", versionData.Id)
	argstr = strings.ReplaceAll(argstr, "${game_directory}", options.GameDirectory)
	argstr = strings.ReplaceAll(argstr, "${assets_root}", path+"/assets")
	argstr = strings.ReplaceAll(argstr, "${assets_index_name}", versionData.Assets)
	if versionData.Assets == "" {
		versionData.Assets = versionData.Id
	}
	argstr = strings.ReplaceAll(argstr, "${auth_uuid}", options.Uuid)
	argstr = strings.ReplaceAll(argstr, "${auth_access_token}", options.Token)
	argstr = strings.ReplaceAll(argstr, "${user_type}", "msa")
	argstr = strings.ReplaceAll(argstr, "${version_type}", versionData.Type)
	argstr = strings.ReplaceAll(argstr, "${user_properties}", "{}")
	argstr = strings.ReplaceAll(argstr, "${resolution_width}", options.ResolutionWidth)
	argstr = strings.ReplaceAll(argstr, "${resolution_height}", options.ResolutionHeight)
	argstr = strings.ReplaceAll(argstr, "${game_assets}", path+"/assets/virtual/legacy")
	argstr = strings.ReplaceAll(argstr, "${auth_session}", options.Token)
	argstr = strings.ReplaceAll(argstr, "${library_directory}", path+"/libraries")
	argstr = strings.ReplaceAll(argstr, "${classpath_separator}", getClasspathSeparator())
	argstr = strings.ReplaceAll(argstr, "${quickPlayPath}", *options.QuickPlayPath)
	argstr = strings.ReplaceAll(argstr, "${quickPlaySingleplayer}", *options.QuickPlaySingleplayer)
	argstr = strings.ReplaceAll(argstr, "${quickPlayMultiplayer}", *options.QuickPlayMultiplayer)
	argstr = strings.ReplaceAll(argstr, "${quickPlayRealms}", *options.QuickPlayRealms)

	return argstr
}

func generateTestOptions() MinecraftOptions {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	
	username := fmt.Sprintf("Player%d", rand.Intn(900)+100) 
	uuidValue := uuid.New().String()                        

	return MinecraftOptions{
		Username: username,
		Uuid:     uuidValue,
		Token:    "", 
	}
}

func getArgumentsString(versionData ClientJson, path string, options MinecraftOptions, classpath string) []string {
	arglist := []string{}

	args := strings.Split(versionData.MinecraftArguments, " ")

	for _, v := range args {
		v = replaceArguments(v, versionData, path, options, classpath)
		arglist = append(arglist, v)
	}

	if options.CustomResolution {
		arglist = append(arglist, "--width", options.ResolutionWidth, "--height", options.ResolutionHeight)
	}

	if options.Demo {
		arglist = append(arglist, "--demo")
	}

	return arglist
}

func getArguments(data []interface{}, versionData ClientJson, path string, options MinecraftOptions, classpath string) []string {
	var arglist []string

	for _, i := range data {
		switch v := i.(type) {
		case string:
			arglist = append(arglist, replaceArguments(v, versionData, path, options, classpath))
		case ClientJsonArgumentRule:
			if !parseRuleList(v.CompatibilityRules, &options) {
				continue
			}

			if !parseRuleList(v.Rules, &options) {
				continue
			}

			if value, ok := v.Value.(string); ok {
				arglist = append(arglist, replaceArguments(value, versionData, path, options, classpath))
			} else if valueList, ok := v.Value.([]string); ok {
				for _, v := range valueList {
					arglist = append(arglist, replaceArguments(v, versionData, path, options, classpath))
				}
			}
		}
	}

	return arglist
}

func (m *Minecraft) GetMinecraftCommand(version string, options MinecraftOptions) ([]string, error) {
	path := m.Config.Directory

	versionDir := filepath.Join(path, "versions", version)
	if _, err := os.Stat(versionDir); os.IsNotExist(err) {
		return nil, errors.New("version not found")
	}

	filePath := filepath.Join(path, "versions", version, version + ".json")
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var data ClientJson
	if err := json.NewDecoder(file).Decode(&data); err != nil {
		return nil, err
	}

	// Обработка inherit
	// if inheritFrom, ok := data.InheritsFrom; ok {
		// Вы можете добавить логику для наследования (если необходимо)
		// Иначе можно пропустить
		// data = inheritJson(data, path)
	// }

	classpath := getLibraries(data, path)
	command := []string{}

	if options.ExecutablePath != "" {
		command = append(command, options.ExecutablePath)
	} else if data.JavaVersion.Component != "" {
		javaPath := getExecutablePath(data.JavaVersion.Component, path)
		if javaPath == "" {
			command = append(command, "java")
		} else {
			command = append(command, javaPath)
		}
	} else {
		command = append(command, options.DefaultExecutablePath)
	}

	command = append(command, options.JvmArguments...)

	if data.Arguments != nil {
		for _, arg := range data.Arguments.Jvm {
			if v, ok := arg.(string); ok {
				command = append(command, v)
			}
		}
	} else {
		command = append(command, "-Djava.library.path=" + options.NativesDirectory)
		command = append(command, "-cp", classpath)
	}

	if options.EnableLoggingConfig {
		loggerFile := filepath.Join(path, "assets", "log_configs", data.Logging.Client.File.Id)
		command = append(command, strings.ReplaceAll(data.Logging.Client.Argument, "${path}", loggerFile))
	}

	command = append(command, data.MainClass)

	if data.Arguments != nil {
		for _, arg := range data.Arguments.Game {
			if v, ok := arg.(string); ok {
				command = append(command, v)
			}
		}
	}

	if options.Server != "" {
		command = append(command, "--server", options.Server)
		if options.Port != "" {
			command = append(command, "--port", options.Port)
		}
	}

	if options.DisableMultiplayer {
		command = append(command, "--disableMultiplayer")
	}
	if options.DisableChat {
		command = append(command, "--disableChat")
	}

	return command, nil
}