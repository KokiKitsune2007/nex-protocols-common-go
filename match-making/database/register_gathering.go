package database

import (
	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	common_globals "github.com/PretendoNetwork/nex-protocols-common-go/v2/globals"
	"github.com/PretendoNetwork/nex-protocols-common-go/v2/match-making/tracking"
	match_making_types "github.com/PretendoNetwork/nex-protocols-go/v2/match-making/types"
)

// RegisterGathering registers a new gathering on the databse. No participants are added
func RegisterGathering(manager *common_globals.MatchmakingManager, pid *types.PID, gathering *match_making_types.Gathering, gatheringType string) (*types.DateTime, *nex.Error) {
	startedTime := types.NewDateTime(0).Now()

	err := manager.Database.QueryRow(`INSERT INTO matchmaking.gatherings (
		owner_pid,
		host_pid,
		min_participants,
		max_participants,
		participation_policy,
		policy_argument,
		flags,
		state,
		description,
		type,
		started_time
	) VALUES (
		$1,
		$2,
		$3,
		$4,
		$5,
		$6,
		$7,
		$8,
		$9,
		$10,
		$11
	) RETURNING id`,
		pid.Value(),
		pid.Value(),
		gathering.MinimumParticipants.Value,
		gathering.MaximumParticipants.Value,
		gathering.ParticipationPolicy.Value,
		gathering.PolicyArgument.Value,
		gathering.Flags.Value,
		gathering.State.Value,
		gathering.Description.Value,
		gatheringType,
		startedTime.Standard(),
	).Scan(&gathering.ID.Value)
	if err != nil {
		return nil, nex.NewError(nex.ResultCodes.Core.Unknown, err.Error())
	}

	nexError := tracking.LogRegisterGathering(manager.Database, pid, gathering.ID.Value)
	if nexError != nil {
		return nil, nexError
	}

	gathering.OwnerPID = pid.Copy().(*types.PID)
	gathering.HostPID = pid.Copy().(*types.PID)

	return startedTime, nil
}
