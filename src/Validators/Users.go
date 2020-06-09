package validator

type Query struct {
	Page        int64 `form:"user" json:"user" xml:"user"  binding:"required"`
	RowsPerPage int64 `form:"password" json:"password" xml:"password" binding:"required"`
}
