// main.go
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

// AnythingLLM API yapılandırması
type AnythingLLMConfig struct {
	BaseURL       string
	APIKey        string
	WorkspaceSlug string
}

// Chat isteği için yapı
type ChatRequest struct {
	Message string `json:"message"`
	Mode    string `json:"mode"`
}

// Chat yanıtı için yapı
type ChatResponse struct {
	TextResponse string        `json:"textResponse"`
	Type         string        `json:"type"`
	Sources      []interface{} `json:"sources"`
	Close        bool          `json:"close"`
	Error        string        `json:"error"`
}

// AnythingLLM istemcisi
type AnythingLLMClient struct {
	config AnythingLLMConfig
	client *http.Client
}

func NewAnythingLLMClient(baseURL, apiKey, workspaceSlug string) *AnythingLLMClient {
	return &AnythingLLMClient{
		config: AnythingLLMConfig{
			BaseURL:       strings.TrimSuffix(baseURL, "/"),
			APIKey:        apiKey,
			WorkspaceSlug: workspaceSlug,
		},
		client: &http.Client{},
	}
}

// Workspace'e mesaj gönder
func (a *AnythingLLMClient) SendMessage(message string) (string, error) {
	url := fmt.Sprintf("%s/api/v1/workspace/%s/chat", a.config.BaseURL, a.config.WorkspaceSlug)

	chatReq := ChatRequest{
		Message: message,
		Mode:    "chat",
	}

	jsonData, err := json.Marshal(chatReq)
	if err != nil {
		return "", fmt.Errorf("JSON oluşturma hatası: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("istek oluşturma hatası: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", a.config.APIKey))

	log.Printf("AnythingLLM'e istek gönderiliyor: %s", url)

	resp, err := a.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("API isteği hatası: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("yanıt okuma hatası: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API hatası (Status: %d): %s", resp.StatusCode, string(body))
	}

	var chatResp ChatResponse
	if err := json.Unmarshal(body, &chatResp); err != nil {
		return "", fmt.Errorf("JSON parse hatası: %v", err)
	}

	if chatResp.Error != "" {
		return "", fmt.Errorf("chat hatası: %s", chatResp.Error)
	}

	return chatResp.TextResponse, nil
}

// Workspace bağlantısını test et
func (a *AnythingLLMClient) TestConnection() error {
	url := fmt.Sprintf("%s/api/v1/workspace/%s", a.config.BaseURL, a.config.WorkspaceSlug)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", a.config.APIKey))

	resp, err := a.client.Do(req)
	if err != nil {
		return fmt.Errorf("bağlantı hatası: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("workspace bulunamadı (Status: %d): %s", resp.StatusCode, string(body))
	}

	return nil
}

func main() {
	// .env dosyasını yükle
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  .env dosyası bulunamadı, ortam değişkenleri kullanılacak")
	}

	// Yapılandırmayı al
	telegramToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	anythingLLMURL := os.Getenv("ANYTHINGLLM_URL")
	anythingLLMKey := os.Getenv("ANYTHINGLLM_API_KEY")
	workspaceSlug := os.Getenv("WORKSPACE_SLUG")

	// Zorunlu değişkenleri kontrol et
	if telegramToken == "" {
		log.Fatal("❌ TELEGRAM_BOT_TOKEN ortam değişkeni ayarlanmamış!")
	}
	if anythingLLMURL == "" {
		log.Fatal("❌ ANYTHINGLLM_URL ortam değişkeni ayarlanmamış!")
	}
	if anythingLLMKey == "" {
		log.Fatal("❌ ANYTHINGLLM_API_KEY ortam değişkeni ayarlanmamış!")
	}
	if workspaceSlug == "" {
		log.Fatal("❌ WORKSPACE_SLUG ortam değişkeni ayarlanmamış!")
	}

	log.Println("🚀 Telegram AnythingLLM Bot başlatılıyor...")

	// AnythingLLM istemcisini oluştur
	llmClient := NewAnythingLLMClient(anythingLLMURL, anythingLLMKey, workspaceSlug)

	// Bağlantıyı test et
	log.Println("🔍 AnythingLLM bağlantısı test ediliyor...")
	if err := llmClient.TestConnection(); err != nil {
		log.Fatalf("❌ AnythingLLM bağlantı hatası: %v", err)
	}
	log.Println("✅ AnythingLLM bağlantısı başarılı!")

	// Telegram botunu başlat
	bot, err := tgbotapi.NewBotAPI(telegramToken)
	if err != nil {
		log.Fatalf("❌ Telegram bot hatası: %v", err)
	}

	bot.Debug = false
	log.Printf("✅ Bot @%s olarak başlatıldı", bot.Self.UserName)

	// Güncellemeleri al
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	log.Println("✅ Bot çalışıyor! Mesajlar bekleniyor...")

	// Mesajları işle
	for update := range updates {
		if update.Message == nil {
			continue
		}

		// Sadece metin mesajlarını işle
		if update.Message.Text == "" {
			continue
		}

		userName := update.Message.From.UserName
		if userName == "" {
			userName = update.Message.From.FirstName
		}

		log.Printf("📨 [%s] %s", userName, update.Message.Text)

		// /start komutu
		if update.Message.Text == "/start" {
			welcomeMsg := fmt.Sprintf(
				"👋 Merhaba %s!\n\n"+
					"Ben AnythingLLM destekli bir yapay zeka botuyum.\n"+
					"Bana istediğiniz soruyu sorabilirsiniz!\n\n"+
					"🤖 Workspace: %s",
				update.Message.From.FirstName,
				workspaceSlug,
			)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, welcomeMsg)
			bot.Send(msg)
			continue
		}

		// /help komutu
		if update.Message.Text == "/help" {
			helpMsg := "📚 Kullanılabilir Komutlar:\n\n" +
				"/start - Botu başlat\n" +
				"/help - Yardım mesajı\n\n" +
				"Bunun dışında bana doğrudan soru sorabilirsiniz!"
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, helpMsg)
			bot.Send(msg)
			continue
		}

		// "Yazıyor..." göster
		typingAction := tgbotapi.NewChatAction(update.Message.Chat.ID, tgbotapi.ChatTyping)
		bot.Send(typingAction)

		// AnythingLLM'e mesaj gönder
		response, err := llmClient.SendMessage(update.Message.Text)
		if err != nil {
			log.Printf("❌ Hata: %v", err)
			errorMsg := "😔 Üzgünüm, bir hata oluştu. Lütfen daha sonra tekrar deneyin."
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, errorMsg)
			bot.Send(msg)
			continue
		}

		log.Printf("✅ Yanıt gönderildi: %d karakter", len(response))

		// Yanıtı gönder
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, response)
		msg.ReplyToMessageID = update.Message.MessageID
		bot.Send(msg)
	}
}
