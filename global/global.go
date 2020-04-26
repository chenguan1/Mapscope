package global

func InitGlobal() error {
	if err := initConfig(); err != nil{
		return err
	}

	if err := initDb(); err != nil{
		return err
	}

	return nil
}

func Clear()  {
	db.Close()
}