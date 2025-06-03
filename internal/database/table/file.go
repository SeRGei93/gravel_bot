package table

type File struct {
	Id       string `gorm:"primaryKey;autoIncrement:false"`
	EntityId uint
	Type     string
}
