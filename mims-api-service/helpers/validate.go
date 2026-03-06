package helpers

import (
	"errors"
	"fmt"
	"reflect"
)

func Msg(v interface{}, param string) string {
	st := reflect.ValueOf(v)
	// st.Elem()
	dataType := st.Kind().String()
	fmt.Println("st.dataType()", dataType)

	dataType2 := st.String()
	fmt.Println("st.String()", dataType2)

	// fmt.Println("ststststst", st, st.Kind(), reflect.String)
	// if string(st.Kind()) == "string" {
	// 	return errors.New("notZZ only validates strings")
	// }
	// if st.String() == "" {
	return "ddd"
	// }
	// return nil
}

func ValidateRequest(errMessageSlice []string) error {
	if len(errMessageSlice) == 0 {
		return nil
	}

	var errMessage string
	for i := 0; i < len(errMessageSlice); i++ {
		if len(errMessageSlice)-1 == i {
			errMessage += errMessageSlice[i]
		} else {
			errMessage += errMessageSlice[i] + ", "
		}
	}

	return errors.New(errMessage)
}
