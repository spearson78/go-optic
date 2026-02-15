package expr

type HtmlChildren struct {
	OpticTypeExpr
}

func (e HtmlChildren) Short() string {
	return "HtmlChildren"
}

func (e HtmlChildren) String() string {
	return "HtmlChildren"
}

type ParseHtml struct {
	OpticTypeExpr
}

func (e ParseHtml) Short() string {
	return "ParseHtml"
}

func (e ParseHtml) String() string {
	return "ParseHtml"
}
