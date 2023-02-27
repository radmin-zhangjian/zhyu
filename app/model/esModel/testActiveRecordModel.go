package esModel

type TestActiveRecordModel struct {
	Controller string `json:"controller"`
	Action     string `json:"action"`
	Token      string `json:"token"`
	ShopCode   string `json:"shopCode"`
	Buid       string `json:"buid"`
	Cuid       string `json:"cuid"`
	Device     string `json:"device"`
	Version    string `json:"version"`
	Day        string `json:"day"`
	Type       string `json:"type"`
	Ip         string `json:"ip"`
	Route      string `json:"route"`
	PostData   string `json:"postData"`
	//UserInfo   map[string]string `json:"userInfo"`
	UserInfo  UserInfo `json:"userInfo"`
	Datetime  string   `json:"datetime"`
	Timestamp string   `json:"timestamp"`
}

type UserInfo struct {
	Id       string `json:"id"`
	UserName string `json:"username"`
	Phone    string `json:"phone"`
	Address  string `json:"address"`
}
