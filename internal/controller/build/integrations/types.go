package integrations

import (
	choreov1 "github.com/choreo-idp/choreo/api/v1"
)

type BuildContext struct {
	Component       *choreov1.Component
	DeploymentTrack *choreov1.DeploymentTrack
	Build           *choreov1.Build
}
