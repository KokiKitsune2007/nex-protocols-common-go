package matchmake_extension

import (
	"github.com/PretendoNetwork/nex-go"
	"github.com/PretendoNetwork/nex-go/types"
	common_globals "github.com/PretendoNetwork/nex-protocols-common-go/globals"
	matchmake_extension "github.com/PretendoNetwork/nex-protocols-go/matchmake-extension"
)

func modifyCurrentGameAttribute(err error, packet nex.PacketInterface, callID uint32, gid *types.PrimitiveU32, attribIndex *types.PrimitiveU32, newValue *types.PrimitiveU32) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		common_globals.Logger.Error(err.Error())
		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, "change_error")
	}

	connection := packet.Sender().(*nex.PRUDPConnection)
	endpoint := connection.Endpoint().(*nex.PRUDPEndPoint)

	session, ok := common_globals.Sessions[gid.Value]
	if !ok {
		return nil, nex.NewError(nex.ResultCodes.RendezVous.SessionVoid, "change_error")
	}

	if session.GameMatchmakeSession.Gathering.OwnerPID.Equals(connection.PID()) {
		return nil, nex.NewError(nex.ResultCodes.RendezVous.PermissionDenied, "change_error")
	}

	index := int(attribIndex.Value)

	if index > session.GameMatchmakeSession.Attributes.Length() {
		return nil, nex.NewError(nex.ResultCodes.Core.InvalidIndex, "change_error")
	}

	session.GameMatchmakeSession.Attributes.SetIndex(index, newValue.Copy().(*types.PrimitiveU32))

	rmcResponse := nex.NewRMCSuccess(endpoint, nil)
	rmcResponse.ProtocolID = matchmake_extension.ProtocolID
	rmcResponse.MethodID = matchmake_extension.MethodModifyCurrentGameAttribute
	rmcResponse.CallID = callID

	return rmcResponse, nil
}
