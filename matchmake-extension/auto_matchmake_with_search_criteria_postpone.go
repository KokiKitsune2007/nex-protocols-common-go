package matchmake_extension

import (
	nex "github.com/PretendoNetwork/nex-go"
	common_globals "github.com/PretendoNetwork/nex-protocols-common-go/globals"
	match_making_types "github.com/PretendoNetwork/nex-protocols-go/match-making/types"
	matchmake_extension "github.com/PretendoNetwork/nex-protocols-go/matchmake-extension"
)

func autoMatchmakeWithSearchCriteria_Postpone(err error, client *nex.Client, callID uint32, lstSearchCriteria []*match_making_types.MatchmakeSessionSearchCriteria, anyGathering *nex.DataHolder, message string) uint32 {
	if commonMatchmakeExtensionProtocol.cleanupMatchmakeSessionSearchCriteriaHandler == nil {
		logger.Warning("MatchmakeExtension::AutoMatchmakeWithSearchCriteria_Postpone missing CleanupMatchmakeSessionSearchCriteriaHandler!")
		return nex.Errors.Core.NotImplemented
	}

	if err != nil {
		logger.Error(err.Error())
		return nex.Errors.Core.InvalidArgument
	}

	commonMatchmakeExtensionProtocol.cleanupMatchmakeSessionSearchCriteriaHandler(lstSearchCriteria)

	server := commonMatchmakeExtensionProtocol.server

	// * A client may disconnect from a session without leaving reliably,
	// * so let's make sure the client is removed from the session
	common_globals.RemoveClientFromAllSessions(client)

	var matchmakeSession *match_making_types.MatchmakeSession
	anyGatheringDataType := anyGathering.TypeName()

	if anyGatheringDataType == "MatchmakeSession" {
		matchmakeSession = anyGathering.ObjectData().(*match_making_types.MatchmakeSession)
	} else {
		logger.Critical("Non-MatchmakeSession DataType?!")
		return nex.Errors.Core.InvalidArgument
	}

	sessions := common_globals.FindSessionsByMatchmakeSessionSearchCriterias(lstSearchCriteria, commonMatchmakeExtensionProtocol.gameSpecificMatchmakeSessionSearchCriteriaChecksHandler)
	var session *common_globals.CommonMatchmakeSession

	if len(sessions) == 0 {
		var errCode uint32
		session, err, errCode = common_globals.CreateSessionBySearchCriteria(matchmakeSession, lstSearchCriteria, client.PID())
		if err != nil {
			logger.Error(err.Error())
			return errCode
		}
	} else {
		session = sessions[0]
	}

	err, errCode := common_globals.AddPlayersToSession(session, []uint32{client.ConnectionID()}, client, message)
	if err != nil {
		logger.Error(err.Error())
		return errCode
	}

	rmcResponseStream := nex.NewStreamOut(server)
	matchmakeDataHolder := nex.NewDataHolder()
	matchmakeDataHolder.SetTypeName("MatchmakeSession")
	matchmakeDataHolder.SetObjectData(session.GameMatchmakeSession)
	rmcResponseStream.WriteDataHolder(matchmakeDataHolder)

	rmcResponseBody := rmcResponseStream.Bytes()

	rmcResponse := nex.NewRMCResponse(matchmake_extension.ProtocolID, callID)
	rmcResponse.SetSuccess(matchmake_extension.MethodAutoMatchmakeWithSearchCriteriaPostpone, rmcResponseBody)

	rmcResponseBytes := rmcResponse.Bytes()

	var responsePacket nex.PacketInterface

	if server.PRUDPVersion() == 0 {
		responsePacket, _ = nex.NewPacketV0(client, nil)
		responsePacket.SetVersion(0)
	} else {
		responsePacket, _ = nex.NewPacketV1(client, nil)
		responsePacket.SetVersion(1)
	}

	responsePacket.SetSource(0xA1)
	responsePacket.SetDestination(0xAF)
	responsePacket.SetType(nex.DataPacket)
	responsePacket.SetPayload(rmcResponseBytes)

	responsePacket.AddFlag(nex.FlagNeedsAck)
	responsePacket.AddFlag(nex.FlagReliable)

	server.Send(responsePacket)

	return 0
}
