// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: route_queries.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const deleteStaleRoutes = `-- name: DeleteStaleRoutes :exec
DELETE FROM route
WHERE 
    route.feed_pk = $1
    AND NOT route.pk = ANY($2::bigint[])
`

type DeleteStaleRoutesParams struct {
	FeedPk          int64
	UpdatedRoutePks []int64
}

func (q *Queries) DeleteStaleRoutes(ctx context.Context, arg DeleteStaleRoutesParams) error {
	_, err := q.db.Exec(ctx, deleteStaleRoutes, arg.FeedPk, arg.UpdatedRoutePks)
	return err
}

const estimateHeadwaysForRoutes = `-- name: EstimateHeadwaysForRoutes :many
WITH per_stop_data AS (
    SELECT
        trip.route_pk route_pk,
        EXTRACT(epoch FROM MAX(trip_stop_time.arrival_time) - MIN(trip_stop_time.arrival_time)) total_diff,
        COUNT(*)-1 num_diffs
    FROM trip_stop_time
        INNER JOIN trip ON trip.pk = trip_stop_time.trip_pk
    WHERE trip.route_pk = ANY($1::bigint[])
        AND NOT trip_stop_time.past
        AND trip_stop_time.arrival_time IS NOT NULL
        AND trip_stop_time.arrival_time >= $2
    GROUP BY trip_stop_time.stop_pk, trip.route_pk
        HAVING COUNT(*) > 1
)
SELECT 
    route_pk,
    COALESCE(ROUND(SUM(total_diff) / (SUM(num_diffs)))::integer, -1)::integer estimated_headway
FROM per_stop_data
GROUP BY route_pk
`

type EstimateHeadwaysForRoutesParams struct {
	RoutePks    []int64
	PresentTime pgtype.Timestamptz
}

type EstimateHeadwaysForRoutesRow struct {
	RoutePk          int64
	EstimatedHeadway int32
}

func (q *Queries) EstimateHeadwaysForRoutes(ctx context.Context, arg EstimateHeadwaysForRoutesParams) ([]EstimateHeadwaysForRoutesRow, error) {
	rows, err := q.db.Query(ctx, estimateHeadwaysForRoutes, arg.RoutePks, arg.PresentTime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []EstimateHeadwaysForRoutesRow
	for rows.Next() {
		var i EstimateHeadwaysForRoutesRow
		if err := rows.Scan(&i.RoutePk, &i.EstimatedHeadway); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getRoute = `-- name: GetRoute :one
SELECT route.pk, route.id, route.system_pk, route.color, route.text_color, route.short_name, route.long_name, route.description, route.url, route.sort_order, route.type, route.agency_pk, route.continuous_drop_off, route.continuous_pickup, route.feed_pk FROM route
    INNER JOIN system ON route.system_pk = system.pk
    WHERE system.pk = $1
    AND route.id = $2
`

type GetRouteParams struct {
	SystemPk int64
	RouteID  string
}

func (q *Queries) GetRoute(ctx context.Context, arg GetRouteParams) (Route, error) {
	row := q.db.QueryRow(ctx, getRoute, arg.SystemPk, arg.RouteID)
	var i Route
	err := row.Scan(
		&i.Pk,
		&i.ID,
		&i.SystemPk,
		&i.Color,
		&i.TextColor,
		&i.ShortName,
		&i.LongName,
		&i.Description,
		&i.Url,
		&i.SortOrder,
		&i.Type,
		&i.AgencyPk,
		&i.ContinuousDropOff,
		&i.ContinuousPickup,
		&i.FeedPk,
	)
	return i, err
}

const insertRoute = `-- name: InsertRoute :one
INSERT INTO route
    (id, system_pk, feed_pk, color, text_color,
     short_name, long_name, description, url, sort_order,
     type, continuous_pickup, continuous_drop_off, agency_pk)
VALUES
    ($1, $2, $3, $4, $5,
     $6, $7, $8, $9, $10,
     $11, $12,$13, $14)
RETURNING pk
`

type InsertRouteParams struct {
	ID                string
	SystemPk          int64
	FeedPk            int64
	Color             string
	TextColor         string
	ShortName         pgtype.Text
	LongName          pgtype.Text
	Description       pgtype.Text
	Url               pgtype.Text
	SortOrder         pgtype.Int4
	Type              string
	ContinuousPickup  string
	ContinuousDropOff string
	AgencyPk          int64
}

func (q *Queries) InsertRoute(ctx context.Context, arg InsertRouteParams) (int64, error) {
	row := q.db.QueryRow(ctx, insertRoute,
		arg.ID,
		arg.SystemPk,
		arg.FeedPk,
		arg.Color,
		arg.TextColor,
		arg.ShortName,
		arg.LongName,
		arg.Description,
		arg.Url,
		arg.SortOrder,
		arg.Type,
		arg.ContinuousPickup,
		arg.ContinuousDropOff,
		arg.AgencyPk,
	)
	var pk int64
	err := row.Scan(&pk)
	return pk, err
}

const listRoutes = `-- name: ListRoutes :many
SELECT pk, id, system_pk, color, text_color, short_name, long_name, description, url, sort_order, type, agency_pk, continuous_drop_off, continuous_pickup, feed_pk FROM route WHERE system_pk = $1 ORDER BY id
`

func (q *Queries) ListRoutes(ctx context.Context, systemPk int64) ([]Route, error) {
	rows, err := q.db.Query(ctx, listRoutes, systemPk)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Route
	for rows.Next() {
		var i Route
		if err := rows.Scan(
			&i.Pk,
			&i.ID,
			&i.SystemPk,
			&i.Color,
			&i.TextColor,
			&i.ShortName,
			&i.LongName,
			&i.Description,
			&i.Url,
			&i.SortOrder,
			&i.Type,
			&i.AgencyPk,
			&i.ContinuousDropOff,
			&i.ContinuousPickup,
			&i.FeedPk,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listRoutesByPk = `-- name: ListRoutesByPk :many
SELECT route.pk, route.id id, route.color, system.id system_id
FROM route
    INNER JOIN system on route.system_pk = system.pk
WHERE route.pk = ANY($1::bigint[])
`

type ListRoutesByPkRow struct {
	Pk       int64
	ID       string
	Color    string
	SystemID string
}

func (q *Queries) ListRoutesByPk(ctx context.Context, routePks []int64) ([]ListRoutesByPkRow, error) {
	rows, err := q.db.Query(ctx, listRoutesByPk, routePks)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListRoutesByPkRow
	for rows.Next() {
		var i ListRoutesByPkRow
		if err := rows.Scan(
			&i.Pk,
			&i.ID,
			&i.Color,
			&i.SystemID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listRoutesInAgency = `-- name: ListRoutesInAgency :many
SELECT route.id, route.color FROM route
WHERE route.agency_pk = $1
`

type ListRoutesInAgencyRow struct {
	ID    string
	Color string
}

func (q *Queries) ListRoutesInAgency(ctx context.Context, agencyPk int64) ([]ListRoutesInAgencyRow, error) {
	rows, err := q.db.Query(ctx, listRoutesInAgency, agencyPk)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListRoutesInAgencyRow
	for rows.Next() {
		var i ListRoutesInAgencyRow
		if err := rows.Scan(&i.ID, &i.Color); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const mapRouteIDToPkInSystem = `-- name: MapRouteIDToPkInSystem :many
SELECT pk, id from route
WHERE
    system_pk = $1
    AND (
        NOT $2::bool
        OR id = ANY($3::text[])
    )
`

type MapRouteIDToPkInSystemParams struct {
	SystemPk        int64
	FilterByRouteID bool
	RouteIds        []string
}

type MapRouteIDToPkInSystemRow struct {
	Pk int64
	ID string
}

func (q *Queries) MapRouteIDToPkInSystem(ctx context.Context, arg MapRouteIDToPkInSystemParams) ([]MapRouteIDToPkInSystemRow, error) {
	rows, err := q.db.Query(ctx, mapRouteIDToPkInSystem, arg.SystemPk, arg.FilterByRouteID, arg.RouteIds)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []MapRouteIDToPkInSystemRow
	for rows.Next() {
		var i MapRouteIDToPkInSystemRow
		if err := rows.Scan(&i.Pk, &i.ID); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateRoute = `-- name: UpdateRoute :exec
UPDATE route SET
    feed_pk = $1,
    color = $2,
    text_color = $3,
    short_name = $4, 
    long_name = $5, 
    description = $6, 
    url = $7, 
    sort_order = $8, 
    type = $9, 
    continuous_pickup = $10, 
    continuous_drop_off = $11, 
    agency_pk = $12
WHERE
    pk = $13
`

type UpdateRouteParams struct {
	FeedPk            int64
	Color             string
	TextColor         string
	ShortName         pgtype.Text
	LongName          pgtype.Text
	Description       pgtype.Text
	Url               pgtype.Text
	SortOrder         pgtype.Int4
	Type              string
	ContinuousPickup  string
	ContinuousDropOff string
	AgencyPk          int64
	Pk                int64
}

func (q *Queries) UpdateRoute(ctx context.Context, arg UpdateRouteParams) error {
	_, err := q.db.Exec(ctx, updateRoute,
		arg.FeedPk,
		arg.Color,
		arg.TextColor,
		arg.ShortName,
		arg.LongName,
		arg.Description,
		arg.Url,
		arg.SortOrder,
		arg.Type,
		arg.ContinuousPickup,
		arg.ContinuousDropOff,
		arg.AgencyPk,
		arg.Pk,
	)
	return err
}
