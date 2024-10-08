package admin

import (
	"jtools/mms"
	"net/http"
	"txscheduler/brix/tools/database/mongo"
	"txscheduler/brix/tools/dbg"
	"txscheduler/brix/tools/jnet/chttp"
	"txscheduler/brix/tools/jtoken"
	"txscheduler/txm/ack"
	"txscheduler/txm/inf"
	"txscheduler/txm/model"
	"txscheduler/txm/pwd"
)

func init() {
	hLogin()
	hRefreshToken()
	hChangePWD()
	hAddSubAdmin()
	hRmvSubAdmin()
}

func adminDB(name, password string, callback func(db mongo.DATABASE, admin *model.Admin)) chttp.CError {
	ackError := chttp.ErrorNone
	model.DB(func(db mongo.DATABASE) {
		c := db.C(inf.COLAdmin)
		selector := mongo.Bson{"name": name}
		admin := model.Admin{}
		if c.Find(selector).One(&admin) != nil {
			ackError = ack.InvalidNick
			return
		}
		clientPWD := pwd.VerifyPWD(admin.Salt, password)
		if admin.Pwd != clientPWD {
			ackError = ack.InvalidPassword
			return
		}

		callback(db, &admin)
		ackError = chttp.ErrorNone
	})
	return ackError
}

func hLogin() {
	method := chttp.POST
	url := "/admin.login"

	type LoginData struct {
		Name string `json:"name"`
		Pwd  string `json:"pwd"`
	}
	type SignResult struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}

	Doc().Comment("[ 로그인 ] 로그인 요청").
		Method(method).URL(url).
		JParam(LoginData{},
			"name", "아이디",
			"pwd", "비번",
		).
		JResultOK(SignResult{},
			"access_token", "접속토큰",
			"refresh_token", "갱신토큰",
		).
		ResultERRR(ack.BadParam).
		ResultERRR(ack.InvalidNick).
		ResultERRR(ack.InvalidPassword).
		ResultERRR(ack.DBJob).
		Apply()

	handle.Append(
		method, url,
		func(w http.ResponseWriter, req *http.Request, ps chttp.Params) {
			cdata := &LoginData{}
			if chttp.BindingJSON(req, cdata) != nil {
				chttp.Fail(w, ack.BadParam)
				return
			}

			authError := adminDB(
				cdata.Name, cdata.Pwd,
				func(db mongo.DATABASE, admin *model.Admin) {
					createAt := mms.Now()
					dbpassword, salt := pwd.MakeDBPWD(cdata.Pwd)
					admin.Pwd = dbpassword
					admin.Salt = salt
					admin.Timestamp = createAt
					admin.YMD = createAt.YMD()

					accessToken := model.MakeTokenForAdmin(admin, jtoken.AccessToken)
					refreshToken := model.MakeTokenForAdmin(admin, jtoken.RefreshToken)
					admin.RefreshToken = refreshToken.ToString()

					admin.UpdateDB(db)

					chttp.OK(w, SignResult{
						AccessToken:  accessToken.ToString(),
						RefreshToken: admin.RefreshToken,
					})
				},
			)
			if authError != chttp.ErrorNone {
				chttp.Fail(w, authError)
			}
		},
	)
}

func hRefreshToken() {
	method := chttp.POST
	url := "/admin.refresh"
	type TokenData struct {
		Name         string `json:"name"`
		RefreshToken string `json:"refresh_token"`
	}
	type SignResult struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}

	Doc().Comment("[ 로그인 ] 토큰 갱신").
		Method(method).URL(url).
		JParam(TokenData{},
			"name", "아이디",
			"refresh_token", "갱신토큰",
		).
		JResultOK(SignResult{},
			"access_token", "접속토큰",
			"refresh_token", "갱신토큰",
		).
		ResultERRR(ack.BadParam).
		ResultERRR(ack.InvalidNick).
		ResultERRR(ack.InvalidToeken, "refresh-token 불일치").
		ResultERRR(ack.DBJob).
		Apply()

	handle.Append(
		method, url,
		func(w http.ResponseWriter, req *http.Request, ps chttp.Params) {
			cdata := &TokenData{}
			if chttp.BindingJSON(req, cdata) != nil {
				chttp.Fail(w, ack.BadParam)
				return
			}

			model.DB(func(db mongo.DATABASE) {
				c := db.C(inf.COLAdmin)

				selector := mongo.Bson{"name": cdata.Name}
				admin := model.Admin{}
				if c.Find(selector).One(&admin) != nil {
					chttp.Fail(w, ack.InvalidNick)
					return
				}
				if admin.RefreshToken != cdata.RefreshToken {
					chttp.Fail(w, ack.InvalidToeken)
					return
				}

				createAt := mms.Now()
				admin.Timestamp = createAt
				admin.YMD = createAt.YMD()

				accessToken := model.MakeTokenForAdmin(&admin, jtoken.AccessToken)
				refreshToken := model.MakeTokenForAdmin(&admin, jtoken.RefreshToken)
				admin.RefreshToken = refreshToken.ToString()

				upQuery := mongo.Bson{"$set": mongo.Bson{
					"timestamp":     admin.Timestamp,
					"refresh_token": admin.RefreshToken,
				}}
				if c.Update(selector, upQuery) != nil {
					chttp.Fail(w, ack.DBJob)
					return
				}

				chttp.OK(w, SignResult{
					AccessToken:  accessToken.ToString(),
					RefreshToken: admin.RefreshToken,
				})
			})

		},
	)
}

func hChangePWD() {
	type CDATA struct {
		Pwd    string `json:"pwd"`
		NewPwd string `json:"new_pwd"`
	}
	type SignResult struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}
	method := chttp.POST
	url := model.V2 + "/admin.change.pw"
	Doc().Comment("[ 비밀번호 변경 ] 관리자 비밀번호 변경").
		Method(method).URL(url).
		JParam(CDATA{},
			"pwd", "기존 비밀번호",
			"new_pwd", "새로운 비밀번호",
		).
		JResultOK(SignResult{},
			"access_token", "접속토큰",
			"refresh_token", "갱신토큰",
		).
		Apply()

	handle.Append(
		method, url,
		func(w http.ResponseWriter, req *http.Request, ps chttp.Params) {
			cdata := CDATA{}
			chttp.BindingJSON(req, &cdata)

			if cdata.NewPwd == "" {
				chttp.Fail(w, ack.BadParam)
				return
			}

			tuser := model.GetTokenAdmin(req)

			authError := adminDB(
				tuser.Name, cdata.Pwd,
				func(db mongo.DATABASE, admin *model.Admin) {
					createAt := mms.Now()
					dbpassword, salt := pwd.MakeDBPWD(cdata.NewPwd)
					admin.Pwd = dbpassword
					admin.Salt = salt
					admin.Timestamp = createAt
					admin.YMD = createAt.YMD()

					accessToken := model.MakeTokenForAdmin(admin, jtoken.AccessToken)
					refreshToken := model.MakeTokenForAdmin(admin, jtoken.RefreshToken)
					admin.RefreshToken = refreshToken.ToString()

					admin.UpdateDB(db)

					chttp.OK(w, SignResult{
						AccessToken:  accessToken.ToString(),
						RefreshToken: admin.RefreshToken,
					})
				},
			)
			if authError != chttp.ErrorNone {
				chttp.Fail(w, authError)
			}
		},
	)
}

func hAddSubAdmin() {
	type CDATA struct {
		Name string `json:"name"`
		Pwd  string `json:"pwd"`
	}

	method := chttp.POST
	url := model.V2 + "/admin.add"
	Doc().Comment("[ 관리자 추가 ] 보조 관리자 추가 (ROOT 권한)").
		Method(method).URL(url).
		JParam(CDATA{},
			"name", "신규 관리자 아이디(보조)",
			"pwd", "비번",
		).
		ResultERRR(ack.BadParam).
		ResultERRR(ack.ExistedName).
		ResultERRR(ack.InvalidRootAdmin).
		Apply()

	handle.Append(
		method, url,
		func(w http.ResponseWriter, req *http.Request, ps chttp.Params) {
			cdata := CDATA{}
			chttp.BindingJSON(req, &cdata)

			dbg.Trim(
				&cdata.Name,
				&cdata.Pwd,
			)
			if cdata.Name == "" || cdata.Pwd == "" {
				chttp.Fail(w, ack.BadParam)
				return
			}

			tuser := model.GetTokenAdmin(req)
			if tuser.IsRoot == false {
				chttp.Fail(w, ack.InvalidRootAdmin)
				return
			}
			if tuser.Name == cdata.Name {
				chttp.Fail(w, ack.ExistedName)
				return
			}

			model.DB(func(db mongo.DATABASE) {

				if cnt, _ := db.C(inf.COLAdmin).Find(mongo.Bson{"name": cdata.Name}).Count(); cnt > 0 {
					chttp.Fail(w, ack.ExistedName)
					return
				}

				admin := model.Admin{
					Name:   cdata.Name,
					IsRoot: false,
				}
				createAt := mms.Now()
				dbpassword, salt := pwd.MakeDBPWD(cdata.Pwd)
				admin.Pwd = dbpassword
				admin.Salt = salt
				admin.CreateAt = createAt
				admin.CreateYMD = createAt.YMD()
				admin.Timestamp = createAt
				admin.YMD = createAt.YMD()

				refreshToken := model.MakeTokenForAdmin(&admin, jtoken.RefreshToken)
				admin.RefreshToken = refreshToken.ToString()

				db.C(inf.COLAdmin).Insert(admin)

				chttp.OK(w, nil)
			})
		},
	)
}

func hRmvSubAdmin() {
	type CDATA struct {
		Name string `json:"name"`
	}

	method := chttp.POST
	url := model.V2 + "/admin.rmv"
	Doc().Comment("[ 관리자 삭제 ] 보조 관리자 삭제 (ROOT 권한)").
		Method(method).URL(url).
		JParam(CDATA{},
			"name", "보조 관리자 아이디",
		).
		ResultERRR(ack.BadParam).
		ResultERRR(ack.NotFoundName).
		ResultERRR(ack.InvalidRootAdmin).
		Apply()

	handle.Append(
		method, url,
		func(w http.ResponseWriter, req *http.Request, ps chttp.Params) {
			cdata := CDATA{}
			chttp.BindingJSON(req, &cdata)

			tuser := model.GetTokenAdmin(req)
			if tuser.IsRoot == false {
				chttp.Fail(w, ack.InvalidRootAdmin)
				return
			}
			if tuser.Name == cdata.Name {
				chttp.Fail(w, ack.BadParam)
				return
			}
			model.DB(func(db mongo.DATABASE) {

				admin := model.Admin{}
				if db.C(inf.COLAdmin).Find(mongo.Bson{"name": cdata.Name}).One(&admin) != nil {
					chttp.Fail(w, ack.NotFoundName)
					return
				}

				db.C(inf.COLAdmin).Remove(admin.Selector())
				chttp.OK(w, nil)
			})
		},
	)
}
