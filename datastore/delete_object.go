package datastore

import (
	"github.com/PretendoNetwork/nex-go"
	common_globals "github.com/PretendoNetwork/nex-protocols-common-go/globals"
	datastore "github.com/PretendoNetwork/nex-protocols-go/datastore"
	datastore_types "github.com/PretendoNetwork/nex-protocols-go/datastore/types"
)

func deleteObject(err error, packet nex.PacketInterface, callID uint32, param *datastore_types.DataStoreDeleteParam) (*nex.RMCMessage, uint32) {
	if commonDataStoreProtocol.getObjectInfoByDataIDHandler == nil {
		common_globals.Logger.Warning("GetObjectInfoByDataID not defined")
		return nil, nex.Errors.Core.NotImplemented
	}

	if commonDataStoreProtocol.deleteObjectByDataIDWithPasswordHandler == nil {
		common_globals.Logger.Warning("DeleteObjectByDataIDWithPassword not defined")
		return nil, nex.Errors.Core.NotImplemented
	}

	if err != nil {
		common_globals.Logger.Error(err.Error())
		return nil, nex.Errors.DataStore.Unknown
	}

	client := packet.Sender()

	metaInfo, errCode := commonDataStoreProtocol.getObjectInfoByDataIDHandler(param.DataID)
	if errCode != 0 {
		return nil, errCode
	}

	errCode = commonDataStoreProtocol.VerifyObjectPermission(metaInfo.OwnerID, client.PID(), metaInfo.DelPermission)
	if errCode != 0 {
		return nil, errCode
	}

	errCode = commonDataStoreProtocol.deleteObjectByDataIDWithPasswordHandler(param.DataID, param.UpdatePassword)
	if errCode != 0 {
		return nil, errCode
	}

	rmcResponse := nex.NewRMCSuccess(nil)
	rmcResponse.ProtocolID = datastore.ProtocolID
	rmcResponse.MethodID = datastore.MethodDeleteObject
	rmcResponse.CallID = callID

	return rmcResponse, 0
}
