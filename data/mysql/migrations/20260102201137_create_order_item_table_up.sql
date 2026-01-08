CREATE TABLE `order_item`
(
    `id`    INT     NOT NULL AUTO_INCREMENT,
    `order_id`  INT NOT NULL,
    `product_id`    INT NOT NULL,
    `amount`    INT NOT NULL,
    `price` INT NOT NULL,
    `size`  INT NOT NULL,
    PRIMARY KEY (`id`),
    FOREIGN KEY (`order_id`) REFERENCES `order` (`id`)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    INDEX `idx_order_id` (`order_id`),
    FOREIGN KEY (`product_id`) REFERENCES `product` (`id`)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    INDEX `idx_product_id` (`product_id`)
)
    ENGINE = InnoDB
    CHARACTER SET = utf8mb4
    COLLATE utf8mb4_unicode_ci
;