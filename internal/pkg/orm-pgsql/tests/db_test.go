package tests

import (
	"context"
	"database/sql/driver"
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	postgrescontainer "github.com/NekKkMirror/go-app/internal/pkg/container/test/postgres"
	"github.com/NekKkMirror/go-app/internal/pkg/orm-pgsql"
	"github.com/NekKkMirror/go-app/internal/pkg/utils/db/pagination"
	"gorm.io/gorm"
)

func TestPaginateWithFilters(t *testing.T) {
	ctx := context.Background()
	DB, m, err := postgrescontainer.Start(ctx, t)
	if err != nil {
		t.Fatalf("expected no error from DB, got %v", err)
	}

	listQuery1 := &pagination.ListQuery{
		Size:    10,
		Page:    1,
		OrderBy: "created_at DESC",
		Filters: []*pagination.FilterModel{
			{Field: "age", Value: "30", Comparison: "="},
		},
	}
	testPaginationWithFilters(t, ctx, listQuery1, DB, m)

	listQuery2 := &pagination.ListQuery{
		Size:    10,
		Page:    1,
		OrderBy: "created_at DESC",
		Filters: []*pagination.FilterModel{
			{Field: "name", Value: "Alice", Comparison: "starts_with"},
		},
	}
	testPaginationWithFilters(t, ctx, listQuery2, DB, m)

	listQuery3 := &pagination.ListQuery{
		Size:    10,
		Page:    1,
		OrderBy: "created_at DESC",
		Filters: []*pagination.FilterModel{
			{Field: "total_spent", Value: "1000", Comparison: ">"},
		},
	}
	testPaginationWithFilters(t, ctx, listQuery3, DB, m)

	listQuery4 := &pagination.ListQuery{
		Size:    10,
		Page:    1,
		OrderBy: "created_at DESC",
		Filters: []*pagination.FilterModel{
			{Field: "total_spent", Value: "500", Comparison: "<"},
		},
	}
	testPaginationWithFilters(t, ctx, listQuery4, DB, m)

	listQuery5 := &pagination.ListQuery{
		Size:    10,
		Page:    1,
		OrderBy: "created_at DESC",
		Filters: []*pagination.FilterModel{
			{Field: "name", Value: "John", Comparison: "!="},
		},
	}
	testPaginationWithFilters(t, ctx, listQuery5, DB, m)

	listQuery6 := &pagination.ListQuery{
		Size:    10,
		Page:    1,
		OrderBy: "created_at DESC",
		Filters: []*pagination.FilterModel{
			{Field: "age", Value: "25,35", Comparison: "between"},
		},
	}
	testPaginationWithFilters(t, ctx, listQuery6, DB, m)

	listQuery7 := &pagination.ListQuery{
		Size:    10,
		Page:    1,
		OrderBy: "created_at DESC",
		Filters: []*pagination.FilterModel{
			{Field: "name", Value: "Doe", Comparison: "ends_with"},
		},
	}
	testPaginationWithFilters(t, ctx, listQuery7, DB, m)

	listQuery8 := &pagination.ListQuery{
		Size:    10,
		Page:    1,
		OrderBy: "created_at DESC",
		Filters: []*pagination.FilterModel{
			{Field: "email", Value: "example.com", Comparison: "ilike"},
		},
	}
	testPaginationWithFilters(t, ctx, listQuery8, DB, m)

	listQuery9 := &pagination.ListQuery{
		Size:    10,
		Page:    1,
		OrderBy: "created_at DESC",
		Filters: []*pagination.FilterModel{
			{Field: "age", Value: "40", Comparison: "is_not_null"},
		},
	}
	testPaginationWithFilters(t, ctx, listQuery9, DB, m)

	listQuery10 := &pagination.ListQuery{
		Size:    10,
		Page:    1,
		OrderBy: "created_at DESC",
		Filters: []*pagination.FilterModel{
			{Field: "total_spent", Value: "2000", Comparison: "not_in"},
		},
	}
	testPaginationWithFilters(t, ctx, listQuery10, DB, m)

	listQuery11 := &pagination.ListQuery{
		Size:    10,
		Page:    1,
		OrderBy: "created_at DESC",
		Filters: []*pagination.FilterModel{
			{Field: "is_active", Comparison: "is_true"},
		},
	}
	testPaginationWithFilters(t, ctx, listQuery11, DB, m)

	listQuery12 := &pagination.ListQuery{
		Size:    10,
		Page:    1,
		OrderBy: "created_at DESC",
		Filters: []*pagination.FilterModel{
			{Field: "is_active", Comparison: "is_false"},
		},
	}
	testPaginationWithFilters(t, ctx, listQuery12, DB, m)

	listQuery13 := &pagination.ListQuery{
		Size:    10,
		Page:    1,
		OrderBy: "created_at DESC",
		Filters: []*pagination.FilterModel{
			{Field: "is_admin", Comparison: "is_not_false"},
		},
	}
	testPaginationWithFilters(t, ctx, listQuery13, DB, m)

	listQuery14 := &pagination.ListQuery{
		Size:    10,
		Page:    1,
		OrderBy: "created_at DESC",
		Filters: []*pagination.FilterModel{
			{Field: "is_active", Comparison: "is_unknown"},
		},
	}
	testPaginationWithFilters(t, ctx, listQuery14, DB, m)

	listQuery15 := &pagination.ListQuery{
		Size:    10,
		Page:    1,
		OrderBy: "created_at DESC",
		Filters: []*pagination.FilterModel{
			{Field: "is_active", Comparison: "is_not_unknown"},
		},
	}
	testPaginationWithFilters(t, ctx, listQuery15, DB, m)

	listQuery16 := &pagination.ListQuery{
		Size:    10,
		Page:    1,
		OrderBy: "created_at DESC",
		Filters: []*pagination.FilterModel{
			{Field: "age", Value: "0", Comparison: "is_positive"},
		},
	}
	testPaginationWithFilters(t, ctx, listQuery16, DB, m)

	listQuery17 := &pagination.ListQuery{
		Size:    10,
		Page:    1,
		OrderBy: "created_at DESC",
		Filters: []*pagination.FilterModel{
			{Field: "age", Value: "0", Comparison: "is_negative"},
		},
	}
	testPaginationWithFilters(t, ctx, listQuery17, DB, m)

	listQuery18 := &pagination.ListQuery{
		Size:    10,
		Page:    1,
		OrderBy: "created_at DESC",
		Filters: []*pagination.FilterModel{
			{Field: "age", Value: "0", Comparison: "is_not_positive"},
		},
	}
	testPaginationWithFilters(t, ctx, listQuery18, DB, m)

	listQuery19 := &pagination.ListQuery{
		Size:    10,
		Page:    1,
		OrderBy: "created_at DESC",
		Filters: []*pagination.FilterModel{
			{Field: "age", Value: "0", Comparison: "is_not_negative"},
		},
	}
	testPaginationWithFilters(t, ctx, listQuery19, DB, m)

	listQuery20 := &pagination.ListQuery{
		Size:    10,
		Page:    1,
		OrderBy: "created_at DESC",
		Filters: []*pagination.FilterModel{
			{Field: "age", Value: "0", Comparison: "is_even"},
		},
	}
	testPaginationWithFilters(t, ctx, listQuery20, DB, m)

	listQuery21 := &pagination.ListQuery{
		Size:    10,
		Page:    1,
		OrderBy: "created_at DESC",
		Filters: []*pagination.FilterModel{
			{Field: "age", Value: "0", Comparison: "is_odd"},
		},
	}
	testPaginationWithFilters(t, ctx, listQuery21, DB, m)

	listQuery22 := &pagination.ListQuery{
		Size:    10,
		Page:    1,
		OrderBy: "created_at DESC",
		Filters: []*pagination.FilterModel{
			{Field: "age", Value: "2", Comparison: "is_divisible_by"},
		},
	}
	testPaginationWithFilters(t, ctx, listQuery22, DB, m)

	listQuery23 := &pagination.ListQuery{
		Size:    10,
		Page:    1,
		OrderBy: "created_at DESC",
		Filters: []*pagination.FilterModel{
			{Field: "is_admin", Comparison: "is_not_true"},
		},
	}
	testPaginationWithFilters(t, ctx, listQuery23, DB, m)
}

func testPaginationWithFilters(t *testing.T, ctx context.Context, listQuery *pagination.ListQuery, DB *gorm.DB, m sqlmock.Sqlmock) {
	var expectedQuery string
	var expectedArgs []driver.Value

	for _, filter := range listQuery.Filters {
		switch filter.Comparison {
		case "is_divisible_by":
			expectedQuery += fmt.Sprintf(" AND %s %s?", filter.Field, "mod")
			expectedArgs = append(expectedArgs, filter.Value, 0)
		case "is_not_divisible_by":
		case "is_negative":
			expectedQuery += fmt.Sprintf(" AND %s < 0", filter.Field)
			expectedArgs = append(expectedArgs, 0)
		case "is_not_negative":
			expectedQuery += fmt.Sprintf(" AND %s >= 0", filter.Field)
			expectedArgs = append(expectedArgs, 0)
		case "is_positive":
			expectedQuery += fmt.Sprintf(" AND %s > 0", filter.Field)
			expectedArgs = append(expectedArgs, 0)
		case "is_not_positive":
			expectedQuery += fmt.Sprintf(" AND %s <= 0", filter.Field)
			expectedArgs = append(expectedArgs, 0)
		case "is_unknown":
			expectedQuery += fmt.Sprintf(" AND %s IS NULL", filter.Field)
			expectedArgs = append(expectedArgs, nil)
		case "is_not_unknown":
			expectedQuery += fmt.Sprintf(" AND %s IS NOT NULL", filter.Field)
			expectedArgs = append(expectedArgs, nil)
		case "is_false":
			expectedQuery += fmt.Sprintf(" AND %s IS FALSE", filter.Field)
			expectedArgs = append(expectedArgs, false)
		case "is_not_false":
			expectedQuery += fmt.Sprintf(" AND %s IS NOT FALSE", filter.Field)
			expectedArgs = append(expectedArgs, false)
		case "is_true":
			expectedQuery += fmt.Sprintf(" AND %s IS TRUE", filter.Field)
			expectedArgs = append(expectedArgs, true)
		case "is_not_true":
			expectedQuery += fmt.Sprintf(" AND %s IS NOT TRUE", filter.Field)
			expectedArgs = append(expectedArgs, true)
		case "is_even":
			expectedQuery += fmt.Sprintf(" AND %s %s 0", filter.Field, "%")
			expectedArgs = append(expectedArgs, 0)
		case "is_odd":
			expectedQuery += fmt.Sprintf(" AND %s %s 1", filter.Field, "%")
			expectedArgs = append(expectedArgs, 1)
		case "=":
			expectedQuery += fmt.Sprintf(" AND %s = ?", filter.Field)
			expectedArgs = append(expectedArgs, filter.Value)
		case "!=":
			expectedQuery += fmt.Sprintf(" AND %s != ?", filter.Field)
			expectedArgs = append(expectedArgs, filter.Value)
		case "starts_with":
			expectedQuery += fmt.Sprintf(" AND %s LIKE ?", filter.Field)
			expectedArgs = append(expectedArgs, filter.Value+"%")
		case "ends_with":
			expectedQuery += fmt.Sprintf(" AND %s LIKE ?", filter.Field)
			expectedArgs = append(expectedArgs, "%"+filter.Value)
		case "contains":
		case ">":
			expectedQuery += fmt.Sprintf(" AND %s > ?", filter.Field)
			expectedArgs = append(expectedArgs, filter.Value)
		case "<":
			expectedQuery += fmt.Sprintf(" AND %s < ?", filter.Field)
			expectedArgs = append(expectedArgs, filter.Value)
		case "like":
			expectedQuery += fmt.Sprintf(" AND %s LIKE ?", filter.Field)
			expectedArgs = append(expectedArgs, "%"+filter.Value+"%")
		case "ilike":
			expectedQuery += fmt.Sprintf(" AND %s ILIKE ?", filter.Field)
			expectedArgs = append(expectedArgs, "%"+filter.Value+"%")
		case "is_null":
			expectedQuery += fmt.Sprintf(" AND %s IS NULL", filter.Field)
			expectedArgs = append(expectedArgs, nil)
		case "is_not_null":
			expectedQuery += fmt.Sprintf(" AND %s IS NOT NULL", filter.Field)
			expectedArgs = append(expectedArgs, nil)
		case "in":
			queryArray := strings.Split(filter.Value, ",")
			expectedQuery += fmt.Sprintf(" AND %s IN (?)", filter.Field)
			expectedArgs = append(expectedArgs, queryArray)
		case "not_in":
			queryArray := strings.Split(filter.Value, ",")
			expectedQuery += fmt.Sprintf(" AND %s NOT IN (?)", filter.Field)
			expectedArgs = append(expectedArgs, queryArray)
		case "between":
			rangeArray := strings.Split(filter.Value, ",")
			minValue, err := strconv.ParseFloat(rangeArray[0], 64)
			if err != nil {
				t.Errorf("failed to parse min value: %v", err)
				return
			}
			maxValue, err := strconv.ParseFloat(rangeArray[1], 64)
			if err != nil {
				t.Errorf("failed to parse max value: %v", err)
				return
			}
			expectedQuery += fmt.Sprintf(" AND %s BETWEEN ? AND ?", filter.Field)
			expectedArgs = append(expectedArgs, minValue, maxValue)
		default:
			t.Errorf("unsupported filter comparison: %s", filter.Comparison)
			return
		}
	}

	m.ExpectQuery("SELECT * FROM user WHERE 1=1" + expectedQuery).
		WithArgs(expectedArgs...).
		WillReturnRows(sqlmock.NewRows([]string{"age"}).AddRow(30))

	result, err := ormpgsql.Paginate[postgrescontainer.User](listQuery, DB)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if result == nil {
		t.Error("pagination result is nil")
	}
}
