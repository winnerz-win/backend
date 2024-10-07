package jlog

import (
	"os"
	"strings"
)

/* yaml 예제
log : {
	file : {
	  is_out_log : true,
	  log_level : "debug", 		# panic > error > warn > info > debug > trace
	  path : "/var/dexlog/",
	  name : "dexserver.log",	# api 서버 또는 core 서버의 경우 core_dexserver.log
	},

	console : {
	  is_out_log : true,
	  log_level : "debug",
	},

	field_base : {
	  server_name : "dexserver",
	  server_version : "0.0.0",
	},
}
*/

type ConfigLogYAML struct {
	File      ConfigFileLog    `yaml:"file"`       // 파일 로그 옵션
	Console   ConfigConsoleLog `yaml:"console"`    // 콘솔 로그 옵션
	FieldBase ConfigFieldBase  `yaml:"field_base"` // 고정된 기본 로그 필드
}
type ConfigFileLog struct {
	IsOutLog bool   `yaml:"is_out_log"` // 로그 사용 여부
	LogLevel string `yaml:"log_level"`  // 로그 레벨
	Path     string `yaml:"path"`       // 저장할 위치
	Name     string `yaml:"name"`       // 파일 이름
}
type ConfigConsoleLog struct {
	IsOutLog bool   `yaml:"is_out_log"` // 로그 사용 여부
	LogLevel string `yaml:"log_level"`  // 로그 레벨
}
type ConfigFieldBase struct {
	ServerName    string `yaml:"server_name"`    // 서버 이름
	ServerVersion string `yaml:"server_version"` // 버전
}

func (my ConfigFileLog) makeFilePath() string {
	file_path := ""
	if strings.HasSuffix(my.Path, "/") {
		file_path = my.Path
	} else {
		file_path = my.Path + "/"
	}
	os.MkdirAll(file_path, os.ModePerm)

	file_path += my.Name

	return file_path
}

// 테스트 등으로 사용 가능한 기본 컨피그
func DefaultConfigLogYAML(is_file_log bool, paths ...string) ConfigLogYAML {
	path := "/var/defaultConfigLog/"
	if len(paths) > 0 && paths[0] != "" {
		path = paths[0]
	}
	return ConfigLogYAML{
		File: ConfigFileLog{
			IsOutLog: is_file_log,
			LogLevel: "trace",
			Path:     path,
			Name:     "JTOOLS",
		},
		Console: ConfigConsoleLog{
			IsOutLog: true,
			LogLevel: "trace",
		},
		FieldBase: ConfigFieldBase{
			ServerName:    "JTOOLS",
			ServerVersion: "0.0.1",
		},
	}
}
