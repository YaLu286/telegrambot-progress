package controllers

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"slices"
	"strings"
	"telegrambot/progress/models"
)

func GetUserState(UserID int64) string {
	ctx := context.Background()
	return models.RedisClient.HGet(ctx, fmt.Sprint(UserID), "state").Val()
}

func SetUserState(UserID int64, State string) {
	ctx := context.Background()
	models.RedisClient.HSet(ctx, fmt.Sprint(UserID), "state", State)
}

func SetUserLocation(UserID int64, Location string) {
	ctx := context.Background()
	models.RedisClient.HSet(ctx, fmt.Sprint(UserID), "location", Location)
}

func GetUserLocation(UserID int64) string {
	ctx := context.Background()
	return models.RedisClient.HGet(ctx, fmt.Sprint(UserID), "location").Val()
}

func SetAdminMode(UserID int64, Mode string) {
	ctx := context.Background()
	models.RedisClient.HSet(ctx, fmt.Sprint(UserID), "admin-mode", Mode)
}

func GetAdminMode(UserID int64) string {
	ctx := context.Background()
	return models.RedisClient.HGet(ctx, fmt.Sprint(UserID), "admin-mode").Val()
}

func CleanUserFilters(UserID int64) {
	ctx := context.Background()
	models.RedisClient.HDel(ctx, fmt.Sprint(UserID), "style")
	models.RedisClient.HDel(ctx, fmt.Sprint(UserID), "brewery")
}

func UpdateFilters(category string, update *tgbotapi.Update, callback *tgbotapi.CallbackConfig) {

	ctx := context.Background()
	filter := update.CallbackQuery.Data
	var new_filters_str string
	current_filters_array := strings.Split(models.RedisClient.HGetAll(ctx, fmt.Sprint(update.CallbackQuery.From.ID)).Val()[category], ",")

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

	models.RedisClient.HSet(ctx, fmt.Sprint(update.CallbackQuery.From.ID), category, new_filters_str)
}

func RemoveStrFromArray(s []string, r string) []string {
	for i, v := range s {
		if v == r {
			return append(s[:i], s[i+1:]...)
		}
	}
	return s
}
