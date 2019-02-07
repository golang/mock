package captor

import "fmt"

//go:generate mockgen -destination mock/mock_dao.go github.com/golang/mock/sample/captor Dao

type Dao interface {
	InsertIDs(ids []int)
}

type realDao struct{}

func (realDao) InsertIDs(ids []int) {
	fmt.Println(fmt.Sprintf("inserting ids %d", ids))
}

func NewDao() Dao {
	return &realDao{}
}
