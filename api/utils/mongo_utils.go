package utils

import (
	"context"
	"math"
	"net/http"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PaginationParams struct {
	Page     int64
	Limit    int64
	SortBy   string
	SortDir  string
	Search   string
	Priority string
	Status   string
}

type DateRange struct {
	StartDate string
	EndDate   string
}

func GetPaginationFromRequest(r *http.Request) PaginationParams {
	query := r.URL.Query()
	page, _ := strconv.ParseInt(query.Get("page"), 10, 64)
	limit, _ := strconv.ParseInt(query.Get("limit"), 10, 64)

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	return PaginationParams{
		Page:     page,
		Limit:    limit,
		SortBy:   query.Get("sort_by"),
		SortDir:  query.Get("sort_dir"),
		Search:   query.Get("search"),
		Priority: query.Get("priority"),
		Status:   query.Get("status"),
	}
}

func BuildSearchFilter(baseFilter bson.M, params PaginationParams, dateRange *DateRange) bson.M {
	filter := baseFilter

	if params.Search != "" {
		filter["$or"] = []bson.M{
			{"title": bson.M{"$regex": params.Search, "$options": "i"}},
			{"description": bson.M{"$regex": params.Search, "$options": "i"}},
		}
	}

	if params.Priority != "" {
		filter["priority"] = params.Priority
	}

	if params.Status != "" {
		filter["status"] = params.Status
	}

	if dateRange != nil {
		dateFilter := bson.M{}
		if dateRange.StartDate != "" {
			startTime, err := time.Parse(time.RFC3339, dateRange.StartDate)
			if err == nil {
				dateFilter["$gte"] = startTime
			}
		}
		if dateRange.EndDate != "" {
			endTime, err := time.Parse(time.RFC3339, dateRange.EndDate)
			if err == nil {
				dateFilter["$lte"] = endTime
			}
		}
		if len(dateFilter) > 0 {
			filter["due_date"] = dateFilter
		}
	}

	return filter
}

func BuildSortOptions(params PaginationParams) bson.M {
	sortOptions := bson.M{"created_at": -1} // default sort
	if params.SortBy != "" {
		sortDirection := 1 // ascending
		if params.SortDir == "desc" {
			sortDirection = -1
		}
		switch params.SortBy {
		case "due_date", "priority", "status", "title":
			sortOptions = bson.M{params.SortBy: sortDirection}
		}
	}
	return sortOptions
}

func ExecutePaginatedQuery(ctx context.Context, collection *mongo.Collection, filter bson.M, params PaginationParams) ([]interface{}, int64, error) {
	skip := (params.Page - 1) * params.Limit

	// Get total count
	total, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	// Execute query
	cursor, err := collection.Find(ctx,
		filter,
		options.Find().
			SetSkip(skip).
			SetLimit(params.Limit).
			SetSort(BuildSortOptions(params)),
	)
	if err != nil {
		return nil, 0, err
	}

	var results []interface{}
	if err = cursor.All(ctx, &results); err != nil {
		return nil, 0, err
	}

	return results, total, nil
}

func CalculateTotalPages(total int64, limit int64) int64 {
	return int64(math.Ceil(float64(total) / float64(limit)))
}
