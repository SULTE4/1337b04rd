CREATE TABLE IF NOT EXISTS Users (
    userID serial PRIMARY KEY NOT NULL,
    username varchar(150) NOT NULL,
    userURL text NOT NULL,
    userToken text NOT NULL
);

CREATE TABLE IF NOT EXISTS Post (
    id serial PRIMARY KEY NOT NULL,
    title varchar(50) NOT NULL,
    content text,
    imageURL varchar(255),
    created timestamp NOT NULL,
    expires timestamp NOT NULL,
    author varchar(150) not null,
    userID int REFERENCES Users(userID)
);

CREATE TABLE IF NOT EXISTS Comment (
    comment_id serial PRIMARY KEY NOT NULL,
    content text,
    imageURL varchar(255),
    created timestamp NOT NULL,
    userID int REFERENCES Users(userID) ON DELETE CASCADE,
    postID int REFERENCES Post(id) ON DELETE CASCADE
);