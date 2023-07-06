package lavalinkbot

import (
	"fmt"
	"os"
	"strings"

	"github.com/disgoorg/disgolink/v3/disgolink"
	"github.com/disgoorg/log"
	"github.com/disgoorg/snowflake/v2"
	"gopkg.in/yaml.v3"
)

var defaultConfig = Config{
	Log: LogConfig{
		Level:     log.LevelInfo,
		AddSource: false,
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
	Level     log.Level `yaml:"level"`
	AddSource bool      `yaml:"add_source"`
}

func (c LogConfig) String() string {
	return fmt.Sprintf("\n  Level: %s\n  AddSource: %t\n",
		c.Level,
		c.AddSource,
	)
}

func (c LogConfig) Flags() int {
	flags := log.LstdFlags
	if c.AddSource {
		flags |= log.Llongfile
	}
	return flags
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
	Dependency string `yaml:"dependency"`
	Repository string `yaml:"repository"`
}

func (c PluginConfig) String() string {
	return fmt.Sprintf("\n   Dependency: %s\n   Repository: %s",
		c.Dependency,
		c.Repository,
	)
}
