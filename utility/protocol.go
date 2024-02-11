package utility

import (
	"github.com/PretendoNetwork/nex-go"
	utility "github.com/PretendoNetwork/nex-protocols-go/utility"
)

type CommonProtocol struct {
	endpoint            nex.EndpointInterface
	protocol            utility.Interface
	GenerateNEXUniqueID func() uint64
}

// NewCommonProtocol returns a new CommonProtocol
func NewCommonProtocol(protocol utility.Interface) *CommonProtocol {
	commonProtocol := &CommonProtocol{
		endpoint: protocol.Endpoint(),
		protocol: protocol,
	}

	protocol.SetHandlerAcquireNexUniqueID(commonProtocol.acquireNexUniqueID)

	return commonProtocol
}
