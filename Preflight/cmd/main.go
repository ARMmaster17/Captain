package main

import (
	"github.com/ARMmaster17/Captain/Shared/framework"
	"github.com/ARMmaster17/Captain/Shared/longjob"
)

func main() {
	preflightFramework := framework.NewFramework("preflight")
	preflightFramework.RegisterCommonAPIRoutes()
	longjob.RegisterLongjobQueue(&preflightFramework, 1, "plane/provision", nil)
	preflightFramework.Start()
}
