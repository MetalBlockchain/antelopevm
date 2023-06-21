package producer

//go:generate msgp
type ProducerSchedule struct {
	Version   uint32
	Producers []ProducerKey
}

func (p ProducerSchedule) Equal(other ProducerSchedule) bool {
	if p.Version != other.Version {
		return false
	}

	if len(p.Producers) != len(other.Producers) {
		return false
	}

	for i := 0; i < len(p.Producers); i++ {
		if !p.Producers[i].Equal(other.Producers[i]) {
			return false
		}
	}

	return true
}
