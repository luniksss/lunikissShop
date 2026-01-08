CREATE TABLE `order`
(
    `id`    INT     NOT NULL AUTO_INCREMENT,
    `user_id` INT   NOT NULL,
    `sales_outlet_id`   INT NOT NULL,
    `created_at` DATETIME DEFAULT NOW(),
    `status_name`   VARCHAR(255)    NOT NULL,
    PRIMARY KEY (`id`),
    FOREIGN KEY (`sales_outlet_id`) REFERENCES `sales_outlet` (`id`)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    INDEX `idx_sales_outlet_id` (`sales_outlet_id`)
)
    ENGINE = InnoDB
    CHARACTER SET = utf8mb4
    COLLATE utf8mb4_unicode_ci
;