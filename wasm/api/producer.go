package api

func init() {
	Functions["get_active_producers"] = getActiveProducers
}

func GetProducerFunctions(context Context) map[string]interface{} {
	functions := make(map[string]interface{})

	functions["get_active_producers"] = getActiveProducers(context)

	return functions
}

func getActiveProducers(context Context) interface{} {
	return func(ptr uint32, length uint32) int32 {
		panic("not implemented")
	}
}
