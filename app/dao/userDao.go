package dao

import (
	"gorm.io/gorm"
	"log"
	"zhyu/app/model"
	"zhyu/utils"
)

// 文档：https://gorm.io/zh_CN/docs/query.html

// UserCreate Create无论什么情况都执行插入
func UserCreate(user *model.User) {
	db := utils.GetDB()
	if result := db.Model(&user).Create(&user); result.Error != nil {
		log.Printf("userCreate: %v", result.Error)
	}
	return
}

// UserSave Save需要插入的数据存在则不进行插入
func UserSave(user *model.User) {
	db := utils.GetDB()
	if result := db.Model(&user).Create(&user); result.Error != nil {
		log.Printf("userCreate: %v", result.Error)
	}
	return
}

// UserCreateBatch 批量插入
func UserCreateBatch(user *[]model.User) {
	db := utils.GetDB()
	db.Model(&user).Create(&user)
	return
}

// UserGetOneNamePass 获取部分参数，例如只获取名字和密码
// SELECT username,password FROM users WHERE xxx = xxx and xx = xx;
func UserGetOneNamePass(username string, password string) (user *model.User, rowsAffected bool) {
	db := utils.GetDB()
	result := db.Model(&user).Select("id", "user_name", "phone", "name", "age", "address").
		Where("user_name = ? and password = ?", username, password).Limit(1).Find(&user)
	rowsAffected = false
	if rows := result.RowsAffected; rows > 0 {
		rowsAffected = true
	}
	return
}

// UserGetFirst 获取第一个，默认查询第一个
// SELECT * FROM users ORDER BY id LIMIT 1;
func UserGetFirst() (user *model.User) {
	db := utils.GetDB()
	db.Model(&user).First(&user)
	return
}

// UserGetLast 获取最后一个
// SELECT * FROM users ORDER BY id DESC LIMIT 1;
func UserGetLast() (user *model.User) {
	db := utils.GetDB()
	db.Model(&user).Last(&user)
	return
}

// UserGetById 通过主键获取
// SELECT * FROM users WHERE id = 1;
func UserGetById(id int64) (user *model.User, rowsAffected bool) {
	db := utils.GetDB()
	result := db.Model(&user).Find(&user, id)
	//db.Model(&user).Where("id = ?",id).Find(&user)
	rowsAffected = false
	if rows := result.RowsAffected; rows > 0 {
		rowsAffected = true
	}
	return
}

// UserGetByIds 通过主键批量查询
// SELECT * FROM users WHERE id IN (1,2,3);
func UserGetByIds(ids []int64) (user *[]model.User) {
	db := utils.GetDB()
	db.Model(&user).Find(&user, ids)
	//db.Model(&user).Where("id in ?",ids).Find(&user)
	return
}

// UserGetSomeParam 获取部分参数，例如只获取名字和密码
// SELECT username,password FROM users WHERE id = 1;
func UserGetSomeParam(id int64) (user *model.User) {
	db := utils.GetDB()
	db.Model(&user).Select("username", "password").Find(&user, id)
	return
}

// UserGetPage 分页查询，可以使用Limit & Offset进行分页查询
// SELECT * FROM users OFFSET 5 LIMIT 10;
func UserGetPage(limit int, offset int) (user *[]model.User) {
	db := utils.GetDB()
	db.Model(&user).Limit(limit).Offset(offset).Order("id desc").Find(&user)
	return
}

// UserGetByOrder order 排序
// SELECT * FROM users ORDER BY id desc, username;
func UserGetByOrder() (user *[]model.User) {
	db := utils.GetDB()
	db.Model(&user).Order("id desc, username").Find(&user)
	//db.Model(&user).Order("id desc").Order("username").Find(&user)
	return
}

// UserUpdateUsername 更新单个字段
// UPDATE users SET username = "lomtom" where id = 1
func UserUpdateUsername(id int64, name string) {
	db := utils.GetDB()
	db.Model(&model.User{}).Where("id = ?", id).Update("name", name)
	return
}

// UserUpdateByUser 全量/多列更新（根据结构体）
// UPDATE `user` SET `id`=14,`user_name`='lomtom',`password`='123456',`create_time`='2021-09-26 14:22:21.271',`update_time`='2021-09-26 14:22:21.271' WHERE id = 14 AND `user`.`deleted` IS NULL
func UserUpdateByUser(user *model.User) {
	db := utils.GetDB()
	db.Model(&model.User{}).Where("id = ?", user.Id).Updates(&user)
	return
}

// UserDeleteByUser 简单删除（根据user里的id进行删除）
// 说明： 结构体未加gorm.DeletedAt标记的字段，直接删除，加了将更新deleted字段，即实现软删除
// DELETE from users where id = 28;
// UPDATE `user` SET `deleted`='2021-09-26 14:25:33.368' WHERE `user`.`id` = 28 AND `user`.`deleted` IS NULL
func UserDeleteByUser(user *model.User) {
	db := utils.GetDB()
	db.Model(&model.User{}).Delete(&user)
	return
}

// UserDeleteById 根据id进行删除
// UPDATE `user` SET `deleted`='2021-09-26 14:29:55.15' WHERE `user`.`id` = 28 AND `user`.`deleted` IS NULL
func UserDeleteById(id int64) {
	db := utils.GetDB()
	db.Model(&model.User{}).Delete(&model.User{}, id)
	return
}

// Transaction 匿名事务
// 可使用db.Transaction匿名方法来表明多个操作在一个事务里面，返回err将回滚，返回nil将提交事务
func Transaction() error {
	db := utils.GetDB()
	err := db.Transaction(func(tx *gorm.DB) error {
		// 在事务中执行一些 db 操作（从这里开始，您应该使用 'tx' 而不是 'db'）
		if err := tx.Create(&model.User{Username: "lomtom"}).Error; err != nil {
			// 返回任何错误都会回滚事务
			return err
		}
		if err := tx.Delete(&model.User{}, 28).Error; err != nil {
			return err
		}
		// 返回 nil 提交事务
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

// Transaction1 手动事务
func Transaction1() error {
	db := utils.GetDB()
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	// 在事务中执行一些 db 操作（从这里开始，您应该使用 'tx' 而不是 'db'）
	if err := tx.Create(&model.User{Username: "lomtom"}).Error; err != nil {
		// 回滚事务
		tx.Rollback()
		return err
	}
	if err := tx.Delete(&model.User{}, 28).Error; err != nil {
		tx.Rollback()
		return err
	}
	// 提交事务
	return tx.Commit().Error
}

// QueryRawDao 原生查询
func QueryRawDao(sql string, values ...interface{}) (user *[]model.User) {
	db := utils.GetDB()
	db.Raw(sql, values...).Scan(&user)
	return
}

// QueryExecDao 一般用于更新不返回数据
// UPDATE users SET money = ? WHERE name = ?
func QueryExecDao(sql string, values ...interface{}) (user *[]model.User) {
	db := utils.GetDB()
	db.Exec(sql, values...)

	// Exec with SQL Expression
	//db.Exec("UPDATE users SET money = ? WHERE name = ?", gorm.Expr("money * ? + ?", 10000, 1), "jinzhu")

	return
}

// UserGetList 分页查询，可以使用Limit & Offset进行分页查询
// SELECT * FROM users OFFSET 5 LIMIT 10;
func UserGetList(limit int, offset int, where string, args ...any) (user *[]model.User, total int64) {
	db := utils.GetDB()
	db.Model(&user).
		Select("id", "user_name", "phone", "name", "age", "address", "photo", "status").
		Where(where, args...).
		Limit(limit).
		Offset(offset).
		Order("id desc").
		Find(&user)

	db.Model(&user).
		Select("id", "user_name", "phone", "name", "age", "address", "photo", "status").
		Where(where, args...).
		Order("id desc").
		Count(&total)
	return
}
