package users

import (
	"context"
	"encoding/json"
	"slices"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/slntopp/core-chatting/cc"
	"github.com/slntopp/core-chatting/pkg/core"
	"github.com/slntopp/core-chatting/pkg/graph"
	settingspb "github.com/slntopp/nocloud-proto/settings"
	"go.uber.org/zap"
)

// Account grades, the same as in the NoCloud admin UI: billing keeps
// account.spending (average real money paid per month) and the grade is
// derived from the "account-grades" setting. Only operators see it.
const accountGradesKey = "account-grades"

const gradesCacheTTL = time.Minute

type gradeDef struct {
	Title string   `json:"title"`
	Min   *float64 `json:"min,omitempty"`
	Color string   `json:"color,omitempty"`
}

// Used until the setting is readable. Thresholds are per month.
var defaultGrades = []gradeDef{
	{Title: "A+", Color: "#8E24AA"},
	{Title: "A", Min: ptr(300.0), Color: "#2E7D32"},
	{Title: "B", Min: ptr(100.0), Color: "#1565C0"},
	{Title: "C", Min: ptr(0.0), Color: "#757575"},
}

func ptr[T any](v T) *T { return &v }

type gradesCache struct {
	mu      sync.Mutex
	grades  []gradeDef
	fetched time.Time
}

// parseGrades reads the setting value; nil when it holds no usable grades.
func parseGrades(raw string) []gradeDef {
	var conf struct {
		Grades []gradeDef `json:"grades"`
	}
	if err := json.Unmarshal([]byte(raw), &conf); err != nil {
		return nil
	}
	var valid []gradeDef
	for _, g := range conf.Grades {
		if g.Title != "" {
			valid = append(valid, g)
		}
	}
	return valid
}

// grades returns the configured grades, read at most once a minute. The setting
// is only read, never written: a failed read keeps the last known grades.
func (s *UsersServer) grades(log *zap.Logger) []gradeDef {
	s.gradesCache.mu.Lock()
	defer s.gradesCache.mu.Unlock()

	if s.gradesCache.grades != nil && time.Since(s.gradesCache.fetched) < gradesCacheTTL {
		return s.gradesCache.grades
	}
	s.gradesCache.fetched = time.Now()
	if s.gradesCache.grades == nil {
		s.gradesCache.grades = defaultGrades
	}
	if s.settings == nil {
		return s.gradesCache.grades
	}

	ctx, cancel := context.WithTimeout(s.settingsCtx, 3*time.Second)
	defer cancel()
	r, err := s.settings.Get(ctx, &settingspb.GetRequest{Keys: []string{accountGradesKey}})
	if err != nil {
		log.Warn("Failed to get account grades setting", zap.Error(err))
		return s.gradesCache.grades
	}
	if v, ok := r.GetFields()[accountGradesKey]; ok {
		if grades := parseGrades(v.GetStringValue()); grades != nil {
			s.gradesCache.grades = grades
		}
	}
	return s.gradesCache.grades
}

// computeGrade mirrors getAccountGrade in the NoCloud admin UI: a manual grade
// wins, otherwise the highest grade whose min the monthly average reaches.
func computeGrade(info graph.AccountGradeInfo, grades []gradeDef, defaultCurrency string) *cc.UserGrade {
	var grade *cc.UserGrade
	override := strings.TrimSpace(info.GradeOverride)

	if override != "" {
		grade = &cc.UserGrade{Title: override, Manual: true}
		for _, g := range grades {
			if g.Title == override {
				grade.Color = g.Color
				break
			}
		}
	} else {
		if info.Spending == nil {
			return nil
		}
		auto := make([]gradeDef, 0, len(grades))
		for _, g := range grades {
			if g.Min != nil {
				auto = append(auto, g)
			}
		}
		sort.SliceStable(auto, func(i, j int) bool { return *auto[i].Min > *auto[j].Min })
		for _, g := range auto {
			if info.Spending.Amount >= *g.Min {
				grade = &cc.UserGrade{Title: g.Title, Color: g.Color}
				break
			}
		}
		if grade == nil {
			return nil
		}
	}

	if sp := info.Spending; sp != nil {
		grade.Amount = ptr(sp.Amount)
		grade.PeriodDays = sp.PeriodDays
		// NCU or no currency means the default currency, as in the admin UI.
		code := sp.Currency.Code
		if code == "" {
			code = sp.Currency.Title
		}
		if (code == "" || code == "NCU") && defaultCurrency != "" {
			code = defaultCurrency
		}
		grade.Currency = code
	}
	return grade
}

// isOperator reports global chat admins and department admins: the people who
// answer tickets. Clients never get grades.
func isOperator(requestor string) bool {
	conf, err := core.Config()
	if err != nil {
		return false
	}
	if slices.Contains(conf.GetAdmins(), requestor) {
		return true
	}
	for _, d := range conf.GetDepartments() {
		if slices.Contains(d.GetAdmins(), requestor) {
			return true
		}
	}
	return false
}

// attachGrades fills User.grade for operators and clears it for everyone else.
func (s *UsersServer) attachGrades(ctx context.Context, log *zap.Logger, requestor string, users []*cc.User) {
	for _, u := range users {
		u.Grade = nil
	}
	if len(users) == 0 || !isOperator(requestor) {
		return
	}

	uuids := make([]string, 0, len(users))
	for _, u := range users {
		uuids = append(uuids, u.GetUuid())
	}
	infos, defaultCurrency, err := s.ctrl.GradeInfo(ctx, uuids)
	if err != nil {
		log.Warn("Failed to get account grades", zap.Error(err))
		return
	}
	grades := s.grades(log)
	for _, u := range users {
		if info, ok := infos[u.GetUuid()]; ok {
			u.Grade = computeGrade(info, grades, defaultCurrency)
		}
	}
}
