package repos

// ListParams represents pagination and sorting options for list queries.
type ListParams struct {
	Page      int
	Limit     int
	SortField string
	SortDir   string
}

// Offset calculates SQL offset.
func (p ListParams) Offset() int {
	if p.Page <= 1 {
		return 0
	}
	return (p.Page - 1) * p.Limit
}

// ValidateSort ensures provided sort options are allowed and returns ORDER BY clause.
func (p ListParams) ValidateSort(allowed map[string]string, defaultField string) (string, error) {
	if len(allowed) == 0 {
		return "", nil
	}

	field := p.SortField
	if field == "" {
		field = defaultField
	}

	column, ok := allowed[field]
	if field != "" && !ok {
		return "", ErrInvalidSortField
	}

	direction := "ASC"
	if p.SortDir != "" {
		switch p.SortDir {
		case "asc", "ASC":
		case "desc", "DESC":
			direction = "DESC"
		default:
			return "", ErrInvalidSortDirection
		}
	}

	return column + " " + direction, nil
}
