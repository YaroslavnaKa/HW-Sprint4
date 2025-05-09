package spentcalories

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	lenStep                    = 0.65 // средняя длина шага.
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

func parseTraining(data string) (int, string, time.Duration, error) {
	words := strings.Split(data, ",")

	if len(words) != 3 {
		return 0, "", 0, fmt.Errorf("invalid data format: expected 3 parts, got %d", len(words))
	}
	step, err := strconv.Atoi(words[0])
	if err != nil {
		return 0, "", 0, fmt.Errorf("convertation steps to int is faild '%d':%w", step, err)
	}
	if step <= 0 {
		return 0, "", 0, fmt.Errorf("invalid steps value '%d':%w", step, err)
	}

	duration, err := time.ParseDuration(words[2])
	if err != nil {
		return 0, "", 0, fmt.Errorf("convertation minutes to int is faild '%d':%w", step, err)
	}
	if duration <= 0 {
		return 0, "", 0, errors.New("invalid data format: minutes is zero")
	}

	return step, words[1], duration, nil
}

func distance(steps int, height float64) float64 {
	lonSt := height * stepLengthCoefficient
	dist := lonSt * float64(steps) / mInKm
	return dist
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}
	speed := distance(steps, height) / duration.Hours()
	return speed
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, trainingType, duration, err := parseTraining(data)
	if err != nil {
		return "", fmt.Errorf("failed to parse data: %w", err)
	}
	if weight <= 0 {
		return "", fmt.Errorf("weight must be > 0, get %.2f", weight)
	}
	if height <= 0 {
		return "", fmt.Errorf("height must be > 0, get %.2f", height)
	}
	var (
		distValue     float64
		speedValue    float64
		caloriesValue float64
		calorieErr    error
	)

	// Рассчитываем общие для всех видов тренировок значения
	distValue = distance(steps, height)
	speedValue = meanSpeed(steps, height, duration)
	switch trainingType {
	case "Бег":
		caloriesValue, calorieErr = RunningSpentCalories(steps, weight, height, duration)
	case "Ходьба":
		caloriesValue, calorieErr = WalkingSpentCalories(steps, weight, height, duration)
	default:
		return "", fmt.Errorf("неизвестный тип тренировки: %s", trainingType)
	}
	if calorieErr != nil {
		return "", fmt.Errorf("ошибка расчета калорий для типа '%s': %w", trainingType, calorieErr)
	}
	resultString := fmt.Sprintf("Тип тренировки: %s\n"+
		"Длительность: %.2f ч.\n"+
		"Дистанция: %.2f км.\n"+
		"Скорость: %.2f км/ч\n"+
		"Сожгли калорий: %.2f\n",
		trainingType,
		duration.Hours(),
		distValue,
		speedValue,
		caloriesValue)

	return resultString, nil

}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 {
		return 0, fmt.Errorf("invalid data format: steps must be greater than zero")
	}
	if weight <= 0 {
		return 0, fmt.Errorf("invalid data format: weight must be greater than zero")
	}
	if height <= 0 {
		return 0, fmt.Errorf("invalid data format: height must be greater than zero")
	}
	if duration <= 0 {
		return 0, fmt.Errorf("invalid data format: duration must be greater than zero")
	}
	meanSpeedValue := meanSpeed(steps, height, duration)
	durationInMinutes := duration.Minutes()
	callor := (weight * meanSpeedValue * durationInMinutes) / minInH
	return callor, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 {
		return 0, fmt.Errorf("invalid data format: steps must be greater than zero")
	}
	if weight <= 0 {
		return 0, fmt.Errorf("invalid data format: weight must be greater than zero")
	}
	if height <= 0 {
		return 0, fmt.Errorf("invalid data format: height must be greater than zero")
	}
	if duration <= 0 {
		return 0, fmt.Errorf("invalid data format: duration must be greater than zero")
	}
	meanSpeedValue := meanSpeed(steps, height, duration)
	durationInMinutes := duration.Minutes()
	baseCalories := (weight * meanSpeedValue * durationInMinutes) / minInH
	callorValue := baseCalories * walkingCaloriesCoefficient
	return callorValue, nil
}
