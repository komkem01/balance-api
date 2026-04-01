package budgets

import (
	"balance/app/modules/entities/ent"
	"balance/internal/config"
	"time"

	"go.opentelemetry.io/otel/trace"
)

type Service struct {
	tracer trace.Tracer
	db     BudgetStore
}

type Config struct{}

type Options struct {
	*config.Config[Config]
	tracer trace.Tracer
	db     BudgetStore
}

func newService(opt *Options) *Service {
	return &Service{tracer: opt.tracer, db: opt.db}
}

func parseBudgetPeriod(value string) (ent.BudgetPeriod, bool) {
	switch value {
	case string(ent.BudgetPeriodDaily):
		return ent.BudgetPeriodDaily, true
	case string(ent.BudgetPeriodWeekly):
		return ent.BudgetPeriodWeekly, true
	case string(ent.BudgetPeriodMonthly):
		return ent.BudgetPeriodMonthly, true
	default:
		return "", false
	}
}

func budgetLocation() *time.Location {
	loc, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		return time.FixedZone("ICT", 7*60*60)
	}
	return loc
}

func dateOnlyInLocation(t time.Time, loc *time.Location) time.Time {
	y, m, d := t.In(loc).Date()
	return time.Date(y, m, d, 0, 0, 0, 0, loc)
}

func periodDateRange(period ent.BudgetPeriod, anchor time.Time, loc *time.Location) (time.Time, time.Time) {
	day := dateOnlyInLocation(anchor, loc)

	switch period {
	case ent.BudgetPeriodDaily:
		return day, day
	case ent.BudgetPeriodWeekly:
		weekday := int(day.Weekday())
		if weekday == 0 {
			weekday = 7
		}
		start := day.AddDate(0, 0, -(weekday - 1))
		end := start.AddDate(0, 0, 6)
		return start, end
	case ent.BudgetPeriodMonthly:
		y, m, _ := day.Date()
		start := time.Date(y, m, 1, 0, 0, 0, 0, loc)
		end := start.AddDate(0, 1, -1)
		return start, end
	default:
		return day, day
	}
}

func resolveBudgetDateRange(period ent.BudgetPeriod, startDate *time.Time, endDate *time.Time) (*time.Time, *time.Time, error) {
	loc := budgetLocation()

	if startDate != nil && endDate != nil {
		s := dateOnlyInLocation(*startDate, loc)
		e := dateOnlyInLocation(*endDate, loc)
		if e.Before(s) {
			return nil, nil, ErrBudgetDateInvalid
		}
		return &s, &e, nil
	}

	anchor := time.Now().In(loc)
	if startDate != nil {
		anchor = *startDate
	}
	if endDate != nil {
		anchor = *endDate
	}

	s, e := periodDateRange(period, anchor, loc)
	return &s, &e, nil
}
