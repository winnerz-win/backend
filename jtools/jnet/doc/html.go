package doc

import (
	"fmt"
	"strings"
)

//HTMLVIEW :
// func HTMLVIEW(w http.ResponseWriter) {
// 	w.Write(HTMLBytes(urlPath, docTItle, docVersion, doclist))
// }

const _style_color = `
	cc_blue {
		font-weight: 900;
		color: blue;
	}
	cc_red {
		font-weight: 900;
		color: red;
	}
	cc_purple{
		font-weight: 900;
		color: #4C0099;
	}	
	cc_green {
		font-weight: 900;
		color: #006600;
	}	
	cc_cyan {
		font-weight: 900;
		color: #006983;
	}
	cc_bold {
		font-weight: 900;
	}`

// HTMLBytes :
func HTMLBytes(urlPath, docTItle, docVersion string, doclist DocStringList) []byte {
	FRAME_WIDTH := "1600" //

	text := `
		<!DOCTYPE html>
		<meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
		
		<style>
			@media screen and (max-width: 700px) {
				body {
					width: 100%;
					margin: auto;
					text-align: center;
				}
			}
		
			body {
				width: 700px;
				margin: auto;
				text-align: center;
			}
			
			` + _style_color + `
			
		</style>
		
		<html>		
			<head>
				<script src="https://code.jquery.com/jquery-1.9.1.min.js"></script>
				<script type="text/javascript">
					$(document).ready(function() {
						console.log("start");
					});
				</script>
		
				<title>[Robin] {.PATH}</title>
			</head>
		
			<body style="width: ` + FRAME_WIDTH + `px; margin: auto; text-align: center;">
		
				<div style="width: 100%; background: linear-gradient( 180deg, #202757, #3C7885 ); position: relative;">
					<dummy id="top"></dummy>
					<!-- title -->
					<div style="width: 100%; text-align: center;">
						
						{.TITLE}
						{.INDEX}
						
						<br>
					</div>

					<!-- body -->
					<div style="display: inline-block; width: 98%; height: auto; margin-left: 1%; margin-right: 1%; margin-bottom: 1%; background: white;">
						<a href="#bottom">인덱스 보기</a>

						{.BODY}

						<br>	
						<br>

						<dummy id="bottom"></dummy>
						{.WARP}
						<a href="#top">맨 위로</a>
					</div>
				</div>

			</body>
		</html>`

	//////
	//////
	divEnd := "</div>"
	divStart := func(bColor, align, fSize, fWeight, fColor string) string {
		return `<div width="640; " ` + bColor + `style="vertical-align:middle; padding:40px 0px 0px 0px; 
		text-align: ` + align + `; 
		font-size: ` + fSize + `px; 
		font-weight: ` + fWeight + `; 
		color: ` + fColor + `; 
		letter-spacing: 0.5px;">`
	}
	_, _ = divStart, divEnd

	apiCount := fmt.Sprintf("&nbsp;[ API-Count : %v ]&nbsp;&nbsp;", doclist.Count())
	docVersion = apiCount + docVersion

	titleTag := ""
	titleTag += divStart("", "center", "40", "bold", "rgb(255, 255, 255)") + docTItle + divEnd
	titleTag += divStart("", "left", "25", "bold", "rgb(255, 255, 255)") + docVersion + divEnd

	bodyTag := ""
	for i, api := range doclist {
		order := fmt.Sprintf("warp_%v", i)
		bodyTag += `<dummy id="` + order + `"></dummy>`
		bodyTag += divStart("", "left", api.Size, api.Weight, api.Color) + api.Text + divEnd //#283A64
		//bodyTag += "―――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――― <br>"
		bodyTag += "―――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――― <br>"
		if i < len(doclist)-1 {
			bodyTag += `<a href="#bottom">인덱스 보기</a>`
		}
	} //for

	warpTag := ""
	for i, api := range doclist {
		order := fmt.Sprintf("warp_%v", i)
		warpTag += `<a href="#` + order + `">` + api.Href + `</a><br>`
	} //for
	//warpTag += "―――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――― <br>"
	warpTag += "―――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――― <br>"

	text = strings.ReplaceAll(text, "{.PATH}", urlPath)
	text = strings.ReplaceAll(text, "{.TITLE}", titleTag)
	text = strings.ReplaceAll(text, "{.BODY}", bodyTag)
	text = strings.ReplaceAll(text, "{.WARP}", warpTag)

	return []byte(text)
}
