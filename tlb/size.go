package tlb

type Size struct {
}
type Size15 Size

func (u Size15) FixedSize() int {
	return 15
}

type Size16 Size

func (u Size16) FixedSize() int {
	return 16
}

type Size19 Size

func (u Size19) FixedSize() int {
	return 19
}

type KeySize32 Size

func (u KeySize32) FixedSize() int {
	return 32
}

type Size96 Size

func (u Size96) FixedSize() int {
	return 96
}

type Size256 Size

func (u Size256) FixedSize() int {
	return 256
}

type Size264 Size

func (u Size264) FixedSize() int {
	return 264
}

type Size320 Size

func (u Size320) FixedSize() int {
	return 320
}
