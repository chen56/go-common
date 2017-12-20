package iox

type NullWriter int

func (NullWriter) Write([]byte) (int, error) { return 0, nil }

var MyNullWriter *NullWriter = new(NullWriter)


