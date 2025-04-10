package lang

// type StrJoiner struct {
// 	Elements  *arraylist.List
// 	CharCount int
// }

// func NewStrJoiner(size int) *StrJoiner {
// 	return &StrJoiner{
// 		Elements:  arraylist.New(),
// 		CharCount: 0,
// 	}
// }

// func (sj *StrJoiner) toString() string {
// 	newBytes := make([]byte, sj.CharCount)
// 	var writerIndex int = 0
// 	sj.Elements.Each(func(index int, value interface{}) {
// 		element := value.(string)
// 		elementBytes := []byte(element)
// 		copy(newBytes[writerIndex:], elementBytes)
// 		writerIndex += len(elementBytes)
// 	})
// 	return string(newBytes)
// }

// func (sj *StrJoiner) toStringWithBrackets() string {
// 	var size int = sj.CharCount + 2
// 	newBytes := make([]byte, size)
// 	newBytes[0] = '['
// 	newBytes[size-1] = ']'
// 	var writerIndex int = 1
// 	sj.Elements.Each(func(index int, value interface{}) {
// 		element := value.(string)
// 		elementBytes := []byte(element)
// 		copy(newBytes[writerIndex:], elementBytes)
// 		writerIndex += len(elementBytes)
// 	})
// 	return string(newBytes)
// }

// func (sj *StrJoiner) add(element string) *StrJoiner {
// 	sj.Elements.Add(element)
// 	sj.CharCount += len(element)
// 	return sj
// }
