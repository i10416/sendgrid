package main

import (
	"context"
	"log"
	"os"

	"github.com/kenzo0107/sendgrid"
)

func main() {
	if err := handler(); err != nil {
		log.Fatal(err)
	}
}

func handler() error {
	apiKey := os.Getenv("SENDGRID_API_KEY")

	c := sendgrid.New(apiKey, sendgrid.OptionDebug(true))

	// メール作成
	mail := sendgrid.NewInputSendMail()

	// 送信者設定
	mail.SetFrom(sendgrid.NewEmail("from@example.com", "送信者名"))
	mail.SetSubject("テストメール")

	// 受信者設定
	p := sendgrid.NewPersonalization()
	p.AddTo(sendgrid.NewEmail("to@example.com", "受信者名"))
	p.AddCc(sendgrid.NewEmail("cc@example.com", "CCユーザー"))
	mail.AddPersonalization(p)

	// コンテンツ追加
	mail.AddContent(sendgrid.NewContent("text/plain", "プレーンテキストの内容です。"))
	mail.AddContent(sendgrid.NewContent("text/html", "<h1>HTMLの内容です</h1><p>SendGrid APIのテストメールです。</p>"))

	// カテゴリ追加
	mail.AddCategory("test")
	mail.AddCategory("api-example")

	// トラッキング設定
	mail.TrackingSettings = &sendgrid.TrackingSettings{
		ClickTracking: &sendgrid.ClickTrackingSetting{
			Enable:     sendgrid.Bool(true),
			EnableText: sendgrid.Bool(true),
		},
		OpenTracking: &sendgrid.OpenTrackingSetting{
			Enable: sendgrid.Bool(true),
		},
	}

	// サンドボックスモード有効化（テスト用）
	mail.MailSettings = &sendgrid.MailSettings{
		SandBoxMode: &sendgrid.Setting{Enable: sendgrid.Bool(true)},
	}

	// メール送信
	result, err := c.SendMail(context.TODO(), mail)
	if err != nil {
		return err
	}

	log.Printf("メール送信成功: MessageID=%s\n", result.MessageID)
	return nil
}
