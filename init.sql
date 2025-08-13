create TABLE if not exists Post(
    id serial primary key not null,
    title varchar(50) not null,
    content text,
    imageURL varchar(255),
    userID int not null,
    created timestamp not null,
    expires timestamp not null
);


create TABLE if not EXISTS Comment(
    comment_id serial primary key not null,
    userID int not null,
    content text,
    imageURL varchar(255),
    created timestamp not null,
    postID int REFERENCES Post(id) ON DELETE CASCADE
);