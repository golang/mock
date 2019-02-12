package captor

func AddIDs(dao Dao, ids []int) {
	dao.InsertIDs(ids)
}
