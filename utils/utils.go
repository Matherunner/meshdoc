package utils

import "reflect"

func ReverseSlice(slice interface{}) {
	s := reflect.ValueOf(slice)
	swap := reflect.Swapper(slice)
	for i := 0; i < s.Len()/2; i++ {
		swap(i, s.Len()-1)
	}
}
