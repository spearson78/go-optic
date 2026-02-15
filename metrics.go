package optic

import (
	"context"
	"fmt"

	"github.com/spearson78/go-optic/expr"
)

// Metrics contains information about the execution of an [Optic].
//
// See:
//   - [WithMetrics] for a combinator that generates metrics.
type Metrics struct {
	Focused   int //Indicates how many elements were focused.
	Access    int //Indicates how many times the source was accessed.
	LengthGet int //Indicates how many LengthGets were performed.
	IxGet     int //Indicates how many index lookups were performed.

	Custom map[string]int //Custom metrics
}

func (m Metrics) String() string {
	return fmt.Sprintf("metrics[F:%v A:%v I:%v L:%v Custom:%v]", m.Focused, m.Access, m.IxGet, m.LengthGet, m.Custom)
}

const metricsKey = "github.com/spearson78/go-optics/metrics"

// The WithMetrics combinator returns an Optic that provides metrics to the given [Metrics] struct
//
// See:
//   - `IncCustomMetric` for a function that enables custom metrics to be published.
func WithMetrics[I, S, T, A, B, RET, RW, DIR, ERR any](o Optic[I, S, T, A, B, RET, RW, DIR, ERR], m *Metrics) Optic[I, S, T, A, B, RET, RW, DIR, ERR] {

	return Omni[I, S, T, A, B, RET, RW, DIR, ERR](
		func(ctx context.Context, source S) (I, A, error) {
			m.Access++
			ctx = context.WithValue(ctx, metricsKey, m)
			i, a, err := o.AsGetter()(ctx, source)
			if err == nil {
				m.Focused++
			}
			return i, a, err
		},
		func(ctx context.Context, focus B, source S) (T, error) {
			m.Access++
			ctx = context.WithValue(ctx, metricsKey, m)
			t, err := o.AsSetter()(ctx, focus, source)
			if err == nil {
				m.Focused++
			}
			return t, err
		},
		func(ctx context.Context, source S) SeqIE[I, A] {
			ctx = context.WithValue(ctx, metricsKey, m)
			return func(yield func(val ValueIE[I, A]) bool) {
				m.Access++
				i := 0
				o.AsIter()(ctx, source)(func(val ValueIE[I, A]) bool {
					index, focus, err := val.Get()
					if err == nil {
						i++
					}
					return yield(ValIE(index, focus, err))
				})
				m.Focused += i
			}
		},
		func(ctx context.Context, source S) (int, error) {
			m.LengthGet++
			ctx = context.WithValue(ctx, metricsKey, m)
			return o.AsLengthGetter()(ctx, source)
		},
		func(ctx context.Context, fmap func(index I, focus A) (B, error), source S) (T, error) {
			ctx = context.WithValue(ctx, metricsKey, m)
			m.Access++
			i := 0
			t, err := o.AsModify()(ctx, func(index I, focus A) (B, error) {
				b, err := fmap(index, focus)
				if err == nil {
					i++
				}
				return b, err
			}, source)
			if err == nil {
				m.Focused += i
			}
			return t, err
		},
		func(ctx context.Context, index I, source S) SeqIE[I, A] {
			ctx = context.WithValue(ctx, metricsKey, m)
			return func(yield func(val ValueIE[I, A]) bool) {
				m.IxGet++
				i := 0
				o.AsIxGetter()(ctx, index, source)(func(val ValueIE[I, A]) bool {
					index, focus, err := val.Get()
					if err == nil {
						i++
					}
					return yield(ValIE(index, focus, err))
				})
				m.Focused += i
			}
		},
		o.AsIxMatch(),
		func(ctx context.Context, focus B) (T, error) {
			ctx = context.WithValue(ctx, metricsKey, m)
			t, err := o.AsReverseGetter()(ctx, focus)
			if err == nil {
				m.Focused++
			}
			return t, err
		},
		ExprDef(
			func(ot expr.OpticTypeExpr) expr.OpticExpression {
				return expr.WithMetrics{
					OpticTypeExpr: ot,
					Optic:         o.AsExpr(),
				}
			},
			o,
		),
	)
}

// IncCustomMetric increments the given custom metric contained within the context.
// Custom metrics will bubble up to the nearest [WithMetric] optic.
func IncCustomMetric(ctx context.Context, name string, val int) {
	if metrics, ok := ctx.Value(metricsKey).(*Metrics); ok && metrics != nil {
		if metrics.Custom == nil {
			metrics.Custom = make(map[string]int)
		}
		curVal, _ := metrics.Custom[name]
		metrics.Custom[name] = curVal + val
	}
}
