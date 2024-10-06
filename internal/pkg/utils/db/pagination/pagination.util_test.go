package pagination

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/NekKkMirror/go-app/internal/pkg/mapper"
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

// Successfully maps ListResult data from TModel to TDTO
func TestListResultToDTOSuccessfulMapping(t *testing.T) {
	type TModel struct {
		ID int
	}
	type TDTO struct {
		ID int
	}

	listResult := &ListResult[TModel]{
		Size:       10,
		Page:       1,
		TotalCount: 100,
		TotalPages: 10,
		Data:       []TModel{{ID: 1}, {ID: 2}},
	}

	err := mapper.CreateMap[TModel, TDTO]()
	if err != nil {
		t.Errorf("expected no error, but got %v", err)
	}

	mappedListResult, err := ListResultToDTO[TModel, TDTO](listResult)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	fmt.Println(mappedListResult)

	// Check if the mappedListResult is not nil
	if mappedListResult == nil {
		t.Fatalf("mappedListResult is nil")
	}

	// Check metadata equality
	if mappedListResult.Size != listResult.Size || mappedListResult.Page != listResult.Page || mappedListResult.TotalCount != listResult.TotalCount || mappedListResult.TotalPages != listResult.TotalPages {
		t.Errorf("expected metadata to be equal, got different values")
	}

	// Check data length equality
	if len(mappedListResult.Data) != len(listResult.Data) {
		t.Errorf("expected data length to be %d, got %d", len(listResult.Data), len(mappedListResult.Data))
	}

	// Check individual data elements
	for i, dto := range mappedListResult.Data {
		if dto.ID != listResult.Data[i].ID {
			t.Errorf("expected ID %d, got %d", listResult.Data[i].ID, dto.ID)
		}
	}
}

// Initializes ListQuery with given size and page
func TestNewListQueryInitialization(t *testing.T) {
	size := 10
	page := 2
	query := NewListQuery(size, page)

	if query.Size != size {
		t.Errorf("expected size %d, got %d", size, query.Size)
	}

	if query.Page != page {
		t.Errorf("expected page %d, got %d", page, query.Page)
	}
}

func TestNewListQueryFromQueryParamsWithInvalidInputs(t *testing.T) {
	sizeStr := "abc"
	pageStr := "xyz"

	query, err := NewListQueryFromQueryParams(sizeStr, pageStr)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if query.Size != defaultSize {
		t.Errorf("expected default size %d, got %d", defaultSize, query.Size)
	}

	if query.Page != defaultPage {
		t.Errorf("expected default page %d, got %d", defaultPage, query.Page)
	}
}
