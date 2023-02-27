package operate

import "github.com/qiniu/uip/db/inf"

type Operate func(data *inf.IpData)

var DefaultOperates = []Operate{
	ReplaceShortage,
	TrimAsnIspDup,
	MergeNearNetwork,
}
