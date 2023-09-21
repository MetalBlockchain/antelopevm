package api

func init() {
	Functions["get_active_producers"] = getActiveProducers
}

func getActiveProducers(context Context) interface{} {
	return func(ptr uint32, length uint32) int32 {
		panic("not implemented")
	}
}
