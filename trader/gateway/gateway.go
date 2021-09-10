package gateway

type BaseGateway interface{
	OnEvent()
	OnTick()
	OnTrade()
	OnOrder()
	OnPosition()
	OnAccount()
	OnQuote()
	OnLog()
	OnContract()
	WriteLog()
	Connect()
	Close()
	Subscribe()
	SendOrder()
	CancelOrder()
	
}