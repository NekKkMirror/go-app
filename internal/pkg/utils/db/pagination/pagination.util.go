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
	Size            int    `json:"size,omitempty"       bson:"size"`
	Page            int    `json:"page,omitempty"       bson:"page"`
	TotalCount      int64  `json:"totalCount,omitempty" bson:"totalCount"`
	TotalPages      int    `json:"totalPages,omitempty" bson:"totalPages"`
	HasPreviousPage bool   `json:"hasPreviousPage,omitempty" bson:"hasPreviousPage"`
	HasNextPage     bool   `json:"hasNextPage,omitempty" bson:"hasNextPage"`
	FirstItemIndex  int    `json:"firstItemIndex,omitempty" bson:"firstItemIndex"`
	LastItemIndex   int    `json:"lastItemIndex,omitempty" bson:"lastItemIndex"`
	IsFirstPage     bool   `json:"isFirstPage,omitempty" bson:"isFirstPage"`
	IsLastPage      bool   `json:"isLastPage,omitempty" bson:"isLastPage"`
	NextPage        int    `json:"nextPage,omitempty" bson:"nextPage"`
	PreviousPage    int    `json:"previousPage,omitempty" bson:"previousPage"`
	IsEmpty         bool   `json:"isEmpty,omitempty" bson:"isEmpty"`
	HasSinglePage   bool   `json:"hasSinglePage,omitempty" bson:"hasSinglePage"`
	HasMorePages    bool   `json:"hasMorePages,omitempty" bson:"hasMorePages"`
	HasLessPages    bool   `json:"hasLessPages,omitempty" bson:"hasLessPages"`
	PaginationInfo  string `json:"paginationInfo,omitempty" bson:"paginationInfo"`
	Data            []T    `json:"data,omitempty"      bson:"data"`
}

// NewListResult creates a new instance of ListResult with the given size, page, totalCount, and data.
// It calculates additional pagination information such as total pages, previous/next page availability,
// first/last item indices, and pagination info string.
//
// Parameters:
// - size: The number of items per page.
// - page: The current page number.
// - totalCount: The total number of items.
// - data: The slice of items to be included in the ListResult.
//
// Returns:
// - A pointer to a new instance of ListResult[T] containing the provided parameters and calculated pagination information.
func NewListResult[T any](size int, page int, totalCount int64, data []T) *ListResult[T] {
	listResult := &ListResult[T]{
		Size:       size,
		Page:       page,
		TotalCount: totalCount,
		Data:       data,
	}

	listResult.TotalPages = getTotalPages(size, totalCount)
	listResult.HasPreviousPage = page > 1
	listResult.HasNextPage = (page * size) < int(totalCount)
	listResult.FirstItemIndex = (page - 1) * size
	listResult.LastItemIndex = page * size
	listResult.IsFirstPage = page == 1
	listResult.IsLastPage = (page * size) >= int(totalCount)
	listResult.NextPage = page + 1
	listResult.PreviousPage = page - 1
	listResult.IsEmpty = len(data) == 0
	listResult.HasSinglePage = listResult.TotalPages == 1
	listResult.HasMorePages = listResult.HasNextPage
	listResult.HasLessPages = listResult.HasPreviousPage
	listResult.PaginationInfo = fmt.Sprintf("Showing data %d to %d of %d", listResult.FirstItemIndex+1, listResult.LastItemIndex, totalCount)

	return listResult
}

// getTotalPages calculates the total number of pages based on the given size and total count.
// It uses the formula: totalPages = ceil(totalCount / size).
//
// Parameters:
// - size: The number of items per page.
// - totalCount: The total number of items.
//
// Returns:
// - The total number of pages.
func getTotalPages(size int, totalCount int64) int {
	d := float64(totalCount) / float64(size)
	return int(math.Ceil(d))
}

type ListQuery struct {
	Size    int            `query:"size"    json:"size,omitempty"`
	Page    int            `query:"page"    json:"page,omitempty"`
	OrderBy string         `query:"orderBy" json:"orderBy,omitempty"`
	Filters []*FilterModel `query:"filters" json:"filters,omitempty"`
}

type FilterModel struct {
	Field      string `query:"field"      json:"field"`
	Value      string `query:"value"      json:"value"`
	Comparison string `query:"comparison" json:"comparison"`
}

// NewListQuery creates a new instance of ListQuery with the given size and page parameters.
//
// Parameters:
// - size: The number of items per page. If size is less than or equal to 0, it defaults to 10.
// - page: The current page number. If page is less than or equal to 0, it defaults to 1.
//
// Returns:
// - A pointer to a new instance of ListQuery containing the provided size and page parameters.
func NewListQuery(size int, page int) *ListQuery {
	return &ListQuery{
		Size: size,
		Page: page,
	}
}

// NewListQueryFromQueryParams creates a new instance of ListQuery based on the provided query parameters.
// It sets the size and page parameters from the given strings, with default values if the conversion fails or the values are zero.
//
// Parameters:
// - size: A string representing the size parameter.
// - page: A string representing the page parameter.
//
// Returns:
// - A pointer to a new instance of ListQuery with the set size and page parameters.

func NewListQueryFromQueryParams(size string, page string) *ListQuery {
	p := &ListQuery{Size: defaultSize, Page: defaultPage}

	if sizeNum, err := strconv.Atoi(size); err != nil && sizeNum != 0 {
		p.Size = sizeNum
	}
	if pageNum, err := strconv.Atoi(page); err != nil && pageNum != 0 {
		p.Size = pageNum
	}

	return p
}

// GetListQueryFromCtx retrieves a ListQuery instance from the provided echo.Context.
// It extracts the query parameters from the context and constructs a new ListQuery instance.
// The function also handles error cases related to query parameter binding and validation.
//
// Parameters:
// - c: An instance of echo.Context, which provides access to the request and response objects.
//
// Returns:
// - A pointer to a new instance of ListQuery containing the extracted query parameters.
// - An error if any issues occur during query parameter binding or validation.
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
// If sizeQuery is an empty string, the size is set to the defaultSize.
// Otherwise, the size is converted from a string to an integer using strconv.Atoi.
// If the conversion fails, an error is returned.
//
// Parameters:
// - sizeQuery: A string representing the size parameter.
//
// Returns:
// - An error if the conversion fails.
func (q *ListQuery) SetSize(sizeQuery string) error {
	if sizeQuery == "" {
		q.Size = defaultSize
		return nil
	}
	sn, err := strconv.Atoi(sizeQuery)
	if err != nil {
		return err
	}
	q.Size = sn

	return nil
}

// SetPage sets the page parameter of the ListQuery instance.
// If pageQuery is an empty string, the page is set to the defaultPage.
// Otherwise, the page is converted from a string to an integer using strconv.Atoi.
// If the conversion fails, an error is returned.
//
// Parameters:
// - pageQuery: A string representing the page parameter.
//
// Returns:
// - An error if the conversion fails.
func (q *ListQuery) SetPage(pageQuery string) error {
	if pageQuery == "" {
		q.Page = defaultPage
		return nil
	}
	pn, err := strconv.Atoi(pageQuery)
	if err != nil {
		return err
	}
	q.Page = pn

	return nil
}

// SetOrderBy sets the order by parameter of the ListQuery instance.
//
// The order by parameter specifies the field to sort the data by.
// It is used in database queries to determine the order in which records are returned.
//
// Parameters:
// - orderByQuery: A string representing the order by parameter for pagination.
//
// The function does not return any value. It directly updates the OrderBy field of the ListQuery instance.
func (q *ListQuery) SetOrderBy(orderByQuery string) {
	q.OrderBy = orderByQuery
}

// GetQueryString generates a query string representation of the ListQuery instance.
// The query string includes the size, page, and orderBy parameters.
//
// Parameters:
// - q: A pointer to the ListQuery instance.
//
// Returns:
//   - A string containing the query string representation of the ListQuery instance.
//     The format of the query string is "size=<size>&page=<page>&orderBy=<orderBy>".
func (q *ListQuery) GetQueryString() string {
	return fmt.Sprintf("size=%v&page=%v&orderBy=%s", q.GetSize(), q.GetPage(), q.GetOrderBy())
}

// GetSize returns the size parameter for pagination.
//
// The size parameter determines the number of items to be displayed per page.
// It is used in conjunction with the page parameter to fetch a specific range of records from a database.
//
// Parameters:
// - q: A pointer to the ListQuery instance.
//
// Returns:
// - An integer representing the size parameter for pagination.
func (q *ListQuery) GetSize() int {
	return q.Size
}

// GetPage returns the current page number for pagination.
//
// The page number is used to determine the subset of data to be displayed.
// It is typically used in conjunction with the size parameter to fetch a specific range of records from a database.
//
// Parameters:
// - q: A pointer to the ListQuery instance.
//
// Returns:
// - An integer representing the current page number for pagination.
func (q *ListQuery) GetPage() int {
	return q.Page
}

// GetOrderBy returns the order by parameter for pagination.
//
// The order by parameter specifies the field to sort the data by.
// It is used in database queries to determine the order in which records are returned.
//
// Parameters:
// - q: A pointer to the ListQuery instance.
//
// Returns:
// - A string representing the order by parameter for pagination.
func (q *ListQuery) GetOrderBy() string {
	return q.OrderBy
}

// GetOffset calculates and returns the offset for pagination based on the current page and size.
// The offset is used to determine the starting position in the dataset when fetching data from a database.
//
// The formula used to calculate the offset is:
// offset = (page - 1) * size
//
// If the current page is 0, the function returns 0 as the offset.
//
// Parameters:
// - q: A pointer to the ListQuery instance.
//
// Returns:
// - An integer representing the offset for pagination.
func (q *ListQuery) GetOffset() int {
	if q.Page == 0 {
		return 0
	}
	return (q.Page - 1) * q.Size
}

// GetLimit calculates and returns the limit for pagination based on the current size.
// The limit is used to determine the maximum number of records to fetch from a database.
//
// The function simply returns the value of the size attribute of the ListQuery instance.
//
// Parameters:
// - q: A pointer to the ListQuery instance.
//
// Returns:
// - An integer representing the limit for pagination.
func (q *ListQuery) GetLimit() int {
	return q.Size
}

// ListResultToDTO converts a ListResult of type TModel to a ListResult of type TDTO.
// It uses the provided mapper to map the data from TModel to TDTO.
//
// Parameters:
// - listResult: A pointer to the ListResult[TModel] instance to be converted.
//
// Returns:
// - A pointer to a new ListResult[TDTO] instance containing the mapped data.
// - An error if any issues occur during the mapping process.
func ListResultToDTO[TDTO any, TModel any](listResult *ListResult[TModel]) (*ListResult[TDTO], error) {
	data, err := mapper.Map[[]TModel, []TDTO](listResult.Data)
	if err != nil {
		return nil, err
	}

	return &ListResult[TDTO]{
		Size:       listResult.Size,
		Page:       listResult.Page,
		TotalCount: listResult.TotalCount,
		TotalPages: listResult.TotalPages,
		Data:       data,
	}, nil
}

// ApplyFilterAction applies a filter action to a GORM query based on the given column, value, and action.
//
// Parameters:
// - query: A pointer to the GORM DB instance to apply the filter action.
// - column: A string representing the column to apply the filter action.
// - value: A string representing the value to compare against the column.
// - action: A string representing the filter action to apply.
//
// Returns:
// - A pointer to the updated GORM DB instance with the applied filter action.
// - An error if any issues occur during the filter action application.
//
// Supported filter actions:
// - equals: Filters records where the column equals the given value.
// - not_equals: Filters records where the column does not equal the given value.
// - greater_than: Filters records where the column is greater than the given value.
// - less_than: Filters records where the column is less than the given value.
// - like: Filters records where the column contains the given value as a substring.
// - in: Filters records where the column is in the given list of values.
// - not_in: Filters records where the column is not in the given list of values.
// - range: Filters records where the column is within the given range of values.
// - is_null: Filters records where the column is NULL.
// - is_not_null: Filters records where the column is not NULL.
// - starts_with: Filters records where the column starts with the given value.
// - ends_with: Filters records where the column ends with the given value.
// - ilike: Filters records where the column contains the given value as a substring (case-insensitive).
// - not_ilike: Filters records where the column does not contain the given value as a substring (case-insensitive).
// - similar_to: Filters records where the column matches the given value using the SIMILAR TO operator.
// - not_similar_to: Filters records where the column does not match the given value using the SIMILAR TO operator.
// - contains: Filters records where the column contains the given value as a substring (array type).
// - contained_in: Filters records where the column is contained in the given list of values (array type).
// - overlap: Filters records where the column overlaps with the given value (array type).
// - distinct_from: Filters records where the column is distinct from the given value.
// - not_distinct_from: Filters records where the column is not distinct from the given value.
// - is_true: Filters records where the column is TRUE.
// - is_not_true: Filters records where the column is not TRUE.
// - is_false: Filters records where the column is FALSE.
// - is_not_false: Filters records where the column is not FALSE.
// - is_unknown: Filters records where the column is UNKNOWN.
// - is_not_unknown: Filters records where the column is not UNKNOWN.
// - is_positive: Filters records where the column is greater than 0.
// - is_negative: Filters records where the column is less than 0.
// - is_not_positive: Filters records where the column is not greater than 0.
// - is_not_negative: Filters records where the column is not less than 0.
// - is_even: Filters records where the column is an even number.
// - is_odd: Filters records where the column is an odd number.
// - is_divisible_by: Filters records where the column is divisible by the given divisor value.
func ApplyFilterAction(query *gorm.DB, column string, value string, action string) (*gorm.DB, error) {
	switch action {
	case "equals":
		whereQuery := fmt.Sprintf("%s = ?", column)
		return query.Where(whereQuery, value), nil
	case "not_equals":
		whereQuery := fmt.Sprintf("%s != ?", column)
		return query.Where(whereQuery, value), nil
	case "greater_than":
		whereQuery := fmt.Sprintf("%s > ?", column)
		return query.Where(whereQuery, value), nil
	case "less_than":
		whereQuery := fmt.Sprintf("%s < ?", column)
		return query.Where(whereQuery, value), nil
	case "like":
		whereQuery := fmt.Sprintf("%s LIKE ?", column)
		return query.Where(whereQuery, "%"+value+"%"), nil
	case "in":
		whereQuery := fmt.Sprintf("%s IN (?)", column)
		queryArray := strings.Split(value, ",")
		return query.Where(whereQuery, queryArray), nil
	case "not_in":
		whereQuery := fmt.Sprintf("%s NOT IN (?)", column)
		queryArray := strings.Split(value, ",")
		return query.Where(whereQuery, queryArray), nil
	case "range":
		whereQuery := fmt.Sprintf("%s BETWEEN ? AND ?", column)
		rangeArray := strings.Split(value, ",")

		minValue, err := strconv.ParseFloat(rangeArray[0], 64)
		if err != nil {
			return nil, errors.Wrap(err, "failed to parse min value")
		}
		maxValue, err := strconv.ParseFloat(rangeArray[1], 64)
		if err != nil {
			return nil, errors.Wrap(err, "failed to parse max value")
		}

		return query.Where(whereQuery, minValue, maxValue), nil
	case "is_null":
		whereQuery := fmt.Sprintf("%s IS NULL", column)
		return query.Where(whereQuery), nil
	case "is_not_null":
		whereQuery := fmt.Sprintf("%s IS NOT NULL", column)
		return query.Where(whereQuery), nil
	case "starts_with":
		whereQuery := fmt.Sprintf("%s LIKE ?", column)
		return query.Where(whereQuery, value+"%"), nil
	case "ends_with":
		whereQuery := fmt.Sprintf("%s LIKE ?", column)
		return query.Where(whereQuery, "%"+value), nil
	case "ilike":
		whereQuery := fmt.Sprintf("%s ILIKE ?", column)
		return query.Where(whereQuery, "%"+value+"%"), nil
	case "not_ilike":
		whereQuery := fmt.Sprintf("%s NOT ILIKE ?", column)
		return query.Where(whereQuery, "%"+value+"%"), nil
	case "similar_to":
		whereQuery := fmt.Sprintf("%s SIMILAR TO ?", column)
		return query.Where(whereQuery, value), nil
	case "not_similar_to":
		whereQuery := fmt.Sprintf("%s NOT SIMILAR TO ?", column)
		return query.Where(whereQuery, value), nil
	case "contains":
		whereQuery := fmt.Sprintf("%s @> ?", column)
		return query.Where(whereQuery, value), nil
	case "contained_in":
		whereQuery := fmt.Sprintf("%s <@ ?", column)
		return query.Where(whereQuery, value), nil
	case "overlap":
		whereQuery := fmt.Sprintf("%s && ?", column)
		return query.Where(whereQuery, value), nil
	case "distinct_from":
		whereQuery := fmt.Sprintf("%s IS DISTINCT FROM ?", column)
		return query.Where(whereQuery, value), nil
	case "not_distinct_from":
		whereQuery := fmt.Sprintf("%s IS NOT DISTINCT FROM ?", column)
		return query.Where(whereQuery, value), nil
	case "is_true":
		whereQuery := fmt.Sprintf("%s IS TRUE", column)
		return query.Where(whereQuery), nil
	case "is_not_true":
		whereQuery := fmt.Sprintf("%s IS NOT TRUE", column)
		return query.Where(whereQuery), nil
	case "is_false":
		whereQuery := fmt.Sprintf("%s IS FALSE", column)
		return query.Where(whereQuery), nil
	case "is_not_false":
		whereQuery := fmt.Sprintf("%s IS NOT FALSE", column)
		return query.Where(whereQuery), nil
	case "is_unknown":
		whereQuery := fmt.Sprintf("%s IS UNKNOWN", column)
		return query.Where(whereQuery), nil
	case "is_not_unknown":
		whereQuery := fmt.Sprintf("%s IS NOT UNKNOWN", column)
		return query.Where(whereQuery), nil
	case "is_positive":
		whereQuery := fmt.Sprintf("%s > 0", column)
		return query.Where(whereQuery), nil
	case "is_negative":
		whereQuery := fmt.Sprintf("%s < 0", column)
		return query.Where(whereQuery), nil
	case "is_not_positive":
		whereQuery := fmt.Sprintf("%s <= 0", column)
		return query.Where(whereQuery), nil
	case "is_not_negative":
		whereQuery := fmt.Sprintf("%s >= 0", column)
		return query.Where(whereQuery), nil
	case "is_even":
		whereQuery := fmt.Sprintf("%s %% 2 = 0", column)
		return query.Where(whereQuery), nil
	case "is_odd":
		whereQuery := fmt.Sprintf("%s %% 2 != 0", column)
		return query.Where(whereQuery), nil
	case "is_divisible_by":
		divisor, err := strconv.Atoi(value)
		if err != nil {
			return nil, errors.Wrap(err, "failed to parse divisor value")
		}
		whereQuery := fmt.Sprintf("%s %% ? = 0", column)
		return query.Where(whereQuery, divisor), nil
	default:
		return nil, errors.New("unsupported filter action")
	}
}
