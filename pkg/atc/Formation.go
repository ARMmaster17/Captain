package atc

type Formation struct {
	ID			int
	BaseName	string
	Domain		string
	TargetCount	int
}

func (f *Formation) GetActualCount() (int, error) {
	// TODO: Active checks or check DB for last health ping
	return 0, nil
}