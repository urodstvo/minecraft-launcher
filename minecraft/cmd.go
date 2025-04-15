package minecraft

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
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
				libstr += filepath.Join(path, "libraries", download.Path) + classpathSeparator
			} else {
				libstr += getLibraryPath(i.Name+"-"+native, path) + classpathSeparator
			}
		}
	}

	if data.Jar != "" {
		libstr += filepath.Join(path, "versions", data.Jar, fmt.Sprintf("%s.jar", data.Jar))
	} else {
		libstr += filepath.Join(path, "versions", data.Id, fmt.Sprintf("%s.jar", data.Id))
	}

	return libstr
}

func replaceArguments(argstr string, versionData ClientJson, path string, options MinecraftOptions, classpath string) string {
	argstr = strings.ReplaceAll(argstr, "${natives_directory}", options.NativesDirectory)
	if options.LauncherName == "" {
		options.LauncherName = "urodstvo-launcher"
	}
	argstr = strings.ReplaceAll(argstr, "${launcher_name}", options.LauncherName)
	if options.LauncherVersion == "" {
		options.LauncherVersion = getLibraryVersion()
	}
	argstr = strings.ReplaceAll(argstr, "${launcher_version}", options.LauncherVersion)
	argstr = strings.ReplaceAll(argstr, "${classpath}", classpath)
	if options.Username == "" {
		options.Username = "{username}"
	}
	argstr = strings.ReplaceAll(argstr, "${auth_player_name}", options.Username)
	argstr = strings.ReplaceAll(argstr, "${version_name}", versionData.Id)
	if options.GameDirectory == "" {
		options.GameDirectory = path
	}
	argstr = strings.ReplaceAll(argstr, "${game_directory}", options.GameDirectory)
	argstr = strings.ReplaceAll(argstr, "${assets_root}", filepath.Join(path,"assets"))
	if versionData.Assets == "" {
		versionData.Assets = versionData.Id
	}
	argstr = strings.ReplaceAll(argstr, "${assets_index_name}", versionData.Assets)
	if options.Uuid == "" {
		options.Uuid = "{uuid}"
	}
	argstr = strings.ReplaceAll(argstr, "${auth_uuid}", options.Uuid)
	if options.Token == "" {
		options.Token = "{token}"
	}
	argstr = strings.ReplaceAll(argstr, "${auth_access_token}", options.Token)
	argstr = strings.ReplaceAll(argstr, "${user_type}", "msa")
	argstr = strings.ReplaceAll(argstr, "${version_type}", versionData.Type)
	argstr = strings.ReplaceAll(argstr, "${user_properties}", "{}")
	if options.ResolutionWidth == "" {
		options.ResolutionWidth = "900"
	}
	argstr = strings.ReplaceAll(argstr, "${resolution_width}", options.ResolutionWidth)
	if options.ResolutionHeight == "" {
		options.ResolutionHeight = "600"
	}
	argstr = strings.ReplaceAll(argstr, "${resolution_height}", options.ResolutionHeight)
	argstr = strings.ReplaceAll(argstr, "${game_assets}", filepath.Join(path,"assets","virtual","legacy"))
	argstr = strings.ReplaceAll(argstr, "${auth_session}", options.Token)
	argstr = strings.ReplaceAll(argstr, "${library_directory}", filepath.Join(path, "libraries"))
	argstr = strings.ReplaceAll(argstr, "${classpath_separator}", getClasspathSeparator())
	if options.QuickPlayPath != nil {
		argstr = strings.ReplaceAll(argstr, "${quickPlayPath}", *options.QuickPlayPath)
	} else {
		argstr = strings.ReplaceAll(argstr, "${quickPlayPath}", "{quickPlayPath}")
	}
	if options.QuickPlaySingleplayer != nil {
		argstr = strings.ReplaceAll(argstr, "${quickPlaySingleplayer}", *options.QuickPlaySingleplayer)
	} else {
		argstr = strings.ReplaceAll(argstr, "${quickPlaySingleplayer}", "{quickPlaySingleplayer}")
	}
	if options.QuickPlayMultiplayer != nil {
		argstr = strings.ReplaceAll(argstr, "${quickPlayMultiplayer}", *options.QuickPlayMultiplayer)
	} else {
		argstr = strings.ReplaceAll(argstr, "${quickPlayMultiplayer}", "{quickPlayMultiplayer}")
	}
	if options.QuickPlayRealms != nil {
		argstr = strings.ReplaceAll(argstr, "${quickPlayRealms}", *options.QuickPlayRealms)
	} else {
		argstr = strings.ReplaceAll(argstr, "${quickPlayRealms}", "{quickPlayRealms}")
	}

	return argstr
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

func convertRulesToClientJsonRules(rules []any) ([]ClientJsonRule, error) {
	var clientRules []ClientJsonRule
	for _, rule := range rules {
		if ruleMap, ok := rule.(map[string]interface{}); ok {
			var clientRule ClientJsonRule
			if action, ok := ruleMap["action"].(string); ok {
				clientRule.Action = action
			}
			if osMap, ok := ruleMap["os"].(map[string]interface{}); ok {
				if name, ok := osMap["name"].(string); ok {
					clientRule.Os.Name = &name
				}
				if arch, ok := osMap["arch"].(string); ok {
					clientRule.Os.Arch = &arch
				}
				if version, ok := osMap["version"].(string); ok {
					clientRule.Os.Version = &version
				}
			}
			clientRules = append(clientRules, clientRule)
		} else {
			return nil, fmt.Errorf("invalid rule type, expected map[string]interface{} but got %s", reflect.TypeOf(rule))
		}
	}
	return clientRules, nil
}

func getArguments(data []any, versionData ClientJson, path string, options MinecraftOptions, classpath string) []string {
	var arglist []string

	for _, i := range data {
		switch v := i.(type) {
		case string:
			arglist = append(arglist, replaceArguments(v, versionData, path, options, classpath))
		case map[string]any:
			if rules, hasRules := v["compatibilityRules"].([]any); hasRules {
				clientRules, err := convertRulesToClientJsonRules(rules)
				if err != nil {
					continue
				}
				if !parseRuleList(clientRules, &options) {
					continue
				}
			}

			if rules, hasRules := v["rules"].([]any); hasRules {
				clientRules, err := convertRulesToClientJsonRules(rules)
				if err != nil {
					continue
				}
				if !parseRuleList(clientRules, &options) {
					continue
				}
			}

			if value, ok := v["value"].(string); ok {
				arglist = append(arglist, replaceArguments(value, versionData, path, options, classpath))
			} else if valueList, ok := v["value"].([]string); ok {
				for _, val := range valueList {
					arglist = append(arglist, replaceArguments(val, versionData, path, options, classpath))					
				}
			}
		}
	}

	return arglist
}

func (m *minecraftConfig) GetMinecraftCommand(version string, options MinecraftOptions) ([]string, error) {
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

	if data.InheritsFrom != "" {
		data, err = inheritJson(data, m.Config.Directory)
		if err != nil {
			return nil, err
		}
	}

	if options.NativesDirectory == "" {
		options.NativesDirectory = filepath.Join(path, "versions", data.Id, "natives")
	}

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

	if len(options.JvmArguments) > 0 {
		command = append(command, options.JvmArguments...)
	}

	if data.Arguments != nil {
		if data.Arguments.Jvm != nil {
			command = append(command, getArguments(data.Arguments.Jvm, data, path, options, classpath)...)
		} else {
			command = append(command, "-Djava.library.path="+options.NativesDirectory, "-cp", classpath)
		}
	}

	if options.EnableLoggingConfig {
		if data.Logging != nil {
			loggerFile := filepath.Join(path, "assets", "log_configs", data.Logging.Client.File.Id)
			command = append(command, strings.Replace(data.Logging.Client.Argument, "${path}", loggerFile, -1))
		}
	}

	command = append(command, data.MainClass)

	if data.MinecraftArguments != "" {
		command = append(command, getArgumentsString(data, path, options, classpath)...)
	} else {
		command = append(command, getArguments(data.Arguments.Game, data, path, options, classpath)...)
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