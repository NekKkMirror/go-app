package pagination

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
)

func TestListQuery_SetSize(t *testing.T) {
	q := &ListQuery{}
	err := q.SetSize("20")
	if err != nil {
		t.Errorf("SetSize failed: %v", err)
	}

	if q.Size != 20 {
		t.Errorf("SetSize did not set the correct size")
	}
}

func TestListQuery_SetPage(t *testing.T) {
	q := &ListQuery{}
	err := q.SetPage("2")
	if err != nil {
		t.Errorf("SetPage failed: %v", err)
	}

	if q.Page != 2 {
		t.Errorf("SetPage did not set the correct page")
	}
}

func TestListQuery_SetOrderBy(t *testing.T) {
	q := &ListQuery{}
	q.SetOrderBy("name")
	if q.OrderBy != "name" {
		t.Errorf("SetOrderBy did not set the correct orderBy value")
	}
}

func TestListQuery_GetQueryString(t *testing.T) {
	q := &ListQuery{Size: 10, Page: 1, OrderBy: "name"}
	expectedQueryString := "size=10&page=1&orderBy=name"
	if q.GetQueryString() != expectedQueryString {
		t.Errorf("GetQueryString returned incorrect value")
	}
}

func TestListQuery_GetOffset(t *testing.T) {
	q := &ListQuery{Size: 10, Page: 2}
	expectedOffset := 10
	if q.GetOffset() != expectedOffset {
		t.Errorf("GetOffset returned incorrect value")
	}
}

func TestListQuery_GetLimit(t *testing.T) {
	q := &ListQuery{Size: 10}
	expectedLimit := 10
	if q.GetLimit() != expectedLimit {
		t.Errorf("GetLimit returned incorrect value")
	}
}

func TestGetListQueryFromCtx(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	q, err := GetListQueryFromCtx(c)
	if err != nil {
		t.Errorf("GetListQueryFromCtx failed: %v", err)
	}

	if len(q.Filters) != 0 {
		t.Errorf("Filters should be empty when no filters are provided")
	}

	if q.Size != 10 {
		t.Errorf("Size should be set to 0 when not provided in the context")
	}
	if q.Page != 1 {
		t.Errorf("Page should be set to 0 when not provided in the context")
	}
	if q.OrderBy != "" {
		t.Errorf("OrderBy should be empty when not provided in the context")
	}
}
