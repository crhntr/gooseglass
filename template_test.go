package gooseglass

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/crhntr/dom/domtest"
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/crhntr/gooseglass/internal/fake"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate
//counterfeiter:generate -o internal/fake/provider.go --fake-name=Provider . migrationProvider

func Test(t *testing.T) {
	type (
		Fakes struct {
			provider *fake.Provider
		}
		Given struct {
			Fakes
		}
		When struct{}
		Then struct {
			Fakes
		}
		Case struct {
			Name  string
			Given func(*testing.T, Given)
			When  func(*testing.T, When) *http.Request
			Then  func(*testing.T, Then, *http.Response)
		}
	)

	newFakes := func() Fakes {
		fakes := Fakes{
			provider: new(fake.Provider),
		}
		return fakes
	}

	run := func(t *testing.T, tc Case) {
		fakes := newFakes()

		if tc.Given != nil {
			tc.Given(t, Given{
				Fakes: fakes,
			})
		}

		mux := http.NewServeMux()
		routes(mux, fakes.provider)

		require.NotNil(t, tc.When)
		req := tc.When(t, When{})
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if tc.Then != nil {
			tc.Then(t, Then{
				Fakes: fakes,
			}, rec.Result())
		}
	}

	for _, tc := range []Case{
		{
			Name: "status list is empty",
			Given: func(t *testing.T, g Given) {
				g.provider.StatusReturns([]*goose.MigrationStatus{}, nil)
			},
			When: func(t *testing.T, when When) *http.Request {
				req := httptest.NewRequest(http.MethodGet, TemplateRoutePaths{}.Status(), nil)
				return req
			},
			Then: func(t *testing.T, when Then, then *http.Response) {
				assert.Equal(t, http.StatusOK, then.StatusCode)
			},
		},
		{
			Name: "status is pending",
			Given: func(t *testing.T, g Given) {
				g.provider.StatusReturns([]*goose.MigrationStatus{
					{State: goose.StatePending, Source: &goose.Source{Type: goose.TypeSQL, Path: "01_init.sql", Version: 1}},
				}, nil)
			},
			When: func(t *testing.T, when When) *http.Request {
				req := httptest.NewRequest(http.MethodGet, TemplateRoutePaths{}.Status(), nil)
				return req
			},
			Then: func(t *testing.T, then Then, response *http.Response) {
				assert.Equal(t, http.StatusOK, response.StatusCode)
				document := domtest.ParseResponseDocument(t, response)

				headerColumns := document.QuerySelectorAll(`#status #status-table thead tr th`)
				require.NotNil(t, headerColumns)

				if tbody := document.QuerySelector(`#status #status-table tbody`); assert.NotNil(t, tbody) {
					assert.Equal(t, 1, tbody.ChildElementCount())
					if tr := tbody.QuerySelector(`tr`); assert.NotNil(t, tr) {
						columns := tr.QuerySelectorAll(`td`)
						assert.Equal(t, headerColumns.Length(), columns.Length())
					}
				}
			},
		},
	} {
		t.Run(tc.Name, func(t *testing.T) { run(t, tc) })
	}
}
