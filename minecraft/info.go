package minecraft

import "encoding/json"

func GetMinecraftNews() (*MinecraftNews, error) {
	resp, err := getRequestsResponseCache("https://launchercontent.mojang.com/news.json")
	if err != nil {
		return nil, err
	}
	var news MinecraftNews
	if err := json.Unmarshal(resp, &news); err != nil {
		return nil, err
	}

	return &news, nil
}

func GetJavaPatchNotes() (*JavaPatchNotes, error) {
	resp, err := getRequestsResponseCache("https://launchercontent.mojang.com/javaPatchNotes.json")
	if err != nil {
		return nil, err
	}

	var notes JavaPatchNotes
	if err := json.Unmarshal(resp, &notes); err != nil {
		return nil, err
	}

	return &notes, nil
}