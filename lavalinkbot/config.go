package lavalinkbot

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/disgoorg/disgolink/v3/disgolink"
	"github.com/disgoorg/snowflake/v2"
	"gopkg.in/yaml.v3"
)

var defaultConfig = Config{
	Log: LogConfig{
		Level:  slog.LevelInfo,
		Format: "text",
	},
}

func ReadConfig(path string) (Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return Config{}, fmt.Errorf("failed to open config: %w", err)
	}
	defer file.Close()

	cfg := defaultConfig
	if err = yaml.NewDecoder(file).Decode(&cfg); err != nil {
		return Config{}, fmt.Errorf("failed to decode config: %w", err)
	}
	return cfg, nil
}

type Config struct {
	Log     LogConfig     `yaml:"log"`
	Bot     BotConfig     `yaml:"bot"`
	GitHub  GitHubConfig  `yaml:"github"`
	Nodes   NodeConfigs   `yaml:"nodes"`
	Plugins PluginConfigs `yaml:"plugins"`
}

func (c Config) String() string {
	return fmt.Sprintf("\n Log: %s\n Bot: %s\n Nodes: %s\n Plugins: %s\n",
		c.Log,
		c.Bot,
		c.Nodes,
		c.Plugins,
	)
}

type LogConfig struct {
	Level     slog.Level `yaml:"level"`
	Format    string     `yaml:"format"`
	AddSource bool       `yaml:"add_source"`
	NoColor   bool       `yaml:"no_color"`
}

func (c LogConfig) String() string {
	return fmt.Sprintf("\n  Level: %s\n  Format: %s\n  AddSource: %t\n  NoColor: %t\n",
		c.Level,
		c.Format,
		c.AddSource,
		c.NoColor,
	)
}

type BotConfig struct {
	Token    string         `yaml:"token"`
	GuildIDs []snowflake.ID `yaml:"guild_ids"`
}

func (c BotConfig) String() string {
	return fmt.Sprintf("\n  Token: %s\n  GuildIDs: %s",
		c.Token,
		c.GuildIDs,
	)
}

type GitHubConfig struct {
	ServerAddr    string                         `yaml:"server_addr"`
	WebhookSecret string                         `yaml:"webhook_secret"`
	Releases      map[string]GithubReleaseConfig `yaml:"releases"`
}

func (c GitHubConfig) String() string {
	var s string
	for repo, cfg := range c.Releases {
		s += fmt.Sprintf("\n %s: %s", repo, cfg)
	}
	return fmt.Sprintf("\n  ServerAddr: %s\n  WebhookSecret: %s\n  Releases: %s",
		c.ServerAddr,
		c.WebhookSecret,
		s,
	)
}

type GithubReleaseConfig struct {
	WebhookID    snowflake.ID `yaml:"webhook_id"`
	WebhookToken string       `yaml:"webhook_token"`
	PingRole     snowflake.ID `yaml:"ping_role"`
}

func (c GithubReleaseConfig) String() string {
	return fmt.Sprintf("\n  WebhookID: %s\n  WebhookToken: %s\n  PingRole: %s",
		c.WebhookID,
		c.WebhookToken,
		c.PingRole,
	)
}

type NodeConfig struct {
	Name      string `yaml:"name"`
	Address   string `yaml:"address"`
	Password  string `yaml:"password"`
	Secure    bool   `yaml:"secure"`
	SessionID string `yaml:"session_id"`
}

func (c NodeConfig) String() string {
	return fmt.Sprintf("\n   Name: %s\n   Address: %s\n   Password: %s\n   Secure: %t\n   SessionID: %s\n",
		c.Name,
		c.Address,
		strings.Repeat("*", len(c.Password)),
		c.Secure,
		c.SessionID,
	)
}

func (c NodeConfig) ToNodeConfig() disgolink.NodeConfig {
	return disgolink.NodeConfig{
		Name:      c.Name,
		Address:   c.Address,
		Password:  c.Password,
		Secure:    c.Secure,
		SessionID: c.SessionID,
	}
}

type NodeConfigs []NodeConfig

func (c NodeConfigs) String() string {
	var str string
	for _, node := range c {
		str += node.String()
	}
	return str
}

type PluginConfigs []PluginConfig

func (c PluginConfigs) String() string {
	var str string
	for _, plugin := range c {
		str += plugin.String()
	}
	return str
}

type PluginConfig struct {
	Name       string `yaml:"name"`
	Dependency string `yaml:"dependency"`
	Repository string `yaml:"repository"`
	Git        string `yaml:"git"`
}

func (c PluginConfig) String() string {
	return fmt.Sprintf("\n   Name: %s\n   Dependency: %s\n   Repository: %s\n   Git: %s\n",
		c.Name,
		c.Dependency,
		c.Repository,
		c.Git,
	)
}
