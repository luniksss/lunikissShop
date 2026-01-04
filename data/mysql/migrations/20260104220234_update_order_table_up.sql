ALTER TABLE `order`
    ADD INDEX `idx_user_id` (`user_id`);

ALTER TABLE `order`
    ADD CONSTRAINT `fk_order_user`
        FOREIGN KEY (`user_id`)
            REFERENCES `user` (`id`)
            ON DELETE CASCADE
            ON UPDATE CASCADE;