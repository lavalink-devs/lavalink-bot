package maven

import (
	"fmt"
	"net/http"
	"testing"
)

func TestMetadata_FetchLatestVersion(t *testing.T) {
	vs := []struct {
		Dependency string
		Repository string
	}{
		{
			"com.github.TopiSenpai.LavaSrc:lavasrc-plugin:x.y.z",
			"https://maven.topi.wtf/releases",
		},
		{
			"com.dunctebot:tts-plugin:VERSION",
			"https://jitpack.io",
		},
		{
			"com.github.TopiSenpai:Sponsorblock-Plugin:x.x.x",
			"https://jitpack.io",
		},
		{
			"com.dunctebot:skybot-lavalink-plugin:VERSION",
			"https://m2.duncte123.dev/releases",
		}, {
			"me.rohank05:lavalink-filter-plugin:x.x.x",
			"https://jitpack.io",
		}, {
			"com.github.esmBot:lava-xm-plugin:vx.x.x",
			"https://jitpack.io",
		},
	}
	for _, v := range vs {
		metadata, err := FetchLatestVersion(http.DefaultClient, v.Dependency, v.Repository)
		if err != nil {
			t.Error(err)
			continue
		}
		fmt.Println(metadata)
	}
}
