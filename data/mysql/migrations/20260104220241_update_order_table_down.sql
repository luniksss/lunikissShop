ALTER TABLE `order`
    DROP FOREIGN KEY `fk_order_user`;

ALTER TABLE `order`
    DROP INDEX `idx_user_id`;