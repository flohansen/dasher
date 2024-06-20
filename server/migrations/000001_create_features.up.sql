create table if not exists features (
    feature_id text not null primary key,
    description text not null,
    enabled boolean not null
);
