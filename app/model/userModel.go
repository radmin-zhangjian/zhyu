package model

type User struct {
	Id       int64  `gorm:"primaryKey;column:id;"`
	Username string `gorm:"column:user_name;type:varchar(30);default:(-)" `
	Password string `gorm:"column:password;type:varchar(100);default:(-)"`
	Phone    string `gorm:"column:phone;type:varchar(20);default:(-)"`
	Name     string `gorm:"column:name;type:varchar(30);default:(-)"`
	Age      uint8  `gorm:"column:age;type:tinyint(1);default:(0)"`
	Address  string `gorm:"column:address;type:varchar(100);default:(-)"`
	Photo    string `gorm:"column:photo;type:varchar(100);default:(-)"`
	Status   uint8  `gorm:"column:status;type:tinyint(1);default:(1)"`
	//Deleted    gorm.DeletedAt `gorm:"column:deleted;type:timestamp;default:(-)"`
	CreatedAt int64 `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt int64 `gorm:"column:updated_at;autoCreateTime"`
}

// TableName 自定义表名
func (*User) TableName() string {
	return "zhyu_user"
}
