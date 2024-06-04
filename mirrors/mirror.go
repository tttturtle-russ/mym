package mirrors

import "Mym/consts"

type MirrorSite struct {
	Url      string
	Location string
	Protocol []consts.Protocol
}

type Mirror interface {
	Timeout() ([]MirrorSite, []float64)
	List() []MirrorSite
	WriteToConfig() error
	IsDeployed() bool
}

type Option func(Mirror)
