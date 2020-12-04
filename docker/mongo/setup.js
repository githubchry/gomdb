// 创建一个名为"chrydb"的DB
db = db.getSiblingDB('chrydb');

// 创建一个名为"chry"的用户，设置密码和权限
//db.createUser({user: "chry", pwd: "chry", roles: [{ role: "dbOwner", db: "chrydb"}]});
//如果是给程序做普通的增删改查操作, 权限不宜太高, 就读写即可
db.createUser({user: "chry", pwd: "chry", roles: ["readWrite"]});

// 在"chry"中创建一个名为"chry"的Collection 因为至少创建一个集合才能auth后执行show dbs看到 可省略
db.createCollection("chry");
