package handlers

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/mattermost/mattermost-server/v6/model"
	"github.com/tarantool/go-tarantool"
)

func generateUUID() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%x", rand.Uint64())
}

func HandleCommand(post *model.Post, client *model.Client4, tnt *tarantool.Connection) {
	args := strings.Fields(post.Message)
	if len(args) < 2 {
		sendResponse(client, post.ChannelId, "Usage: /vote [create|vote|results|end|delete]")
		return
	}

	switch args[1] {
	case "create":
		createPoll(client, tnt, post, args[2:])
	case "vote":
		processVote(client, tnt, post, args[2:])
	case "results":
		showResults(client, tnt, post, args[2:])
	case "end":
		endPoll(client, tnt, post, args[2:])
	case "delete":
		deletePoll(client, tnt, post, args[2:])
	default:
		sendResponse(client, post.ChannelId, "Неизвестная команда")
	}
}

func createPoll(client *model.Client4, tnt *tarantool.Connection, post *model.Post, args []string) {
	if len(args) < 2 {
		sendResponse(client, post.ChannelId, "Usage: /vote create <Вопрос?> <Вариант1,Вариант2>")
		return
	}

	pollID := generateUUID()
	options := strings.Split(args[1], ",")
	optionsMap := make(map[string]int)
	for i := range options {
		optionsMap[strconv.Itoa(i+1)] = 0
	}

	_, err := tnt.Insert("polls", []interface{}{
		pollID,
		args[0],
		optionsMap,
		post.UserId,
		true,
	})

	if err != nil {
		log.Printf("Ошибка создания голосования: %v", err)
		sendResponse(client, post.ChannelId, "Ошибка создания голосования")
		return
	}

	response := fmt.Sprintf("**Новое голосование**\nID: `%s`\nВопрос: %s\nВарианты:\n", pollID, args[0])
	for i, opt := range options {
		response += fmt.Sprintf("%d. %s\n", i+1, strings.TrimSpace(opt))
	}
	sendResponse(client, post.ChannelId, response)
}

func processVote(client *model.Client4, tnt *tarantool.Connection, post *model.Post, args []string) {
	if len(args) < 2 {
		sendResponse(client, post.ChannelId, "Usage: /vote vote <ID> <номер_варианта>")
		return
	}

	pollID := args[0]
	choice := args[1]

	res, err := tnt.Select("polls", "primary", 0, 1, tarantool.IterEq, []interface{}{pollID})
	if err != nil || len(res.Data) == 0 {
		sendResponse(client, post.ChannelId, "Голосование не найдено")
		return
	}

	poll := res.Data[0].([]interface{})
	if !poll[4].(bool) {
		sendResponse(client, post.ChannelId, "Голосование завершено")
		return
	}

	options := poll[2].(map[interface{}]interface{})
	if _, exists := options[choice]; !exists {
		sendResponse(client, post.ChannelId, "Неверный вариант")
		return
	}

	_, err = tnt.Update("polls", "primary", []interface{}{pollID}, []interface{}{
		[]interface{}{"=", 2, map[string]interface{}{
			choice: options[choice].(int) + 1,
		}},
	})

	sendResponse(client, post.ChannelId, fmt.Sprintf("Голос за вариант %s учтен!", choice))
}

func showResults(client *model.Client4, tnt *tarantool.Connection, post *model.Post, args []string) {
	if len(args) < 1 {
		sendResponse(client, post.ChannelId, "Usage: /vote results <ID>")
		return
	}

	pollID := args[0]
	res, err := tnt.Select("polls", "primary", 0, 1, tarantool.IterEq, []interface{}{pollID})
	if err != nil || len(res.Data) == 0 {
		sendResponse(client, post.ChannelId, "Голосование не найдено")
		return
	}

	poll := res.Data[0].([]interface{})
	response := fmt.Sprintf("**Результаты голосования**\n %s\n", poll[1])
	for k, v := range poll[2].(map[interface{}]interface{}) {
		response += fmt.Sprintf("- %s: %d голосов\n", k, v)
	}
	sendResponse(client, post.ChannelId, response)
}

func endPoll(client *model.Client4, tnt *tarantool.Connection, post *model.Post, args []string) {
	if len(args) < 1 {
		sendResponse(client, post.ChannelId, "Usage: /vote end <ID>")
		return
	}

	pollID := args[0]
	res, err := tnt.Select("polls", "primary", 0, 1, tarantool.IterEq, []interface{}{pollID})
	if err != nil || len(res.Data) == 0 {
		sendResponse(client, post.ChannelId, "Голосование не найдено")
		return
	}

	poll := res.Data[0].([]interface{})
	if poll[3].(string) != post.UserId {
		sendResponse(client, post.ChannelId, "Только создатель может завершить голосование")
		return
	}

	_, err = tnt.Update("polls", "primary", []interface{}{pollID}, []interface{}{
		[]interface{}{"=", 4, false},
	})
	sendResponse(client, post.ChannelId, "Голосование завершено!")
}

func deletePoll(client *model.Client4, tnt *tarantool.Connection, post *model.Post, args []string) {
	if len(args) < 1 {
		sendResponse(client, post.ChannelId, "Usage: /vote delete <ID>")
		return
	}

	pollID := args[0]
	_, err := tnt.Delete("polls", "primary", []interface{}{pollID})
	if err != nil {
		sendResponse(client, post.ChannelId, "Ошибка при удалении голосования")
		return
	}
	sendResponse(client, post.ChannelId, "Голосование удалено!")
}

func sendResponse(client *model.Client4, channelID, message string) {
	post := &model.Post{
		ChannelId: channelID,
		Message:   message,
	}
	_, _, err := client.CreatePost(post) // Исправлено на 3 возвращаемых значения
	if err != nil {
		log.Printf("Ошибка отправки сообщения: %v", err)
	}
}
