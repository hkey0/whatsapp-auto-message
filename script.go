package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"model/model"
	"os"
	"time"

	whatsapp "github.com/Rhymen/go-whatsapp"
	qrcodeTerminal "github.com/mdp/qrterminal/v3"
)

func main() {
	wac, _ := whatsapp.NewConn(1 * time.Second)
	wac.SetClientVersion(2, 2123, 7)

	conf := read_config()

	// login with qr
	login(wac)
	if conf.Type == 0 {
		send_messages(wac, conf.Numbers, conf.Type, "", conf.Message, conf.Sleep)
	}
	if conf.Type == 1 {
		send_messages(wac, conf.Numbers, conf.Type, conf.Image, conf.Message, conf.Sleep)
	}

}

func send_messages(wac *whatsapp.Conn, file_name string, tip int, image_path string, message string, sleep_time int) {
	f, _ := os.Open(file_name)

	defer f.Close()

	scanner := bufio.NewScanner(f)

	if tip == 0 {
		for scanner.Scan() {

			send_text(wac, scanner.Text(), message)
			<-time.After(time.Duration(sleep_time) * time.Second)
		}
	}
	if tip == 1 {
		for scanner.Scan() {
			send_image(wac, scanner.Text(), image_path, message)
			<-time.After(time.Duration(sleep_time) * time.Second)
		}
	}
}

func read_config() model.Config {
	jsonFile, _ := os.Open("config.json")
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var config model.Config

	json.Unmarshal(byteValue, &config)
	return config
}

func login(wac *whatsapp.Conn) {
	qrChan := make(chan string)
	go func() {
		config := qrcodeTerminal.Config{
			Level:     qrcodeTerminal.L,
			Writer:    os.Stdout,
			BlackChar: qrcodeTerminal.BLACK,
			WhiteChar: qrcodeTerminal.WHITE,
			QuietZone: 1,
		}
		qrcodeTerminal.GenerateWithConfig(<-qrChan, config)
	}()

	_, err := wac.Login(qrChan)
	if err != nil {
		fmt.Print("err")
	}

}

func send_image(wac *whatsapp.Conn, phone_number string, image_file string, message string) {
	img, _ := os.Open(image_file) // hkey
	msg := whatsapp.ImageMessage{
		Info: whatsapp.MessageInfo{
			RemoteJid: phone_number + "@s.whatsapp.net",
		},
		Type:    "image/jpeg",
		Caption: message,
		Content: img,
	}

	msgId, err := wac.Send(msg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Mesaj gönderirken hata oluştu:  %v", err)
		os.Exit(1)
	} else {
		fmt.Println(phone_number + " - Mesaj gönderildi, mesaj ID'si:  " + msgId)
	}
}

func send_text(wac *whatsapp.Conn, phone_number string, message string) {
	text := whatsapp.TextMessage{
		Info: whatsapp.MessageInfo{
			RemoteJid: phone_number + "@s.whatsapp.net",
		},
		Text: message,
	}

	msgId, err := wac.Send(text)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Mesaj gönderirken hata oluştu:  %v", err)
		os.Exit(1)
	} else {
		fmt.Println(phone_number + " - Mesaj gönderildi, mesaj ID'si:  " + msgId)
	}
}
