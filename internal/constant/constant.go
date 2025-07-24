package constant

type CtxKey string

const(
	NewRelicTransactionCtx CtxKey = "newRelicTransaction"

	Drafted   = 1
	Published = 2
)

var (
	ArticleStatus = map[int]string{
		Drafted:   "drafted",
		Published: "published",
	}
)