package controllers

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/redis/go-redis/v9"
	"slices"
	"strings"
)

var RedisClient *redis.Client

func ConnectRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}

func GetUserState(UserID int64) string {
	ctx := context.Background()
	return RedisClient.HGet(ctx, fmt.Sprint(UserID), "state").Val()
}

func SetUserState(UserID int64, State string) {
	ctx := context.Background()
	RedisClient.HSet(ctx, fmt.Sprint(UserID), "state", State)
}

func UpdateFilters(category string, update *tgbotapi.Update, callback *tgbotapi.CallbackConfig) {

	ctx := context.Background()
	filter := update.CallbackQuery.Data
	var new_filters_str string
	current_filters_array := strings.Split(RedisClient.HGetAll(ctx, fmt.Sprint(update.CallbackQuery.From.ID)).Val()[category], ",")

	if !slices.Contains(current_filters_array, filter) {
		current_filters_str := strings.Join(current_filters_array, ",")
		current_filters_array = RemoveStrFromArray(current_filters_array, "")
		new_filters_str = strings.Join([]string{current_filters_str, filter}, ",")
		callback.Text = "Добавлен фильтр: " + filter
	} else {
		current_filters_array = RemoveStrFromArray(current_filters_array, filter)
		new_filters_str = strings.Join(current_filters_array, ",")
		callback.Text = "Удалён фильтр: " + filter
	}

	RedisClient.HSet(ctx, fmt.Sprint(update.CallbackQuery.From.ID), category, new_filters_str)
}

func RemoveStrFromArray(s []string, r string) []string {
	for i, v := range s {
		if v == r {
			return append(s[:i], s[i+1:]...)
		}
	}
	return s
}
