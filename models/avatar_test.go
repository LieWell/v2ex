package models

import (
	"fmt"
	"testing"
)

func TestSaveAvatar(t *testing.T) {
	prepareMysql()
	avatar := &Avatar{
		Name:   "test",
		Avatar: "test1",
	}

	id, err := SaveAvatar(avatar)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Printf("insert avatar[%v] success!\n", id)

	allAvatar, err := FindAllAvatar()
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Printf("query all avatar success!\n")
	for _, a := range allAvatar {
		fmt.Printf("\t%+v\n", a)
	}

	err = DeleteAvatar(id)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Printf("delete avatar[%v] success!\n", id)
}
