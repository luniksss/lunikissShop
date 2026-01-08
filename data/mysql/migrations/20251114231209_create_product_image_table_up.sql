CREATE TABLE `product_image`
(
    `id`         INT            NOT NULL AUTO_INCREMENT,
    `product_id`    INT   NOT NULL,
    `image_path`    VARCHAR(400)   NOT NULL,
    PRIMARY KEY (`id`),
    FOREIGN KEY (`product_id`) REFERENCES `product` (`id`)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    INDEX `idx_product_id` (`product_id`)
)
    ENGINE = InnoDB
    CHARACTER SET = utf8mb4
    COLLATE utf8mb4_unicode_ci
;