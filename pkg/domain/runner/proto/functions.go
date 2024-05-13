package proto

func (ps *TaskParams) Get(key string) string {
	for _, p := range ps.Params {
		if p.Key == key {
			return p.Value
		}
	}
	return ""
}
