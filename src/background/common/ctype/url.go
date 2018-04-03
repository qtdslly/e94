package ctype

/*
var bindAddr string

func InitBindAddr(addr string) {
	bindAddr = addr
}

type Url string

func (v *Url) MarshalJSON() ([]byte, error) {
	return json.Marshal(bindAddr + string(v))
}

func (v *Url) UnmarshalJSON(data []byte) error {
	url := strings.TrimPrefix(string(data), bindAddr)
	if err := json.Unmarshal(url, v); err != nil {
		return err
	}
	return nil
}
*/
