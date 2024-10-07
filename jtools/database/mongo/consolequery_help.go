package mongo

//ConsoleQueryHelp2 :
func ConsoleQueryHelp(simple ...bool) string {
	v := `
{
	databasenames / show
	db.drop() / db.dropDatabase()
	db.collections() / db.show() 
	db.col.indexes()
	db.col.dropindexAll()
	db.col.ensureindex( key, 1 , false )
	db.col.dropindexname( name_1 )
	db.col.dropcollection()
	db.col.removeAll({query})
	db.col.remove({query})
	db.col.insert({bson})
	db.col.updateAll({bson} , {bson})
	db.col.update({bson} , {bson})
	db.col.upsert({bson} , {bson})
	db.col.aggregate([{bson}...]).one()
	db.col.aggregate([{bson}...]).all()
	db.col.find({bson}).limit(1).sort(key,-no).skip(20).count() / one() / all() / sum(field)
	db.col.findProjection( {bson} , {Dson} )
	db.col.FindOneAndUpdate({query} , {bson})  --> {before_document}
	db.col.FindOneAndDelete({query})           --> {before_document}
	db.col.FindOneAndReplace({query} , {bson}) --> {before_document}
}
`
	if len(simple) > 0 && simple[0] {
		return v
	}
	return "\033[1;33m" + //Yellow
		v +
		"\033[0m" //End
}
