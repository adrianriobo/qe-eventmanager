package rules

import (
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

// TODO create separate PR to handle hot config changes
// viper.OnConfigChange(func(e fsnotify.Event) {
// 	fmt.Println("Config file changed:", e.Name)
// 	// err = viper.Unmarshal(&umbConfig)
// 	// send for channel to manage reload
// })
// viper.WatchConfig()
// Extend to multiple files each one rule?
func LoadRules(configFilePath string) (*Rule, error) {
	viper.SetConfigType(filepath.Ext(configFilePath))
	viper.SetConfigName(
		strings.TrimSuffix(
			filepath.Base(configFilePath),
			filepath.Ext(configFilePath)))
	viper.AddConfigPath(filepath.Dir(configFilePath))
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}
	var rule *Rule
	err = viper.Unmarshal(rule)
	if err != nil {
		return nil, err
	}
	return rule, nil
}
