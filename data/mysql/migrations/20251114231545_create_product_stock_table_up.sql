CREATE TABLE `product_stock`
(
    `sales_outlet_id`   INT     NOT NULL,
    `product_id`        INT     NOT NULL,
    `size`              INT     NOT NULL,
    `amount`            INT     NOT NULL,
    FOREIGN KEY (`product_id`) REFERENCES `product` (`id`)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    FOREIGN KEY (`sales_outlet_id`) REFERENCES `sales_outlet` (`id`)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    INDEX `idx_product_id` (`product_id`),
    INDEX `idx_sales_outlet_id` (`sales_outlet_id`)
)
    ENGINE = InnoDB
    CHARACTER SET = utf8mb4
    COLLATE utf8mb4_unicode_ci
;