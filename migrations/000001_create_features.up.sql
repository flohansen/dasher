begin;

create table if not exists features (
    feature_id text not null primary key,
    description text,
    enabled tinyint
);

commit;
