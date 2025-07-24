package config

import (
	"fmt"
	"os"
	"reflect"
	"regexp"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	dir, _ := os.Getwd()
	fmt.Println("debug macth", dir)
	dir += "/"
	match, _ := regexp.Match(`/internal/domain/.+/(delivery)`, []byte(dir))
	if match {
		dir += "../../../../"
	}
	match, _ = regexp.Match(`/pkg/(mailer)/.+`, []byte(dir))
	if match {
		dir += "../../../"
	}
	match, _ = regexp.Match(`/cmd/api/middleware`, []byte(dir))
	if match {
		dir += "../../../"
	}
	match, _ = regexp.Match(`/internal/domain/.+/(repository)/.+`, []byte(dir))
	if match {
		dir += "../../../../../"
	}

	match, _ = regexp.Match(`/internal/domain/.+/(usecase)`, []byte(dir))
	if match {
		dir += "../../../../"
	}

	envFile := fmt.Sprintf("%v%v", dir, ".env")

	// fmt.Println("debug", envFile)
	godotenv.Load(envFile)
}

func MarshalEnv(d interface{}) {
	v := reflect.Indirect(reflect.ValueOf(d))
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		v.FieldByName(t.Field(i).Name).SetString(os.Getenv(t.Field(i).Tag.Get("env")))
	}
}
