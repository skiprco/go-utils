package mapper

import (
	"encoding/json"

	log "github.com/sirupsen/logrus"
)

func MapLoose(source interface{}, dest interface{}) error {
	b, err := json.Marshal(source)
	if err != nil {
		log.WithFields(log.Fields{
			"Method":        "MapLoose",
			"ErrorDetails":  err.Error(),
		}).Errorf("fail to marshal source")
		return err
	}
	err = json.Unmarshal(b, dest)
	if err != nil {
		log.WithFields(log.Fields{
			"Method":        "MapLoose",
			"ErrorDetails":  err.Error(),
		}).Errorf("Fail to unmarshal to dest")
		return err
	}
	return nil
}
