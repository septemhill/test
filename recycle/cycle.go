package recycle

func NewCycle(grps ...Group) {
	if len(grps) <= 1 {
		return
	}

	for i := 0; i < len(grps); i++ {
		instCh := grps[i].Digested()
		grps[i].Digest(instCh)
	}
}
