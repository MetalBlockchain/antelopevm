package api

import (
	"github.com/MetalBlockchain/antelopevm/chain/time"
)

func init() {
	Functions["current_time"] = currentTime
	Functions["publication_time"] = publicationTime
	Functions["is_feature_activated"] = isFeatureActivated
	Functions["get_sender"] = getSender
}

func currentTime(context Context) interface{} {
	return func() uint64 {
		currentTime := time.Now().TimeSinceEpoch().Count()

		return uint64(currentTime)
	}
}

func publicationTime(context Context) interface{} {
	return func() uint64 {
		published := context.GetTransactionContext().GetPublicationTime().TimeSinceEpoch().Count()

		return uint64(published)
	}
}

func isFeatureActivated(context Context) interface{} {
	return func(ptr uint32) uint32 {
		return 0
	}
}

func getSender(context Context) interface{} {
	return func() uint64 {
		if sender, err := context.GetApplyContext().GetSender(); err == nil {
			if sender == nil {
				return 0
			}

			return uint64(*sender)
		}

		return 0
	}
}
