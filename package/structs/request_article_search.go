package structs

type RequestSearchArticle struct {
	Keyword string `json:"keyword"`
	Page    int    `json:"page"`
	Limit   int    `json:"limit"`
	OrderBy string `json:"order_by"`
}

func (p *RequestSearchArticle) NewQuerySearchArticle() map[string]interface{} {
	// Default pagination values
	if p.Page <= 0 {
		p.Page = 1
	}
	if p.Limit <= 0 {
		p.Limit = 10
	}

	from := (p.Page - 1) * p.Limit

	// Base bool query
	boolQuery := make(map[string]interface{})
	mustQueries := make([]map[string]interface{}, 0)

	// Keyword search
	if p.Keyword != "" {
		mustQueries = append(mustQueries, map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":  p.Keyword,
				"fields": []string{"article.title", "article.body"},
			},
		})
	}

	if len(mustQueries) > 0 {
		boolQuery["must"] = mustQueries
	}

	// Sort default
	sortField := "created_at"
	sortOrder := "desc"
	if p.OrderBy != "" {
		sortOrder = p.OrderBy // Pastikan hanya "asc" atau "desc"
	}

	// Final query
	query := map[string]interface{}{
		"from": from,
		"size": p.Limit,
		"query": map[string]interface{}{
			"bool": boolQuery,
		},
		"sort": []interface{}{
			map[string]interface{}{
				sortField: map[string]interface{}{
					"order":         sortOrder,     // "asc" or "desc"
					"unmapped_type": "date",        // Opsional, hanya jika field tidak selalu ada
				},
			},
		},
	}

	return query
}
