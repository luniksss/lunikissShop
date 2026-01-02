CREATE TABLE `product`
(
    `id`         INT            NOT NULL AUTO_INCREMENT,
    `name`    VARCHAR(255)   NOT NULL,
    `description`    VARCHAR(400)   DEFAULT NULL,
    `price` INT NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY (`name`)
)
    ENGINE = InnoDB
    CHARACTER SET = utf8mb4
    COLLATE utf8mb4_unicode_ci
;