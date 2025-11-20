# Telegram AnythingLLM Bot

AnythingLLM workspace'inizi Telegram üzerinden kullanmanızı sağlayan Go tabanlı bir bot.

## 🚀 Özellikler

- ✅ AnythingLLM workspace'i ile doğrudan entegrasyon
- ✅ Telegram üzerinden yapay zeka ile sohbet
- ✅ Kolay kurulum ve yapılandırma
- ✅ Hata yönetimi ve loglama
- ✅ Systemd servisi desteği

## 📋 Gereksinimler

- Ubuntu 22.04 (veya herhangi bir Linux dağıtımı)
- Go 1.21 veya üzeri
- AnythingLLM kurulu ve çalışır durumda
- Telegram Bot Token (BotFather'dan)
- AnythingLLM API Key

## 🔧 Kurulum

### 1. Go Kurulumu

```bash
# Go'nun kurulu olup olmadığını kontrol edin
go version

# Eğer kurulu değilse:
wget https://go.dev/dl/go1.21.5.linux-amd64.tar.gz
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.21.5.linux-amd64.tar.gz

# PATH'e ekleyin
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# Kontrol edin
go version


2. Projeyi Klonlayın
bash
Copy
cd ~
git clone https://github.com/0xmtnslk/telegram-anythingllm-bot.git
cd telegram-anythingllm-bot
3. Bağımlılıkları Yükleyin
bash
Copy
go mod download
4. Telegram Bot Oluşturun
Telegram'da @BotFather ile konuşun
/newbot komutunu gönderin
Bot için bir isim verin (örn: "My AnythingLLM Bot")
Bot için bir kullanıcı adı verin (örn: "my_anythingllm_bot")
Size verilen token'ı kaydedin
5. AnythingLLM API Key Alın
AnythingLLM arayüzünüze gidin: http://116.203.233.131:3001
Sağ üst köşeden ayarlara gidin
"API Keys" bölümüne gidin
"Generate New API Key" butonuna tıklayın
Oluşturulan API key'i kaydedin
6. Workspace Slug'ını Kontrol Edin
Workspace slug'ınız genellikle workspace adınızın küçük harfli halidir:

Workspace adı: TelegramBot
Workspace slug: telegrambot veya TelegramBot
AnythingLLM arayüzünde workspace'inize girdiğinizde URL'de görebilirsiniz:

http://116.203.233.131:3001/workspace/TelegramBot
                                      ^^^^^^^^^^^
                                      Bu kısım slug
7. Yapılandırma Dosyasını Oluşturun
bash
Copy
# .env.example dosyasını kopyalayın
cp .env.example .env

# .env dosyasını düzenleyin
nano .env
.env dosyasını aşağıdaki gibi doldurun:

bash
Copy
TELEGRAM_BOT_TOKEN=123456789:ABCdefGHIjklMNOpqrsTUVwxyz
ANYTHINGLLM_URL=http://116.203.233.131:3001
ANYTHINGLLM_API_KEY=ANYLLM-XXXXXXXXXXXXXXXXXXXXXXXX
WORKSPACE_SLUG=TelegramBot
Kaydedin ve çıkın (Ctrl+X, Y, Enter)

8. Botu Test Edin
bash
Copy
# Botu çalıştırın
go run main.go
Çıktıda şunları görmelisiniz:

🚀 Telegram AnythingLLM Bot başlatılıyor...
🔍 AnythingLLM bağlantısı test ediliyor...
✅ AnythingLLM bağlantısı başarılı!
✅ Bot @your_bot_username olarak başlatıldı
✅ Bot çalışıyor! Mesajlar bekleniyor...
Telegram'da botunuza /start göndererek test edin!

🔄 Systemd Servisi Olarak Çalıştırma
Botu arka planda sürekli çalışır halde tutmak için:

1. Botu Derleyin
bash
Copy
cd ~/telegram-anythingllm-bot
go build -o telegram-bot main.go
2. Systemd Servis Dosyası Oluşturun
bash
Copy
sudo nano /etc/systemd/system/telegram-anythingllm-bot.service
Aşağıdaki içeriği yapıştırın (USER_NAME yerine kendi kullanıcı adınızı yazın):

ini
Copy
[Unit]
Description=Telegram AnythingLLM Bot
After=network.target

[Service]
Type=simple
User=USER_NAME
WorkingDirectory=/home/USER_NAME/telegram-anythingllm-bot
ExecStart=/home/USER_NAME/telegram-anythingllm-bot/telegram-bot
Restart=always
RestartSec=10
StandardOutput=journal
StandardError=journal

[Install]
WantedBy=multi-user.target
3. Servisi Başlatın
bash
Copy
# Systemd'yi yeniden yükle
sudo systemctl daemon-reload

# Servisi etkinleştir (sistem başlangıcında otomatik başlat)
sudo systemctl enable telegram-anythingllm-bot

# Servisi başlat
sudo systemctl start telegram-anythingllm-bot

# Durumu kontrol et
sudo systemctl status telegram-anythingllm-bot
4. Servis Komutları
bash
Copy
# Servisi durdur
sudo systemctl stop telegram-anythingllm-bot

# Servisi yeniden başlat
sudo systemctl restart telegram-anythingllm-bot

# Logları görüntüle
sudo journalctl -u telegram-anythingllm-bot -f

# Son 100 log satırını göster
sudo journalctl -u telegram-anythingllm-bot -n 100
📱 Kullanım
Telegram'da botunuza gidin ve şu komutları kullanın:

/start - Botu başlat ve hoş geldin mesajı al
/help - Yardım mesajını görüntüle
Herhangi bir mesaj - AnythingLLM ile sohbet et
