package face

type Codec interface {
	Check(val interface{}) bool
	Marshal(val interface{}) ([]byte, error)
	Unmarshal(val interface{}, data []byte) error
}
