begin;

create table if not exists projects
(
    id                  integer primary key,
    name                varchar(50) not null,
    name_with_namespace varchar(100) not null,
    created_at          timestamp    not null,
    last_activity_at    timestamp    not null
);

create table if not exists events
(
    project_id      integer     not null references projects (id),
    action_name     varchar(20) not null,
    target_id       integer,
    target_type     varchar(20),
    author_id       integer     not null,
    author_username varchar(30) not null,
    created_at      timestamp   not null,
    primary key (created_at, author_id, action_name)
);

commit;
