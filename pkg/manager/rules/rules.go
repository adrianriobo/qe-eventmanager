package rules

import (
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
func LoadRules(filename, path string) (*Rule, error) {
	viper.SetConfigType("yaml")
	viper.SetConfigName(filename)
	viper.AddConfigPath(path)
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
