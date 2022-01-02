package types

//go:generate msgp
type SshieldClearanceToken struct {
	CreateTime    int64
	ClearanceTime int64
	Cleared       bool
}

func (z *SshieldClearanceToken) Reset() {
	z.CreateTime = 0
	z.ClearanceTime = 0
	z.Cleared = false
}
