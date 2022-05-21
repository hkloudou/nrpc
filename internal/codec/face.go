package codec

type Codec interface {
	Check(val interface{}) bool
	Marshal(val interface{}) ([]byte, error)
	Unmarshal(data []byte, val interface{}) error
}

var coders map[string]Codec

func init() {
	coders = make(map[string]Codec)
}

func RegisteCodec(key string, codec Codec) {
	coders[key] = codec
}

func GetCodec(val interface{}) Codec {
	for _, codec := range coders {
		if codec.Check(val) {
			return codec
		}
	}
	return defaultCodec
}
