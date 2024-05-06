-- 念のため、使用したSQLをメモしたファイルを作成しておきました

CREATE TABLE `products` (
    `id` INT AUTO_INCREMENT PRIMARY KEY,
    `name` VARCHAR(255) NOT NULL,
    `brand_id` INT NOT NULL,
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    FOREIGN KEY (`brand_id`) REFERENCES `brands`(`id`) ON DELETE RESTRICT ON UPDATE CASCADE
);

CREATE TABLE `images` (
    `id` INT AUTO_INCREMENT PRIMARY KEY ,
    `product_id` INT NOT NULL,
    `path` VARCHAR(255) NOT NULL,
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    FOREIGN KEY (`product_id`) REFERENCES products(`id`) ON DELETE CASCADE
);

CREATE TABLE `brands` (
    `id` INT AUTO_INCREMENT PRIMARY KEY,
    `name` VARCHAR(255) NOT NULL UNIQUE,
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    INDEX (`name`)
);


-- ブランドのデータを追加
INSERT INTO `brands` (`name`) VALUES ('KENZO');

-- 商品のデータを追加
INSERT INTO `products` (`name`, `brand_id`) VALUES ("KENZO 'TIGER CREST' POLO SHIRT", 1);

-- 画像のデータを追加
INSERT INTO `images` (`product_id`, `path`) VALUES (1, "https://stok.store/cdn/shop/files/20220304040105603_E52---kenzo---FA65PO0014PU01B_4_M1.jpg");

-- 全件取得
SELECT
    p.id AS product_id,
    p.name AS product_name,
    b.name AS brand_name,
    i.path AS image_path,
    p.created_at AS product_created_at,
    p.updated_at AS product_updated_at
FROM
    products p
JOIN
    brands b ON p.brand_id = b.id
LEFT JOIN
    images i ON p.id = i.product_id;