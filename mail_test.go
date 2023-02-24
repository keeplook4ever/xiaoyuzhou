package main

import (
	"testing"
	"xiaoyuzhou/pkg/mail"
)

func Test_SendMail(t *testing.T) {
	receiverList := []string{

		"test-u3cmgv0fm@srv1.mail-tester.com",
	}

	subject := "Test Golang"
	txt := "Welcome to mailgun and kiminouchuu.com"
	t.Log(receiverList)

	err := mail.SendMail(receiverList, subject, txt)
	if err != nil {
		t.Error(err)
	}
}

func Test_Mailgun(t *testing.T) {
	id, err := mail.SendSimpleMessage("mail.kiminouchuu.com", "106188044e56b94497503302610953ad-52d193a0-f576a233")
	t.Log(id)
	t.Log(err)
}
