package config

func Read(out interface{}, filename string) error{
	return NewIni().Read( out, filename)
}