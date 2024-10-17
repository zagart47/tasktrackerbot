create table if not exists public.tasks
(
    id            serial primary key,
    user_id       bigint    not null,
    text          text      not null,
    created_at    timestamp default now(),
    expiration    timestamp not null,
    duration      bigint    not null,
    reminder_sent boolean   default false,
    chat_id       bigint    not null,
    msg_id        int    not null
)
