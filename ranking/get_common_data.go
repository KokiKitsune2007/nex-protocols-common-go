package ranking

import (
	"github.com/PretendoNetwork/nex-go"
	"github.com/PretendoNetwork/nex-go/types"
	common_globals "github.com/PretendoNetwork/nex-protocols-common-go/globals"
	ranking "github.com/PretendoNetwork/nex-protocols-go/ranking"
)

func getCommonData(err error, packet nex.PacketInterface, callID uint32, uniqueID *types.PrimitiveU64) (*nex.RMCMessage, *nex.Error) {
	if commonProtocol.GetCommonData == nil {
		common_globals.Logger.Warning("Ranking::GetCommonData missing GetCommonData!")
		return nil, nex.NewError(nex.ResultCodes.Core.NotImplemented, "change_error")
	}

	if err != nil {
		common_globals.Logger.Error(err.Error())
		return nil, nex.NewError(nex.ResultCodes.Ranking.InvalidArgument, "change_error")
	}

	connection := packet.Sender()
	endpoint := connection.Endpoint()

	commonData, err := commonProtocol.GetCommonData(uniqueID)
	if err != nil {
		return nil, nex.NewError(nex.ResultCodes.Ranking.NotFound, "change_error")
	}

	rmcResponseStream := nex.NewByteStreamOut(endpoint.LibraryVersions(), endpoint.ByteStreamSettings())

	commonData.WriteTo(rmcResponseStream)

	rmcResponseBody := rmcResponseStream.Bytes()

	rmcResponse := nex.NewRMCSuccess(endpoint, rmcResponseBody)
	rmcResponse.ProtocolID = ranking.ProtocolID
	rmcResponse.MethodID = ranking.MethodGetCommonData
	rmcResponse.CallID = callID

	return rmcResponse, nil
}
