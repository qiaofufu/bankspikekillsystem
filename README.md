# bankspikekillsystem
大二服创参赛作品(国赛三等奖)

一个银行理财产品秒杀系统
## 主要功能:
  1. 秒杀活动模块
  2. 理财产品管理模块
  3. 过滤规则配置模块
  4. 订单管理模块
  5. 用户管理模块
  6. 管理员权限分级模块

## 相关技术
  - 利用Radis缓存用户预订单信息, 预加载订单数量
  - 利用Mysql实现乐观锁
  - 采用Gin Web Framework
  - 采用Gorm ,对象关系映射
  - 采用JSON Web Token , 解决了认证相关问题
  - 采用RabbitMQ 消息队列，异步处理订单生成，有效进行消峰处理
  - 采用RESTFUL设计风格和开发方式
