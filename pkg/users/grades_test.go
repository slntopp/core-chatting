package users

import (
	"testing"

	"github.com/slntopp/core-chatting/pkg/graph"
)

func spending(amount float64, code string) *graph.AccountSpending {
	s := &graph.AccountSpending{Amount: amount, PeriodDays: 395}
	s.Currency.Code = code
	return s
}

func TestComputeGrade(t *testing.T) {
	cases := []struct {
		name   string
		info   graph.AccountGradeInfo
		title  string
		manual bool
		cur    string
	}{
		{"A from 300", graph.AccountGradeInfo{Spending: spending(300, "BYN")}, "A", false, "BYN"},
		{"B between", graph.AccountGradeInfo{Spending: spending(299.99, "BYN")}, "B", false, "BYN"},
		{"C below 100", graph.AccountGradeInfo{Spending: spending(5, "BYN")}, "C", false, "BYN"},
		{"NCU shown as default", graph.AccountGradeInfo{Spending: spending(150, "NCU")}, "B", false, "PLN"},
		{"manual wins", graph.AccountGradeInfo{Spending: spending(5, "BYN"), GradeOverride: " A+ "}, "A+", true, "BYN"},
		{"manual without spending", graph.AccountGradeInfo{GradeOverride: "VIP"}, "VIP", true, ""},
		{"nothing counted", graph.AccountGradeInfo{}, "", false, ""},
	}
	for _, c := range cases {
		g := computeGrade(c.info, defaultGrades, "PLN")
		if c.title == "" {
			if g != nil {
				t.Errorf("%s: want no grade, got %v", c.name, g)
			}
			continue
		}
		if g == nil || g.GetTitle() != c.title || g.GetManual() != c.manual || g.GetCurrency() != c.cur {
			t.Errorf("%s: got %v", c.name, g)
		}
	}
}

func TestParseGrades(t *testing.T) {
	g := parseGrades(`{"grades":[{"title":"A","min":300,"color":"red"},{"min":1},{"title":"A+"}]}`)
	if len(g) != 2 || g[0].Title != "A" || g[0].Min == nil || *g[0].Min != 300 || g[1].Min != nil {
		t.Errorf("got %+v", g)
	}
	if parseGrades("not json") != nil || parseGrades(`{"grades":[]}`) != nil {
		t.Error("want nil for unusable values")
	}
}
