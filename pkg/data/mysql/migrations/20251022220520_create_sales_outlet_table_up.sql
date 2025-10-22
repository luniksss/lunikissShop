CREATE TABLE sales_outlet (
    id VARCHAR(36) PRIMARY KEY,
    address STRING NOT NULL
);

CREATE TABLE product_stock (
     sales_outlet_id VARCHAR(36) NOT NULL,
     product_id VARCHAR(36) NOT NULL,
     size INT NULL DEFAULT NULL,
     amount INT NOT NULL DEFAULT 0,
     FOREIGN KEY (sales_outlet_id) REFERENCES sales_outlet(id) ON DELETE CASCADE,
     FOREIGN KEY (product_id) REFERENCES product(id) ON DELETE CASCADE,
     UNIQUE KEY unique_outlet_product (sales_outlet_id, product_id)
);

CREATE INDEX idx_product_stock_outlet ON product_stock(sales_outlet_id);
CREATE INDEX idx_product_stock_product ON product_stock(product_id);