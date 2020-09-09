# gomdb
```
sudo apt install mongodb
mongo -version
service mongodb start

进入mongo shell
mongo

查看dbs列表
> show dbs
admin   0.000GB
config  0.000GB
local   0.000GB
test    0.000GB

切换到其中一个db
> use test
switched to db test

创建collection(集合)
> db.createCollection("trainer")
{ "ok" : 1 }
> db.createCollection("student")
{ "ok" : 1 }

显示所有collection
>  show collections
student
trainer

删除指定数据集
> db.trainer.drop()
true

插入一条文档
> db.student.insertOne({name:"小王子",age:18});
{
        "acknowledged" : true,
        "insertedId" : ObjectId("5f5848ee38ad7149539434e9")
}

插入多条文档
> db.student.insertMany([
... {name:"张三",age:20},
... {name:"李四",age:25}
... ]);
{
	"acknowledged" : true,
	"insertedIds" : [
		ObjectId("5f58491d38ad7149539434ea"),
		ObjectId("5f58491d38ad7149539434eb")
	]
}

查询所有文档：
> db.student.find()
{ "_id" : ObjectId("5f5848ee38ad7149539434e9"), "name" : "小王子", "age" : 18 }
{ "_id" : ObjectId("5f58491d38ad7149539434ea"), "name" : "张三", "age" : 20 }
{ "_id" : ObjectId("5f58491d38ad7149539434eb"), "name" : "李四", "age" : 25 }

查询age>20岁的文档：
> db.student.find(
... {age:{$gt:20}}
... )
{ "_id" : ObjectId("5f58491d38ad7149539434eb"), "name" : "李四", "age" : 25 }

更新文档：
> db.student.update({name:"小王子"},{name:"老王子",age:98})
WriteResult({ "nMatched" : 1, "nUpserted" : 0, "nModified" : 1 })

删除文档：
db.student.deleteOne({name:"李四"});
```