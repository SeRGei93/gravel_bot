package table

type File struct {
	ID       string `gorm:"primaryKey;autoIncrement:false"`
	Type     string
	EntityId uint
}
