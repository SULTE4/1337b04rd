create TABLE Form(
    id serial primary key not null,
    title varchar(50) not null,
    content text,
    image varchar(255),
    userID int not null,
    created timestamp not null,
    expires timestamp not null
);
