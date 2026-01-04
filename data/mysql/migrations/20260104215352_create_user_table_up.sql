CREATE TABLE `user`
(
    `id`    INT     NOT NULL AUTO_INCREMENT,
    `name` VARCHAR(255)   NOT NULL,
    `surname`   VARCHAR(255) NOT NULL,
    `email` VARCHAR(255) NOT NULL,
    `password`   VARCHAR(255)    NOT NULL,
    `role` VARCHAR(255) NOT NULL,
    `phone` VARCHAR(255) DEFAULT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY (`email`)
)
    ENGINE = InnoDB
    CHARACTER SET = utf8mb4
    COLLATE utf8mb4_unicode_ci
;