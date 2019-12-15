package server

func Run()  {
	r := setupRouter()
	r.Run(":8088")
}
