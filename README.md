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
