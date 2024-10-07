package mongo

/*
	https://www.mongodb.com/docs/v5.0/reference/operator/query/regex/?_ga=2.4136432.1374970675.1661476989-840893440.1612774695#std-label-regex-multiline-match

	: 검색글자로 시작하는 단어
	Bson{
		"name" : Bson{
			"$regex" : "^Coffee",
			"$options" : "m",
		},
	}

	: 검색글자가 중간에 포함하는
	Bson{
		"name" : Bson{
			"$regex" : ".Coffee.",
		},
	}

*/
