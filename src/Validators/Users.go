package validatores

type Query struct {
	Page        uint64 `form:"user" json:"user" xml:"user"  binding:"required"`
	RowsPerPage uint64 `form:"password" json:"password" xml:"password" binding:"required"`
	Total       uint64
	Table       []interface{}
}
