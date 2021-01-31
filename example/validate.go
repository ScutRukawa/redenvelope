package example

import (
	"fmt"
	"reflect"

	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"gopkg.in/go-playground/validator.v9"
	vtzh "gopkg.in/go-playground/validator.v9/translations/zh"
)

type User struct {
	FirstName string `validate:"required"`
	LastName  string `validate:"required"`
	Age       int64  `validate:"gte=0,lte=130"`
	Email     string `validate:"required,email"`
}

type MyType struct {
	kkk int
}

func (m *MyType) Error() string {
	return "testxxxxxxxxxxxxxx"
}

var interfaceEx error = &MyType{}
var kkkk = MyType{}

func test() {
	user := &User{
		FirstName: "firstName",
		LastName:  "lastName",
		Age:       136,
		Email:     "fl163.com",
	}
	validate := validator.New()
	//创建消息国际化通用翻译器
	cn := zh.New()
	uni := ut.New(cn, cn)
	translator, found := uni.GetTranslator("zh")
	if found {
		err := vtzh.RegisterDefaultTranslations(validate, translator)
		if err != nil {
			fmt.Println(err)
		}
	} else {
		fmt.Println("not found")
	}
	err := validate.Struct(user)
	if err != nil {
		//_, ok := err.(*MyType)
		//xxx, _ := interfaceEx.(*MyType) //当为指针接受者实现接口时，只有指向这个类型的指针才被认为实现了该接口。当为值接收者时，普通接口变量 以及指针类型的接口变量 都被认为实现了该接口
		//_, ok := err.(validator.InvalidValidationError)
		fmt.Println("reflect TypeOf err:", reflect.TypeOf(err))
		//fmt.Println("reflect TypeOf xxx:", reflect.TypeOf(xxx))
		// if !ok {
		// 	fmt.Println("err", err)
		// }
		errs, ok := err.(validator.ValidationErrors)
		if ok {
			for _, err := range errs {
				fmt.Println(err.Translate(translator))
			}
		}

	}
}
