create table if not exists users (
    id         serial primary key,
    name       varchar(255) not null,
    email      varchar(255) not null unique,
    age        int          not null default 0,
    phone      varchar(50)  not null default '',
    created_at timestamp    not null default now()
);

insert into users (name, email, age, phone) values ('John Doe', 'john@example.com', 30, '+77001112233');
