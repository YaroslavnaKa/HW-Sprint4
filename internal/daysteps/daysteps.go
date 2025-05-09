package daysteps

import (
	"fmt"
	"github.com/Yandex-Practicum/tracker/internal/spentcalories"
	"log"
	"strconv"
	"strings"
	"time"
)

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

func parsePackage(data string) (int, time.Duration, error) {
	words := strings.Split(data, ",")
	if len(words) != 2 {
		return 0, 0, fmt.Errorf("invalid data format: expected 2 parts separated by comma, got %d", len(words))
	}
	step, err := strconv.Atoi(words[0])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid steps value %d:%w", step, err)
	}
	if step <= 0 {
		return 0, 0, fmt.Errorf("invalid steps value '%d':%w", step, err)
	}

	duration, err := time.ParseDuration(words[1])
	if err != nil {
		return 0, 0, fmt.Errorf("converting string to duration: %w", err)

	}
	if duration <= 0 {
		return 0, 0, fmt.Errorf("invalid minutes value '%d':%w", duration, err)

	}
	return step, duration, nil

}

func DayActionInfo(data string, weight, height float64) string {
	steps, duration, err := parsePackage(data)
	if err != nil {
		log.Println("data processing error:", err)
		return ""
	}
	if steps <= 0 {
		log.Println("data processing error: steps must be greater than zero")
		return ""
	}

	distKm := float64(steps) * stepLength / float64(mInKm)
	calories, err := spentcalories.WalkingSpentCalories(steps, weight, height, duration)
	if err != nil {
		fmt.Println("Ошибка расчета калорий:", err)
		return ""
	}
	return fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n", steps, distKm, calories)
}
