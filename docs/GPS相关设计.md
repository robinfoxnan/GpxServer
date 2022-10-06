## 一、需求与功能分析



## 二、服务接口设计

接口汇总：

| 接口                | 功能               | 备注                                 |
| ------------------- | ------------------ | ------------------------------------ |
| /v1/regist          | 注册               |                                      |
| /v1/login           | 请求登录           | 使用密码或者验证码                   |
| /v1/loginCheckcode  | 请求验证码         | 这个功能需要后端开启短信服务         |
| /v1/addfriends      | 添加好友           |                                      |
| /v1/searchByPhone   | 搜索手机号         | 其他的搜索功能需要ES数据库支持       |
| /v1/listFriends     | 列出当前好友       |                                      |
| /v1/setFridendInfo  | 设置好友的相关权限 |                                      |
| /v1/blockFriend     | 屏蔽某人           |                                      |
| /v1/delFriend       | 删除好友关系       |                                      |
| /v1/gpx             | 单点上报           |                                      |
| /v1/gpxlist         | 批量上报           |                                      |
| /v1/track           | 查询一段时间轨迹   | 数据已经压缩                         |
| /v1/position        | 查询最后的活动点   | 这是最常用的一个功能                 |
| /v1/addtoGroup      | 创建群组或者添加人 |                                      |
| /v1/dismissGroup    | 解散               |                                      |
| /v1/SetGroupInfo    | 设置群组信息       |                                      |
| /v1/SetGroupMemInfo | 设置群组内个人信息 |                                      |
| /v1/listGroupMems   | 获取群组人员列表   | 获取此列表后，根据人员再检查动态位置 |
|                     |                    |                                      |



### 1.注册

功能定义：

注册新用户，在用户静态信息表中添加相关信息，成功后返回用户唯一ID；（userId为64位整数）

| 类型  | 字段       | 备注 |
| ----- | ---------- | ---- |
| 地址  | /v1/regist |      |
| 参数1 | phonenum   |      |
| 参数2 | pwd        |      |
| 参数3 | phonetype  |      |



| http code | body                            | 备注 |
| --------- | ------------------------------- | ---- |
| 200       | {"status:"ok", "userId":0001, } |      |
|           |                                 |      |
|           |                                 |      |



### 2.注销



### 3.登录



### 4 上报数据

功能：先查询当日已经上报数据的时间点，并上报剩余的点；服务端返回争取存储的最后时间点；

如果多日未上线，则上报7天内的记录数据；数据编码使用JSON编码

这里列举gpx格式中点的部分如下：

```
<trkpt lat="40.00533616" lon="116.25281886">
<ele>58.2550048828125</ele>
<time>2022-08-29T23:46:44.000Z</time>
<speed>0.0</speed>
<geoidheight>-8.0</geoidheight>
<src>gps</src>
<sat>3</sat>
<hdop>1.0</hdop>
<vdop>0.8</vdop>
<pdop>1.3</pdop>
</trkpt>
```

5 查询位置

功能：本人位置信息通过GPS获得，主要是查询群组中所有人位置，以及查询某个好友位置；

查询好友需要检查好友的设置，是否屏蔽（是否需要这个功能）









查询轨迹

返回的信息为了简化，不使用json，而使用固定的格式：

```
tm64|40.00533616|116.25281886|58.2550048828125|0.0
```







## 三、缓存设计

1 用户静态信息表

```sql

```

2. 用户动态信息表

```

  
```

3 好友关系表



4 群组表

定义群组的成员，以及每个人的权限；

备注：最后登录时间，位置信息如果通过关联效率低，这部分动态信息应该存在redis中；定期保存到数据库中；

```
CREATE TABLE `grouplist` (
  `groupID` varchar(32) NOT NULL,
  `userID` varchar(32) NOT NULL,
  `role`   int(11) DEFAULT NULL,
  `sequence` int(11) DEFAULT NULL,
  `lastUpdateTime` int(10) unsigned DEFAULT NULL,
  `creationTime` int(10) unsigned DEFAULT NULL,
  `lat`  double DEFAULT 0,
  `lon`  double DEFAULT 0,
  `ele`  double DEFAULT 0,
  `speed` double DEFAULT 0,
  PRIMARY KEY (`groupID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
```



