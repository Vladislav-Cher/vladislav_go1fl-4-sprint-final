package daysteps

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/spentcalories"
)

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

func parsePackage(data string) (int, time.Duration, error) {
	dataInSlice := strings.Split(data, ",")
	if len(dataInSlice) != 2 {
		return 0, time.Duration(0), fmt.Errorf("требуется 2 элемента данных (получено %d)", len(dataInSlice))
	}
	//Cделал раздельно по две проверки для каждой переменной,
	//чтобы в случае нулевых значений не возвращалсь нули с "nil" ошибкой.
	steps, err := strconv.Atoi(dataInSlice[0])
	if err != nil {
		return 0, time.Duration(0), err
	}
	if steps <= 0 {
		return 0, time.Duration(0), fmt.Errorf("некорректное количество шагов (получено %d)", steps)
	}

	duration, err := time.ParseDuration(dataInSlice[1])
	if err != nil {
		return 0, time.Duration(0), err
	}
	if duration <= 0 {
		return 0, time.Duration(0), fmt.Errorf("некорректная продолжительность тренировки")
	}

	return steps, duration, nil
}

func DayActionInfo(data string, weight, height float64) string {
	steps, duration, err := parsePackage(data)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	distInM := float64(steps) * stepLength
	distInKm := distInM / mInKm

	calories, err := spentcalories.WalkingSpentCalories(steps, weight, height, duration)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	return fmt.Sprintf(
		"Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n",
		steps,
		distInKm,
		calories,
	)
}
