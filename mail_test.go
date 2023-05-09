package main

import (
	"fmt"
	"testing"
	"xiaoyuzhou/pkg/mail"
)

func Test_SendMail(t *testing.T) {
	receiverList := []string{

		"charles_ln@163.com",
	}

	subject := "TTTT"
	txt := fmt.Sprintf(`<!DOCTYPE html>\n<html>\n<head> <meta charset=\"utf-8\">\n</head>\n\n小さな宇宙、大きな感じ。今一番流行っている占いページ！\nここでは―この世界、宇宙、そして自分の運命を知る。今日の運勢から彼との恋愛まで。\n<br><br>\n\nこんにちは、小さな宇宙 です。\n<br><br>\n\nこの度はタロットアンサーをご利用いただきありがとうございます。\n以下、注文情報です。\n<br><br>\n<h3>▶ 送信情報</h3>\nタロットの回答提出日：%s\n注文番号：%s\n<br>\n\n<h3>▶ お支払いについて</h3>\n支払方法：%s\n支払いメールです：%s\n<br><br>\n\n<a href=\"%s\">小さな宇宙-タロット解答</a>のページでは、タロット解答の投稿状況を確認することができ、疑問点を把握することができます。\n<br>\n\n<h3>▶ 私たちについて・小さな宇宙</h3>\n2023 年人気絶頂のタロット占いサイト、今日の運勢を占う\n恋愛占いとタロットカード\n毎日一発の抽選で幸運を得ます\n占いや占いが好き。 あなたの人生のリーダーになりましょう！ 幸運を手に入れよう\nこれまで2万人以上の実績\n<br><br>\n\n今回は、願いを叶える！\n<br><br><br>\n\nCopyright ⓒ 2022 - 2023 All rights reserved. kiminouchuu.com\n\n</html>\n`, "ss", "ss", "ss", "ss", "ss")
	t.Log(receiverList)

	err := mail.SendMail(receiverList, subject, txt)
	if err != nil {
		t.Error(err)
	}
}
