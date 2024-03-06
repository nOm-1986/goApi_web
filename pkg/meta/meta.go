package meta

type Meta struct {
	TotalCount int `json:"total_count"`
}

func NewMeta(total int) (*Meta, error) {
	return &Meta{
		TotalCount: total,
	}, nil
}
