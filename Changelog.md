## Change Log
### 2024-07-12
**feature:**
- 新增过期类型，可使用类型redis的ttl逻辑做过期
- 新增使用Value()方法时，判断key是否过期
- 新增带expire function的add key方法
- 新增对key update的计数

**optimization:**
- 优化过期key的清理逻辑，转为使用ticker清理
- go version 1.15 -> 1.22
