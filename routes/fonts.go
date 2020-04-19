package routes

import (
	"fmt"
	"github.com/kataras/iris/v12/context"
	"io/ioutil"
)

/*
https://docs.mapbox.com/api/maps/#fonts
Fonts API errors
Response body message	HTTP status code	Description
Invalid Range	400	Check the font glyph range. This error can also occur with empty fonts.
Maximum of 10 font faces permitted	400	Too many font faces in the query.
Not Authorized - No Token	401	No token was used in the query.
Not Authorized - Invalid Token	401	Check the access token you used in the query.
Not Found	404	Check the font name or names you used in the query.
*/

func FontGlypRange(ctx context.Context) {
	user := ctx.Params().Get("username")
	font := ctx.Params().Get("font")
	pbf := ctx.Params().Get("rangepbf") // 0-255.pbf

	str := fmt.Sprintf("%v,%v,%v",user,font,pbf)
	ctx.WriteString(str)
}

func FontList(ctx context.Context) {
	user := ctx.Params().Get("username")
	ftlist := make([]string,0)
	ftlist = append(ftlist,user+"'s font")
	ctx.JSON(ftlist)
}

func FontAdd(ctx context.Context) {
	user := ctx.Params().Get("username")

	body,err := ctx.GetBody()
	if err != nil{
		ctx.JSON(err)
		return
	}

	ctx.Application().Logger().Println(len(body))

	err = ioutil.WriteFile("./data/uploads/myfont.ttf",body,0777)
	if err != nil{
		ctx.JSON(err)
		return
	}

	ctx.WriteString(user)
}