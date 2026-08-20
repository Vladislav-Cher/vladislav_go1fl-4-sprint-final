package spentcalories

import (
	"fmt"
	"log"
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
	dataInSlice := strings.Split(data, ",")
	if len(dataInSlice) != 3 {
		return 0, "", time.Duration(0), fmt.Errorf("требуется 3 элемента данных (получено %d)", len(dataInSlice))
	}

	//Cделал раздельно по две проверки для каждой переменной,
	//чтобы в случае нулевых значений не возвращалсь нули с "nil" ошибкой.
	steps, err := strconv.Atoi(dataInSlice[0])
	if err != nil {
		return 0, "", time.Duration(0), err
	}
	if steps <= 0 {
		return 0, "", time.Duration(0), fmt.Errorf("некорректное количество шагов (получено %d)", steps)
	}

	duration, err := time.ParseDuration(dataInSlice[2])
	if err != nil {
		return 0, "", time.Duration(0), err
	}
	if duration <= 0 {
		return 0, "", time.Duration(0), fmt.Errorf("продолжительность тренировки = 0")
	}

	return steps, dataInSlice[1], duration, nil
}

func distance(steps int, height float64) float64 {
	// TODO: реализовать функцию
	stepLength := height * stepLengthCoefficient

	distInM := float64(steps) * stepLength
	distInKm := distInM / float64(mInKm)

	return distInKm
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}

	return distance(steps, height) / duration.Hours()
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	// TODO: реализовать функцию
	steps, activity, duration, err := parseTraining(data)
	if err != nil {
		log.Println(err)
		return "", nil
	}
	switch activity {
	case "Бег":
		dist := distance(steps, height)
		meanS := meanSpeed(steps, height, duration)
		calories, err := RunningSpentCalories(steps, weight, height, duration)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf(
			"Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км."+
				"\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n",
			activity,
			duration.Hours(),
			dist,
			meanS,
			calories,
		), nil
	case "Ходьба":
		dist := distance(steps, height)
		meanS := meanSpeed(steps, height, duration)
		calories, err := WalkingSpentCalories(steps, weight, height, duration)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf(
			"Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км."+
				"\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n",
			activity,
			duration.Hours(),
			dist,
			meanS,
			calories,
		), nil
	default:
		return "", fmt.Errorf("неизвестный тип тренировки")
	}

}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 {
		return 0, fmt.Errorf("некорректное количество шагов (получено %d)", steps)
	}

	if weight <= 0 {
		return 0, fmt.Errorf("некорректный вес (получено %.1f кг)", weight)
	}

	if height <= 0 {
		return 0, fmt.Errorf("некорректный рост (получено %.2f м)", height)
	}

	if duration <= 0 {
		return 0, fmt.Errorf("продолжительность тренировки = 0")
	}

	meanS := meanSpeed(steps, height, duration)
	durationInMinutes := duration.Minutes()

	return (weight * meanS * durationInMinutes) / minInH, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 {
		return 0, fmt.Errorf("некорректное количество шагов (получено %d)", steps)
	}

	if weight <= 0 {
		return 0, fmt.Errorf("некорректный вес (получено %.1f кг)", weight)
	}

	if height <= 0 {
		return 0, fmt.Errorf("некорректный рост (получено %.2f м)", height)
	}

	if duration <= 0 {
		return 0, fmt.Errorf("продолжительность тренировки = 0")
	}

	meanS := meanSpeed(steps, height, duration)
	durationInMinutes := duration.Minutes()

	caloriesIfRunning := (weight * meanS * durationInMinutes) / minInH

	return caloriesIfRunning * walkingCaloriesCoefficient, nil
}
