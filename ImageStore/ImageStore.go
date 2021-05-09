package ImageStore

import (
	"fmt"
	"github.com/spf13/viper"
)

func GetProviderSpecificImageConfiguration(driverYamlTag string, imagetype string) (string, error) {
	if !viper.IsSet(fmt.Sprintf("config.images.%s", imagetype)) {
		return "", fmt.Errorf("image type %s is not configured in config.yaml", imagetype)
	}
	if !viper.IsSet(fmt.Sprintf("config.images.%s.%s", imagetype, driverYamlTag)) {
		return "", fmt.Errorf("image type %s not found from driver %s in config.yaml", imagetype, driverYamlTag)
	}
	return viper.GetString(fmt.Sprintf("config.images.%s.%s", imagetype, driverYamlTag)), nil
}
