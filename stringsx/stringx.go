package stringsx

func CutRune(str string,minLength int)string{
	r:=[]rune(str)
	size:=len(r)
	if minLength>size{
		minLength=size
	}
	return string(r[0:minLength])
}

func GetOr(value string, default_ string) string {
	if value==""{
		return default_
	}
	return value
}
