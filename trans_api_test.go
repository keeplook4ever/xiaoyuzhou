package main
import (
	"context"
	"fmt"
	"testing"

	"cloud.google.com/go/translate"
	"golang.org/x/text/language"
)

func TranslateText(targetLanguage, text string) (string, error) {
	// text := "The Go Gopher is cute"
	ctx := context.Background()

	lang, err := language.Parse(targetLanguage)
	if err != nil {
		return "", fmt.Errorf("language.Parse: %v", err)
	}

	client, err := translate.NewClient(ctx)
	if err != nil {
		return "", err
	}
	defer client.Close()

	resp, err := client.Translate(ctx, []string{text}, lang, nil)
	if err != nil {
		return "", fmt.Errorf("Translate: %v", err)
	}
	if len(resp) == 0 {
		return "", fmt.Errorf("Translate returned empty response to text: %s", text)
	}
	return resp[0].Text, nil
}

func TestTrans(t *testing.T) {
	resp, err := TranslateText("ja", "你好，我是略略略")
	if err != nil {
		t.Error(err)
	}
	t.Log(resp)
}