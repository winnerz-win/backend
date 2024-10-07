package mongo

/*
	https://m.blog.naver.com/ijoos/221312444591
*/

func (my *cCollection) Aggregate(pipeline interface{}) Iterator {
	cur, err := my.col.Aggregate(
		my.session.ctx(),
		pipeline,
	)
	if err != nil {
		//cc.RedItalic(err)
		return &cIterator{
			err: err,
		}
	}

	cursor := &cIterator{
		ctx:    my.session.ctx(),
		cursor: cur,
	}

	return cursor
}
