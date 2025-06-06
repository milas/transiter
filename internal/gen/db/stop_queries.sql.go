// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: stop_queries.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jamespfennell/transiter/db/types"
)

const deleteStaleStops = `-- name: DeleteStaleStops :exec
DELETE FROM stop
WHERE
    stop.feed_pk = $1
    AND NOT stop.pk = ANY($2::bigint[])
`

type DeleteStaleStopsParams struct {
	FeedPk         int64
	UpdatedStopPks []int64
}

func (q *Queries) DeleteStaleStops(ctx context.Context, arg DeleteStaleStopsParams) error {
	_, err := q.db.Exec(ctx, deleteStaleStops, arg.FeedPk, arg.UpdatedStopPks)
	return err
}

const deleteWheelchairBoardingForSystem = `-- name: DeleteWheelchairBoardingForSystem :exec
UPDATE stop
SET wheelchair_boarding = NULL
WHERE system_pk = $1
`

func (q *Queries) DeleteWheelchairBoardingForSystem(ctx context.Context, systemPk int64) error {
	_, err := q.db.Exec(ctx, deleteWheelchairBoardingForSystem, systemPk)
	return err
}

const getStop = `-- name: GetStop :one
SELECT stop.pk, stop.id, stop.system_pk, stop.parent_stop_pk, stop.name, stop.url, stop.code, stop.description, stop.platform_code, stop.timezone, stop.type, stop.wheelchair_boarding, stop.zone_id, stop.feed_pk, stop.location FROM stop
    INNER JOIN system ON stop.system_pk = system.pk
    WHERE system.id = $1
    AND stop.id = $2
`

type GetStopParams struct {
	SystemID string
	StopID   string
}

func (q *Queries) GetStop(ctx context.Context, arg GetStopParams) (Stop, error) {
	row := q.db.QueryRow(ctx, getStop, arg.SystemID, arg.StopID)
	var i Stop
	err := row.Scan(
		&i.Pk,
		&i.ID,
		&i.SystemPk,
		&i.ParentStopPk,
		&i.Name,
		&i.Url,
		&i.Code,
		&i.Description,
		&i.PlatformCode,
		&i.Timezone,
		&i.Type,
		&i.WheelchairBoarding,
		&i.ZoneID,
		&i.FeedPk,
		&i.Location,
	)
	return i, err
}

const insertStop = `-- name: InsertStop :one
INSERT INTO stop
    (id, system_pk, feed_pk, name, location,
     url, code, description, platform_code, timezone, type,
     wheelchair_boarding, zone_id)
VALUES
    ($1, $2, $3, $4,
     $5::geography,
     $6, $7, $8, $9,
     $10, $11, $12, $13)
RETURNING pk
`

type InsertStopParams struct {
	ID                 string
	SystemPk           int64
	FeedPk             int64
	Name               pgtype.Text
	Location           types.Geography
	Url                pgtype.Text
	Code               pgtype.Text
	Description        pgtype.Text
	PlatformCode       pgtype.Text
	Timezone           pgtype.Text
	Type               string
	WheelchairBoarding pgtype.Bool
	ZoneID             pgtype.Text
}

func (q *Queries) InsertStop(ctx context.Context, arg InsertStopParams) (int64, error) {
	row := q.db.QueryRow(ctx, insertStop,
		arg.ID,
		arg.SystemPk,
		arg.FeedPk,
		arg.Name,
		arg.Location,
		arg.Url,
		arg.Code,
		arg.Description,
		arg.PlatformCode,
		arg.Timezone,
		arg.Type,
		arg.WheelchairBoarding,
		arg.ZoneID,
	)
	var pk int64
	err := row.Scan(&pk)
	return pk, err
}

const listStops = `-- name: ListStops :many
SELECT pk, id, system_pk, parent_stop_pk, name, url, code, description, platform_code, timezone, type, wheelchair_boarding, zone_id, feed_pk, location FROM stop
WHERE system_pk = $1
  AND id >= $2
  AND (
    NOT $3::bool OR
    id = ANY($4::text[])
  )
  AND (
    NOT $5::bool OR
    type = ANY($6::text[])
  )
ORDER BY id
LIMIT $7
`

type ListStopsParams struct {
	SystemPk     int64
	FirstStopID  string
	FilterByID   bool
	StopIds      []string
	FilterByType bool
	Types        []string
	NumStops     int32
}

func (q *Queries) ListStops(ctx context.Context, arg ListStopsParams) ([]Stop, error) {
	rows, err := q.db.Query(ctx, listStops,
		arg.SystemPk,
		arg.FirstStopID,
		arg.FilterByID,
		arg.StopIds,
		arg.FilterByType,
		arg.Types,
		arg.NumStops,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Stop
	for rows.Next() {
		var i Stop
		if err := rows.Scan(
			&i.Pk,
			&i.ID,
			&i.SystemPk,
			&i.ParentStopPk,
			&i.Name,
			&i.Url,
			&i.Code,
			&i.Description,
			&i.PlatformCode,
			&i.Timezone,
			&i.Type,
			&i.WheelchairBoarding,
			&i.ZoneID,
			&i.FeedPk,
			&i.Location,
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

const listStopsByPk = `-- name: ListStopsByPk :many
SELECT stop.pk, stop.id stop_id, stop.name, system.id system_id
FROM stop
    INNER JOIN system on stop.system_pk = system.pk
WHERE stop.pk = ANY($1::bigint[])
`

type ListStopsByPkRow struct {
	Pk       int64
	StopID   string
	Name     pgtype.Text
	SystemID string
}

func (q *Queries) ListStopsByPk(ctx context.Context, stopPks []int64) ([]ListStopsByPkRow, error) {
	rows, err := q.db.Query(ctx, listStopsByPk, stopPks)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListStopsByPkRow
	for rows.Next() {
		var i ListStopsByPkRow
		if err := rows.Scan(
			&i.Pk,
			&i.StopID,
			&i.Name,
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

const listStops_Geographic = `-- name: ListStops_Geographic :many
WITH distance AS (
    SELECT
        stop.pk stop_pk,
        stop.location <-> $4::geography distance
    FROM stop
    WHERE stop.location IS NOT NULL
    AND (
        NOT $5::bool OR
        type = ANY($6::text[])
    )
)
SELECT stop.pk, stop.id, stop.system_pk, stop.parent_stop_pk, stop.name, stop.url, stop.code, stop.description, stop.platform_code, stop.timezone, stop.type, stop.wheelchair_boarding, stop.zone_id, stop.feed_pk, stop.location FROM stop
INNER JOIN distance ON stop.pk = distance.stop_pk
WHERE stop.system_pk = $1
    AND distance.distance <= 1000 * $2::float
ORDER by distance.distance
LIMIT $3
`

type ListStops_GeographicParams struct {
	SystemPk     int64
	MaxDistance  float64
	MaxResults   int32
	Base         types.Geography
	FilterByType bool
	Types        []string
}

func (q *Queries) ListStops_Geographic(ctx context.Context, arg ListStops_GeographicParams) ([]Stop, error) {
	rows, err := q.db.Query(ctx, listStops_Geographic,
		arg.SystemPk,
		arg.MaxDistance,
		arg.MaxResults,
		arg.Base,
		arg.FilterByType,
		arg.Types,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Stop
	for rows.Next() {
		var i Stop
		if err := rows.Scan(
			&i.Pk,
			&i.ID,
			&i.SystemPk,
			&i.ParentStopPk,
			&i.Name,
			&i.Url,
			&i.Code,
			&i.Description,
			&i.PlatformCode,
			&i.Timezone,
			&i.Type,
			&i.WheelchairBoarding,
			&i.ZoneID,
			&i.FeedPk,
			&i.Location,
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

const listTripStopTimesByStops = `-- name: ListTripStopTimesByStops :many
SELECT trip_stop_time.pk, trip_stop_time.stop_pk, trip_stop_time.trip_pk, trip_stop_time.arrival_time, trip_stop_time.arrival_delay, trip_stop_time.arrival_uncertainty, trip_stop_time.departure_time, trip_stop_time.departure_delay, trip_stop_time.departure_uncertainty, trip_stop_time.stop_sequence, trip_stop_time.track, trip_stop_time.headsign, trip_stop_time.past,
       trip.pk, trip.id, trip.route_pk, trip.direction_id, trip.started_at, trip.gtfs_hash, trip.feed_pk, vehicle.id vehicle_id,
       vehicle.location::geography vehicle_location,
       vehicle.bearing vehicle_bearing,
       vehicle.updated_at vehicle_updated_at,
       COALESCE(scheduled_trip_stop_time.headsign, scheduled_trip.headsign) scheduled_trip_headsign
    FROM trip_stop_time
    INNER JOIN trip ON trip_stop_time.trip_pk = trip.pk
    LEFT JOIN vehicle ON vehicle.trip_pk = trip.pk
    LEFT JOIN scheduled_trip ON scheduled_trip.id = trip.id AND scheduled_trip.route_pk = trip.route_pk
    LEFT JOIN scheduled_trip_stop_time ON scheduled_trip_stop_time.trip_pk = scheduled_trip.pk AND
                                          scheduled_trip_stop_time.stop_pk = trip_stop_time.stop_pk AND
                                          scheduled_trip_stop_time.stop_sequence = trip_stop_time.stop_sequence
    WHERE trip_stop_time.stop_pk = ANY($1::bigint[])
    AND NOT trip_stop_time.past
    ORDER BY COALESCE(trip_stop_time.arrival_time, trip_stop_time.departure_time)
`

type ListTripStopTimesByStopsRow struct {
	Pk                    int64
	StopPk                int64
	TripPk                int64
	ArrivalTime           pgtype.Timestamptz
	ArrivalDelay          pgtype.Int4
	ArrivalUncertainty    pgtype.Int4
	DepartureTime         pgtype.Timestamptz
	DepartureDelay        pgtype.Int4
	DepartureUncertainty  pgtype.Int4
	StopSequence          int32
	Track                 pgtype.Text
	Headsign              pgtype.Text
	Past                  bool
	Pk_2                  int64
	ID                    string
	RoutePk               int64
	DirectionID           pgtype.Bool
	StartedAt             pgtype.Timestamptz
	GtfsHash              string
	FeedPk                int64
	VehicleID             pgtype.Text
	VehicleLocation       types.Geography
	VehicleBearing        pgtype.Float4
	VehicleUpdatedAt      pgtype.Timestamptz
	ScheduledTripHeadsign pgtype.Text
}

func (q *Queries) ListTripStopTimesByStops(ctx context.Context, stopPks []int64) ([]ListTripStopTimesByStopsRow, error) {
	rows, err := q.db.Query(ctx, listTripStopTimesByStops, stopPks)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListTripStopTimesByStopsRow
	for rows.Next() {
		var i ListTripStopTimesByStopsRow
		if err := rows.Scan(
			&i.Pk,
			&i.StopPk,
			&i.TripPk,
			&i.ArrivalTime,
			&i.ArrivalDelay,
			&i.ArrivalUncertainty,
			&i.DepartureTime,
			&i.DepartureDelay,
			&i.DepartureUncertainty,
			&i.StopSequence,
			&i.Track,
			&i.Headsign,
			&i.Past,
			&i.Pk_2,
			&i.ID,
			&i.RoutePk,
			&i.DirectionID,
			&i.StartedAt,
			&i.GtfsHash,
			&i.FeedPk,
			&i.VehicleID,
			&i.VehicleLocation,
			&i.VehicleBearing,
			&i.VehicleUpdatedAt,
			&i.ScheduledTripHeadsign,
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

const mapStopIDAndPkToStationPk = `-- name: MapStopIDAndPkToStationPk :many
WITH RECURSIVE
ancestor AS (
    SELECT
    id stop_id,
    pk stop_pk,
    pk station_pk,
    parent_stop_pk,
    (type = 'STATION') is_station
    FROM stop
        WHERE stop.system_pk = $1
        AND (
            NOT $2::bool
            OR stop.pk = ANY($3::bigint[])
        )
    UNION
    SELECT
    child.stop_id stop_id,
    child.stop_pk stop_pk,
    parent.pk station_pk,
    parent.parent_stop_pk,
    (parent.type = 'STATION') is_station
        FROM stop parent
        INNER JOIN ancestor child
    ON child.parent_stop_pk = parent.pk
    AND NOT child.is_station
)
SELECT stop_id, stop_pk, station_pk
  FROM ancestor
  WHERE parent_stop_pk IS NULL
  OR is_station
`

type MapStopIDAndPkToStationPkParams struct {
	SystemPk       int64
	FilterByStopPk bool
	StopPks        []int64
}

type MapStopIDAndPkToStationPkRow struct {
	StopID    string
	StopPk    int64
	StationPk int64
}

func (q *Queries) MapStopIDAndPkToStationPk(ctx context.Context, arg MapStopIDAndPkToStationPkParams) ([]MapStopIDAndPkToStationPkRow, error) {
	rows, err := q.db.Query(ctx, mapStopIDAndPkToStationPk, arg.SystemPk, arg.FilterByStopPk, arg.StopPks)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []MapStopIDAndPkToStationPkRow
	for rows.Next() {
		var i MapStopIDAndPkToStationPkRow
		if err := rows.Scan(&i.StopID, &i.StopPk, &i.StationPk); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const mapStopIDToPk = `-- name: MapStopIDToPk :many
SELECT pk, id from stop
WHERE
    system_pk = $1
    AND (
        NOT $2::bool
        OR id = ANY($3::text[])
    )
`

type MapStopIDToPkParams struct {
	SystemPk       int64
	FilterByStopID bool
	StopIds        []string
}

type MapStopIDToPkRow struct {
	Pk int64
	ID string
}

func (q *Queries) MapStopIDToPk(ctx context.Context, arg MapStopIDToPkParams) ([]MapStopIDToPkRow, error) {
	rows, err := q.db.Query(ctx, mapStopIDToPk, arg.SystemPk, arg.FilterByStopID, arg.StopIds)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []MapStopIDToPkRow
	for rows.Next() {
		var i MapStopIDToPkRow
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

const mapStopPkToChildPks = `-- name: MapStopPkToChildPks :many
SELECT parent_stop_pk parent_pk, pk child_pk
FROM stop
WHERE stop.parent_stop_pk = ANY($1::bigint[])
`

type MapStopPkToChildPksRow struct {
	ParentPk pgtype.Int8
	ChildPk  int64
}

func (q *Queries) MapStopPkToChildPks(ctx context.Context, stopPks []int64) ([]MapStopPkToChildPksRow, error) {
	rows, err := q.db.Query(ctx, mapStopPkToChildPks, stopPks)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []MapStopPkToChildPksRow
	for rows.Next() {
		var i MapStopPkToChildPksRow
		if err := rows.Scan(&i.ParentPk, &i.ChildPk); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const mapStopPkToDescendentPks = `-- name: MapStopPkToDescendentPks :many
WITH RECURSIVE descendent AS (
    SELECT
    stop.pk root_stop_pk,
    stop.pk descendent_stop_pk
    FROM stop
        WHERE stop.pk = ANY($1::bigint[])
    UNION
    SELECT
    descendent.root_stop_pk root_stop_pk,
    child.pk descendent_stop_pk
        FROM stop child
        INNER JOIN descendent
    ON child.parent_stop_pk = descendent.descendent_stop_pk
)
SELECT root_stop_pk, descendent_stop_pk FROM descendent
`

type MapStopPkToDescendentPksRow struct {
	RootStopPk       int64
	DescendentStopPk int64
}

func (q *Queries) MapStopPkToDescendentPks(ctx context.Context, stopPks []int64) ([]MapStopPkToDescendentPksRow, error) {
	rows, err := q.db.Query(ctx, mapStopPkToDescendentPks, stopPks)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []MapStopPkToDescendentPksRow
	for rows.Next() {
		var i MapStopPkToDescendentPksRow
		if err := rows.Scan(&i.RootStopPk, &i.DescendentStopPk); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateStop = `-- name: UpdateStop :exec
UPDATE stop SET
    feed_pk = $1,
    name = $2,
    location = $3::geography,
    url = $4,
    code = $5,
    description = $6,
    platform_code = $7,
    timezone = $8,
    type = $9,
    wheelchair_boarding = CASE WHEN $10::boolean THEN $11 ELSE wheelchair_boarding END,
    zone_id = $12,
    parent_stop_pk = NULL
WHERE
    pk = $13
`

type UpdateStopParams struct {
	FeedPk                   int64
	Name                     pgtype.Text
	Location                 types.Geography
	Url                      pgtype.Text
	Code                     pgtype.Text
	Description              pgtype.Text
	PlatformCode             pgtype.Text
	Timezone                 pgtype.Text
	Type                     string
	UpdateWheelchairBoarding bool
	WheelchairBoarding       pgtype.Bool
	ZoneID                   pgtype.Text
	Pk                       int64
}

func (q *Queries) UpdateStop(ctx context.Context, arg UpdateStopParams) error {
	_, err := q.db.Exec(ctx, updateStop,
		arg.FeedPk,
		arg.Name,
		arg.Location,
		arg.Url,
		arg.Code,
		arg.Description,
		arg.PlatformCode,
		arg.Timezone,
		arg.Type,
		arg.UpdateWheelchairBoarding,
		arg.WheelchairBoarding,
		arg.ZoneID,
		arg.Pk,
	)
	return err
}

const updateStop_Parent = `-- name: UpdateStop_Parent :exec
UPDATE stop SET
    parent_stop_pk = $1
WHERE
    pk = $2
`

type UpdateStop_ParentParams struct {
	ParentStopPk pgtype.Int8
	Pk           int64
}

func (q *Queries) UpdateStop_Parent(ctx context.Context, arg UpdateStop_ParentParams) error {
	_, err := q.db.Exec(ctx, updateStop_Parent, arg.ParentStopPk, arg.Pk)
	return err
}

const updateWheelchairBoardingForStop = `-- name: UpdateWheelchairBoardingForStop :exec
UPDATE stop
SET wheelchair_boarding = $1
WHERE pk = $2
`

type UpdateWheelchairBoardingForStopParams struct {
	WheelchairBoarding pgtype.Bool
	StopPk             int64
}

func (q *Queries) UpdateWheelchairBoardingForStop(ctx context.Context, arg UpdateWheelchairBoardingForStopParams) error {
	_, err := q.db.Exec(ctx, updateWheelchairBoardingForStop, arg.WheelchairBoarding, arg.StopPk)
	return err
}
