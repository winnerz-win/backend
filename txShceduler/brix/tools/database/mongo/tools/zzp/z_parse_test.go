package zzp_test

import (
	"testing"
	"txscheduler/brix/tools/database/mongo/tools/cc"
	"txscheduler/brix/tools/database/mongo/tools/zzp"
)

func get_col_node() []interface{} {
	node, cmd, brunch := zzp.NodeMaker()
	list := []interface{}{}
	find := node(
		cmd(
			"find",
			zzp.TagSmall,
			zzp.ArgJson,
		),

		brunch(
			cmd(
				"skip",
				zzp.TagSmall,
				zzp.ArgText,
			),
			cmd(
				"sort",
				zzp.TagSmall,
				zzp.ArgTextComma,
			),
			cmd(
				"limit",
				zzp.TagSmall,
				zzp.ArgText,
			),
		),

		node(
			cmd(
				"count",
				zzp.TagSmall,
			),
		),
		node(
			cmd(
				"one",
				zzp.TagSmall,
			),
		),
		node(
			cmd(
				"all",
				zzp.TagSmall,
			),
		),
		node(
			cmd(
				"sum",
				zzp.TagSmall,
				zzp.ArgText,
			),
		),
	)
	update := node(
		cmd(
			"update",
			zzp.TagSmall,
			zzp.ArgJson,
			zzp.ArgJson,
		),
	)
	updateAll := node(
		cmd(
			"updateAll",
			zzp.TagSmall,
			zzp.ArgJson,
			zzp.ArgJson,
		),
	)
	upsert := node(
		cmd(
			"upsert",
			zzp.TagSmall,
			zzp.ArgJson,
			zzp.ArgJson,
		),
	)
	insert := node(
		cmd(
			"insert",
			zzp.TagSmall,
			zzp.ArgJson,
		),
	)
	remove := node(
		cmd(
			"remove",
			zzp.TagSmall,
			zzp.ArgJson,
		),
	)
	removeAll := node(
		cmd(
			"removeAll",
			zzp.TagSmall,
			zzp.ArgJson,
		),
	)
	aggregate := node(
		cmd(
			"aggregate",
			zzp.TagSmall,
			zzp.ArgJsonArray,
		),

		node(
			cmd(
				"one",
				zzp.TagSmall,
			),
		),
		node(
			cmd(
				"all",
				zzp.TagSmall,
			),
		),
	)
	dropcollection := node(
		cmd(
			"dropcollection",
			zzp.TagSmall,
		),
	)
	dropindexname := node(
		cmd(
			"dropindexname",
			zzp.TagSmall,
			zzp.ArgTextComma,
		),
	)
	ensureindex := node(
		cmd(
			"ensureindex",
			zzp.TagSmall,
			zzp.ArgTextComma,
		),
	)
	dropindexAll := node(
		cmd(
			"dropindexAll",
			zzp.TagSmall,
		),
	)
	indexes := node(
		cmd(
			"indexes",
			zzp.TagSmall,
		),
	)
	list = append(list,
		find,
		update,
		upsert,
		updateAll,
		insert,
		remove,
		removeAll,
		aggregate,
		dropcollection,
		dropindexname,
		ensureindex,
		dropindexAll,
		indexes,
	)
	return list
}

func TestPar(t *testing.T) {
	node, cmd, brunch := zzp.NodeMaker()
	_ = brunch

	db_node := node(
		node(cmd("collections", zzp.TagSmall)),
		node(cmd("show", zzp.TagSmall)),
		node(cmd("drop", zzp.TagSmall)),
		node(cmd("dropDatabase", zzp.TagSmall)),
		node(get_col_node()...),
	)
	// for _, node := range getColNodes() {
	// 	db_node.AddFork(node.(*zzp.Node))
	// }

	pack := db_node

	//dbg.Green(pack)

	queryText := `nft.nft_main.find( {"key":2} ).skip(2).sort(key,-a).limit(11).one()`
	queryText = `nft.nft_main.updateAll( {"key":2} , {"xxxx":3})`
	queryText = `nft.dropDatabase()`

	queryText = `db.col.aggregate( [{} , {}]   ).one()`
	queryText = `db.col.dropindexname( name_1 , key_2   )`
	queryText = `nft.nft_sell.indexes()`
	queryText = `db.col.upsert({"key":1} , {"key":1 ,"data":"sdfsdf  0000" } )`
	queryText = `nft.nft_sell.updateAll({"z_filter.file_type":"unknown"}, {"$set":{"z_filter.file_type" : "image"}})`
	// queryText = `db.col.ensureindex( name , -1 , true )`
	// queryText = `db.col.dropindexAll()`
	// queryText = `db.col.find(nil).sum(data.float)`

	r := pack.Parse(queryText) //

	cc.Red(r)
}
