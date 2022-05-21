package face

type Codec interface {
	Check(val interface{}) bool
	Marshal(val interface{}) ([]byte, error)
	Unmarshal(data []byte, val interface{}) error
}
