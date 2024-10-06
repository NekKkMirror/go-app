package pagination

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/NekKkMirror/go-app/internal/pkg/mapper"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

const (
	defaultSize = 10
	defaultPage = 1
)

type ListResult[T interface{}] struct {
	Size            int    `json:"size,omitempty"            bson:"size"`
	Page            int    `json:"page,omitempty"            bson:"page"`
	TotalCount      int64  `json:"totalCount,omitempty"      bson:"totalCount"`
	TotalPages      int    `json:"totalPages,omitempty"      bson:"totalPages"`
	HasPreviousPage bool   `json:"hasPreviousPage,omitempty" bson:"hasPreviousPage"`
	HasNextPage     bool   `json:"hasNextPage,omitempty"     bson:"hasNextPage"`
	FirstItemIndex  int    `json:"firstItemIndex,omitempty"  bson:"firstItemIndex"`
	LastItemIndex   int    `json:"lastItemIndex,omitempty"   bson:"lastItemIndex"`
	IsFirstPage     bool   `json:"isFirstPage,omitempty"     bson:"isFirstPage"`
	IsLastPage      bool   `json:"isLastPage,omitempty"      bson:"isLastPage"`
	NextPage        int    `json:"nextPage,omitempty"        bson:"nextPage"`
	PreviousPage    int    `json:"previousPage,omitempty"    bson:"previousPage"`
	IsEmpty         bool   `json:"isEmpty,omitempty"         bson:"isEmpty"`
	HasSinglePage   bool   `json:"hasSinglePage,omitempty"   bson:"hasSinglePage"`
	HasMorePages    bool   `json:"hasMorePages,omitempty"    bson:"hasMorePages"`
	HasLessPages    bool   `json:"hasLessPages,omitempty"    bson:"hasLessPages"`
	PaginationInfo  string `json:"paginationInfo,omitempty"  bson:"paginationInfo"`
	Data            []T    `json:"data,omitempty"            bson:"data"`
}

// NewListResult constructs an instance of ListResult with calculated pagination details.
func NewListResult[T any](size, page int, totalCount int64, data []T) *ListResult[T] {
	totalPages := calculateTotalPages(size, totalCount)
	firstItemIndex := (page - 1) * size
	lastItemIndex := page * size

	return &ListResult[T]{
		Size:            size,
		Page:            page,
		TotalCount:      totalCount,
		TotalPages:      totalPages,
		FirstItemIndex:  firstItemIndex,
		LastItemIndex:   lastItemIndex,
		IsFirstPage:     page == 1,
		IsLastPage:      lastItemIndex >= int(totalCount),
		HasPreviousPage: page > 1,
		HasNextPage:     lastItemIndex < int(totalCount),
		NextPage:        page + 1,
		PreviousPage:    page - 1,
		IsEmpty:         len(data) == 0,
		HasSinglePage:   totalPages == 1,
		HasMorePages:    lastItemIndex < int(totalCount),
		HasLessPages:    page > 1,
		PaginationInfo:  fmt.Sprintf("Showing data %d to %d of %d", firstItemIndex+1, lastItemIndex, totalCount),
		Data:            data,
	}
}

// calculateTotalPages determines the number of pages given the size and total count.
func calculateTotalPages(size int, totalCount int64) int {
	return int(math.Ceil(float64(totalCount) / float64(size)))
}

// ListQuery represents the query parameters for pagination and filtering.
type ListQuery struct {
	Size    int            `query:"size"    json:"size,omitempty"`
	Page    int            `query:"page"    json:"page,omitempty"`
	OrderBy string         `query:"orderBy" json:"orderBy,omitempty"`
	Filters []*FilterModel `query:"filters" json:"filters,omitempty"`
}

// FilterModel represents the filtering model with field, value, and comparison parameters.
type FilterModel struct {
	Field      string `query:"field"      json:"field"`
	Value      string `query:"value"      json:"value"`
	Comparison string `query:"comparison" json:"comparison"`
}

// NewListQuery creates a new instance of ListQuery with the given size and page parameters.
func NewListQuery(size, page int) *ListQuery {
	return &ListQuery{
		Size: size,
		Page: page,
	}
}

// NewListQueryFromQueryParams creates a new instance of ListQuery based on the provided query parameters.
func NewListQueryFromQueryParams(sizeStr, pageStr string) (*ListQuery, error) {
	size, err := strconv.Atoi(sizeStr)
	if err != nil || size == 0 {
		size = defaultSize
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil || page == 0 {
		page = defaultPage
	}

	return &ListQuery{
		Size: size,
		Page: page,
	}, nil
}

// GetListQueryFromCtx retrieves a ListQuery instance from the provided echo.Context.
func GetListQueryFromCtx(c echo.Context) (*ListQuery, error) {
	q := &ListQuery{}
	var page, size, orderBy string

	err := echo.QueryParamsBinder(c).
		CustomFunc("filters", func(values []string) []error {
			for _, v := range values {
				if v == "" {
					continue
				}
				f := &FilterModel{}
				if err := c.Bind(f); err != nil {
					return []error{err}
				}
				q.Filters = append(q.Filters, f)
			}
			return nil
		}).
		String("size", &size).
		String("page", &page).
		String("orderBy", &orderBy).
		BindError()

	if err != nil {
		return nil, err
	}

	if err = q.SetSize(size); err != nil {
		return nil, err
	}
	if err = q.SetPage(page); err != nil {
		return nil, err
	}
	q.SetOrderBy(orderBy)

	return q, nil
}

// SetSize sets the size parameter of the ListQuery instance.
func (q *ListQuery) SetSize(sizeQuery string) error {
	if sizeQuery == "" {
		q.Size = defaultSize
		return nil
	}
	size, err := strconv.Atoi(sizeQuery)
	if err != nil {
		return errors.Wrap(err, "invalid size parameter")
	}
	q.Size = size
	return nil
}

// SetPage sets the page parameter of the ListQuery instance.
func (q *ListQuery) SetPage(pageQuery string) error {
	if pageQuery == "" {
		q.Page = defaultPage
		return nil
	}
	page, err := strconv.Atoi(pageQuery)
	if err != nil {
		return errors.Wrap(err, "invalid page parameter")
	}
	q.Page = page
	return nil
}

// SetOrderBy sets the order by parameter of the ListQuery instance.
func (q *ListQuery) SetOrderBy(orderByQuery string) {
	q.OrderBy = orderByQuery
}

// GetQueryString generates a query string representation of the ListQuery instance.
func (q *ListQuery) GetQueryString() string {
	return fmt.Sprintf("size=%d&page=%d&orderBy=%s", q.GetSize(), q.GetPage(), q.GetOrderBy())
}

// GetSize returns the size parameter for pagination.
func (q *ListQuery) GetSize() int {
	return q.Size
}

// GetPage returns the current page number for pagination.
func (q *ListQuery) GetPage() int {
	return q.Page
}

// GetOrderBy returns the order by parameter for pagination.
func (q *ListQuery) GetOrderBy() string {
	return q.OrderBy
}

// GetOffset calculates and returns the offset for pagination based on the current page and size.
func (q *ListQuery) GetOffset() int {
	if q.Page == 0 {
		return 0
	}
	return (q.Page - 1) * q.Size
}

// GetLimit calculates and returns the limit for pagination based on the current size.
func (q *ListQuery) GetLimit() int {
	return q.Size
}

// ListResultToDTO converts a ListResult of type TModel to a ListResult of type TDTO.
func ListResultToDTO[TModel any, TDTO any](listResult *ListResult[TModel]) (*ListResult[TDTO], error) {
	if listResult == nil {
		return nil, errors.New("ListResult is nil")
	}

	var dataDTO []TDTO
	for _, model := range listResult.Data {
		dto, err := mapper.Map[TModel, TDTO](model)
		if err != nil {
			return nil, err
		}
		dataDTO = append(dataDTO, dto)
	}

	dtoResult := &ListResult[TDTO]{
		Size:            listResult.Size,
		Page:            listResult.Page,
		TotalCount:      listResult.TotalCount,
		TotalPages:      listResult.TotalPages,
		HasPreviousPage: listResult.HasPreviousPage,
		HasNextPage:     listResult.HasNextPage,
		FirstItemIndex:  listResult.FirstItemIndex,
		LastItemIndex:   listResult.LastItemIndex,
		IsFirstPage:     listResult.IsFirstPage,
		IsLastPage:      listResult.IsLastPage,
		NextPage:        listResult.NextPage,
		PreviousPage:    listResult.PreviousPage,
		IsEmpty:         listResult.IsEmpty,
		HasSinglePage:   listResult.HasSinglePage,
		HasMorePages:    listResult.HasMorePages,
		HasLessPages:    listResult.HasLessPages,
		PaginationInfo:  listResult.PaginationInfo,
		Data:            dataDTO,
	}

	return dtoResult, nil
}

// ApplyFilterAction applies the filters defined in ListQuery to the gorm.DB instance.
func ApplyFilterAction(db *gorm.DB, filters []*FilterModel, fieldsNotAllowed map[string]bool) (*gorm.DB, error) {
	for _, filter := range filters {
		if len(fieldsNotAllowed) > 0 && fieldsNotAllowed[filter.Field] {
			return nil, fmt.Errorf("filter field %s is not allowed", filter.Field)
		}

		condition, value, err := buildCondition(filter)
		if err != nil {
			return nil, err
		}

		db = db.Where(condition, value...)
	}
	return db, nil
}

// buildCondition builds the SQL condition string based on the FilterModel.
func buildCondition(filter *FilterModel) (string, []interface{}, error) {
	var condition string
	var value []interface{}

	// Convert comparison to lowercase to handle case insensitivity
	comparison := strings.ToLower(filter.Comparison)

	switch comparison {
	case "eq", "=":
		condition = fmt.Sprintf("%s = ?", filter.Field)
		value = []interface{}{filter.Value}
	case "ne", "!=", "<>":
		condition = fmt.Sprintf("%s <> ?", filter.Field)
		value = []interface{}{filter.Value}
	case "lt", "<":
		condition = fmt.Sprintf("%s < ?", filter.Field)
		value = []interface{}{filter.Value}
	case "lte", "<=":
		condition = fmt.Sprintf("%s <= ?", filter.Field)
		value = []interface{}{filter.Value}
	case "gt", ">":
		condition = fmt.Sprintf("%s > ?", filter.Field)
		value = []interface{}{filter.Value}
	case "gte", ">=":
		condition = fmt.Sprintf("%s >= ?", filter.Field)
		value = []interface{}{filter.Value}
	case "like":
		condition = fmt.Sprintf("%s LIKE ?", filter.Field)
		value = []interface{}{filter.Value}
	case "ilike":
		condition = fmt.Sprintf("%s ILIKE ?", filter.Field) // Case-insensitive like
		value = []interface{}{filter.Value}
	case "similar_to":
		condition = fmt.Sprintf("%s SIMILAR TO ?", filter.Field)
		value = []interface{}{filter.Value}
	case "not_similar_to":
		condition = fmt.Sprintf("%s NOT SIMILAR TO ?", filter.Field)
		value = []interface{}{filter.Value}
	case "ends_with":
		condition = fmt.Sprintf("%s LIKE ?", filter.Field)
		value = []interface{}{fmt.Sprintf("%%%s", filter.Value)}
	case "starts_with":
		condition = fmt.Sprintf("%s LIKE ?", filter.Field)
		value = []interface{}{fmt.Sprintf("%s%%", filter.Value)}
	case "in":
		condition = fmt.Sprintf("%s IN (?)", filter.Field)
		value = []interface{}{filter.Value}
	case "not_in":
		condition = fmt.Sprintf("%s NOT IN (?)", filter.Field)
		value = []interface{}{filter.Value}
	case "is_null":
		condition = fmt.Sprintf("%s IS NULL", filter.Field)
		value = nil
	case "is_not_null":
		condition = fmt.Sprintf("%s IS NOT NULL", filter.Field)
		value = nil
	case "between":
		parts := strings.Split(filter.Value, ",")
		if len(parts) != 2 {
			return "", nil, errors.New("invalid value for between operator, expected two values separated by a comma")
		}
		lowerBound := parts[0]
		upperBound := parts[1]
		condition = fmt.Sprintf("%s BETWEEN ? AND ?", filter.Field)
		value = []interface{}{lowerBound, upperBound}
	case "contains":
		condition = fmt.Sprintf("%s @> ?", filter.Field)
		value = []interface{}{filter.Value}
	case "contained_in":
		condition = fmt.Sprintf("%s <@ ?", filter.Field)
		value = []interface{}{filter.Value}
	case "overlap":
		condition = fmt.Sprintf("%s && ?", filter.Field)
		value = []interface{}{filter.Value}
	case "distinct_from":
		condition = fmt.Sprintf("%s IS DISTINCT FROM ?", filter.Field)
		value = []interface{}{filter.Value}
	case "not_distinct_from":
		condition = fmt.Sprintf("%s IS NOT DISTINCT FROM ?", filter.Field)
		value = []interface{}{filter.Value}
	case "is_true":
		condition = fmt.Sprintf("%s IS TRUE", filter.Field)
		value = nil
	case "is_not_true":
		condition = fmt.Sprintf("%s IS NOT TRUE", filter.Field)
		value = nil
	case "is_false":
		condition = fmt.Sprintf("%s IS FALSE", filter.Field)
		value = nil
	case "is_not_false":
		condition = fmt.Sprintf("%s IS NOT FALSE", filter.Field)
		value = nil
	case "is_unknown":
		condition = fmt.Sprintf("%s IS UNKNOWN", filter.Field)
		value = nil
	case "is_not_unknown":
		condition = fmt.Sprintf("%s IS NOT UNKNOWN", filter.Field)
		value = nil
	case "is_positive":
		condition = fmt.Sprintf("%s > 0", filter.Field)
		value = nil
	case "is_negative":
		condition = fmt.Sprintf("%s < 0", filter.Field)
		value = nil
	case "is_not_positive":
		condition = fmt.Sprintf("%s <= 0", filter.Field)
		value = nil
	case "is_not_negative":
		condition = fmt.Sprintf("%s >= 0", filter.Field)
		value = nil
	case "is_even":
		condition = fmt.Sprintf("%s %% 2 = 0", filter.Field)
		value = nil
	case "is_odd":
		condition = fmt.Sprintf("%s %% 2 != 0", filter.Field)
		value = nil
	case "is_divisible_by":
		divisor, err := strconv.Atoi(filter.Value)
		if err != nil {
			return "", nil, fmt.Errorf("failed to parse divisor value: %w", err)
		}
		condition = fmt.Sprintf("%s %% ? = 0", filter.Field)
		value = []interface{}{divisor}
	default:
		return "", nil, fmt.Errorf("unsupported comparison operator: %s", filter.Comparison)
	}

	return condition, value, nil
}
