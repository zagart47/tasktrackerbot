create table if not exists public.tasks
(
    id            serial                  primary key,
    user_id       bigint                  not null,
    text          text                    not null,
    created_at    timestamp default now() not null,
    reminder_time timestamp with time zone,
    reminder_sent boolean
)
