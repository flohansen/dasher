-- name: GetAllFeatures :many
select * from features;

-- name: UpsertFeature :exec
insert into features (feature_id, description, enabled)
values (?, ?, ?)
on conflict (feature_id) do update set
    description = excluded.description,
    enabled     = excluded.enabled
returning *;

-- name: DeleteFeature :exec
delete from features where feature_id = ?;
