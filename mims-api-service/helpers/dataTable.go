package helpers

func QueryOffset(offset int) int {
	if offset == 0 {
		offset = 0
	}
	return offset
}

func QueryLimit(limit int) int {
	if limit == 0 {
		limit = 25
	}
	return limit
}

func QueryOrder(column, dir string, columnOrder map[string]string) string {
	orderString := columnOrder[column] + " " + dir
	return orderString
}
