create TABLE "Post"(
    id serial primary key not null,
    title varchar(50) not null,
    content text,
    imageURL varchar(255),
    userID int not null,
    created timestamp not null,
    expires timestamp not null
);
