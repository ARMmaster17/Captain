package main

import (
	"github.com/ARMmaster17/Captain/Shared/framework"
	"github.com/ARMmaster17/Captain/Shared/longjob"
)

func main() {
	builderFramework := framework.NewFramework("builder")
	builderFramework.RegisterCommonAPIRoutes()
	longjob.RegisterLongjobQueue(&builderFramework, 1, "plane/build", nil)
	longjob.RegisterLongjobQueue(&builderFramework, 1, "plane/destroy", nil)
	builderFramework.Start()
}
