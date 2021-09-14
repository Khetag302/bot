package telegram

import (
	"TelegramBot"
	"fmt"
	owm "github.com/briandowns/openweathermap"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"

	"time"
)

const (
	commandStart = "start"
	weather      = "weather"
	weatherCycle = "weatherCycle"
)

func (b *Bot) handleCommand(message *tgbotapi.Message, city string) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, "Неизвестная команда")

	switch message.Command() {
	case commandStart:
		msg.Text = "Поехали"
		_, err := b.bot.Send(msg)
		return err
	case weather:
		err := b.handleWeatherMess(message, city)
		return err
	case weatherCycle:
		err:= b.handleWeatherCycleMess(message, city)
		return err

	default:
		_, err := b.bot.Send(msg)
		return err
	}
	return nil
}

func (b *Bot) handleWeatherMess(message *tgbotapi.Message, city string) error {
	weather, err := owm.NewCurrent("c", "RU", TelegramBot.WeatherToken)
	if err != nil {
		return err
	}

	err = weather.CurrentByName(city)
	if err != nil {
		return err
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("Город: %s\n"+
		"Погода: %s\n"+
		"Влажность: %d%%\n"+
		"Температура: %d°C\n"+
		"Ощущается как:%d°C",
		city,
		weather.Weather[0].Description,
		int8(weather.Main.Humidity),
		int8(weather.Main.Temp),
		int8(weather.Main.FeelsLike)))

	mess, err := b.bot.Send(msg)
	if err != nil {
		return err
	}
	time.Sleep(5 * time.Second)
	_, err = b.bot.DeleteMessage(tgbotapi.DeleteMessageConfig{ChatID: message.Chat.ID, MessageID: message.MessageID})
	if err != nil {
		return err
	}
	_, err = b.bot.DeleteMessage(tgbotapi.DeleteMessageConfig{ChatID: message.Chat.ID, MessageID: mess.MessageID})
	if err != nil {
		return err
	}
	return nil
}

func (b *Bot) handleWeatherCycleMess(message *tgbotapi.Message, city string) error {
	for {
		weather, err := owm.NewCurrent("c", "RU", TelegramBot.WeatherToken)
		if err != nil {
			return err
		}

		err = weather.CurrentByName(city)
		if err != nil {
			return err
		}

		msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("Город: %s\n"+
			"Погода: %s\n"+
			"Влажность: %d%%\n"+
			"Температура: %d°C\n"+
			"Ощущается как:%d°C",
			city,
			weather.Weather[0].Description,
			int8(weather.Main.Humidity),
			int8(weather.Main.Temp),
			int8(weather.Main.FeelsLike)))

		 b.bot.Send(msg)
	}
}


func (b *Bot) handleMessage(message *tgbotapi.Message) {
	log.Printf("[%s] %s", message.From.UserName, message.Text)

	msg := tgbotapi.NewMessage(message.Chat.ID, message.Text)

	b.bot.Send(msg)
}
