package functions

// SearchFields cria a where com os campos passados
func SearchFields(search string, fields []string) string {
	query := ""
	if search != "" {
		query = " WHERE "
		for i, field := range fields {
			if len(fields) != i+1 {
				query += " upper(" + field + ") like upper('%" + search + "%') OR"
			} else {
				query += " upper(" + field + ") like upper('%" + search + "%') "
			}
		}
	}
	return query
}

// SelectFields cria o selec com os campos passados
func SelectFields(fields []string) string {
	query := "*"
	if len(fields) != 0 {
		query = " "
		for i, field := range fields {
			if len(fields) != i+1 {
				query += field + ","
			} else {
				query += field
			}
		}
	}
	return query
}
