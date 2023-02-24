package api

func GetProducerFunctions(context Context) map[string]interface{} {
	functions := make(map[string]interface{})

	functions["get_active_producers"] = getActiveProducers(context)

	return functions
}

func getActiveProducers(context Context) func(uint32, uint32) int32 {
	return func(ptr uint32, length uint32) int32 {
		panic("not implemented")
	}
}
