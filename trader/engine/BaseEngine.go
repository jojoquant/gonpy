package engine


type BaseEngineInterface interface{
	close()
}

type BaseEngine struct{
	name string
}
