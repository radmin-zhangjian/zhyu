# zhyu
****
说明: gin的中间件堪称它的精髓  
1.根据gin框架搭建的脚手架 自定义Context头 定义开发结构  
2.实现了动态路由 (api开头) 支持版本  
3.静态路由也可以匹配动态方法 支持版本  
4.添加了一些基础中间件 (logger & JWT & contextKeys & limit & cors等)  
5.setting 配置信息设置  
6.开发结构分层 (api控制器层、service服务层、dao数据获取逻辑层、model层)  
7.jwt-go 实例 & 简单用户验证  
8.实现了自定义异步日志 & gorm自定义日志 & TraceId & 链路日志
9.自定义框架 劫持ServeHTTP方法 实现路由和中间件功能 (zhyu目录) 仅供学习   

### 这里我设置成自动模式
go env -w GO111MODULE=auto

### 设置代理
go env -w GOPROXY=https://goproxy.io,direct

### 安装gin 及 用到的组件
go get -u github.com/gin-gonic/gin  
go get -u github.com/petermattis/goid  
go get -u github.com/didip/tollbooth  
go get -u gopkg.in/yaml.v2  
go get -u github.com/sony/sonyflake  
go get -u gorm.io/gorm  
go get -u gorm.io/driver/mysql  
go get -u gorm.io/gorm/logger  
go get -u gorm.io/gorm/schema  
go get -u github.com/dgrijalva/jwt-go  
go get -u github.com/go-redis/redis/v8  
#### 安装验证器
go get -u github.com/go-playground/locales/zh  
go get -u github.com/go-playground/universal-translator  
go get -u github.com/go-playground/validator/v10  
go get -u github.com/go-playground/validator/v10/translations/zh  

#### mysql user 数据结构
CREATE TABLE `zhyu_user` (  
`id` mediumint(6) unsigned NOT NULL AUTO_INCREMENT,  
`username` varchar(20) DEFAULT NULL,  
`password` varchar(32) DEFAULT NULL,  
`roleid` smallint(5) DEFAULT '0',  
`encrypt` varchar(6) DEFAULT NULL,  
`lastloginip` varchar(15) DEFAULT NULL,  
`lastlogintime` int(10) unsigned DEFAULT '0',  
`email` varchar(40) DEFAULT NULL,  
`realname` varchar(50) NOT NULL DEFAULT '',  
PRIMARY KEY (`userid`),  
KEY `username` (`username`) USING BTREE  
) ENGINE=MyISAM AUTO_INCREMENT=1 DEFAULT CHARSET=gbk  

### 测试用例 动态路由
http://localhost:9090/api/v1/test  
http://localhost:9090/api/v2/test  

### 测试用例 静态路由
http://localhost:9090/v1/test  
http://localhost:9090/v2/test  
