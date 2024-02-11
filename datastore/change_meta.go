package datastore

import (
	"github.com/PretendoNetwork/nex-go"
	common_globals "github.com/PretendoNetwork/nex-protocols-common-go/globals"
	datastore "github.com/PretendoNetwork/nex-protocols-go/datastore"
	datastore_types "github.com/PretendoNetwork/nex-protocols-go/datastore/types"
)

func changeMeta(err error, packet nex.PacketInterface, callID uint32, param *datastore_types.DataStoreChangeMetaParam) (*nex.RMCMessage, *nex.Error) {
	if commonProtocol.GetObjectInfoByDataID == nil {
		common_globals.Logger.Warning("GetObjectInfoByDataID not defined")
		return nil, nex.NewError(nex.ResultCodes.Core.NotImplemented, "change_error")
	}

	if commonProtocol.UpdateObjectPeriodByDataIDWithPassword == nil {
		common_globals.Logger.Warning("UpdateObjectPeriodByDataIDWithPassword not defined")
		return nil, nex.NewError(nex.ResultCodes.Core.NotImplemented, "change_error")
	}

	if commonProtocol.UpdateObjectMetaBinaryByDataIDWithPassword == nil {
		common_globals.Logger.Warning("UpdateObjectMetaBinaryByDataIDWithPassword not defined")
		return nil, nex.NewError(nex.ResultCodes.Core.NotImplemented, "change_error")
	}

	if commonProtocol.UpdateObjectDataTypeByDataIDWithPassword == nil {
		common_globals.Logger.Warning("UpdateObjectDataTypeByDataIDWithPassword not defined")
		return nil, nex.NewError(nex.ResultCodes.Core.NotImplemented, "change_error")
	}

	if err != nil {
		common_globals.Logger.Error(err.Error())
		return nil, nex.NewError(nex.ResultCodes.DataStore.Unknown, "change_error")
	}

	connection := packet.Sender()
	endpoint := connection.Endpoint()

	metaInfo, errCode := commonProtocol.GetObjectInfoByDataID(param.DataID)
	if errCode != nil {
		return nil, errCode
	}

	// TODO - Is this the right permission?
	errCode = commonProtocol.VerifyObjectPermission(metaInfo.OwnerID, connection.PID(), metaInfo.DelPermission)
	if errCode != nil {
		return nil, errCode
	}

	if param.ModifiesFlag.PAND(0x08) != 0 {
		errCode = commonProtocol.UpdateObjectPeriodByDataIDWithPassword(param.DataID, param.Period, param.UpdatePassword)
		if errCode != nil {
			return nil, errCode
		}
	}

	if param.ModifiesFlag.PAND(0x10) != 0 {
		errCode = commonProtocol.UpdateObjectMetaBinaryByDataIDWithPassword(param.DataID, param.MetaBinary, param.UpdatePassword)
		if errCode != nil {
			return nil, errCode
		}
	}

	if param.ModifiesFlag.PAND(0x80) != 0 {
		errCode = commonProtocol.UpdateObjectDataTypeByDataIDWithPassword(param.DataID, param.DataType, param.UpdatePassword)
		if errCode != nil {
			return nil, errCode
		}
	}

	rmcResponse := nex.NewRMCSuccess(endpoint, nil)
	rmcResponse.ProtocolID = datastore.ProtocolID
	rmcResponse.MethodID = datastore.MethodChangeMeta
	rmcResponse.CallID = callID

	return rmcResponse, nil
}
