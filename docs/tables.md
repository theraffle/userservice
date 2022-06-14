# Tables

This guide contain infomation of tables used in `UserService`.

## Table List
```mysql
mysql> show tables;
+------------------+
| Tables_in_raffle |
+------------------+
| project          |
| user             |
| user_project     |
| user_wallet      |
+------------------+
```

### Table1. user
```mysql
mysql> show columns from user;
+----------+-------------+------+-----+---------+----------------+
| Field    | Type        | Null | Key | Default | Extra          |
+----------+-------------+------+-----+---------+----------------+
| id       | int         | NO   | PRI | NULL    | auto_increment |
| type     | varchar(10) | NO   |     | NULL    |                |
| telegram | varchar(50) | YES  |     |         |                |
| discord  | varchar(50) | YES  |     |         |                |
| twitter  | varchar(50) | YES  |     |         |                |
+----------+-------------+------+-----+---------+----------------+
```

### Table2. project
```mysql
mysql> show columns from project;
+-----------------+--------------+------+-----+---------+----------------+
| Field           | Type         | Null | Key | Default | Extra          |
+-----------------+--------------+------+-----+---------+----------------+
| id              | int          | NO   | PRI | NULL    | auto_increment |
| name            | varchar(50)  | NO   |     | NULL    |                |
| chain_id        | int          | NO   |     | NULL    |                |
| raffle_contract | varchar(100) | NO   |     | NULL    |                |
+-----------------+--------------+------+-----+---------+----------------+
```

### Table3. user_project
```mysql
mysql> show columns from user_project;
+------------+--------------+------+-----+---------+-------+
| Field      | Type         | Null | Key | Default | Extra |
+------------+--------------+------+-----+---------+-------+
| user_id    | int          | NO   |     | NULL    |       |
| project_id | int          | NO   | PRI | NULL    |       |
| chain_id   | int          | NO   | PRI | NULL    |       |
| address    | varchar(200) | NO   | PRI | NULL    |       |
+------------+--------------+------+-----+---------+-------+
```

### Table3. user_wallet
```mysql
mysql> show columns from user_wallet;
+----------+--------------+------+-----+---------+-------+
| Field    | Type         | Null | Key | Default | Extra |
+----------+--------------+------+-----+---------+-------+
| user_id  | int          | NO   |     | NULL    |       |
| chain_id | int          | NO   | PRI | NULL    |       |
| address  | varchar(200) | NO   | PRI | NULL    |       |
+----------+--------------+------+-----+---------+-------+
```