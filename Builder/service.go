package Builder

import (
	"fmt"
	"github.com/ARMmaster17/Captain/Shared/framework"
	"github.com/ARMmaster17/Captain/Shared/longjob"
	"github.com/rs/zerolog/log"
)

func NewBuilder() (*framework.Framework, error) {
	builderFramework, err := framework.NewFramework("builder")
	if err != nil {
		log.Error().Err(err).Stack().Msg("unable to initialize framework")
		return nil, fmt.Errorf("unable to create builder service")
	}
	log.Trace().Msg("registering api routes")
	builderFramework.RegisterCommonAPIRoutes()
	log.Trace().Msg("initializing job routes")
	longjob.RegisterLongjobQueue(builderFramework, 1, "plane/build", nil)
	longjob.RegisterLongjobQueue(builderFramework, 1, "plane/destroy", nil)
	return builderFramework, nil
}
