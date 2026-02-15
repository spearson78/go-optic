package data

import . "github.com/spearson78/go-optic"

func (s *lBlogPost[I, S, T, RET, RW, DIR, ERR]) MeanRating() Optic[Void, S, S, int, int, ReturnMany, ReadOnly, UniDir, CompositionTree[ERR, Pure]] {
	return Reduce(
		s.Ratings().Traverse().Stars(),
		Mean[int](),
	)
}
