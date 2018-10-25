package rmgo

//Keeper stores Filter in memory
type Keeper interface {
	Insert(*Condition, Filter)
	Select(Message) []Filter
	Remove(*Condition)
	Destroy()
}

type MessageKeeperFactory struct {
}

func (f MessageKeeperFactory) Create() Keeper {
	return &MessageKeeper{
		cfMap: map[*Condition]Filter{},
	}
}

//MessageKeeper is a Keeper implementation.
type MessageKeeper struct {
	cfMap map[*Condition]Filter
}

//Insert condition to keeper
func (dk *MessageKeeper) Insert(conditionPtr *Condition, filter Filter) {
	dk.cfMap[conditionPtr] = filter
}

//Select condition by Message
func (dk *MessageKeeper) Select(msg Message) []Filter {
	result := []Filter{}
	for k, filter := range dk.cfMap {
		condition := *k
		//append
		if condition.Match(msg) {
			result = append(result, filter)
		}
	}
	return result
}

func (dk *MessageKeeper) Remove(conditionPtr *Condition) {
	delete(dk.cfMap, conditionPtr)
}

func (dk *MessageKeeper) Destroy() {
	m := dk.cfMap
	for k := range m {
		delete(m, k)
	}
}
