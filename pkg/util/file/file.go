package file

import (
	"path/filepath"
	"strings"

	"github.com/adrianriobo/qe-eventmanager/pkg/util/logging"
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
func LoadFileAsStruct(filePath string, structured interface{}) error {
	viper.SetConfigType(strings.ReplaceAll(filepath.Ext(filePath), ".", ""))
	logging.Debugf("%v", strings.ReplaceAll(filepath.Ext(filePath), ".", ""))
	viper.SetConfigName(
		strings.TrimSuffix(
			filepath.Base(filePath),
			filepath.Ext(filePath)))
	logging.Debugf("%v", strings.TrimSuffix(
		filepath.Base(filePath),
		filepath.Ext(filePath)))
	viper.AddConfigPath(filepath.Dir(filePath))
	logging.Debugf("%v", filepath.Dir(filePath))
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}
	err = viper.Unmarshal(structured)
	if err != nil {
		return err
	}
	return nil
}
