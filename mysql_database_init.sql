
CREATE TABLE db(
database_id  INT PRIMARY KEY AUTO_INCREMENT,dbdb
database_name VARCHAR(255) NOT NULL UNIQUE
);

CREATE TABLE tble(
	table_id INT PRIMARY KEY AUTO_INCREMENT,
    table_name VARCHAR(255),
    db_id INT ,
    FOREIGN KEY (db_id) REFERENCES db(database_id)
);

CREATE TABLE attribute(
	attribute_id INT PRIMARY KEY AUTO_INCREMENT,
    attribute_name VARCHAR(255) NOT NULL, 
    attribute_type VARCHAR(255) NOT NULL,
    table_id INT,
    FOREIGN KEY (table_id) REFERENCES tble(table_id)
);

ALTER TABLE tble 
ADD CONSTRAINT unique_table UNIQUE(db_id, table_name);

ALTER TABLE attribute 
ADD CONSTRAINT unique_attribute UNIQUE(attribute_name, table_id);

ALTER TABLE attribute 
ADD CONSTRAINT chk_type CHECK (attribute_type IN ('INT', 'CHAR', 'FLOAT', 'VARCHAR', 'BINARY', 'BOOL'));

ALTER TABLE attribute 
ADD CONSTRAINT chk_type CHECK (attribute_type IN ('INT', 'CHAR', 'FLOAT', 'VARCHAR', 'BOOL'));

CREATE TABLE constraints(
	constraint_id INT PRIMARY KEY AUTO_INCREMENT,
    constraint_name VARCHAR(255) NOT NULL,
    constraint_type VARCHAR(255) NOT NULL,
    constraint_detail VARCHAR(255),
    attribute_id INT,
    FOREIGN KEY (attribute_id) REFERENCES attribute(attribute_id),
    CHECK(constraint_type in ('NOT NULL', 'PRIMARY KEY', 'FOREIGN KEY', 'UNIQUE'))
);