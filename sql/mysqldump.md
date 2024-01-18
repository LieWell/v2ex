## 备份数据

mysqldump 无法直接忽略某个表的字段，因此还需要对导出的语句进行处理。

一般是处理 `id` 字段,可以全局替换为 `null`。

```
mysqldump -uroot -p v2ex member --skip-extended-insert > member.sql
```

| 参数名称                   | 参数解释                            | 备注                    |
|------------------------|---------------------------------|-----------------------|
| --skip-add-drop-table  | 不输出 drop table 语句               |                       |
| --no-create-db         | 不输出 create database 语句          |                       |
| --no-create-info       | 不输出 create table 语句             |                       |
| --databases            | 指定数据库名称                         |                       |
| --tables               | 指定表名称                           | 此参数讲覆盖 --databases 配置 |
| --skip-extended-insert | 不实用多 values 的 insert 语法         | 就是每条数据对应一个 insert 语句  |
| --skip-comments        | 不导出注释                           |                       |
| --insert-ignore        | 使用 insert ignore 语句代替 insert 语句 |                       |

## 导入数据
```shell
mysql -uroot -p -D v2ex < ./member.sql
```