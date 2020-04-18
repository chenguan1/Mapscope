package routes

import (
	"fmt"
	"github.com/kataras/iris/v12/context"
)

/*
https://docs.mapbox.com/api/maps/#datasets
Datasets API errors
Unauthorized	401	Check the access token you used in the query.
No such user	404	Check the username you used in the query.
Not found	404	The access token used in the query needs the datasets:list scope.
Invalid start key	422	Check the start key used in the query.
No dataset	422	Check the dataset ID used in the query (or, if you are retrieving a feature, check the feature ID).
*/

func FontGlypRange(ctx context.Context) {
	user := ctx.Params().Get("username")
	font := ctx.Params().Get("font")
	start := ctx.Params().Get("start")
	end := ctx.Params().Get("end")

	str := fmt.Sprintf("%v,%v,%v,%v",user,font,start,end)
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
	ctx.WriteString(user)
}