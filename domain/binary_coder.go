package domain

type BinaryCoder interface {
	Encode(m interface{}) ([]byte, error)
	Decode(str string) (interface{}, error)
}