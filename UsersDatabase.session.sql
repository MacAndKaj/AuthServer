CREATE TABLE IF NOT EXISTS users ( 
    id              INT AUTO_INCREMENT NOT NULL, 
    nickname        VARCHAR(64) NOT NULL, 
    first_name      VARCHAR(128) NOT NULL, 
    email           VARCHAR(256) NOT NULL, 
    creation_date   TIMESTAMP, 
    password        VARCHAR(256) NOT NULL, 
    PRIMARY KEY (`id`) 
);

DROP TABLE tokens;

CREATE TABLE tokens (
    user_id         INT NOT NULL,
    hash            VARCHAR(256),
    expiration_date TIMESTAMP,
    permissions     VARCHAR(2),
    PRIMARY KEY(`user_id`)
);
