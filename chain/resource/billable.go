package resource

const billableAlignment uint64 = 16

type BillableSize = uint64

func NewBillableSize(value uint64) BillableSize {
	return BillableSize(((value + billableAlignment - 1) / billableAlignment) * billableAlignment)
}
